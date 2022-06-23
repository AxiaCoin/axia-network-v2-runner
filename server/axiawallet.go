// Copyright (C) 2019-2022, Axia Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package server

import (
	"context"
	"time"

	"github.com/axiacoin/axia-network-v2/ids"
	"github.com/axiacoin/axia-network-v2/utils/constants"
	"github.com/axiacoin/axia-network-v2/vms/avm"
	"github.com/axiacoin/axia-network-v2/vms/platformvm"
	"github.com/axiacoin/axia-network-v2/vms/secp256k1fx"
	"github.com/axiacoin/axia-network-v2/axiawallet/chain/core"
	"github.com/axiacoin/axia-network-v2/axiawallet/chain/swap"
	"github.com/axiacoin/axia-network-v2/axiawallet/allychain/primary"
)

const defaultTimeout = time.Minute

func createDefaultCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, defaultTimeout)
}

type refreshableAXIAWallet struct {
	primary.AXIAWallet
	kc *secp256k1fx.Keychain

	pBackend core.Backend
	pBuilder core.Builder
	pSigner  core.Signer

	xBackend swap.Backend
	xBuilder swap.Builder
	xSigner  swap.Signer

	httpRPCEp string
}

// Creates a new axiawallet to work around the case where the new axiawallet object
// is not able to find previous transactions in the cache.
// TODO: support tx backfilling in upstream axiawallet SDK.
func createRefreshableAXIAWallet(ctx context.Context, httpRPCEp string, kc *secp256k1fx.Keychain) (*refreshableAXIAWallet, error) {
	cctx, cancel := createDefaultCtx(ctx)
	pCTX, xCTX, utxos, err := primary.FetchState(cctx, httpRPCEp, kc.Addrs)
	cancel()
	if err != nil {
		return nil, err
	}

	pUTXOs := primary.NewChainUTXOs(constants.CoreChainID, utxos)
	pTXs := make(map[ids.ID]*platformvm.Tx)
	pBackend := core.NewBackend(pCTX, pUTXOs, pTXs)
	pBuilder := core.NewBuilder(kc.Addrs, pBackend)
	pSigner := core.NewSigner(kc, pBackend)

	// need updates when reconnected
	pClient := platformvm.NewClient(httpRPCEp)
	pw := core.NewAXIAWallet(pBuilder, pSigner, pClient, pBackend)

	swapChainID := xCTX.BlockchainID()
	xUTXOs := primary.NewChainUTXOs(swapChainID, utxos)
	xBackend := swap.NewBackend(xCTX, swapChainID, xUTXOs)
	xBuilder := swap.NewBuilder(kc.Addrs, xBackend)
	xSigner := swap.NewSigner(kc, xBackend)

	// need updates when reconnected
	xClient := avm.NewClient(httpRPCEp, "Swap")
	xw := swap.NewAXIAWallet(xBuilder, xSigner, xClient, xBackend)

	return &refreshableAXIAWallet{
		AXIAWallet: primary.NewAXIAWallet(pw, xw),
		kc:     kc,

		pBackend: pBackend,
		pBuilder: pBuilder,
		pSigner:  pSigner,

		xBackend: xBackend,
		xBuilder: xBuilder,
		xSigner:  xSigner,

		httpRPCEp: httpRPCEp,
	}, nil
}

// Refreshes the txs and utxos in case of extended disconnection/restarts.
// TODO: should be "primary.FetchState" again?
// here we assume there's no contending axiawallet user, so just cache everything...
func (w *refreshableAXIAWallet) refresh(httpRPCEp string) {
	// need updates when reconnected
	pClient := platformvm.NewClient(httpRPCEp)
	pw := core.NewAXIAWallet(w.pBuilder, w.pSigner, pClient, w.pBackend)

	// need updates when reconnected
	xClient := avm.NewClient(httpRPCEp, "Swap")
	xw := swap.NewAXIAWallet(w.xBuilder, w.xSigner, xClient, w.xBackend)

	w.AXIAWallet = primary.NewAXIAWallet(pw, xw)
	w.httpRPCEp = httpRPCEp
}
