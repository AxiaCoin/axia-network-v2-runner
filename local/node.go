package local

import (
	"context"
	"crypto"
	"fmt"
	"net"
	"os/exec"
	"syscall"
	"time"

	"github.com/axiacoin/axia-network-runner/api"
	"github.com/axiacoin/axia-network-runner/network/node"
	"github.com/axiacoin/axia-network-v2/ids"
	"github.com/axiacoin/axia-network-v2/message"
	"github.com/axiacoin/axia-network-v2/network/peer"
	"github.com/axiacoin/axia-network-v2/network/throttling"
	"github.com/axiacoin/axia-network-v2/snow/networking/router"
	"github.com/axiacoin/axia-network-v2/snow/networking/tracker"
	"github.com/axiacoin/axia-network-v2/snow/validators"
	"github.com/axiacoin/axia-network-v2/staking"
	"github.com/axiacoin/axia-network-v2/utils/constants"
	"github.com/axiacoin/axia-network-v2/utils/ips"
	"github.com/axiacoin/axia-network-v2/utils/logging"
	"github.com/axiacoin/axia-network-v2/utils/math/meter"
	"github.com/axiacoin/axia-network-v2/utils/resource"
	"github.com/axiacoin/axia-network-v2/version"
	"github.com/prometheus/client_golang/prometheus"
)

// interface compliance
var (
	_ node.Node   = (*localNode)(nil)
	_ NodeProcess = (*nodeProcessImpl)(nil)
	_ getConnFunc = defaultGetConnFunc
)

type getConnFunc func(context.Context, node.Node) (net.Conn, error)

// NodeProcess as an interface so we can mock running
// Axia binaries in tests
type NodeProcess interface {
	// Start this process
	Start() error
	// Send a SIGTERM to this process
	Stop() error
	// Returns when the process finishes exiting
	Wait() error
}

const (
	peerMsgQueueBufferSize      = 1024
	peerResourceTrackerDuration = 10 * time.Second
)

type nodeProcessImpl struct {
	cmd *exec.Cmd
}

func (p *nodeProcessImpl) Start() error {
	return p.cmd.Start()
}

func (p *nodeProcessImpl) Wait() error {
	return p.cmd.Wait()
}

func (p *nodeProcessImpl) Stop() error {
	return p.cmd.Process.Signal(syscall.SIGTERM)
}

// Gives access to basic node info, and to most axia apis
type localNode struct {
	// Must be unique across all nodes in this network.
	name string
	// [nodeID] is this node's Axia Node ID.
	// Set in network.AddNode
	nodeID ids.NodeID
	// The ID of the network this node exists in
	networkID uint32
	// Allows user to make API calls to this node.
	client api.Client
	// The process running this node.
	process NodeProcess
	// The API port
	apiPort uint16
	// The P2P (staking) port
	p2pPort uint16
	// Returns a connection to this node
	getConnFunc getConnFunc
	// The db dir of the node
	dbDir string
	// The logs dir of the node
	logsDir string
	// The node config
	config node.Config
}

func defaultGetConnFunc(ctx context.Context, node node.Node) (net.Conn, error) {
	dialer := net.Dialer{}
	return dialer.DialContext(ctx, constants.NetworkType, net.JoinHostPort(node.GetURL(), fmt.Sprintf("%d", node.GetP2PPort())))
}

// AttachPeer: see Network
func (node *localNode) AttachPeer(ctx context.Context, router router.InboundHandler) (peer.Peer, error) {
	tlsCert, err := staking.NewTLSCert()
	if err != nil {
		return nil, err
	}
	tlsConfg := peer.TLSConfig(*tlsCert)
	clientUpgrader := peer.NewTLSClientUpgrader(tlsConfg)
	conn, err := node.getConnFunc(ctx, node)
	if err != nil {
		return nil, err
	}
	mc, err := message.NewCreator(
		prometheus.NewRegistry(),
		true,
		"",
		10*time.Second,
	)
	if err != nil {
		return nil, err
	}

	metrics, err := peer.NewMetrics(
		logging.NoLog{},
		"",
		prometheus.NewRegistry(),
	)
	if err != nil {
		return nil, err
	}
	ip := ips.IPPort{
		IP:   net.IPv6zero,
		Port: 0,
	}
	resourceTracker, err := tracker.NewResourceTracker(
		prometheus.NewRegistry(),
		resource.NoUsage,
		meter.ContinuousFactory{},
		peerResourceTrackerDuration,
	)
	if err != nil {
		return nil, err
	}
	config := &peer.Config{
		Metrics:             metrics,
		MessageCreator:      mc,
		Log:                 logging.NoLog{},
		InboundMsgThrottler: throttling.NewNoInboundThrottler(),
		Network: peer.NewTestNetwork(
			mc,
			node.networkID,
			ip,
			version.CurrentApp,
			tlsCert.PrivateKey.(crypto.Signer),
			ids.Set{},
			100,
		),
		Router:               router,
		VersionCompatibility: version.GetCompatibility(node.networkID),
		VersionParser:        version.DefaultApplicationParser,
		MyAllychains:            ids.Set{},
		Beacons:              validators.NewSet(),
		NetworkID:            node.networkID,
		PingFrequency:        constants.DefaultPingFrequency,
		PongTimeout:          constants.DefaultPingPongTimeout,
		MaxClockDifference:   time.Minute,
		ResourceTracker:      resourceTracker,
	}
	_, conn, cert, err := clientUpgrader.Upgrade(conn)
	if err != nil {
		return nil, err
	}

	p := peer.Start(
		config,
		conn,
		cert,
		ids.NodeIDFromCert(tlsCert.Leaf),
		peer.NewBlockingMessageQueue(
			config.Metrics,
			logging.NoLog{},
			peerMsgQueueBufferSize,
		),
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// See node.Node
func (node *localNode) GetName() string {
	return node.name
}

// See node.Node
func (node *localNode) GetNodeID() ids.NodeID {
	return node.nodeID
}

// See node.Node
func (node *localNode) GetAPIClient() api.Client {
	return node.client
}

// See node.Node
func (node *localNode) GetURL() string {
	return "127.0.0.1"
}

// See node.Node
func (node *localNode) GetP2PPort() uint16 {
	return node.p2pPort
}

// See node.Node
func (node *localNode) GetAPIPort() uint16 {
	return node.apiPort
}

// See node.Node
func (node *localNode) GetBinaryPath() string {
	return node.config.BinaryPath
}

// See node.Node
func (node *localNode) GetDbDir() string {
	return node.dbDir
}

// See node.Node
func (node *localNode) GetLogsDir() string {
	return node.logsDir
}

// See node.Node
func (node *localNode) GetConfigFile() string {
	return node.config.ConfigFile
}
