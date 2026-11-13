package batcher

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/require"
)

func TestChannelBuilderPayloadSize(t *testing.T) {
	// Match op-geth miner.maxRollupPayloadTxBytes. Include a real L1-info
	// deposit and signed transactions to exercise both batch encodings.
	const producerBudget = 8*1024*1024 - 1_000_000
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	recipient := crypto.PubkeyToAddress(key.PublicKey)
	for _, batchType := range []uint{derive.SingularBatchType, derive.SpanBatchType} {
		for _, count := range []int{56, 100} {
			t.Run(fmt.Sprintf("batch-%d/txs-%d", batchType, count), func(t *testing.T) {
				base := newMiniL2BlockWithNumberParent(count, big.NewInt(1), [32]byte{})
				txs := base.Transactions()
				signer := types.LatestSignerForChainID(defaultTestRollupConfig.L2ChainID)
				var txBytes uint64
				for i := 1; i < len(txs); i++ {
					txs[i] = types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
						ChainID: defaultTestRollupConfig.L2ChainID, Nonce: uint64(i - 1),
						To:        &recipient,
						Gas:       params.TxGas + params.TxDataZeroGas*130_000,
						GasFeeCap: big.NewInt(params.InitialBaseFee), GasTipCap: big.NewInt(0),
						Data: make([]byte, 130_000),
					})
				}
				for _, tx := range txs {
					txBytes += tx.Size()
				}
				block := types.NewBlock(base.Header(), txs, nil, nil, trie.NewStackTrie(nil))
				batch, _, err := derive.BlockToSingularBatch(&defaultTestRollupConfig, block)
				require.NoError(t, err)
				var encoded bytes.Buffer
				require.NoError(t, rlp.Encode(&encoded, derive.NewBatchData(batch)))
				cfg := defaultTestChannelConfig()
				cfg.BatchType = batchType
				cb, err := NewChannelBuilder(cfg, defaultTestRollupConfig, latestL1BlockOrigin)
				require.NoError(t, err)
				_, err = cb.AddBlock(block)
				if count == 56 {
					require.LessOrEqual(t, txBytes, uint64(producerBudget))
					require.Less(t, uint64(producerBudget)-txBytes, uint64(131_072))
					require.Less(t, encoded.Len(), rollup.SafeMaxRLPBytesPerChannel)
					require.NoError(t, err)
				} else {
					require.Greater(t, encoded.Len(), rollup.SafeMaxRLPBytesPerChannel)
					require.ErrorIs(t, err, derive.ErrTooManyRLPBytes)
				}
				t.Logf("tx_bytes=%d singular_rlp_bytes=%d", txBytes, encoded.Len())
			})
		}
	}
}
