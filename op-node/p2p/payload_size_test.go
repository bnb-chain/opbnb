package p2p

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/snappy"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

func TestBlockValidatorPayloadSize(t *testing.T) {
	// Match op-geth miner.maxRollupPayloadTxBytes; channel and gossip limits
	// must be able to carry the sequencer's largest transaction set.
	const producerBudget = 8*1024*1024 - 1_000_000
	cfg := &rollup.Config{L2ChainID: big.NewInt(204)}
	secrets, err := e2eutils.DefaultMnemonicConfig.Secrets()
	require.NoError(t, err)
	runCfg := &testutils.MockRuntimeConfig{P2PSeqAddress: crypto.PubkeyToAddress(secrets.SequencerP2P.PublicKey)}
	signer := &PreparedSigner{Signer: NewLocalSigner(secrets.SequencerP2P)}

	for _, tc := range []struct {
		name  string
		count int
		want  pubsub.ValidationResult
	}{
		{"producer-budget", 56, pubsub.ValidationAccept},
		{"oversized-report", 100, pubsub.ValidationReject},
	} {
		t.Run(tc.name, func(t *testing.T) {
			zero := uint64(0)
			beaconHash := common.Hash{1}
			envelope := createEnvelope(&beaconHash, types.Withdrawals{}, &zero, &zero)
			for i := 0; i < tc.count; i++ {
				tx := types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &common.Address{}, Gas: 541_000, Data: make([]byte, 130_000)})
				encoded, err := tx.MarshalBinary()
				require.NoError(t, err)
				envelope.ExecutionPayload.Transactions = append(envelope.ExecutionPayload.Transactions, encoded)
			}
			// Hash the real transaction encodings, so the accepted case exercises
			// signature, decoding and block-hash validation as well as the size gate.
			hash, _ := envelope.CheckBlockHash()
			envelope.ExecutionPayload.BlockHash = hash
			data, err := createSignedP2Payload(envelope, signer, cfg.L2ChainID)
			require.NoError(t, err)
			decodedLen, err := snappy.DecodedLen(data)
			require.NoError(t, err)
			require.Less(t, len(data), maxGossipSize)
			if tc.want == pubsub.ValidationAccept {
				var txBytes int
				for _, tx := range envelope.ExecutionPayload.Transactions {
					txBytes += len(tx)
				}
				require.LessOrEqual(t, txBytes, producerBudget)
				require.Less(t, producerBudget-txBytes, 131_072)
				require.Less(t, decodedLen, maxGossipSize)
			} else {
				require.Greater(t, decodedLen, maxGossipSize)
			}
			validator := BuildBlocksValidator(testlog.Logger(t, log.LevelCrit), cfg, runCfg, eth.BlockV3)
			message := &pubsub.Message{Message: &pubsubpb.Message{Data: data}}
			require.Equal(t, tc.want, validator(context.Background(), peer.ID("payload-size"), message))
			t.Logf("decoded_bytes=%d compressed_bytes=%d result=%s", decodedLen, len(data), validationResultString(tc.want))
		})
	}
}
