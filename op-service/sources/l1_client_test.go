package sources

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGoOrUpdatePreFetchReceipts(t *testing.T) {
	t.Run("handleReOrg", func(t *testing.T) {
		m := new(mockRPC)
		ctx := context.Background()
		clientLog := testlog.Logger(t, log.LvlDebug)
		latestHead := &RPCHeader{
			ParentHash:      randHash(),
			UncleHash:       common.Hash{},
			Coinbase:        common.Address{},
			Root:            types.EmptyRootHash,
			TxHash:          types.EmptyTxsHash,
			ReceiptHash:     types.EmptyReceiptsHash,
			Bloom:           eth.Bytes256{},
			Difficulty:      hexutil.Big{},
			Number:          100,
			GasLimit:        0,
			GasUsed:         0,
			Time:            0,
			Extra:           nil,
			MixDigest:       common.Hash{},
			Nonce:           types.BlockNonce{},
			BaseFee:         nil,
			WithdrawalsRoot: nil,
			Hash:            randHash(),
		}
		m.On("CallContext", ctx, new(*RPCHeader),
			"eth_getBlockByNumber", []any{"latest", false}).Run(func(args mock.Arguments) {
			*args[1].(**RPCHeader) = latestHead
		}).Return([]error{nil})
		for i := 81; i <= 90; i++ {
			currentHead := &RPCHeader{
				ParentHash:      randHash(),
				UncleHash:       common.Hash{},
				Coinbase:        common.Address{},
				Root:            types.EmptyRootHash,
				TxHash:          types.EmptyTxsHash,
				ReceiptHash:     types.EmptyReceiptsHash,
				Bloom:           eth.Bytes256{},
				Difficulty:      hexutil.Big{},
				Number:          hexutil.Uint64(i),
				GasLimit:        0,
				GasUsed:         0,
				Time:            0,
				Extra:           nil,
				MixDigest:       common.Hash{},
				Nonce:           types.BlockNonce{},
				BaseFee:         nil,
				WithdrawalsRoot: nil,
				Hash:            randHash(),
			}
			currentBlock := &RPCBlock{
				RPCHeader:    *currentHead,
				Transactions: []*types.Transaction{},
			}
			m.On("CallContext", ctx, new(*RPCHeader),
				"eth_getBlockByNumber", []any{numberID(i).Arg(), false}).Once().Run(func(args mock.Arguments) {
				*args[1].(**RPCHeader) = currentHead
			}).Return([]error{nil})
			m.On("CallContext", ctx, new(*RPCBlock),
				"eth_getBlockByHash", []any{currentHead.Hash, true}).Once().Run(func(args mock.Arguments) {
				*args[1].(**RPCBlock) = currentBlock
			}).Return([]error{nil})
		}
		for i := 91; i <= 100; i++ {
			currentHead := &RPCHeader{
				ParentHash:      randHash(),
				UncleHash:       common.Hash{},
				Coinbase:        common.Address{},
				Root:            types.EmptyRootHash,
				TxHash:          types.EmptyTxsHash,
				ReceiptHash:     types.EmptyReceiptsHash,
				Bloom:           eth.Bytes256{},
				Difficulty:      hexutil.Big{},
				Number:          hexutil.Uint64(i),
				GasLimit:        0,
				GasUsed:         0,
				Time:            0,
				Extra:           nil,
				MixDigest:       common.Hash{},
				Nonce:           types.BlockNonce{},
				BaseFee:         nil,
				WithdrawalsRoot: nil,
				Hash:            randHash(),
			}
			m.On("CallContext", ctx, new(*RPCHeader),
				"eth_getBlockByNumber", []any{numberID(i).Arg(), false}).Once().Run(func(args mock.Arguments) {
				*args[1].(**RPCHeader) = currentHead
			}).Return([]error{nil})
			currentBlock := &RPCBlock{
				RPCHeader:    *currentHead,
				Transactions: []*types.Transaction{},
			}
			m.On("CallContext", ctx, new(*RPCBlock),
				"eth_getBlockByHash", []any{currentHead.Hash, true}).Once().Run(func(args mock.Arguments) {
				*args[1].(**RPCBlock) = currentBlock
			}).Return([]error{nil})
		}
		var lastParentHeader common.Hash
		var real100Hash common.Hash
		for i := 76; i <= 100; i++ {
			currentHead := &RPCHeader{
				ParentHash:      lastParentHeader,
				UncleHash:       common.Hash{},
				Coinbase:        common.Address{},
				Root:            types.EmptyRootHash,
				TxHash:          types.EmptyTxsHash,
				ReceiptHash:     types.EmptyReceiptsHash,
				Bloom:           eth.Bytes256{},
				Difficulty:      hexutil.Big{},
				Number:          hexutil.Uint64(i),
				GasLimit:        0,
				GasUsed:         0,
				Time:            0,
				Extra:           nil,
				MixDigest:       common.Hash{},
				Nonce:           types.BlockNonce{},
				BaseFee:         nil,
				WithdrawalsRoot: nil,
				Hash:            randHash(),
			}
			if i == 100 {
				real100Hash = currentHead.Hash
			}
			lastParentHeader = currentHead.Hash
			m.On("CallContext", ctx, new(*RPCHeader),
				"eth_getBlockByNumber", []any{numberID(i).Arg(), false}).Once().Run(func(args mock.Arguments) {
				*args[1].(**RPCHeader) = currentHead
			}).Return([]error{nil})
			currentBlock := &RPCBlock{
				RPCHeader:    *currentHead,
				Transactions: []*types.Transaction{},
			}
			m.On("CallContext", ctx, new(*RPCBlock),
				"eth_getBlockByHash", []any{currentHead.Hash, true}).Once().Run(func(args mock.Arguments) {
				*args[1].(**RPCBlock) = currentBlock
			}).Return([]error{nil})
		}
		s, err := NewL1Client(m, clientLog, nil, L1ClientDefaultConfig(&rollup.Config{SeqWindowSize: 1000}, true, RPCKindBasic))
		require.NoError(t, err)
		err2 := s.GoOrUpdatePreFetchReceipts(ctx, 81)
		require.NoError(t, err2)
		time.Sleep(1 * time.Second)
		pair, ok := s.recProvider.GetReceiptsCache().Get(100, false)
		require.True(t, ok, "100 cache miss")
		require.Equal(t, real100Hash, pair.blockHash, "block 100 hash is different,want:%s,but:%s", real100Hash, pair.blockHash)
		_, ok2 := s.recProvider.GetReceiptsCache().Get(76, false)
		require.True(t, ok2, "76 cache miss")
	})
}
