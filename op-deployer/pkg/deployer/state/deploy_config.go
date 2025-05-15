package state

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/standard"

	op_service "github.com/ethereum-optimism/optimism/op-service"

	"github.com/ethereum-optimism/optimism/op-service/jsonutil"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	l2GenesisBlockBaseFeePerGas = hexutil.Big(*(big.NewInt(1000000000)))
)

func CombineDeployConfig(intent *Intent, chainIntent *ChainIntent, state *State, chainState *ChainState) (genesis.DeployConfig, error) {
	upgradeSchedule := standard.DefaultHardforkScheduleForTag(intent.L1ContractsLocator.Tag)
	if intent.UseInterop {
		if upgradeSchedule.L2GenesisIsthmusTimeOffset == nil {
			return genesis.DeployConfig{}, errors.New("expecting isthmus fork to be enabled for interop deployments")
		}
		upgradeSchedule.UseInterop = true
	}

	cfg := genesis.DeployConfig{
		L1DependenciesConfig: genesis.L1DependenciesConfig{
			L1StandardBridgeProxy:       chainState.L1StandardBridgeProxy,
			L1CrossDomainMessengerProxy: chainState.L1CrossDomainMessengerProxy,
			L1ERC721BridgeProxy:         chainState.L1Erc721BridgeProxy,
			SystemConfigProxy:           chainState.SystemConfigProxy,
			OptimismPortalProxy:         chainState.OptimismPortalProxy,
			ProtocolVersionsProxy:       state.SuperchainDeployment.ProtocolVersionsProxy,
		},
		L2InitializationConfig: genesis.L2InitializationConfig{
			DevDeployConfig: genesis.DevDeployConfig{
				FundDevAccounts: intent.FundDevAccounts,
			},
			L2GenesisBlockDeployConfig: genesis.L2GenesisBlockDeployConfig{
				L2GenesisBlockGasLimit:      60_000_000,
				L2GenesisBlockBaseFeePerGas: &l2GenesisBlockBaseFeePerGas,
			},
			L2VaultsDeployConfig: genesis.L2VaultsDeployConfig{
				BaseFeeVaultWithdrawalNetwork:            "local",
				L1FeeVaultWithdrawalNetwork:              "local",
				SequencerFeeVaultWithdrawalNetwork:       "local",
				SequencerFeeVaultMinimumWithdrawalAmount: standard.VaultMinWithdrawalAmount,
				BaseFeeVaultMinimumWithdrawalAmount:      standard.VaultMinWithdrawalAmount,
				L1FeeVaultMinimumWithdrawalAmount:        standard.VaultMinWithdrawalAmount,
				BaseFeeVaultRecipient:                    chainIntent.BaseFeeVaultRecipient,
				L1FeeVaultRecipient:                      chainIntent.L1FeeVaultRecipient,
				SequencerFeeVaultRecipient:               chainIntent.SequencerFeeVaultRecipient,
			},
			GovernanceDeployConfig: genesis.GovernanceDeployConfig{
				EnableGovernance:      false,
				GovernanceTokenSymbol: "OP",
				GovernanceTokenName:   "Optimism",
				GovernanceTokenOwner:  common.HexToAddress("0xDeaDDEaDDeAdDeAdDEAdDEaddeAddEAdDEAdDEad"),
			},
			GasPriceOracleDeployConfig: genesis.GasPriceOracleDeployConfig{
				GasPriceOracleBaseFeeScalar:       1368,
				GasPriceOracleBlobBaseFeeScalar:   810949,
				GasPriceOracleOperatorFeeScalar:   chainIntent.OperatorFeeScalar,
				GasPriceOracleOperatorFeeConstant: chainIntent.OperatorFeeConstant,
			},
			EIP1559DeployConfig: genesis.EIP1559DeployConfig{
				EIP1559Denominator:       chainIntent.Eip1559Denominator,
				EIP1559DenominatorCanyon: 250,
				EIP1559Elasticity:        chainIntent.Eip1559Elasticity,
			},

			// STOP! This struct sets the _default_ upgrade schedule for all chains.
			// Any upgrades you enable here will be enabled for all new deployments.
			// In-development hardforks should never be activated here. Instead, they
			// should be specified as overrides.
			UpgradeScheduleDeployConfig: *upgradeSchedule,
			L2CoreDeployConfig: genesis.L2CoreDeployConfig{
				L1ChainID:                 intent.L1ChainID,
				L2ChainID:                 chainState.ID.Big().Uint64(),
				L2BlockTime:               2,
				FinalizationPeriodSeconds: 12,
				MaxSequencerDrift:         600,
				SequencerWindowSize:       3600,
				ChannelTimeoutBedrock:     300,
				SystemConfigStartBlock:    0,
				BatchInboxAddress:         calculateBatchInboxAddr(chainState.ID),
			},
			OperatorDeployConfig: genesis.OperatorDeployConfig{
				BatchSenderAddress:  chainIntent.Roles.Batcher,
				P2PSequencerAddress: chainIntent.Roles.UnsafeBlockSigner,
			},
			OwnershipDeployConfig: genesis.OwnershipDeployConfig{
				ProxyAdminOwner:  chainIntent.Roles.L2ProxyAdminOwner,
				FinalSystemOwner: chainIntent.Roles.L1ProxyAdminOwner,
			},
		},
		FaultProofDeployConfig: genesis.FaultProofDeployConfig{
			UseFaultProofs:                  true,
			FaultGameWithdrawalDelay:        604800,
			PreimageOracleMinProposalSize:   126000,
			PreimageOracleChallengePeriod:   86400,
			ProofMaturityDelaySeconds:       604800,
			DisputeGameFinalityDelaySeconds: 302400,
		},
	}

	if intent.UseInterop {
		cfg.L2InitializationConfig.UpgradeScheduleDeployConfig.L2GenesisInteropTimeOffset = op_service.U64UtilPtr(0)
	}

	if chainState.StartBlock == nil {
		// These are dummy variables - see below for rationale.
		num := rpc.LatestBlockNumber
		cfg.L1StartingBlockTag = &genesis.MarshalableRPCBlockNumberOrHash{
			BlockNumber: &num,
		}
	} else {
		startHash := chainState.StartBlock.Hash
		cfg.L1StartingBlockTag = &genesis.MarshalableRPCBlockNumberOrHash{
			BlockHash: &startHash,
		}
	}

	if chainIntent.DangerousAltDAConfig.UseAltDA {
		cfg.AltDADeployConfig = chainIntent.DangerousAltDAConfig
		cfg.L1DependenciesConfig.DAChallengeProxy = chainState.AltDAChallengeProxy
	}

	// The below dummy variables are set in order to allow the deploy
	// config to pass validation. The validation checks are useful to
	// ensure that the L2 is properly configured. They are not used by
	// the L2 genesis script itself.

	cfg.L1BlockTime = 12
	dummyAddr := common.Address{19: 0x01}
	cfg.SuperchainL1DeployConfig = genesis.SuperchainL1DeployConfig{
		SuperchainConfigGuardian: dummyAddr,
	}
	cfg.OutputOracleDeployConfig = genesis.OutputOracleDeployConfig{
		L2OutputOracleSubmissionInterval: 1,
		L2OutputOracleStartingTimestamp:  1,
		L2OutputOracleProposer:           dummyAddr,
		L2OutputOracleChallenger:         dummyAddr,
	}
	// End of dummy variables

	// Apply overrides after setting the main values.
	var err error
	if len(intent.GlobalDeployOverrides) > 0 {
		cfg, err = jsonutil.MergeJSON(cfg, intent.GlobalDeployOverrides)
		if err != nil {
			return genesis.DeployConfig{}, fmt.Errorf("error merging global L2 overrides: %w", err)

		}
	}

	if len(chainIntent.DeployOverrides) > 0 {
		cfg, err = jsonutil.MergeJSON(cfg, chainIntent.DeployOverrides)
		if err != nil {
			return genesis.DeployConfig{}, fmt.Errorf("error merging chain L2 overrides: %w", err)
		}
	}

	if err := cfg.Check(log.New(log.DiscardHandler())); err != nil {
		return cfg, fmt.Errorf("combined deploy config failed validation: %w", err)
	}

	return cfg, nil
}

func calculateBatchInboxAddr(chainID common.Hash) common.Address {
	var out common.Address
	copy(out[1:], crypto.Keccak256(chainID[:])[:19])
	return out
}
