package metrics

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum-optimism/optimism/op-service/metrics"
)

type NoopTxMetrics struct {
	metrics.NoopRPCMetrics
}

func (*NoopTxMetrics) RecordNonce(uint64)                {}
func (*NoopTxMetrics) RecordPendingTx(int64)             {}
func (*NoopTxMetrics) RecordGasBumpCount(int)            {}
func (*NoopTxMetrics) RecordTxConfirmationLatency(int64) {}
func (*NoopTxMetrics) TxConfirmed(*types.Receipt)        {}
func (*NoopTxMetrics) TxPublished(string)                {}
func (*NoopTxMetrics) RecordBaseFee(*big.Int)            {}
func (*NoopTxMetrics) RecordBlobBaseFee(*big.Int)        {}
func (*NoopTxMetrics) RecordTipCap(*big.Int)             {}
func (*NoopTxMetrics) RPCError()                         {}
func (m *NoopTxMetrics) RecordL1UrlSwitchEvt(url string) {}
func (m *NoopTxMetrics) RecordBlobsNumber(number int)    {}
