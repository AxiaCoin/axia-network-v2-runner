// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	api "github.com/axiacoin/axia-network-v2-runner/api"
	admin "github.com/axiacoin/axia-network-v2/api/admin"

	avm "github.com/axiacoin/axia-network-v2/vms/avm"

	evm "github.com/axiacoin/axia-network-v2-coreth/plugin/evm"

	health "github.com/axiacoin/axia-network-v2/api/health"

	indexer "github.com/axiacoin/axia-network-v2/indexer"

	info "github.com/axiacoin/axia-network-v2/api/info"

	ipcs "github.com/axiacoin/axia-network-v2/api/ipcs"

	keystore "github.com/axiacoin/axia-network-v2/api/keystore"

	mock "github.com/stretchr/testify/mock"

	platformvm "github.com/axiacoin/axia-network-v2/vms/platformvm"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// AdminAPI provides a mock function with given fields:
func (_m *Client) AdminAPI() admin.Client {
	ret := _m.Called()

	var r0 admin.Client
	if rf, ok := ret.Get(0).(func() admin.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(admin.Client)
		}
	}

	return r0
}

// AXChainAPI provides a mock function with given fields:
func (_m *Client) AXChainAPI() evm.Client {
	ret := _m.Called()

	var r0 evm.Client
	if rf, ok := ret.Get(0).(func() evm.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(evm.Client)
		}
	}

	return r0
}

// AXChainEthAPI provides a mock function with given fields:
func (_m *Client) AXChainEthAPI() api.EthClient {
	ret := _m.Called()

	var r0 api.EthClient
	if rf, ok := ret.Get(0).(func() api.EthClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.EthClient)
		}
	}

	return r0
}

// AXChainIndexAPI provides a mock function with given fields:
func (_m *Client) AXChainIndexAPI() indexer.Client {
	ret := _m.Called()

	var r0 indexer.Client
	if rf, ok := ret.Get(0).(func() indexer.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(indexer.Client)
		}
	}

	return r0
}

// HealthAPI provides a mock function with given fields:
func (_m *Client) HealthAPI() health.Client {
	ret := _m.Called()

	var r0 health.Client
	if rf, ok := ret.Get(0).(func() health.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(health.Client)
		}
	}

	return r0
}

// InfoAPI provides a mock function with given fields:
func (_m *Client) InfoAPI() info.Client {
	ret := _m.Called()

	var r0 info.Client
	if rf, ok := ret.Get(0).(func() info.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(info.Client)
		}
	}

	return r0
}

// IpcsAPI provides a mock function with given fields:
func (_m *Client) IpcsAPI() ipcs.Client {
	ret := _m.Called()

	var r0 ipcs.Client
	if rf, ok := ret.Get(0).(func() ipcs.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ipcs.Client)
		}
	}

	return r0
}

// KeystoreAPI provides a mock function with given fields:
func (_m *Client) KeystoreAPI() keystore.Client {
	ret := _m.Called()

	var r0 keystore.Client
	if rf, ok := ret.Get(0).(func() keystore.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(keystore.Client)
		}
	}

	return r0
}

// CoreChainAPI provides a mock function with given fields:
func (_m *Client) CoreChainAPI() platformvm.Client {
	ret := _m.Called()

	var r0 platformvm.Client
	if rf, ok := ret.Get(0).(func() platformvm.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(platformvm.Client)
		}
	}

	return r0
}

// CoreChainIndexAPI provides a mock function with given fields:
func (_m *Client) CoreChainIndexAPI() indexer.Client {
	ret := _m.Called()

	var r0 indexer.Client
	if rf, ok := ret.Get(0).(func() indexer.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(indexer.Client)
		}
	}

	return r0
}

// SwapChainAPI provides a mock function with given fields:
func (_m *Client) SwapChainAPI() avm.Client {
	ret := _m.Called()

	var r0 avm.Client
	if rf, ok := ret.Get(0).(func() avm.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(avm.Client)
		}
	}

	return r0
}

// SwapChainWalletAPI provides a mock function with given fields:
func (_m *Client) SwapChainWalletAPI() avm.WalletClient {
	ret := _m.Called()

	var r0 avm.WalletClient
	if rf, ok := ret.Get(0).(func() avm.WalletClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(avm.WalletClient)
		}
	}

	return r0
}
