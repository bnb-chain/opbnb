package client

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func MultiUrlParse(url string) (isMultiUrl bool, urlList []string) {
	if strings.Contains(url, ",") {
		return true, strings.Split(url, ",")
	}
	return false, []string{}
}

type FallbackClient struct {
	firstClient       EthClient
	urlList           []string
	clientInitFunc    func(url string) (EthClient, error)
	lastMinuteFail    atomic.Int64
	currentClient     EthClient
	currentIndex      int
	mx                sync.Mutex
	log               log.Logger
	isInFallbackState bool
}

func NewFallbackClient(rpc EthClient, urlList []string, log log.Logger, clientInitFunc func(url string) (EthClient, error)) EthClient {
	fallbackClient := &FallbackClient{firstClient: rpc, urlList: urlList, log: log, clientInitFunc: clientInitFunc, currentClient: rpc, currentIndex: 0}
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

func (l *FallbackClient) BlockNumber(ctx context.Context) (uint64, error) {
	number, err := l.currentClient.BlockNumber(ctx)
	if err != nil {
		l.handleErr(err)
	}
	return number, err
}

func (l *FallbackClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := l.currentClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		l.handleErr(err)
	}
	return receipt, err
}

func (l *FallbackClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	err := l.currentClient.SendTransaction(ctx, tx)
	if err != nil {
		l.handleErr(err)
	}
	return err
}

func (l *FallbackClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	tipCap, err := l.currentClient.SuggestGasTipCap(ctx)
	if err != nil {
		l.handleErr(err)
	}
	return tipCap, err
}

func (l *FallbackClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	at, err := l.currentClient.PendingNonceAt(ctx, account)
	if err != nil {
		l.handleErr(err)
	}
	return at, err
}

func (l *FallbackClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	estimateGas, err := l.currentClient.EstimateGas(ctx, msg)
	if err != nil {
		l.handleErr(err)
	}
	return estimateGas, err
}

func (l *FallbackClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	contract, err := l.currentClient.CallContract(ctx, call, blockNumber)
	if err != nil {
		l.handleErr(err)
	}
	return contract, err
}

func (l *FallbackClient) Close() {
	l.currentClient.Close()
	if l.currentClient != l.firstClient {
		l.firstClient.Close()
	}
}

func (l *FallbackClient) ChainID(ctx context.Context) (*big.Int, error) {
	id, err := l.currentClient.ChainID(ctx)
	if err != nil {
		l.handleErr(err)
	}
	return id, err
}

func (l *FallbackClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	balanceAt, err := l.currentClient.BalanceAt(ctx, account, blockNumber)
	if err != nil {
		l.handleErr(err)
	}
	return balanceAt, err
}

func (l *FallbackClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	headerByNumber, err := l.currentClient.HeaderByNumber(ctx, number)
	if err != nil {
		l.handleErr(err)
	}
	return headerByNumber, err
}

func (l *FallbackClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	storageAt, err := l.currentClient.StorageAt(ctx, account, key, blockNumber)
	if err != nil {
		l.handleErr(err)
	}
	return storageAt, err
}

func (l *FallbackClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	codeAt, err := l.currentClient.CodeAt(ctx, account, blockNumber)
	if err != nil {
		l.handleErr(err)
	}
	return codeAt, err
}

func (l *FallbackClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	nonceAt, err := l.currentClient.NonceAt(ctx, account, blockNumber)
	if err != nil {
		l.handleErr(err)
	}
	return nonceAt, err
}

func (l *FallbackClient) handleErr(err error) {
	if err == rpc.ErrNoResult {
		return
	}
	errCount := l.lastMinuteFail.Add(1)
	if errCount > 10 {
		l.switchCurrentClient()
	}
}

func (l *FallbackClient) switchCurrentClient() {
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
	newClient, err := l.clientInitFunc(url)
	if err != nil {
		log.Error("fallback client switch current client fail", "url", url, "err", err)
		return
	}
	if l.currentClient != l.firstClient {
		l.currentClient.Close()
	}
	l.lastMinuteFail.Store(0)
	l.currentClient = newClient
	log.Info("switch current client new url", "url", url)
	if !l.isInFallbackState {
		l.isInFallbackState = true
		l.recoverIfFirstRpcHealth()
	}
}

func (l *FallbackClient) recoverIfFirstRpcHealth() {
	go func() {
		count := 0
		for {
			_, err := l.firstClient.ChainID(context.Background())
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
		l.currentClient.Close()
		l.lastMinuteFail.Store(0)
		l.currentClient = l.firstClient
		l.currentIndex = 0
		l.isInFallbackState = false
		log.Info("recover current client to first client", "url", l.urlList[0])
	}()
}
