package pipeline

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/backend/depset"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/types"
)

func GenerateInteropDepset(ctx context.Context, pEnv *Env, globalIntent *state.Intent, st *state.State) error {
	lgr := pEnv.Logger.New("stage", "generate-interop-depset")

	if !globalIntent.UseInterop {
		lgr.Warn("interop not enabled - skipping interop depset generation")
		return nil
	}

	lgr.Info("creating interop dependency set...")
	deps := make(map[eth.ChainID]*depset.StaticConfigDependency)
	for i, chain := range globalIntent.Chains {
		id := eth.ChainIDFromBytes32(chain.ID)
		deps[id] = &depset.StaticConfigDependency{ChainIndex: types.ChainIndex(i)}
	}

	interopDepSet, err := depset.NewStaticConfigDependencySet(deps)
	if err != nil {
		return fmt.Errorf("failed to create interop dependency set: %w", err)
	}
	st.InteropDepSet = interopDepSet

	if err := pEnv.StateWriter.WriteState(st); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}
