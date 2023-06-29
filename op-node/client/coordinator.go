package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
)

// CoordinatorClient is a client for the coordinator RPC.
type CoordinatorClient struct {
	sequencerId string
	rpc         *rpc.Client
}

// NewCoordinatorClient creates a new client for the coordinator RPC.
func NewCoordinatorClient(url string, sequencerId string) (*CoordinatorClient, error) {
	rpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &CoordinatorClient{
		sequencerId: sequencerId,
		rpc:         rpc,
	}, nil
}

// RequestBuildingBlock is called by the sequencer to request a building block when using coordinator-mode.
func (c *CoordinatorClient) RequestBuildingBlock() error {
	var respErr error
	err := c.rpc.Call(respErr, "coordinator_requestBuildingBlock", c.sequencerId)
	if err != nil {
		return fmt.Errorf("failed to call coordinator_requestBuildingBlock: %w", err)
	}
	if respErr != nil {
		return fmt.Errorf("coordinator_requestBuildingBlock refused request: %w", respErr)
	}
	return nil
}
