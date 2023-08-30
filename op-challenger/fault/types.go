package fault

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNegativeIndex = errors.New("index cannot be negative")
	ErrIndexTooLarge = errors.New("index is larger than the maximum index")
)

// StepCallData encapsulates the data needed to perform a step.
type StepCallData struct {
	StateIndex uint64
	ClaimIndex uint64
	IsAttack   bool
	StateData  []byte
	Proof      []byte
}

// TraceProvider is a generic way to get a claim value at a specific
// step in the trace.
// The [AlphabetProvider] is a minimal implementation of this interface.
type TraceProvider interface {
	Get(i uint64) (common.Hash, error)
	GetPreimage(i uint64) ([]byte, error)
}

// ClaimData is the core of a claim. It must be unique inside a specific game.
type ClaimData struct {
	Value common.Hash
	Position
}

func (c *ClaimData) ValueBytes() [32]byte {
	responseBytes := c.Value.Bytes()
	var responseArr [32]byte
	copy(responseArr[:], responseBytes[:32])
	return responseArr
}

// Claim extends ClaimData with information about the relationship between two claims.
// It uses ClaimData to break cyclicity without using pointers.
// If the position of the game is Depth 0, IndexAtDepth 0 it is the root claim
// and the Parent field is empty & meaningless.
type Claim struct {
	ClaimData
	Parent ClaimData
	// Location of the claim & it's parent inside the contract. Does not exist
	// for claims that have not made it to the contract.
	ContractIndex       int
	ParentContractIndex int
}

// IsRoot returns true if this claim is the root claim.
func (c *Claim) IsRoot() bool {
	return c.Position.IsRootPosition()
}

// DefendsParent returns true if the the claim is a defense (i.e. goes right) of the
// parent. It returns false if the claim is an attack (i.e. goes left) of the parent.
func (c *Claim) DefendsParent() bool {
	return (c.IndexAtDepth() >> 1) != c.Parent.IndexAtDepth()
}

// Responder takes a response action & executes.
// For full op-challenger this means executing the transaction on chain.
type Responder interface {
	Respond(ctx context.Context, response Claim) error
	Step(ctx context.Context, stepData StepCallData) error
}
