package sources

import (
	"context"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func makeTestBSCBlobSidecar(blockHash common.Hash, blobs []eth.Blob) ([]eth.IndexedBlobHash, *eth.BSCBlobSidecar) {
	commitments := []eth.Bytes48{}
	proofs := []eth.Bytes48{}
	ibhs := []eth.IndexedBlobHash{}
	for i, blob := range blobs {
		commit, _ := kzg4844.BlobToCommitment(kzg4844.Blob(blob))
		proof, _ := kzg4844.ComputeBlobProof(kzg4844.Blob(blob), commit)
		hash := eth.KZGToVersionedHash(commit)
		commitments = append(commitments, eth.Bytes48(commit))
		proofs = append(proofs, eth.Bytes48(proof))
		ibhs = append(ibhs, eth.IndexedBlobHash{
			Index: uint64(i),
			Hash:  hash,
		})
	}

	sidecar := eth.BSCBlobSidecar{
		BlockHash:   blockHash,
		BlockNumber: &hexutil.Big{},
		BSCBlobTxSidecar: eth.BSCBlobTxSidecar{
			Blobs:       blobs,
			Commitments: commitments,
			Proofs:      proofs,
		},
	}
	return ibhs, &sidecar
}

func TestValidateBlobSidecars(t *testing.T) {
	blockHash := common.BytesToHash([]byte{1})
	blobs := []eth.Blob{}
	blob1 := eth.Blob{}
	blob1[0] = 1
	blob2 := eth.Blob{}
	blob2[0] = 2
	blobs = append(blobs, blob1)
	blobs = append(blobs, blob2)
	ibhs, sidecar := makeTestBSCBlobSidecar(blockHash, blobs)

	sidecars := eth.BSCBlobSidecars{sidecar}
	ref := eth.L1BlockRef{
		Hash: blockHash,
	}
	validatedSidecars, err := validateBlobSidecars(sidecars, ref)
	require.NoError(t, err)
	vBlob1, ok := validatedSidecars[ibhs[0].Hash]
	require.Equal(t, *vBlob1, blob1)
	require.Equal(t, ok, true)
	vBlob2, ok := validatedSidecars[ibhs[1].Hash]
	require.Equal(t, *vBlob2, blob2)
	require.Equal(t, ok, true)
	_, ok = validatedSidecars[common.Hash{}]
	require.Equal(t, ok, false)

	// mangle block hash to make sure it's detected
	ref = eth.L1BlockRef{}
	_, err = validateBlobSidecars(sidecars, ref)
	require.ErrorIs(t, err, ethereum.NotFound)
	// mangle blob to make sure it's detected
	sidecars[0].BSCBlobTxSidecar.Blobs[0][11]++
	_, err = validateBlobSidecars(sidecars, ref)
	require.Error(t, err)
	// mangle commitment to make sure it's detected
	sidecars[0].BSCBlobTxSidecar.Commitments[0][11]++
	_, err = validateBlobSidecars(sidecars, ref)
	require.Error(t, err)
	// mangle proof to make sure it's detected
	sidecars[0].BSCBlobTxSidecar.Proofs[0][11]++
	_, err = validateBlobSidecars(sidecars, ref)
	require.Error(t, err)
}

func TestBSCBlobClient(t *testing.T) {
	blockHash := common.BytesToHash([]byte{1})
	blobs := []eth.Blob{}
	blob1 := eth.Blob{}
	blob1[0] = 1
	blob2 := eth.Blob{}
	blob2[0] = 2
	blobs = append(blobs, blob1)
	blobs = append(blobs, blob2)
	ibhs, sidecar := makeTestBSCBlobSidecar(blockHash, blobs)
	sidecars := eth.BSCBlobSidecars{sidecar}
	ref := eth.L1BlockRef{
		Hash: blockHash,
	}

	m := new(mockRPC)
	ctx := context.Background()
	m.On("CallContext", ctx, new(eth.BSCBlobSidecars),
		"eth_getBlobSidecars", []any{"0x0"}).Run(func(args mock.Arguments) {
		*args[1].(*eth.BSCBlobSidecars) = sidecars
	}).Return([]error{nil})
	bscBlobClient := NewBSCBlobClient([]client.RPC{m})

	gotBlobs, err := bscBlobClient.GetBlobs(ctx, ref, ibhs)
	require.NoError(t, err)
	require.Equal(t, len(gotBlobs), 2)
	require.Equal(t, *gotBlobs[0], blob1)
	require.Equal(t, *gotBlobs[1], blob2)

	// mangle block hash to make sure it's detected
	_, err = bscBlobClient.GetBlobs(ctx, eth.L1BlockRef{}, ibhs)
	println(err)
	require.ErrorIs(t, err, ethereum.NotFound)

	// mangle blob hash to make sure it's detected
	ibhs[0].Hash[10]++
	_, err = bscBlobClient.GetBlobs(ctx, eth.L1BlockRef{}, ibhs)
	require.ErrorIs(t, err, ethereum.NotFound)
}
