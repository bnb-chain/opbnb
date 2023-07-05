package sources

import (
	"context"
	"github.com/ethereum-optimism/optimism/op-node/client"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"sync"
	"sync/atomic"
	"time"
)

type FallbackClient struct {
	firstRpc          client.RPC
	urlList           []string
	rpcInitFunc       func(url string) (client.RPC, error)
	lastMinuteFail    atomic.Int64
	currentRpc        client.RPC
	currentIndex      int
	mx                sync.Mutex
	log               log.Logger
	isInFallbackState bool
	subscribeFunc     func() (event.Subscription, error)
	l1HeadsSub        *ethereum.Subscription
}

func NewFallbackClient(rpc client.RPC, urlList []string, log log.Logger, rpcInitFunc func(url string) (client.RPC, error)) client.RPC {
	fallbackClient := &FallbackClient{firstRpc: rpc, urlList: urlList, log: log, rpcInitFunc: rpcInitFunc, currentRpc: rpc, currentIndex: 0}
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				fallbackClient.lastMinuteFail.Store(0)
			}
		}
	}()
	return fallbackClient
}

func (l *FallbackClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	return l.EthSubscribe(ctx, ch, "newHeads")
}

func (l *FallbackClient) Close() {
	l.currentRpc.Close()
	if l.currentRpc != l.firstRpc {
		l.firstRpc.Close()
	}
}

func (l *FallbackClient) CallContext(ctx context.Context, result any, method string, args ...any) error {
	err := l.currentRpc.CallContext(ctx, result, method, args...)
	if err != nil {
		l.handleErr(err)
	}
	return err
}

func (l *FallbackClient) handleErr(err error) {
	if err == rpc.ErrNoResult {
		return
	}
	errCount := l.lastMinuteFail.Add(1)
	if errCount > 10 {
		l.switchCurrentRpc()
	}
}

func (l *FallbackClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	err := l.currentRpc.BatchCallContext(ctx, b)
	if err != nil {
		l.handleErr(err)
	}
	return err
}

func (l *FallbackClient) EthSubscribe(ctx context.Context, channel any, args ...any) (ethereum.Subscription, error) {
	subscribe, err := l.currentRpc.EthSubscribe(ctx, channel, args...)
	if err != nil {
		l.handleErr(err)
	}
	return subscribe, err
}

func (l *FallbackClient) switchCurrentRpc() {
	l.mx.Lock()
	defer l.mx.Unlock()
	if l.lastMinuteFail.Load() <= 10 {
		return
	}
	l.currentIndex++
	if l.currentIndex >= len(l.urlList) {
		log.Error("fallback client has tried all urls")
		return
	}
	url := l.urlList[l.currentIndex]
	newRpc, err := l.rpcInitFunc(url)
	if err != nil {
		log.Error("fallback client switch current RPC fail", "url", url, "err", err)
		return
	}
	if l.currentRpc != l.firstRpc {
		l.currentRpc.Close()
	}
	l.lastMinuteFail.Store(0)
	l.currentRpc = newRpc
	if l.subscribeFunc != nil {
		l.reSubscribeNewRpc(url)
	}
	log.Info("switch current rpc new url", "url", url)
	if !l.isInFallbackState {
		l.isInFallbackState = true
		l.recoverIfFirstRpcHealth()
	}
}

func (l *FallbackClient) reSubscribeNewRpc(url string) {
	(*l.l1HeadsSub).Unsubscribe()
	subscriptionNew, err := l.subscribeFunc()
	if err != nil {
		l.log.Error("can not subscribe new url", "url", url, "err", err)
	} else {
		*l.l1HeadsSub = subscriptionNew
	}
}

func (l *FallbackClient) recoverIfFirstRpcHealth() {
	go func() {
		count := 0
		for {
			var id hexutil.Big
			err := l.firstRpc.CallContext(context.Background(), &id, "eth_chainId")
			if err != nil {
				count = 0
				time.Sleep(5 * time.Second)
				continue
			}
			count++
			if count >= 3 {
				break
			}
		}
		l.mx.Lock()
		defer l.mx.Unlock()
		if !l.isInFallbackState {
			return
		}
		l.currentRpc.Close()
		l.lastMinuteFail.Store(0)
		l.currentRpc = l.firstRpc
		l.currentIndex = 0
		l.isInFallbackState = false
		if l.subscribeFunc != nil {
			l.reSubscribeNewRpc(l.urlList[0])
		}
		log.Info("recover current rpc to first rpc", "url", l.urlList[0])
	}()
}

func (l *FallbackClient) RegisterSubscribeFunc(f func() (event.Subscription, error), l1HeadsSub *ethereum.Subscription) {
	l.subscribeFunc = f
	l.l1HeadsSub = l1HeadsSub
}
