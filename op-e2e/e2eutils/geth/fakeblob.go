package geth

import (
	"context"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	geth "github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/rpc"
)

type BlobService struct {
	beacon  Beacon
	backend *geth.EthAPIBackend
}

func blobAPIs(beacon Beacon, backend *geth.EthAPIBackend) []rpc.API {
	// Append all the local APIs and return
	return []rpc.API{
		{
			Namespace: "eth",
			Service:   NewBlobService(beacon, backend),
		},
	}
}

func NewBlobService(beacon Beacon, backend *geth.EthAPIBackend) *BlobService {
	return &BlobService{beacon: beacon, backend: backend}
}

func (b *BlobService) GetBlobSidecars(ctx context.Context, blockNumber rpc.BlockNumber) (*eth.BSCBlobSidecars, error) {
	//slot := (envelope.ExecutionPayload.Timestamp - f.eth.BlockChain().Genesis().Time()) / f.blockTime
	header, err := b.backend.HeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	slot := (header.Time - b.backend.Genesis().Time()) / 2
	bundle, err := b.beacon.LoadBlobsBundle(slot)
	if err != nil {
		return nil, err
	}
	if len(bundle.Blobs) == 0 {
		return nil, nil
	}
	var sidecars eth.BSCBlobSidecars
	bn := hexutil.Big(*big.NewInt(blockNumber.Int64()))
	oneSidecar := eth.BSCBlobSidecar{
		BlockNumber: &bn,
		BlockHash:   header.Hash(),
		TxIndex:     nil,
		TxHash:      common.Hash{},
	}
	blobSize := len(bundle.Blobs)
	oneSidecar.BSCBlobTxSidecar.Blobs = make([]eth.Blob, blobSize)
	oneSidecar.BSCBlobTxSidecar.Proofs = make([]eth.Bytes48, blobSize)
	oneSidecar.BSCBlobTxSidecar.Commitments = make([]eth.Bytes48, blobSize)
	for i, blob := range bundle.Blobs {
		copy(oneSidecar.BSCBlobTxSidecar.Blobs[i][:], blob)
		copy(oneSidecar.BSCBlobTxSidecar.Proofs[i][:], bundle.Proofs[i])
		copy(oneSidecar.BSCBlobTxSidecar.Commitments[i][:], bundle.Commitments[i])
	}
	sidecars = append(sidecars, &oneSidecar)

	return &sidecars, nil
}
