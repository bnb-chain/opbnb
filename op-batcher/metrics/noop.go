package metrics

import (
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-batcher/flags"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	txmetrics "github.com/ethereum-optimism/optimism/op-service/txmgr/metrics"
)

type noopMetrics struct {
	opmetrics.NoopRefMetrics
	txmetrics.NoopTxMetrics
}

var NoopMetrics Metricer = new(noopMetrics)

func (*noopMetrics) Document() []opmetrics.DocumentedMetric { return nil }

func (*noopMetrics) RecordInfo(version string) {}
func (*noopMetrics) RecordUp()                 {}

func (*noopMetrics) RecordLatestL1Block(l1ref eth.L1BlockRef)               {}
func (*noopMetrics) RecordL2BlocksLoaded(eth.L2BlockRef)                    {}
func (*noopMetrics) RecordChannelOpened(derive.ChannelID, int)              {}
func (*noopMetrics) RecordL2BlocksAdded(eth.L2BlockRef, int, int, int, int) {}
func (*noopMetrics) RecordL2BlockInPendingQueue(*types.Block)               {}
func (*noopMetrics) RecordL2BlockInChannel(*types.Block)                    {}

func (*noopMetrics) RecordChannelClosed(derive.ChannelID, int, int, int, int, error) {}

func (*noopMetrics) RecordChannelFullySubmitted(derive.ChannelID) {}
func (*noopMetrics) RecordChannelTimedOut(derive.ChannelID)       {}

func (*noopMetrics) RecordBatchTxSubmitted() {}
func (*noopMetrics) RecordBatchTxSuccess()   {}
func (*noopMetrics) RecordBatchTxFailed()    {}
func (*noopMetrics) RecordBlobUsedBytes(int) {}
func (*noopMetrics) StartBalanceMetrics(log.Logger, ethereum.ChainStateReader, common.Address) io.Closer {
	return nil
}
func (*noopMetrics) RecordBlobsNumber(number int) {}

func (*noopMetrics) RecordAutoChoosedDAType(daType flags.DataAvailabilityType) {}
func (*noopMetrics) RecordEconomicAutoSwitchCount()                            {}
func (*noopMetrics) RecordReservedErrorSwitchCount()                           {}
func (*noopMetrics) RecordAutoSwitchTimeDuration(duration time.Duration)       {}
func (*noopMetrics) RecordEstimatedCalldataTypeFee(fee *big.Int)               {}
func (*noopMetrics) RecordEstimatedBlobTypeFee(fee *big.Int)                   {}

func (m *noopMetrics) RecordL1UrlSwitchEvt(url string) {
}
