package fault

import (
	"errors"
)

var (
	// ErrClaimExists is returned when a claim already exists in the game state.
	ErrClaimExists = errors.New("claim exists in game state")

	// ErrClaimNotFound is returned when a claim does not exist in the game state.
	ErrClaimNotFound = errors.New("claim not found in game state")
)

// Game is an interface that represents the state of a dispute game.
type Game interface {
	// Put adds a claim into the game state.
	Put(claim Claim) error

	// PutAll adds a list of claims into the game state.
	PutAll(claims []Claim) error

	// Claims returns all of the claims in the game.
	Claims() []Claim

	// IsDuplicate returns true if the provided [Claim] already exists in the game state.
	IsDuplicate(claim Claim) bool

	// PreStateClaim gets the claim which commits to the pre-state of this specific claim.
	// This will return an error if it is called with a non-leaf claim.
	PreStateClaim(claim Claim) (Claim, error)

	// PostStateClaim gets the claim which commits to the post-state of this specific claim.
	// This will return an error if it is called with a non-leaf claim.
	PostStateClaim(claim Claim) (Claim, error)
}

type extendedClaim struct {
	self          Claim
	contractIndex int
	children      []ClaimData
}

// gameState is a struct that represents the state of a dispute game.
// The game state implements the [Game] interface.
type gameState struct {
	root   ClaimData
	claims map[ClaimData]*extendedClaim
	depth  uint64
}

// NewGameState returns a new game state.
// The provided [Claim] is used as the root node.
func NewGameState(root Claim, depth uint64) *gameState {
	claims := make(map[ClaimData]*extendedClaim)
	claims[root.ClaimData] = &extendedClaim{
		self:          root,
		contractIndex: 0,
		children:      make([]ClaimData, 0),
	}
	return &gameState{
		root:   root.ClaimData,
		claims: claims,
		depth:  depth,
	}
}

// PutAll adds a list of claims into the [Game] state.
// If any of the claims already exist in the game state, an error is returned.
func (g *gameState) PutAll(claims []Claim) error {
	for _, claim := range claims {
		if err := g.Put(claim); err != nil {
			return err
		}
	}
	return nil
}

// Put adds a claim into the game state.
func (g *gameState) Put(claim Claim) error {
	if claim.IsRoot() || g.IsDuplicate(claim) {
		return ErrClaimExists
	}
	if parent, ok := g.claims[claim.Parent]; !ok {
		return errors.New("no parent claim")
	} else {
		parent.children = append(parent.children, claim.ClaimData)
	}
	g.claims[claim.ClaimData] = &extendedClaim{
		self:          claim,
		contractIndex: claim.ContractIndex,
		children:      make([]ClaimData, 0),
	}
	return nil
}

func (g *gameState) IsDuplicate(claim Claim) bool {
	_, ok := g.claims[claim.ClaimData]
	return ok
}

func (g *gameState) Claims() []Claim {
	queue := []ClaimData{g.root}
	var out []Claim
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		queue = append(queue, g.getChildren(item)...)
		out = append(out, g.claims[item].self)
	}
	return out
}

func (g *gameState) getChildren(c ClaimData) []ClaimData {
	return g.claims[c].children
}

func (g *gameState) getParent(claim Claim) (Claim, error) {
	if claim.IsRoot() {
		return Claim{}, ErrClaimNotFound
	}
	if parent, ok := g.claims[claim.Parent]; !ok {
		return Claim{}, ErrClaimNotFound
	} else {
		return parent.self, nil
	}
}

func (g *gameState) PreStateClaim(claim Claim) (Claim, error) {
	// Do checks in PreStateClaim because these do not hold while walking the tree
	if claim.Depth() != int(g.depth) {
		return Claim{}, errors.New("Only leaf claims have pre or post state")
	}
	// If the claim is the far left most claim, the pre-state is pulled from the contracts & we can supply at contract index.
	if claim.IndexAtDepth() == 0 {
		return Claim{
			ContractIndex: -1,
		}, nil
	}
	return g.preStateClaim(claim)
}

// preStateClaim is the internal tree walker which does not do error handling
func (g *gameState) preStateClaim(claim Claim) (Claim, error) {
	parent, _ := g.getParent(claim)
	if claim.DefendsParent() {
		return parent, nil
	} else {
		return g.preStateClaim(parent)
	}
}

func (g *gameState) PostStateClaim(claim Claim) (Claim, error) {
	// Do checks in PostStateClaim because these do not hold while walking the tree
	if claim.Depth() != int(g.depth) {
		return Claim{}, errors.New("Only leaf claims have pre or post state")
	}
	return g.postStateClaim(claim)
}

// postStateClaim is the internal tree walker which does not do error handling
func (g *gameState) postStateClaim(claim Claim) (Claim, error) {
	parent, _ := g.getParent(claim)
	if claim.DefendsParent() {
		return g.postStateClaim(parent)
	} else {
		return parent, nil
	}
}
