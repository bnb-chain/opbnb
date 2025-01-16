package zk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ProofAccessor struct {
	dir    string
	l1Head eth.BlockID
}

func (a *ProofAccessor) GetProof(baseBlock uint64, challengeBlock uint64) ([]byte, error) {
	if err := os.MkdirAll(a.dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create proofs directory %v: %w", a.dir, err)
	}
	proofFile := filepath.Join(a.dir, "proofs")
	err := mockFile(proofFile)
	if err != nil {
		return nil, fmt.Errorf("could not mock file,err: %w", err)
	}
	data, err := os.ReadFile(proofFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("read config file (%v): %w", proofFile, err)
	}
	return data, nil
}

func mockFile(filePath string) error {
	tmpFile := filePath + ".tmp"
	file, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open file (%v) for writing: %w", tmpFile, err)
	}
	defer file.Close() // Ensure file is closed even if write or sync fails
	data := hexutil.MustDecode("0x1")
	if _, err = file.Write(data); err != nil {
		return fmt.Errorf("write new config to temp file (%v): %w", tmpFile, err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync new config temp file (%v): %w", tmpFile, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close new config temp file (%v): %w", tmpFile, err)
	}
	// Rename to replace the previous file
	if err := os.Rename(tmpFile, filePath); err != nil {
		return fmt.Errorf("rename temp config file to final destination: %w", err)
	}
	return nil
}

func NewProofAccessor(dir string, l1Head eth.BlockID) (*ProofAccessor, error) {
	return &ProofAccessor{
		dir:    dir,
		l1Head: l1Head,
	}, nil
}
