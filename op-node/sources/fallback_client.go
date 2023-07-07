package sources

import (
	"context"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-node/client"
	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
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
	l1ChainId         *big.Int
	l1Block           eth.BlockID
	ctx               context.Context
}

func NewFallbackClient(ctx context.Context, rpc client.RPC, urlList []string, log log.Logger, l1ChainId *big.Int, l1Block eth.BlockID, rpcInitFunc func(url string) (client.RPC, error)) client.RPC {
	fallbackClient := &FallbackClient{
		ctx:          ctx,
		firstRpc:     rpc,
		urlList:      urlList,
		log:          log,
		rpcInitFunc:  rpcInitFunc,
		currentRpc:   rpc,
		currentIndex: 0,
		l1ChainId:    l1ChainId,
		l1Block:      l1Block,
	}
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				log.Debug("FallbackClient clear lastMinuteFail 0")
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
	for {
		l.currentIndex++
		err := l.switchCurrentRpcLogic()
		if err != nil {
			l.log.Warn("fallback client switch current rpc fail", "err", err)
		} else {
			break
		}
	}
}

func (l *FallbackClient) switchCurrentRpcLogic() error {
	if l.currentIndex >= len(l.urlList) {
		return fmt.Errorf("fallback client has tried all urls")
	}
	url := l.urlList[l.currentIndex]
	newRpc, err := l.rpcInitFunc(url)
	if err != nil {
		return fmt.Errorf("fallback client init RPC fail,url:%s, err:%v", url, err)
	}
	vErr := l.validateRpc(newRpc)
	if vErr != nil {
		return vErr
	}
	lastRpc := l.currentRpc
	l.currentRpc = newRpc
	if lastRpc != l.firstRpc {
		lastRpc.Close()
	}
	l.lastMinuteFail.Store(0)
	if l.subscribeFunc != nil {
		l.reSubscribeNewRpc(url)
	}
	l.log.Info("switch current rpc new url", "url", url)
	if !l.isInFallbackState {
		l.isInFallbackState = true
		l.recoverIfFirstRpcHealth()
	}
	return nil
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
			err := l.firstRpc.CallContext(l.ctx, &id, "eth_chainId")
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
		lastRpc := l.currentRpc
		l.currentRpc = l.firstRpc
		lastRpc.Close()
		l.lastMinuteFail.Store(0)
		l.currentIndex = 0
		l.isInFallbackState = false
		if l.subscribeFunc != nil {
			l.reSubscribeNewRpc(l.urlList[0])
		}
		l.log.Info("recover current rpc to first rpc", "url", l.urlList[0])
	}()
}

func (l *FallbackClient) RegisterSubscribeFunc(f func() (event.Subscription, error), l1HeadsSub *ethereum.Subscription) {
	l.subscribeFunc = f
	l.l1HeadsSub = l1HeadsSub
}

func (l *FallbackClient) validateRpc(newRpc client.RPC) error {
	chainID, err := l.ChainID(l.ctx, newRpc)
	if err != nil {
		return err
	}
	if l.l1ChainId.Cmp(chainID) != 0 {
		return fmt.Errorf("incorrect L1 RPC chain id %d, expected %d", chainID, l.l1ChainId)
	}
	l1GenesisBlockRef, err := l.l1BlockRefByNumber(l.ctx, l.l1Block.Number, newRpc)
	if err != nil {
		return err
	}
	if l1GenesisBlockRef.Hash != l.l1Block.Hash {
		return fmt.Errorf("incorrect L1 genesis block hash %s, expected %s", l1GenesisBlockRef.Hash, l.l1Block.Hash)
	}
	return nil
}

func (l *FallbackClient) ChainID(ctx context.Context, rpc client.RPC) (*big.Int, error) {
	var id hexutil.Big
	err := rpc.CallContext(ctx, &id, "eth_chainId")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&id), nil
}

func (l *FallbackClient) l1BlockRefByNumber(ctx context.Context, number uint64, newRpc client.RPC) (*rpcHeader, error) {
	var header *rpcHeader
	err := newRpc.CallContext(ctx, &header, "eth_getBlockByNumber", numberID(number).Arg(), false) // headers are just blocks without txs
	if err != nil {
		return nil, err
	}
	return header, nil
}
