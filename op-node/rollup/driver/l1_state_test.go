package driver

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-node/metrics"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

func TestL1StateFinalized(t *testing.T) {
	state := NewL1State(log.NewLogger(log.DiscardHandler()), metrics.NoopMetrics)
	require.Equal(t, eth.L1BlockRef{}, state.L1Finalized())

	finalized := eth.L1BlockRef{
		Hash:       common.Hash{0x11},
		Number:     100,
		ParentHash: common.Hash{0x22},
		Time:       1_000,
		MilliTime:  111,
	}
	state.HandleNewL1FinalizedBlock(finalized)
	require.Equal(t, finalized, state.L1Finalized())
}

func TestL1StateFinalizedConcurrentAccess(t *testing.T) {
	state := NewL1State(log.NewLogger(log.DiscardHandler()), metrics.NoopMetrics)
	refs := [2]eth.L1BlockRef{
		{Hash: common.Hash{0x11}, Number: 100, ParentHash: common.Hash{0x22}, Time: 1_000, MilliTime: 111},
		{Hash: common.Hash{0xff}, Number: 101, ParentHash: common.Hash{0xee}, Time: 2_000, MilliTime: 999},
	}
	state.HandleNewL1FinalizedBlock(refs[0])

	const iterations = 10_000
	const readers = 4
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1 + readers)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < iterations; i++ {
			state.HandleNewL1FinalizedBlock(refs[i%len(refs)])
		}
	}()
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				ref := state.L1Finalized()
				if ref != refs[0] && ref != refs[1] {
					t.Errorf("read inconsistent finalized block snapshot: %+v", ref)
					return
				}
			}
		}()
	}
	close(start)
	wg.Wait()

	require.Equal(t, refs[1], state.L1Finalized())
}
