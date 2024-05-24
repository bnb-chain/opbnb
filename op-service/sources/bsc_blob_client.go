package sources

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

type BSCBlobClient struct {
	// BSCBlobClient will rotate client.RPC in pool whenever a client runs into an error or return nil while fetching blobs
	pool *ClientPool[client.RPC]
}

func NewBSCBlobClient(clients []client.RPC) *BSCBlobClient {
	return &BSCBlobClient{
		pool: NewClientPool[client.RPC](clients...),
	}
}

func (s *BSCBlobClient) GetBlobs(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.Blob, error) {
	if len(hashes) == 0 {
		return []*eth.Blob{}, nil
	}

	blobSidecars, err := s.GetBlobSidecars(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob sidecars for L1BlockRef %s: %w", ref, err)
	}

	validatedBlobs, err := validateBlobSidecars(blobSidecars, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to validate blob sidecars for L1BlockRef %s: %w", ref, err)
	}

	blobs := make([]*eth.Blob, len(hashes))
	for i, indexedBlobHash := range hashes {
		blob, ok := validatedBlobs[indexedBlobHash.Hash]
		if !ok {
			return nil, fmt.Errorf("blob sidecars fetched from rpc mismatched with expected hash %s for L1BlockRef %s :%w", indexedBlobHash.Hash, ref, ethereum.NotFound)
		}
		blobs[i] = blob
	}
	return blobs, nil
}

func (s *BSCBlobClient) GetBlobSidecars(ctx context.Context, ref eth.L1BlockRef) (eth.BSCBlobSidecars, error) {
	var errs []error
	for i := 0; i < s.pool.Len(); i++ {
		var blobSidecars eth.BSCBlobSidecars

		f := s.pool.Get()
		err := f.CallContext(ctx, &blobSidecars, "eth_getBlobSidecars", numberID(ref.Number).Arg())
		if err != nil {
			s.pool.MoveToNext()
			errs = append(errs, err)
		} else {
			if len(blobSidecars) == 0 {
				err = ethereum.NotFound
				errs = append(errs, err)
				s.pool.MoveToNext()
			} else {
				return blobSidecars, nil
			}
		}
	}
	return nil, errors.Join(errs...)
}

func validateBlobSidecars(blobSidecars eth.BSCBlobSidecars, ref eth.L1BlockRef) (map[common.Hash]*eth.Blob, error) {
	if len(blobSidecars) == 0 {
		return nil, fmt.Errorf("invalidate api response, blob sidecars of block %s are empty: %w", ref.Hash, ethereum.NotFound)
	}
	blobsMap := make(map[common.Hash]*eth.Blob)
	for _, blobSidecar := range blobSidecars {
		if blobSidecar.BlockNumber.ToInt().Cmp(big.NewInt(0).SetUint64(ref.Number)) != 0 {
			return nil, fmt.Errorf("invalidate api response of tx %s, expect block number %d, got %d", blobSidecar.TxHash, ref.Number, blobSidecar.BlockNumber.ToInt().Uint64())
		}
		if blobSidecar.BlockHash.Cmp(ref.Hash) != 0 {
			return nil, fmt.Errorf("invalidate api response of tx %s, expect block hash %s, got %s :%w", blobSidecar.TxHash, ref.Hash, blobSidecar.BlockHash, ethereum.NotFound)
		}
		if len(blobSidecar.Blobs) == 0 || len(blobSidecar.Blobs) != len(blobSidecar.Commitments) || len(blobSidecar.Blobs) != len(blobSidecar.Proofs) {
			return nil, fmt.Errorf("invalidate api response of tx %s,idx:%d, len of blobs(%d)/commitments(%d)/proofs(%d) is not equal or is 0", blobSidecar.TxHash, blobSidecar.TxIndex, len(blobSidecar.Blobs), len(blobSidecar.Commitments), len(blobSidecar.Proofs))
		}

		for i := 0; i < len(blobSidecar.Blobs); i++ {
			// confirm blob data is valid by verifying its proof against the commitment
			if err := eth.VerifyBlobProof(&blobSidecar.Blobs[i], kzg4844.Commitment(blobSidecar.Commitments[i]), kzg4844.Proof(blobSidecar.Proofs[i])); err != nil {
				return nil, fmt.Errorf("blob of tx %s at index %d failed verification: %w", blobSidecar.TxHash, i, err)
			}
			// the blob's kzg commitment hashes
			hash := eth.KZGToVersionedHash(kzg4844.Commitment(blobSidecar.Commitments[i]))
			blobsMap[hash] = &blobSidecar.Blobs[i]
		}
	}
	return blobsMap, nil
}
