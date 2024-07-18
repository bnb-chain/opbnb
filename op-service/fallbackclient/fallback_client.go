package fallbackclient

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

func MultiUrlParse(url string) (isMultiUrl bool, urlList []string) {
	if strings.Contains(url, ",") {
		return true, strings.Split(url, ",")
	}
	return false, []string{}
}

type FallbackClientMetricer interface {
	RecordL1UrlSwitchEvt(url string)
}

type FallbackClientMetrics struct {
	l1UrlSwitchEvt opmetrics.EventVec
}

func (f *FallbackClientMetrics) RecordL1UrlSwitchEvt(url string) {
	f.l1UrlSwitchEvt.Record(url)
}

func NewFallbackClientMetrics(ns string, factory opmetrics.Factory) *FallbackClientMetrics {
	return &FallbackClientMetrics{
		l1UrlSwitchEvt: opmetrics.NewEventVec(factory, ns, "", "l1_url_switch", "l1 url switch", []string{"url_idx"}),
	}
}

// FallbackClient is an Client, it can automatically switch to the next l1 endpoint
// when there is a problem with the current l1 endpoint
// and automatically switch back after the first l1 endpoint recovers.
type FallbackClient struct {
	// firstClient is created by the first of the l1 urls, it should be used first in a healthy state
	firstClient       Client
	urlList           []string
	clientInitFunc    func(url string) (Client, error)
	lastMinuteFail    atomic.Int64
	currentClient     atomic.Pointer[Client]
	currentIndex      int
	mx                sync.Mutex
	log               log.Logger
	isInFallbackState bool
	// fallbackThreshold specifies how many errors have occurred in the past 1 minute to trigger the switching logic
	fallbackThreshold int64
	isClose           chan struct{}
	metrics           FallbackClientMetricer
}

var _ Client = (*FallbackClient)(nil)

// NewFallbackClient returns a new FallbackClient.
func NewFallbackClient(rpc Client, urlList []string, log log.Logger, fallbackThreshold int64, m FallbackClientMetricer, clientInitFunc func(url string) (Client, error)) Client {
	fallbackClient := &FallbackClient{
		firstClient:       rpc,
		urlList:           urlList,
		log:               log,
		clientInitFunc:    clientInitFunc,
		currentIndex:      0,
		fallbackThreshold: fallbackThreshold,
		metrics:           m,
		isClose:           make(chan struct{}),
	}
	fallbackClient.currentClient.Store(&rpc)
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				log.Debug("FallbackClient clear lastMinuteFail 0")
				fallbackClient.lastMinuteFail.Store(0)
			case <-fallbackClient.isClose:
				return
			default:
				if fallbackClient.lastMinuteFail.Load() >= fallbackClient.fallbackThreshold {
					fallbackClient.switchCurrentClient()
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return fallbackClient
}

func (l *FallbackClient) BlockNumber(ctx context.Context) (uint64, error) {
	number, err := (*l.currentClient.Load()).BlockNumber(ctx)
	if err != nil {
		l.handleErr(err, "BlockNumber")
	}
	return number, err
}

func (l *FallbackClient) NetworkID(ctx context.Context) (*big.Int, error) {
	id, err := (*l.currentClient.Load()).NetworkID(ctx)
	if err != nil {
		l.handleErr(err, "NetworkID")
	}
	return id, err
}

func (l *FallbackClient) PeerCount(ctx context.Context) (uint64, error) {
	count, err := (*l.currentClient.Load()).PeerCount(ctx)
	if err != nil {
		l.handleErr(err, "PeerCount")
	}
	return count, err
}

func (l *FallbackClient) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	count, err := (*l.currentClient.Load()).TransactionCount(ctx, blockHash)
	if err != nil {
		l.handleErr(err, "TransactionCount")
	}
	return count, err
}

func (l *FallbackClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	count, err := (*l.currentClient.Load()).PendingTransactionCount(ctx)
	if err != nil {
		l.handleErr(err, "PendingTransactionCount")
	}
	return count, err
}

func (l *FallbackClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	progress, err := (*l.currentClient.Load()).SyncProgress(ctx)
	if err != nil {
		l.handleErr(err, "SyncProgress")
	}
	return progress, err
}

func (l *FallbackClient) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := (*l.currentClient.Load()).TransactionByHash(ctx, hash)
	if err != nil {
		l.handleErr(err, "TransactionByHash")
	}
	return tx, isPending, err
}

func (l *FallbackClient) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	tx, err := (*l.currentClient.Load()).TransactionInBlock(ctx, blockHash, index)
	if err != nil {
		l.handleErr(err, "TransactionInBlock")
	}
	return tx, err
}

func (l *FallbackClient) TransactionSender(ctx context.Context, tx *types.Transaction, block common.Hash, index uint) (common.Address, error) {
	sender, err := (*l.currentClient.Load()).TransactionSender(ctx, tx, block, index)
	if err != nil {
		l.handleErr(err, "TransactionSender")
	}
	return sender, err
}

func (l *FallbackClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := (*l.currentClient.Load()).TransactionReceipt(ctx, txHash)
	if err != nil {
		l.handleErr(err, "TransactionReceipt")
	}
	return receipt, err
}

func (l *FallbackClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	err := (*l.currentClient.Load()).SendTransaction(ctx, tx)
	if err != nil {
		l.handleErr(err, "SendTransaction")
	}
	return err
}

func (l *FallbackClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	tipCap, err := (*l.currentClient.Load()).SuggestGasTipCap(ctx)
	if err != nil {
		l.handleErr(err, "SuggestGasTipCap")
	}
	return tipCap, err
}

func (l *FallbackClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	price, err := (*l.currentClient.Load()).SuggestGasPrice(ctx)
	if err != nil {
		l.handleErr(err, "SuggestGasPrice")
	}
	return price, err
}

func (l *FallbackClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	at, err := (*l.currentClient.Load()).PendingNonceAt(ctx, account)
	if err != nil {
		l.handleErr(err, "PendingNonceAt")
	}
	return at, err
}

func (l *FallbackClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	estimateGas, err := (*l.currentClient.Load()).EstimateGas(ctx, msg)
	if err != nil {
		l.handleErr(err, "EstimateGas")
	}
	return estimateGas, err
}

func (l *FallbackClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	contract, err := (*l.currentClient.Load()).CallContract(ctx, call, blockNumber)
	if err != nil {
		l.handleErr(err, "CallContract")
	}
	return contract, err
}

func (l *FallbackClient) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	contract, err := (*l.currentClient.Load()).PendingCallContract(ctx, call)
	if err != nil {
		l.handleErr(err, "PendingCallContract")
	}
	return contract, err
}

func (l *FallbackClient) CallContractAtHash(ctx context.Context, call ethereum.CallMsg, blockHash common.Hash) ([]byte, error) {
	contract, err := (*l.currentClient.Load()).CallContractAtHash(ctx, call, blockHash)
	if err != nil {
		l.handleErr(err, "CallContractAtHash")
	}
	return contract, err
}

func (l *FallbackClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	block, err := (*l.currentClient.Load()).BlockByNumber(ctx, number)
	if err != nil {
		l.handleErr(err, "BlockByNumber")
	}
	return block, err
}

func (l *FallbackClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	block, err := (*l.currentClient.Load()).BlockByHash(ctx, hash)
	if err != nil {
		l.handleErr(err, "BlockByHash")
	}
	return block, err
}

func (l *FallbackClient) Close() {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.isClose <- struct{}{}
	currentClient := *l.currentClient.Load()
	currentClient.Close()
	if currentClient != l.firstClient {
		l.firstClient.Close()
	}
}

func (l *FallbackClient) ChainID(ctx context.Context) (*big.Int, error) {
	id, err := (*l.currentClient.Load()).ChainID(ctx)
	if err != nil {
		l.handleErr(err, "ChainID")
	}
	return id, err
}

func (l *FallbackClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	balanceAt, err := (*l.currentClient.Load()).BalanceAt(ctx, account, blockNumber)
	if err != nil {
		l.handleErr(err, "BalanceAt")
	}
	return balanceAt, err
}

func (l *FallbackClient) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	balanceAt, err := (*l.currentClient.Load()).PendingBalanceAt(ctx, account)
	if err != nil {
		l.handleErr(err, "PendingBalanceAt")
	}
	return balanceAt, err
}

func (l *FallbackClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	headerByNumber, err := (*l.currentClient.Load()).HeaderByNumber(ctx, number)
	if err != nil {
		l.handleErr(err, "HeaderByNumber")
	}
	return headerByNumber, err
}

func (l *FallbackClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	header, err := (*l.currentClient.Load()).HeaderByHash(ctx, hash)
	if err != nil {
		l.handleErr(err, "HeaderByHash")
	}
	return header, err
}

func (l *FallbackClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	storageAt, err := (*l.currentClient.Load()).StorageAt(ctx, account, key, blockNumber)
	if err != nil {
		l.handleErr(err, "StorageAt")
	}
	return storageAt, err
}

func (l *FallbackClient) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	storageAt, err := (*l.currentClient.Load()).PendingStorageAt(ctx, account, key)
	if err != nil {
		l.handleErr(err, "PendingStorageAt")
	}
	return storageAt, err
}

func (l *FallbackClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	codeAt, err := (*l.currentClient.Load()).CodeAt(ctx, account, blockNumber)
	if err != nil {
		l.handleErr(err, "CodeAt")
	}
	return codeAt, err
}

func (l *FallbackClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	codeAt, err := (*l.currentClient.Load()).PendingCodeAt(ctx, account)
	if err != nil {
		l.handleErr(err, "PendingCodeAt")
	}
	return codeAt, err
}

func (l *FallbackClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	nonceAt, err := (*l.currentClient.Load()).NonceAt(ctx, account, blockNumber)
	if err != nil {
		l.handleErr(err, "NonceAt")
	}
	return nonceAt, err
}

func (l *FallbackClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	logs, err := (*l.currentClient.Load()).FilterLogs(ctx, q)
	if err != nil {
		l.handleErr(err, "FilterLogs")
	}
	return logs, err
}

func (l *FallbackClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	sub, err := (*l.currentClient.Load()).SubscribeFilterLogs(ctx, q, ch)
	if err != nil {
		l.handleErr(err, "SubscribeFilterLogs")
	}
	return sub, err
}

func (l *FallbackClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	sub, err := (*l.currentClient.Load()).SubscribeNewHead(ctx, ch)
	if err != nil {
		l.handleErr(err, "SubscribeNewHead")
	}
	return sub, err
}

func (l *FallbackClient) handleErr(err error, methodName string) {
	if errors.Is(err, rpc.ErrNoResult) {
		return
	}
	if errors.Is(err, ethereum.NotFound) {
		return
	}
	var targetErr rpc.Error
	if errors.As(err, &targetErr) {
		return
	}
	log.Debug("fallback client fail count+1", "err", err, "methodName", methodName)
	l.lastMinuteFail.Add(1)
}

func (l *FallbackClient) switchCurrentClient() {
	l.mx.Lock()
	defer l.mx.Unlock()
	if l.lastMinuteFail.Load() <= l.fallbackThreshold {
		return
	}
	//Use defer to ensure that recoverIfFirstRpcHealth will always be executed regardless of the circumstances.
	defer func() {
		if !l.isInFallbackState {
			l.isInFallbackState = true
			l.recoverIfFirstRpcHealth()
		}
	}()
	l.currentIndex++
	if l.currentIndex >= len(l.urlList) {
		l.log.Error("the fallback client has tried all urls")
		return
	}
	l.metrics.RecordL1UrlSwitchEvt(strconv.Itoa(l.currentIndex))
	url := l.urlList[l.currentIndex]
	newClient, err := l.clientInitFunc(url)
	if err != nil {
		l.log.Error("the fallback client failed to switch the current client", "url", url, "err", err)
		return
	}
	lastClient := *l.currentClient.Load()
	l.currentClient.Store(&newClient)
	if lastClient != l.firstClient {
		lastClient.Close()
	}
	l.lastMinuteFail.Store(0)
	l.log.Info("switched current rpc to new url", "url", url)
}

func (l *FallbackClient) recoverIfFirstRpcHealth() {
	go func() {
		count := 0
		for {
			_, err := l.firstClient.ChainID(context.Background())
			if err != nil {
				count = 0
				time.Sleep(3 * time.Second)
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
		lastClient := *l.currentClient.Load()
		l.currentClient.Store(&l.firstClient)
		if lastClient != l.firstClient {
			lastClient.Close()
		}
		l.lastMinuteFail.Store(0)
		l.currentIndex = 0
		l.isInFallbackState = false
		l.log.Info("recover the current client to the first client", "url", l.urlList[0])
	}()
}
