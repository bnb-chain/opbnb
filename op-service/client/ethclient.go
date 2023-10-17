package client

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// DialEthClientWithTimeout attempts to dial the L1 provider using the provided
// URL. If the dial doesn't complete within defaultDialTimeout seconds, this
// method will return an error.
func DialEthClientWithTimeout(ctx context.Context, url string, timeout time.Duration) (*ethclient.Client, error) {
	ctxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return ethclient.DialContext(ctxt, url)
}

const BatcherFallbackThreshold int64 = 10
const ProposerFallbackThreshold int64 = 3
const TxmgrFallbackThreshold int64 = 3

// DialEthClientWithTimeoutAndFallback will try to dial within the timeout period and create an EthClient.
// If the URL is a multi URL, then a fallbackClient will be created to add the fallback capability to the client
func DialEthClientWithTimeoutAndFallback(ctx context.Context, url string, timeout time.Duration, l log.Logger, fallbackThreshold int64, m FallbackClientMetricer) (EthClient, error) {
	ctxt, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	isMultiUrl, urlList := MultiUrlParse(url)
	if isMultiUrl {
		firstEthClient, err := ethclient.DialContext(ctxt, urlList[0])
		if err != nil {
			return nil, err
		}
		fallbackClient := NewFallbackClient(firstEthClient, urlList, l, fallbackThreshold, m, func(url string) (EthClient, error) {
			ctxtIn, cancelIn := context.WithTimeout(ctx, timeout)
			defer cancelIn()
			ethClientNew, err := ethclient.DialContext(ctxtIn, url)
			if err != nil {
				return nil, err
			}
			return ethClientNew, nil
		})
		return fallbackClient, nil
	}

	return ethclient.DialContext(ctxt, url)
}

type EthClient interface {
	ChainID(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	BlockNumber(ctx context.Context) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	Close()
}
