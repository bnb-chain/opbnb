package opbnb_upgrades

import (
	"context"
	"fmt"
	oldBindings "github.com/ethereum-optimism/optimism/op-chain-ops/opbnb-upgrades/old-contracts/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// ContractVersions
type ContractVersions struct {
	L1CrossDomainMessenger       string `yaml:"l1_cross_domain_messenger"`
	L1ERC721Bridge               string `yaml:"l1_erc721_bridge"`
	L1StandardBridge             string `yaml:"l1_standard_bridge"`
	L2OutputOracle               string `yaml:"l2_output_oracle,omitempty"`
	OptimismMintableERC20Factory string `yaml:"optimism_mintable_erc20_factory"`
	OptimismPortal               string `yaml:"optimism_portal"`
	SystemConfig                 string `yaml:"system_config"`
	// Superchain-wide contracts:
	ProtocolVersions string `yaml:"protocol_versions"`
	SuperchainConfig string `yaml:"superchain_config,omitempty"`
}

// GetProxyContractVersions will fetch the versions of all of the contracts.
func GetProxyContractVersions(ctx context.Context, addresses map[string]common.Address, backend bind.ContractBackend) (ContractVersions, error) {
	var versions ContractVersions
	var err error

	versions.L1CrossDomainMessenger, err = getVersion(ctx, addresses["L1CrossDomainMessengerProxy"], backend)
	if err != nil {
		return versions, fmt.Errorf("L1CrossDomainMessenger: %w", err)
	}
	versions.L1ERC721Bridge, err = getVersion(ctx, addresses["L1ERC721BridgeProxy"], backend)
	if err != nil {
		return versions, fmt.Errorf("L1ERC721Bridge: %w", err)
	}
	versions.L1StandardBridge, err = getVersion(ctx, addresses["L1StandardBridgeProxy"], backend)
	if err != nil {
		return versions, fmt.Errorf("L1StandardBridge: %w", err)
	}
	versions.L2OutputOracle, err = getVersion(ctx, addresses["L2OutputOracleProxy"], backend)
	if err != nil {
		return versions, fmt.Errorf("L2OutputOracle: %w", err)
	}
	versions.OptimismMintableERC20Factory, err = getVersion(ctx, addresses["OptimismMintableERC20FactoryProxy"], backend)
	if err != nil {
		return versions, fmt.Errorf("OptimismMintableERC20Factory: %w", err)
	}
	versions.OptimismPortal, err = getVersion(ctx, addresses["OptimismPortalProxy"], backend)
	if err != nil {
		return versions, fmt.Errorf("OptimismPortal: %w", err)
	}
	versions.SystemConfig, err = getVersion(ctx, addresses["SystemConfigProxy"], backend)
	if err != nil {
		return versions, fmt.Errorf("SystemConfig: %w", err)
	}
	return versions, err
}

// GetImplContractVersions will fetch the versions of all of the contracts.
func GetImplContractVersions(ctx context.Context, addresses map[string]common.Address, backend bind.ContractBackend) (ContractVersions, error) {
	var versions ContractVersions
	var err error

	versions.L1CrossDomainMessenger, err = getVersion(ctx, addresses["L1CrossDomainMessenger"], backend)
	if err != nil {
		return versions, fmt.Errorf("L1CrossDomainMessenger: %w", err)
	}
	versions.L1ERC721Bridge, err = getVersion(ctx, addresses["L1ERC721Bridge"], backend)
	if err != nil {
		return versions, fmt.Errorf("L1ERC721Bridge: %w", err)
	}
	versions.L1StandardBridge, err = getVersion(ctx, addresses["L1StandardBridge"], backend)
	if err != nil {
		return versions, fmt.Errorf("L1StandardBridge: %w", err)
	}
	versions.L2OutputOracle, err = getVersion(ctx, addresses["L2OutputOracle"], backend)
	if err != nil {
		return versions, fmt.Errorf("L2OutputOracle: %w", err)
	}
	versions.OptimismMintableERC20Factory, err = getVersion(ctx, addresses["OptimismMintableERC20Factory"], backend)
	if err != nil {
		return versions, fmt.Errorf("OptimismMintableERC20Factory: %w", err)
	}
	versions.OptimismPortal, err = getVersion(ctx, addresses["OptimismPortal"], backend)
	if err != nil {
		return versions, fmt.Errorf("OptimismPortal: %w", err)
	}
	versions.SystemConfig, err = getVersion(ctx, addresses["SystemConfig"], backend)
	if err != nil {
		return versions, fmt.Errorf("SystemConfig: %w", err)
	}
	return versions, err
}

// getVersion will get the version of a contract at a given address.
func getVersion(ctx context.Context, addr common.Address, backend bind.ContractBackend) (string, error) {
	semver, err := oldBindings.NewSemver(addr, backend)
	if err != nil {
		return "", fmt.Errorf("%s: %w", addr, err)
	}
	version, err := semver.Version(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", addr, err)
	}
	return version, nil
}
