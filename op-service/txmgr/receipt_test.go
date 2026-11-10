package txmgr

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

// receiptBackend scripts independent receipt, latest, finalized and canonical
// header responses, including fallback switches during a confirmation check.
type receiptBackend struct {
	ETHBackend
	receipts         []*types.Receipt
	receiptErr       error
	tip              *types.Header
	tipErr           error
	block            *types.Header
	blockErr         error
	finalized        *types.Header
	finalizedErr     error
	finalizedHeads   []*types.Header
	generation       uint64
	switchAt         string
	switches         uint64
	canonicalReads   int
	waitForCanonical bool
}

func (b *receiptBackend) TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error) {
	if b.receiptErr != nil {
		return nil, b.receiptErr
	}
	receipt := b.receipts[0]
	if len(b.receipts) > 1 {
		b.receipts = b.receipts[1:]
	}
	return receipt, nil
}

func (b *receiptBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if number == nil {
		b.switchProvider("tip")
		return b.tip, b.tipErr
	}
	if number.Int64() == int64(rpc.FinalizedBlockNumber) {
		b.switchProvider("finalized")
		if len(b.finalizedHeads) > 0 {
			h := b.finalizedHeads[0]
			b.finalizedHeads = b.finalizedHeads[1:]
			return h, nil
		}
		return b.finalized, b.finalizedErr
	}
	b.canonicalReads++
	b.switchProvider("canonical")
	if b.waitForCanonical {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return b.block, b.blockErr
}

func TestQueryReceiptCanonicalBlock(t *testing.T) {
	t.Parallel()

	block := &types.Header{Number: big.NewInt(100), Extra: []byte("canonical")}
	orphan := &types.Header{Number: big.NewInt(100), Extra: []byte("orphan")}
	for _, test := range []struct {
		name           string
		receiptBlock   *types.Header
		canonicalBlock *types.Header
		blockErr       error
		tipErr         error
		tipHeight      int64
		confirmations  uint64
		wantReceipt    bool
		wantWaiting    bool
	}{
		{name: "canonical", receiptBlock: block, canonicalBlock: block, tipHeight: 109, confirmations: 10, wantReceipt: true, wantWaiting: true},
		{name: "one_confirmation", receiptBlock: block, canonicalBlock: block, tipHeight: 100, confirmations: 1, wantReceipt: true, wantWaiting: true},
		{name: "orphan", receiptBlock: orphan, canonicalBlock: block, tipHeight: 109, confirmations: 10},
		{name: "orphan_one_confirmation", receiptBlock: orphan, canonicalBlock: block, tipHeight: 100, confirmations: 1},
		{name: "missing_block", receiptBlock: block, blockErr: ethereum.NotFound, tipHeight: 109, confirmations: 10},
		{name: "temporary_error", receiptBlock: block, blockErr: errors.New("RPC unavailable"), tipHeight: 109, confirmations: 10, wantWaiting: true},
		{name: "empty_header", receiptBlock: block, tipHeight: 109, confirmations: 10, wantWaiting: true},
		{name: "below_confirmation_depth", receiptBlock: block, tipHeight: 108, confirmations: 10, wantWaiting: true},
		{name: "tip_error", receiptBlock: block, canonicalBlock: block, tipErr: errors.New("RPC unavailable"), tipHeight: 109, confirmations: 10, wantWaiting: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := configWithNumConfs(test.confirmations)
			cfg.NetworkTimeout = time.Second
			h := newTestHarnessWithConfig(t, cfg)
			txHash := common.Hash{0x42}
			receipt := &types.Receipt{
				TxHash: txHash, BlockNumber: test.receiptBlock.Number,
				BlockHash: test.receiptBlock.Hash(), Status: types.ReceiptStatusSuccessful,
			}
			h.mgr.backend = &receiptBackend{
				ETHBackend: h.backend, receipts: []*types.Receipt{receipt},
				tip:    &types.Header{Number: big.NewInt(test.tipHeight), BaseFee: big.NewInt(1)},
				tipErr: test.tipErr, block: test.canonicalBlock, blockErr: test.blockErr,
				finalized: &types.Header{Number: big.NewInt(109)},
			}
			state := testSendState()
			state.TxMined(txHash)
			got := h.mgr.queryReceipt(context.Background(), txHash, state)
			if test.wantReceipt {
				require.Same(t, receipt, got)
			} else {
				require.Nil(t, got, "must not confirm an unverified receipt")
			}
			require.Equal(t, test.wantWaiting, state.IsWaitingForConfirmation())
		})
	}
}

func TestSendTxWaitsForCanonicalReceipt(t *testing.T) {
	t.Parallel()

	cfg := configWithNumConfs(10)
	cfg.NetworkTimeout = time.Second
	cfg.ReceiptQueryInterval = time.Millisecond
	h := newTestHarnessWithConfig(t, cfg)
	h.backend.setTxSender(func(context.Context, *types.Transaction) error { return nil })
	tx := types.NewTx(&types.DynamicFeeTx{GasTipCap: big.NewInt(1), GasFeeCap: big.NewInt(2)})
	orphanBlock := &types.Header{Number: big.NewInt(100), Extra: []byte("orphan")}
	canonicalBlock := &types.Header{Number: big.NewInt(100), Extra: []byte("canonical")}
	orphan := &types.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(100), BlockHash: orphanBlock.Hash()}
	canonical := &types.Receipt{TxHash: tx.Hash(), BlockNumber: big.NewInt(100), BlockHash: canonicalBlock.Hash()}
	h.mgr.backend = &receiptBackend{
		ETHBackend: h.backend, receipts: []*types.Receipt{orphan, canonical},
		tip: &types.Header{Number: big.NewInt(109), BaseFee: big.NewInt(1)}, block: canonicalBlock, finalized: canonicalBlock,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	receipt, err := h.mgr.sendTx(ctx, tx)
	require.NoError(t, err)
	require.Equal(t, canonical, receipt, "sendTx must wait for canonical reinclusion")
}

func (b *receiptBackend) RPCGeneration() uint64 { return b.generation }
func (b *receiptBackend) switchProvider(phase string) {
	if b.switchAt == phase {
		b.generation += b.switches
		b.switchAt = ""
	}
}

func TestQueryReceiptFinalized(t *testing.T) {
	block := &types.Header{Number: big.NewInt(100)}
	orphan := &types.Header{Number: big.NewInt(100), Extra: []byte("orphan")}
	for _, test := range []struct {
		name                     string
		finalized, canonical     *types.Header
		finalizedErr, blockErr   error
		switchAt                 string
		switches                 uint64
		wantReceipt, wantWaiting bool
	}{
		{name: "not_finalized", finalized: &types.Header{Number: big.NewInt(99)}, canonical: block, wantWaiting: true},
		{name: "finalized_error", finalizedErr: errors.New("unsupported finalized tag"), canonical: block, wantWaiting: true},
		{name: "nil_finalized", canonical: block, wantWaiting: true},
		{name: "nil_finalized_number", finalized: &types.Header{}, canonical: block, wantWaiting: true},
		{name: "finalized", finalized: &types.Header{Number: big.NewInt(105)}, canonical: block, wantReceipt: true, wantWaiting: true},
		{name: "orphan", finalized: block, canonical: orphan},
		{name: "missing_canonical", finalized: block, blockErr: ethereum.NotFound},
		{name: "canonical_error", finalized: block, blockErr: errors.New("unavailable"), wantWaiting: true},
		{name: "nil_canonical", finalized: block, wantWaiting: true},
		{name: "wrong_canonical_height", finalized: block, canonical: &types.Header{Number: big.NewInt(101)}, wantWaiting: true},
		{name: "switch_before_finalized", finalized: block, canonical: block, switchAt: "tip", switches: 1, wantWaiting: true},
		{name: "switch_at_finalized", finalized: block, canonical: block, switchAt: "finalized", switches: 1, wantWaiting: true},
		{name: "switch_at_canonical", finalized: block, canonical: block, switchAt: "canonical", switches: 1, wantWaiting: true},
		{name: "switch_away_and_back", finalized: block, canonical: block, switchAt: "canonical", switches: 2, wantWaiting: true},
		{name: "switch_with_orphan", finalized: block, canonical: orphan, switchAt: "canonical", switches: 1, wantWaiting: true},
		{name: "switch_with_not_found", finalized: block, blockErr: ethereum.NotFound, switchAt: "canonical", switches: 1, wantWaiting: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := configWithNumConfs(10)
			h := newTestHarnessWithConfig(t, cfg)
			hash := common.Hash{0x42}
			receipt := &types.Receipt{TxHash: hash, BlockNumber: block.Number, BlockHash: block.Hash()}
			backend := &receiptBackend{ETHBackend: h.backend, receipts: []*types.Receipt{receipt}, tip: &types.Header{Number: big.NewInt(109), BaseFee: big.NewInt(1)}, block: test.canonical, blockErr: test.blockErr, finalized: test.finalized, finalizedErr: test.finalizedErr, switchAt: test.switchAt, switches: test.switches}
			h.mgr.backend = backend
			state := testSendState()
			got := h.mgr.queryReceipt(context.Background(), hash, state)
			if test.wantReceipt {
				require.Same(t, receipt, got)
			} else {
				require.Nil(t, got)
			}
			require.Equal(t, test.wantWaiting, state.IsWaitingForConfirmation())
			if test.finalized == nil || test.finalized.Number == nil || test.finalized.Number.Uint64() < 100 {
				require.Zero(t, backend.canonicalReads)
			}
			if test.switches > 0 && test.blockErr == nil && test.canonical == block {
				require.Same(t, receipt, h.mgr.queryReceipt(context.Background(), hash, state))
			}
		})
	}
}

func TestSendTxWaitsForFinality(t *testing.T) {
	cfg := configWithNumConfs(1)
	cfg.ReceiptQueryInterval = time.Millisecond
	h := newTestHarnessWithConfig(t, cfg)
	h.backend.setTxSender(func(context.Context, *types.Transaction) error { return nil })
	tx := types.NewTx(&types.DynamicFeeTx{GasTipCap: big.NewInt(1), GasFeeCap: big.NewInt(2)})
	block := &types.Header{Number: big.NewInt(100), BaseFee: big.NewInt(1)}
	receipt := &types.Receipt{TxHash: tx.Hash(), BlockNumber: block.Number, BlockHash: block.Hash()}
	backend := &receiptBackend{ETHBackend: h.backend, receipts: []*types.Receipt{receipt}, tip: block, block: block, finalizedHeads: []*types.Header{{Number: big.NewInt(99)}, block}}
	h.mgr.backend = backend
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, err := h.mgr.sendTx(ctx, tx)
	require.NoError(t, err)
	require.Same(t, receipt, got)
	require.Empty(t, backend.finalizedHeads)
}

func TestQueryReceiptFinalizedTimeout(t *testing.T) {
	cfg := configWithNumConfs(1)
	cfg.NetworkTimeout = 10 * time.Millisecond
	h := newTestHarnessWithConfig(t, cfg)
	block := &types.Header{Number: big.NewInt(100), BaseFee: big.NewInt(1)}
	hash := common.Hash{0x42}
	h.mgr.backend = &receiptBackend{ETHBackend: h.backend, receipts: []*types.Receipt{{TxHash: hash, BlockNumber: block.Number, BlockHash: block.Hash()}}, tip: block, block: block, finalized: block, waitForCanonical: true}
	state := testSendState()
	require.Nil(t, h.mgr.queryReceipt(context.Background(), hash, state))
	require.True(t, state.IsWaitingForConfirmation())
}
