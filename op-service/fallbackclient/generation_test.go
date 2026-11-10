package fallbackclient

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

type generationClient struct {
	Client
	recover chan struct{}
}

func (c *generationClient) ChainID(context.Context) (*big.Int, error) {
	<-c.recover
	return big.NewInt(56), nil
}
func (*generationClient) Close() {}

type generationMetrics struct{}

func (generationMetrics) RecordL1UrlSwitchEvt(string) {}

func TestRPCGenerationTracksFallbackAndRecovery(t *testing.T) {
	ready := make(chan struct{})
	first := Client(&generationClient{recover: ready})
	next := Client(&generationClient{})
	client := &FallbackClient{
		firstClient: first, urlList: []string{"first", "next"}, log: log.New(),
		metrics: generationMetrics{}, fallbackThreshold: 1,
		clientInitFunc: func(string) (Client, error) { return next, nil },
	}
	client.currentClient.Store(&first)
	observer, ok := any(client).(interface{ RPCGeneration() uint64 })
	require.True(t, ok, "fallback must expose a monotonic generation")
	if !ok {
		return
	}
	require.Zero(t, observer.RPCGeneration())
	client.lastMinuteFail.Store(2)
	client.switchCurrentClient()
	require.EqualValues(t, 1, observer.RPCGeneration())
	close(ready)
	require.Eventually(t, func() bool { return observer.RPCGeneration() == 2 }, time.Second, time.Millisecond)
	require.Equal(t, first, *client.currentClient.Load(), "A -> B -> A must still change the generation")
}

func TestInstrumentedClientPreservesRPCGeneration(t *testing.T) {
	first := Client(&generationClient{})
	client := &FallbackClient{}
	client.currentClient.Store(&first)
	wrapped := NewInstrumentedClient(client, nil)
	observer, ok := any(wrapped).(interface{ RPCGeneration() uint64 })
	require.True(t, ok, "instrumentation must not hide generation changes")
	if !ok {
		return
	}
	require.Zero(t, observer.RPCGeneration())
	client.generationMu.Lock()
	client.generation++
	client.generationMu.Unlock()
	require.EqualValues(t, 1, observer.RPCGeneration())
	// A slow dial holds mx, but reading the published generation must not wait
	// for network I/O or exceed the txmgr's NetworkTimeout.
	client.mx.Lock()
	read := make(chan uint64, 1)
	go func() { read <- observer.RPCGeneration() }()
	select {
	case generation := <-read:
		require.EqualValues(t, 1, generation)
	case <-time.After(time.Second):
		t.Error("generation read blocked on dialing lock")
	}
	client.mx.Unlock()
	plain := NewInstrumentedClient(first, nil)
	require.Zero(t, any(plain).(interface{ RPCGeneration() uint64 }).RPCGeneration())
}
