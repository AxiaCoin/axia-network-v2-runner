package api

import (
	"fmt"

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

// interface compliance
var (
	_ Client        = (*APIClient)(nil)
	_ NewAPIClientF = NewAPIClient
)

// APIClient gives access to most axia apis (or suitable wrappers)
type APIClient struct {
	platform     platformvm.Client
	swapChain       avm.Client
	swapChainAXIAWallet avm.AXIAWalletClient
	axChain       evm.Client
	axChainEth    EthClient
	info         info.Client
	health       health.Client
	ipcs         ipcs.Client
	keystore     keystore.Client
	admin        admin.Client
	pindex       indexer.Client
	cindex       indexer.Client
}

// Returns a new API client for a node at [ipAddr]:[port].
type NewAPIClientF func(ipAddr string, port uint16) Client

// NewAPIClient initialize most of axia apis
func NewAPIClient(ipAddr string, port uint16) Client {
	uri := fmt.Sprintf("http://%s:%d", ipAddr, port)
	return &APIClient{
		platform:     platformvm.NewClient(uri),
		swapChain:       avm.NewClient(uri, "Swap"),
		swapChainAXIAWallet: avm.NewAXIAWalletClient(uri, "Swap"),
		axChain:       evm.NewAXChainClient(uri),
		axChainEth:    NewEthClient(ipAddr, uint(port)), // wrapper over ethclient.Client
		info:         info.NewClient(uri),
		health:       health.NewClient(uri),
		ipcs:         ipcs.NewClient(uri),
		keystore:     keystore.NewClient(uri),
		admin:        admin.NewClient(uri),
		pindex:       indexer.NewClient(uri),
		cindex:       indexer.NewClient(uri),
	}
}

func (c APIClient) CoreChainAPI() platformvm.Client {
	return c.platform
}

func (c APIClient) SwapChainAPI() avm.Client {
	return c.swapChain
}

func (c APIClient) SwapChainAXIAWalletAPI() avm.AXIAWalletClient {
	return c.swapChainAXIAWallet
}

func (c APIClient) AXChainAPI() evm.Client {
	return c.axChain
}

func (c APIClient) AXChainEthAPI() EthClient {
	return c.axChainEth
}

func (c APIClient) InfoAPI() info.Client {
	return c.info
}

func (c APIClient) HealthAPI() health.Client {
	return c.health
}

func (c APIClient) IpcsAPI() ipcs.Client {
	return c.ipcs
}

func (c APIClient) KeystoreAPI() keystore.Client {
	return c.keystore
}

func (c APIClient) AdminAPI() admin.Client {
	return c.admin
}

func (c APIClient) CoreChainIndexAPI() indexer.Client {
	return c.pindex
}

func (c APIClient) AXChainIndexAPI() indexer.Client {
	return c.cindex
}
