package proposer

import (
	"context"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type OutputRootCacheHandler struct {
	ctx                 context.Context
	log                 log.Logger
	isStart             *atomic.Bool
	readyChan           chan []eth.Bytes32
	outputRootFetchFunc func(ctx context.Context, block *big.Int) (*eth.OutputResponse, bool, error)
}

func newOutputRootCacheHandler(ctx context.Context, log log.Logger, outputRootFetchFunc func(ctx context.Context, block *big.Int) (*eth.OutputResponse, bool, error)) *OutputRootCacheHandler {
	return &OutputRootCacheHandler{
		ctx:                 ctx,
		log:                 log,
		isStart:             &atomic.Bool{},
		readyChan:           make(chan []eth.Bytes32, 1),
		outputRootFetchFunc: outputRootFetchFunc,
	}
}

func (h *OutputRootCacheHandler) startFrom(parentGame *GameInformation) {
	if !h.isStart.CompareAndSwap(false, true) {
		return
	}
	go h.loop(parentGame)
}

func (h *OutputRootCacheHandler) loop(parentGame *GameInformation) {
	currentBlockNumber := new(big.Int).Add(parentGame.extraData.endL2BlockNumber, big.NewInt(1))
	endBlockNumber := new(big.Int).Add(parentGame.extraData.endL2BlockNumber, big.NewInt(3600))
	outputRootResult := make([]eth.Bytes32, 0, 1200)
	for {
		outputRootResponse, shouldPropose, err := h.outputRootFetchFunc(h.ctx, currentBlockNumber)
		if err != nil {
			h.log.Error("failed to fetch outputRoot", "err", err, "currentBlockNumber", currentBlockNumber)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if !shouldPropose {
			h.log.Warn("the current block number cannot submit the output root yet; it needs to wait for a while", "currentBlockNumber", currentBlockNumber)
			time.Sleep(1 * time.Second)
			continue
		}
		outputRootResult = append(outputRootResult, outputRootResponse.OutputRoot)
		currentBlockNumber.Add(currentBlockNumber, big.NewInt(1))
		if currentBlockNumber.Cmp(endBlockNumber) > 0 {
			break
		}
	}
	h.readyChan <- outputRootResult
}
