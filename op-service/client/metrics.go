package client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// InstrumentedClient is an Ethereum client that tracks
// Prometheus metrics for each call.
type InstrumentedClient struct {
	c EthClient
	m Metricer
}

type Metricer interface {
	RecordRPCClientRequest(method string) func(err error)
}

// NewInstrumentedClient creates a new instrumented client. It takes
// a concrete *rpc.Client to prevent people from passing in an already
// instrumented client.
func NewInstrumentedClient(c EthClient, m Metricer) EthClient {
	return &InstrumentedClient{
		c: c,
		m: m,
	}
}

func (ic *InstrumentedClient) Close() {
	ic.c.Close()
}

func (ic *InstrumentedClient) ChainID(ctx context.Context) (*big.Int, error) {
	return instrument2[*big.Int](ic.m, "eth_chainId", func() (*big.Int, error) {
		return ic.c.ChainID(ctx)
	})
}

func (ic *InstrumentedClient) BlockNumber(ctx context.Context) (uint64, error) {
	return instrument2[uint64](ic.m, "eth_blockNumber", func() (uint64, error) {
		return ic.c.BlockNumber(ctx)
	})
}

func (ic *InstrumentedClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return instrument2[*types.Header](ic.m, "eth_getHeaderByNumber", func() (*types.Header, error) {
		return ic.c.HeaderByNumber(ctx, number)
	})
}

func (ic *InstrumentedClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return instrument2[*types.Receipt](ic.m, "eth_getTransactionReceipt", func() (*types.Receipt, error) {
		return ic.c.TransactionReceipt(ctx, txHash)
	})
}

func (ic *InstrumentedClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return instrument2[*big.Int](ic.m, "eth_getBalance", func() (*big.Int, error) {
		return ic.c.BalanceAt(ctx, account, blockNumber)
	})
}

func (ic *InstrumentedClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	return instrument2[[]byte](ic.m, "eth_getStorageAt", func() ([]byte, error) {
		return ic.c.StorageAt(ctx, account, key, blockNumber)
	})
}

func (ic *InstrumentedClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return instrument2[[]byte](ic.m, "eth_getCode", func() ([]byte, error) {
		return ic.c.CodeAt(ctx, account, blockNumber)
	})
}

func (ic *InstrumentedClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return instrument2[uint64](ic.m, "eth_getTransactionCount", func() (uint64, error) {
		return ic.c.NonceAt(ctx, account, blockNumber)
	})
}

func (ic *InstrumentedClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return instrument2[uint64](ic.m, "eth_getTransactionCount", func() (uint64, error) {
		return ic.c.PendingNonceAt(ctx, account)
	})
}

func (ic *InstrumentedClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return instrument2[[]byte](ic.m, "eth_call", func() ([]byte, error) {
		return ic.c.CallContract(ctx, msg, blockNumber)
	})
}

func (ic *InstrumentedClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return instrument2[uint64](ic.m, "eth_estimateGas", func() (uint64, error) {
		return ic.c.EstimateGas(ctx, msg)
	})
}

func (ic *InstrumentedClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return instrument1(ic.m, "eth_sendRawTransaction", func() error {
		return ic.c.SendTransaction(ctx, tx)
	})
}

func (ic *InstrumentedClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return instrument2[*big.Int](ic.m, "eth_maxPriorityFeePerGas", func() (*big.Int, error) {
		return ic.c.SuggestGasTipCap(ctx)
	})
}

func (ic *InstrumentedClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return instrument2[*types.Block](ic.m, "eth_getBlockByNumber", func() (*types.Block, error) {
		return ic.c.BlockByNumber(ctx, number)
	})
}

func instrument1(m Metricer, name string, cb func() error) error {
	record := m.RecordRPCClientRequest(name)
	err := cb()
	record(err)
	return err
}

func instrument2[O any](m Metricer, name string, cb func() (O, error)) (O, error) {
	record := m.RecordRPCClientRequest(name)
	res, err := cb()
	record(err)
	return res, err
}
