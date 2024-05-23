package sources

import (
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

func makeTestBSCBlobSidecar(blockHash common.Hash, blobs []eth.Blob) ([]eth.IndexedBlobHash, *eth.BSCBlobSidecar) {
	commitments := []eth.Bytes48{}
	proofs := []eth.Bytes48{}
	hashs := []common.Hash{}
	ibhs := []eth.IndexedBlobHash{}
	for i, blob := range blobs {
		commit, _ := kzg4844.BlobToCommitment(kzg4844.Blob(blob))
		proof, _ := kzg4844.ComputeBlobProof(kzg4844.Blob(blob), commit)
		hash := eth.KZGToVersionedHash(commit)
		commitments = append(commitments, eth.Bytes48(commit))
		proofs = append(proofs, eth.Bytes48(proof))
		hashs = append(hashs, hash)
		ibhs = append(ibhs, eth.IndexedBlobHash{
			Index: uint64(i),
			Hash: hash,
		})
	}

	sidecar := eth.BSCBlobSidecar{
		BlockHash: blockHash,
		BlockNumber: &hexutil.Big{},
		BSCBlobTxSidecar: eth.BSCBlobTxSidecar{
			Blobs:          blobs,
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
	_, sidecar := makeTestBSCBlobSidecar(blockHash, blobs)

	sidecars := eth.BSCBlobSidecars{sidecar}
	ref := eth.L1BlockRef{
		Hash: blockHash,
	}
	validateBlobSidecars(sidecars, ref)


}
