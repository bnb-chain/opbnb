package dial

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// DefaultDialTimeout is a default timeout for dialing a client.
const DefaultDialTimeout = 1 * time.Minute
const defaultRetryCount = 30
const defaultRetryTime = 2 * time.Second

const BatcherFallbackThreshold int64 = 10
const ProposerFallbackThreshold int64 = 3
const TxmgrFallbackThreshold int64 = 3

// DialEthClientWithTimeout attempts to dial the L1 provider using the provided
// URL. If the dial doesn't complete within defaultDialTimeout seconds, this
// method will return an error.
func DialEthClientWithTimeout(ctx context.Context, timeout time.Duration, log log.Logger, url string) (client.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	c, err := dialRPCClientWithBackoff(ctx, log, url)
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(c), nil
}

// DialRollupClientWithTimeout attempts to dial the RPC provider using the provided URL.
// If the dial doesn't complete within timeout seconds, this method will return an error.
func DialRollupClientWithTimeout(ctx context.Context, timeout time.Duration, log log.Logger, url string) (*sources.RollupClient, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rpcCl, err := dialRPCClientWithBackoff(ctx, log, url)
	if err != nil {
		return nil, err
	}

	return sources.NewRollupClient(client.NewBaseRPCClient(rpcCl)), nil
}

// Dials a JSON-RPC endpoint repeatedly, with a backoff, until a client connection is established. Auth is optional.
func dialRPCClientWithBackoff(ctx context.Context, log log.Logger, addr string) (*rpc.Client, error) {
	bOff := retry.Fixed(defaultRetryTime)
	return retry.Do(ctx, defaultRetryCount, bOff, func() (*rpc.Client, error) {
		if !client.IsURLAvailable(addr) {
			log.Warn("failed to dial address, but may connect later", "addr", addr)
			return nil, fmt.Errorf("address unavailable (%s)", addr)
		}
		client, err := rpc.DialOptions(ctx, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to dial address (%s): %w", addr, err)
		}
		return client, nil
	})
}

// DialEthClientWithTimeoutAndFallback will try to dial within the timeout period and create an EthClient.
// If the URL is a multi URL, then a fallbackClient will be created to add the fallback capability to the client
func DialEthClientWithTimeoutAndFallback(ctx context.Context, url string, timeout time.Duration, l log.Logger, fallbackThreshold int64, m client.FallbackClientMetricer) (client.Client, error) {
	isMultiUrl, urlList := client.MultiUrlParse(url)
	if isMultiUrl {
		firstEthClient, err := DialEthClientWithTimeout(ctx, timeout, l, urlList[0])
		if err != nil {
			return nil, err
		}
		fallbackClient := client.NewFallbackClient(firstEthClient, urlList, l, fallbackThreshold, m, func(url string) (client.Client, error) {
			ethClientNew, err := DialEthClientWithTimeout(ctx, timeout, l, url)
			if err != nil {
				return nil, err
			}
			return ethClientNew, nil
		})
		return fallbackClient, nil
	}

	return DialEthClientWithTimeout(ctx, timeout, l, url)
}
