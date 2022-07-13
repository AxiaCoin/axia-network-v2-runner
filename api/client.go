package api

import (
	"github.com/axiacoin/axia-network-v2/api/admin"
	"github.com/axiacoin/axia-network-v2/api/health"
	"github.com/axiacoin/axia-network-v2/api/info"
	"github.com/axiacoin/axia-network-v2/api/ipcs"
	"github.com/axiacoin/axia-network-v2/api/keystore"
	"github.com/axiacoin/axia-network-v2/indexer"
	"github.com/axiacoin/axia-network-v2/vms/avm"
	"github.com/axiacoin/axia-network-v2/vms/platformvm"
	"github.com/axiacoin/axia-network-v2-coreth/plugin/evm"
)

// Issues API calls to a node
// TODO: byzantine api. check if appropiate. improve implementation.
type Client interface {
	CoreChainAPI() platformvm.Client
	SwapChainAPI() avm.Client
	SwapChainAxiaWalletAPI() avm.AxiaWalletClient
	AXChainAPI() evm.Client
	AXChainEthAPI() EthClient // ethclient websocket wrapper that adds mutexed calls, and lazy conn init (on first call)
	InfoAPI() info.Client
	HealthAPI() health.Client
	IpcsAPI() ipcs.Client
	KeystoreAPI() keystore.Client
	AdminAPI() admin.Client
	CoreChainIndexAPI() indexer.Client
	AXChainIndexAPI() indexer.Client
	// TODO add methods
}
