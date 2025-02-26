package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockID struct {
	Hash   common.Hash `json:"hash"`
	Number uint64      `json:"number"`
}

func (id BlockID) String() string {
	return fmt.Sprintf("%s:%d", id.Hash.String(), id.Number)
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (id BlockID) TerminalString() string {
	return fmt.Sprintf("%s:%d", id.Hash.TerminalString(), id.Number)
}

func ReceiptBlockID(r *types.Receipt) BlockID {
	return BlockID{Number: r.BlockNumber.Uint64(), Hash: r.BlockHash}
}

func HeaderBlockID(h *types.Header) BlockID {
	return BlockID{Number: h.Number.Uint64(), Hash: h.Hash()}
}

type L2BlockRef struct {
	Hash           common.Hash `json:"hash"`
	Number         uint64      `json:"number"`
	ParentHash     common.Hash `json:"parentHash"`
	Time           uint64      `json:"timestamp"`
	MilliTime      uint64      `json:"millitimestamp"`
	L1Origin       BlockID     `json:"l1origin"`
	SequenceNumber uint64      `json:"sequenceNumber"` // distance to first block of epoch
}

func (id L2BlockRef) String() string {
	return fmt.Sprintf("%s:%d", id.Hash.String(), id.Number)
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (id L2BlockRef) TerminalString() string {
	return fmt.Sprintf("%s:%d", id.Hash.TerminalString(), id.Number)
}

// TODO:
func (id L2BlockRef) MilliTimestamp() uint64 {
	if id.MilliTime != 0 {
		return id.MilliTime
	}
	return id.Time * 1000
}

func (id L2BlockRef) Timestamp() uint64 {
	if id.Time != 0 {
		return id.Time
	}
	return id.MilliTime / 1000
}

type L1BlockRef struct {
	Hash       common.Hash `json:"hash"`
	Number     uint64      `json:"number"`
	ParentHash common.Hash `json:"parentHash"`
	Time       uint64      `json:"timestamp"`
}

func (id L1BlockRef) String() string {
	return fmt.Sprintf("%s:%d", id.Hash.String(), id.Number)
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (id L1BlockRef) TerminalString() string {
	return fmt.Sprintf("%s:%d", id.Hash.TerminalString(), id.Number)
}

func (id L1BlockRef) ID() BlockID {
	return BlockID{
		Hash:   id.Hash,
		Number: id.Number,
	}
}

func (id L1BlockRef) ParentID() BlockID {
	n := id.ID().Number
	// Saturate at 0 with subtraction
	if n > 0 {
		n -= 1
	}
	return BlockID{
		Hash:   id.ParentHash,
		Number: n,
	}
}

func (id L2BlockRef) ID() BlockID {
	return BlockID{
		Hash:   id.Hash,
		Number: id.Number,
	}
}

func (id L2BlockRef) ParentID() BlockID {
	n := id.ID().Number
	// Saturate at 0 with subtraction
	if n > 0 {
		n -= 1
	}
	return BlockID{
		Hash:   id.ParentHash,
		Number: n,
	}
}

// IndexedBlobHash represents a blob hash that commits to a single blob confirmed in a block.  The
// index helps us avoid unnecessary blob to blob hash conversions to find the right content in a
// sidecar.
type IndexedBlobHash struct {
	Index uint64      // absolute index in the block, a.k.a. position in sidecar blobs array
	Hash  common.Hash // hash of the blob, used for consistency checks
}
