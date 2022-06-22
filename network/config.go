package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/axiacoin/axia-network-runner/network/node"
	"github.com/axiacoin/axia-network-runner/utils"
	"github.com/axiacoin/axia-network-v2/genesis"
	"github.com/axiacoin/axia-network-v2/ids"
	"github.com/axiacoin/axia-network-v2/utils/constants"
	"github.com/axiacoin/axia-network-v2/utils/formatting/address"
	"github.com/axiacoin/axia-network-v2/utils/units"
)

var axChainConfig map[string]interface{}

const (
	validatorStake         = units.MegaAxc
	defaultAXChainConfigStr = "{\"config\":{\"chainId\":43115,\"homesteadBlock\":0,\"daoForkBlock\":0,\"daoForkSupport\":true,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"apricotPhase1BlockTimestamp\":0,\"apricotPhase2BlockTimestamp\":0,\"apricotPhase3BlockTimestamp\":0,\"apricotPhase4BlockTimestamp\":0,\"apricotPhase5BlockTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}"
)

func init() {
	if err := json.Unmarshal([]byte(defaultAXChainConfigStr), &axChainConfig); err != nil {
		panic(err)
	}
}

// AddrAndBalance holds both an address and its balance
type AddrAndBalance struct {
	Addr    ids.ShortID
	Balance uint64
}

// Config that defines a network when it is created.
type Config struct {
	// Must not be empty
	Genesis string `json:"genesis"`
	// May have length 0
	// (i.e. network may have no nodes on creation.)
	NodeConfigs []node.Config `json:"nodeConfigs"`
	// Flags that will be passed to each node in this network.
	// It can be empty.
	// Config flags may also be passed in a node's config struct
	// or config file.
	// The precedence of flags handling is, from highest to lowest:
	// 1. Flags defined in a node's node.Config
	// 2. Flags defined in a network's network.Config
	// 3. Flags defined in a node's config file
	// For example, if a network.Config has flag W set to X,
	// and a node within that network has flag W set to Y,
	// and the node's config file has flag W set to Z,
	// then the node will be started with flag W set to Y.
	Flags map[string]interface{} `json:"flags"`
}

// Validate returns an error if this config is invalid
func (c *Config) Validate() error {
	var someNodeIsBeacon bool
	switch {
	case len(c.Genesis) == 0:
		return errors.New("no genesis given")
	}
	networkID, err := utils.NetworkIDFromGenesis([]byte(c.Genesis))
	if err != nil {
		return fmt.Errorf("couldn't get network ID from genesis: %w", err)
	}
	for i, nodeConfig := range c.NodeConfigs {
		if err := nodeConfig.Validate(networkID); err != nil {
			var nodeName string
			if len(nodeConfig.Name) > 0 {
				nodeName = nodeConfig.Name
			} else {
				nodeName = strconv.Itoa(i)
			}
			return fmt.Errorf("node %q config failed validation: %w", nodeName, err)
		}
		if nodeConfig.IsBeacon {
			someNodeIsBeacon = true
		}
	}
	if len(c.NodeConfigs) > 0 && !someNodeIsBeacon {
		return errors.New("beacon nodes not given")
	}
	return nil
}

// Return a genesis JSON where:
// The nodes in [genesisVdrs] are validators.
// The AXChain and SwapChain balances are given by
// [axChainBalances] and [swapChainBalances].
// Note that many of the genesis fields (i.e. reward addresses)
// are randomly generated or hard-coded.
func NewAxiaGenesis(
	networkID uint32,
	swapChainBalances []AddrAndBalance,
	axChainBalances []AddrAndBalance,
	genesisVdrs []ids.NodeID,
) ([]byte, error) {
	switch networkID {
	case constants.TestnetID, constants.MainnetID, constants.LocalID:
		return nil, errors.New("network ID can't be mainnet, testnet or local network ID")
	}
	switch {
	case len(genesisVdrs) == 0:
		return nil, errors.New("no genesis validators provided")
	case len(swapChainBalances)+len(axChainBalances) == 0:
		return nil, errors.New("no genesis balances given")
	}

	// Address that controls stake doesn't matter -- generate it randomly
	genesisVdrStakeAddr, _ := address.Format(
		"Swap",
		constants.GetHRP(networkID),
		ids.GenerateTestShortID().Bytes(),
	)
	config := genesis.UnparsedConfig{
		NetworkID: networkID,
		Allocations: []genesis.UnparsedAllocation{
			{
				ETHAddr:       "0x0000000000000000000000000000000000000000",
				AXCAddr:      genesisVdrStakeAddr, // Owner doesn't matter
				InitialAmount: 0,
				UnlockSchedule: []genesis.LockedAmount{ // Provides stake to validators
					{
						Amount: uint64(len(genesisVdrs)) * validatorStake,
					},
				},
			},
		},
		StartTime:                  uint64(time.Now().Unix()),
		InitialStakedFunds:         []string{genesisVdrStakeAddr},
		InitialStakeDuration:       31_536_000, // 1 year
		InitialStakeDurationOffset: 5_400,      // 90 minutes
		Message:                    "hello world",
	}

	for _, swapChainBal := range swapChainBalances {
		swapChainAddr, _ := address.Format("Swap", constants.GetHRP(networkID), swapChainBal.Addr[:])
		config.Allocations = append(
			config.Allocations,
			genesis.UnparsedAllocation{
				ETHAddr:       "0x0000000000000000000000000000000000000000",
				AXCAddr:      swapChainAddr,
				InitialAmount: swapChainBal.Balance,
				UnlockSchedule: []genesis.LockedAmount{
					{
						Amount:   validatorStake * uint64(len(genesisVdrs)), // Stake
						Locktime: uint64(time.Now().Add(7 * 24 * time.Hour).Unix()),
					},
				},
			},
		)
	}

	// Set initial AXChain balances.
	axChainAllocs := map[string]interface{}{}
	for _, axChainBal := range axChainBalances {
		addrHex := fmt.Sprintf("0x%s", axChainBal.Addr.Hex())
		balHex := fmt.Sprintf("0x%x", axChainBal.Balance)
		axChainAllocs[addrHex] = map[string]interface{}{
			"balance": balHex,
		}
	}
	// avoid modifying original axChainConfig
	localAXChainConfig := map[string]interface{}{}
	for k, v := range axChainConfig {
		localAXChainConfig[k] = v
	}
	localAXChainConfig["alloc"] = axChainAllocs
	axChainConfigBytes, _ := json.Marshal(localAXChainConfig)
	config.AXChainGenesis = string(axChainConfigBytes)

	// Set initial validators.
	// Give staking rewards to random address.
	rewardAddr, _ := address.Format("Swap", constants.GetHRP(networkID), ids.GenerateTestShortID().Bytes())
	for _, genesisVdr := range genesisVdrs {
		config.InitialStakers = append(
			config.InitialStakers,
			genesis.UnparsedStaker{
				NodeID:        genesisVdr,
				RewardAddress: rewardAddr,
				DelegationFee: 10_000,
			},
		)
	}

	// TODO add validation (from Axia's function validateConfig?)
	return json.Marshal(config)
}
