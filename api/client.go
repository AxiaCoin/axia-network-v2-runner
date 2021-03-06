package api

import (
	"github.com/axiacoin/axia/api/admin"
	"github.com/axiacoin/axia/api/health"
	"github.com/axiacoin/axia/api/info"
	"github.com/axiacoin/axia/api/ipcs"
	"github.com/axiacoin/axia/api/keystore"
	"github.com/axiacoin/axia/indexer"
	"github.com/axiacoin/axia/vms/avm"
	"github.com/axiacoin/axia/vms/platformvm"
	"github.com/axiacoin/coreth/plugin/evm"
)

// Issues API calls to a node
// TODO: byzantine api. check if appropiate. improve implementation.
type Client interface {
	PChainAPI() platformvm.Client
	XChainAPI() avm.Client
	XChainWalletAPI() avm.WalletClient
	CChainAPI() evm.Client
	CChainEthAPI() EthClient // ethclient websocket wrapper that adds mutexed calls, and lazy conn init (on first call)
	InfoAPI() info.Client
	HealthAPI() health.Client
	IpcsAPI() ipcs.Client
	KeystoreAPI() keystore.Client
	AdminAPI() admin.Client
	PChainIndexAPI() indexer.Client
	CChainIndexAPI() indexer.Client
	// TODO add methods
}
