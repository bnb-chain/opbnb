package proposer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type LastGameInfo struct {
	Address *common.Address `json:"address,omitempty"`
	Idx     *big.Int        `json:"idx,omitempty"`
}

type LastGamePersistenceCache struct {
	lock sync.Mutex
	file string
}

func (c *LastGamePersistenceCache) loadFile() (*LastGameInfo, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	data, err := os.ReadFile(c.file)
	if errors.Is(err, os.ErrNotExist) {
		// persistedState.SequencerStarted == nil: SequencerState() will return StateUnset if no state is found
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("read config file (%v): %w", c.file, err)
	}
	var config LastGameInfo
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&config); err != nil {
		return nil, fmt.Errorf("invalid config file (%v): %w", c.file, err)
	}
	if config.Address == nil {
		return nil, fmt.Errorf("missing address value in config file (%v)", c.file)
	}
	return &config, nil
}

func (c *LastGamePersistenceCache) cacheFile(address *common.Address, idx *big.Int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	data, err := json.Marshal(&LastGameInfo{Address: address, Idx: idx})
	if err != nil {
		return fmt.Errorf("marshall new config: %w", err)
	}
	dir := filepath.Dir(c.file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir (%v): %w", c.file, err)
	}
	// Write the new content to a temp file first, then rename into place
	// Avoids corrupting the content if the disk is full or there are IO errors
	tmpFile := c.file + ".tmp"
	file, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open file (%v) for writing: %w", tmpFile, err)
	}
	defer file.Close() // Ensure file is closed even if write or sync fails
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
	if err := os.Rename(tmpFile, c.file); err != nil {
		return fmt.Errorf("rename temp config file to final destination: %w", err)
	}
	return nil
}

func NewLastGamePersistenceCache(file string) *LastGamePersistenceCache {
	return &LastGamePersistenceCache{file: file}
}
