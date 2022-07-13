// Copyright (C) 2019-2022, Axia Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/axiacoin/axia-network-v2-runner/api"
	"github.com/axiacoin/axia-network-v2-runner/local"
	"github.com/axiacoin/axia-network-v2-runner/network"
	"github.com/axiacoin/axia-network-v2-runner/network/node"
	"github.com/axiacoin/axia-network-v2-runner/pkg/color"
	"github.com/axiacoin/axia-network-v2-runner/rpcpb"
	"github.com/axiacoin/axia-network-v2/utils/constants"
	"github.com/axiacoin/axia-network-v2/utils/logging"
)

type localNetwork struct {
	logger logging.Logger

	binPath string
	cfg     network.Config

	nw network.Network

	nodeNames []string
	nodes     map[string]node.Node
	nodeInfos map[string]*rpcpb.NodeInfo

	apiClis map[string]api.Client

	readyc          chan struct{} // closed when local network is ready/healthy
	readycCloseOnce sync.Once

	stopc chan struct{}
	donec chan struct{}
	errc  chan error

	stopOnce sync.Once
}

func newNetwork(execPath string, rootDataDir string, whitelistedAllychains string, logLevel string) (*localNetwork, error) {
	lcfg := logging.DefaultConfig
	// if err != nil {
	// 	return nil, err
	// }
	lcfg.Directory = rootDataDir
	logFactory := logging.NewFactory(lcfg)
	logger, _ := logFactory.Make("main")
	// if err != nil {
	// 	return nil, err
	// }

	if logLevel == "" {
		logLevel = "INFO"
	}

	nodeInfos := make(map[string]*rpcpb.NodeInfo)
	cfg := local.NewDefaultConfig(execPath)
	nodeNames := make([]string, len(cfg.NodeConfigs))
	for i := range cfg.NodeConfigs {
		nodeName := fmt.Sprintf("node%d", i+1)
		logDir := filepath.Join(rootDataDir, nodeName, "log")
		dbDir := filepath.Join(rootDataDir, nodeName, "db-dir")

		nodeNames[i] = nodeName
		cfg.NodeConfigs[i].Name = nodeName

		// need to whitelist allychain ID to create custom VM chain
		// ref. vms/platformvm/createChain
		cfg.NodeConfigs[i].ConfigFile = fmt.Sprintf(`{
	"network-peer-list-gossip-frequency":"250ms",
	"network-max-reconnect-delay":"1s",
	"public-ip":"127.0.0.1",
	"health-check-frequency":"2s",
	"api-admin-enabled":true,
	"api-ipcs-enabled":true,
	"index-enabled":true,
	"log-display-level":"INFO",
	"log-level":"%s",
	"log-dir":"%s",
	"db-dir":"%s",
	"whitelisted-allychains":"%s"
}`,
			strings.ToUpper(logLevel),
			logDir,
			dbDir,
			whitelistedAllychains,
		)
		cfg.NodeConfigs[i].ImplSpecificConfig = json.RawMessage(fmt.Sprintf(`{"binaryPath":"%s","redirectStdout":true,"redirectStderr":true}`, execPath))

		nodeInfos[nodeName] = &rpcpb.NodeInfo{
			Name:                  nodeName,
			ExecPath:              execPath,
			Uri:                   "",
			Id:                    "",
			LogDir:                logDir,
			DbDir:                 dbDir,
			WhitelistedAllychains: whitelistedAllychains,
			Config:                []byte(cfg.NodeConfigs[i].ConfigFile),
		}
	}

	return &localNetwork{
		logger: logger,

		binPath: execPath,
		cfg:     cfg,

		nodeNames: nodeNames,
		nodeInfos: nodeInfos,
		apiClis:   make(map[string]api.Client),

		readyc: make(chan struct{}),

		stopc: make(chan struct{}),
		donec: make(chan struct{}),
		errc:  make(chan error, 1),
	}, nil
}

func (lc *localNetwork) start() {
	defer func() {
		close(lc.donec)
	}()

	color.Outf("{{blue}}{{bold}}create and run local network{{/}}\n")
	nw, err := local.NewNetwork(lc.logger, lc.cfg, os.TempDir())
	if err != nil {
		lc.errc <- err
		return
	}
	lc.nw = nw

	if err := lc.waitForHealthy(); err != nil {
		lc.errc <- err
		return
	}
}

const healthyWait = 2 * time.Minute

var errAborted = errors.New("aborted")

func (lc *localNetwork) waitForHealthy() error {
	color.Outf("{{blue}}{{bold}}waiting for all nodes to report healthy...{{/}}\n")

	ctx, cancel := context.WithTimeout(context.Background(), healthyWait)
	defer cancel()
	hc := lc.nw.Healthy(ctx)
	select {
	case <-lc.stopc:
		return errAborted
	case <-ctx.Done():
		return ctx.Err()
	case err := <-hc:
		if err != nil {
			return err
		}
	}

	nodes, err := lc.nw.GetAllNodes()
	if err != nil {
		return err
	}
	lc.nodes = nodes

	for name, node := range nodes {
		uri := fmt.Sprintf("http://%s:%d", node.GetURL(), node.GetAPIPort())
		nodeID := node.GetNodeID().PrefixedString(constants.NodeIDPrefix)

		lc.nodeInfos[name].Uri = uri
		lc.nodeInfos[name].Id = nodeID

		lc.apiClis[name] = node.GetAPIClient()
		color.Outf("{{cyan}}%s: node ID %q, URI %q{{/}}\n", name, nodeID, uri)
	}

	lc.readycCloseOnce.Do(func() {
		close(lc.readyc)
	})
	return nil
}

func (lc *localNetwork) stop() {
	lc.stopOnce.Do(func() {
		close(lc.stopc)
		serr := lc.nw.Stop(context.Background())
		<-lc.donec
		color.Outf("{{red}}{{bold}}terminated network{{/}} (error %v)\n", serr)
	})
}
