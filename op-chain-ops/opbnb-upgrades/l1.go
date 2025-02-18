package opbnb_upgrades

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	newBindings "github.com/ethereum-optimism/optimism/op-chain-ops/opbnb-upgrades/new-contracts/bindings"
	oldBindings "github.com/ethereum-optimism/optimism/op-chain-ops/opbnb-upgrades/old-contracts/bindings"
	"github.com/ethereum-optimism/optimism/op-chain-ops/safe"
	"github.com/ethereum-optimism/optimism/op-chain-ops/upgrades/bindings"
)

const (
	// upgradeAndCall represents the signature of the upgradeAndCall function
	// on the ProxyAdmin contract.
	upgradeAndCall = "upgradeAndCall(address,address,bytes)"

	method = "setBytes32"
)

// L1 will add calls for upgrading each of the L1 contracts.
func L1(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend, chainId *big.Int) error {
	if err := SuperChainConfig(batch, proxyAddresses, implAddresses, backend); err != nil {
		return fmt.Errorf("upgrading SuperChainConfig: %w", err)
	}
	if err := L1CrossDomainMessenger(batch, proxyAddresses, implAddresses, backend); err != nil {
		return fmt.Errorf("upgrading L1CrossDomainMessenger: %w", err)
	}
	if err := L1ERC721Bridge(batch, proxyAddresses, implAddresses, backend); err != nil {
		return fmt.Errorf("upgrading L1ERC721Bridge: %w", err)
	}
	if err := L1StandardBridge(batch, proxyAddresses, implAddresses, backend); err != nil {
		return fmt.Errorf("upgrading L1StandardBridge: %w", err)
	}
	if err := L2OutputOracle(batch, proxyAddresses, implAddresses, backend); err != nil {
		return fmt.Errorf("upgrading L2OutputOracle: %w", err)
	}
	if err := OptimismMintableERC20Factory(batch, proxyAddresses, implAddresses, backend); err != nil {
		return fmt.Errorf("upgrading OptimismMintableERC20Factory: %w", err)
	}
	if err := OptimismPortal(batch, proxyAddresses, implAddresses, backend); err != nil {
		return fmt.Errorf("upgrading OptimismPortal: %w", err)
	}
	if err := SystemConfig(batch, proxyAddresses, implAddresses, backend, chainId); err != nil {
		return fmt.Errorf("upgrading SystemConfig: %w", err)
	}
	return nil
}

// SuperChainConfig will add a call to the batch that upgrades the SuperChainConfig.
func SuperChainConfig(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	// can remove upgrade storageSetter superChainConfig first upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L11-L13
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return err
		}
		args := []any{
			proxyAddresses["SuperChainConfigProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	optimismPortal, err := oldBindings.NewOptimismPortalCaller(proxyAddresses["OptimismPortalProxy"], backend)
	if err != nil {
		return err
	}
	guardian, err := optimismPortal.GUARDIAN(&bind.CallOpts{})
	if err != nil {
		return err
	}
	paused, err := optimismPortal.Paused(&bind.CallOpts{})
	if err != nil {
		return err
	}

	superChainConfigABI, err := newBindings.SuperChainConfigMetaData.GetAbi()
	if err != nil {
		return err
	}

	calldata, err := superChainConfigABI.Pack("initialize", guardian, paused)
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["SuperChainConfigProxy"],
		implAddresses["SuperChainConfig"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}

// L1CrossDomainMessenger will add a call to the batch that upgrades the L1CrossDomainMessenger.
func L1CrossDomainMessenger(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L11-L13
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return err
		}
		args := []any{
			proxyAddresses["L1CrossDomainMessengerProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	l1CrossDomainMessenger, err := oldBindings.NewL1CrossDomainMessengerCaller(proxyAddresses["L1CrossDomainMessengerProxy"], backend)
	if err != nil {
		return err
	}
	optimismPortal, err := l1CrossDomainMessenger.PORTAL(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l1CrossDomainMessengerABI, err := newBindings.L1CrossDomainMessengerMetaData.GetAbi()
	if err != nil {
		return err
	}

	calldata, err := l1CrossDomainMessengerABI.Pack("initialize", proxyAddresses["SuperChainConfigProxy"], optimismPortal, proxyAddresses["SystemConfigProxy"])
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["L1CrossDomainMessengerProxy"],
		implAddresses["L1CrossDomainMessenger"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}

// L1ERC721Bridge will add a call to the batch that upgrades the L1ERC721Bridge.
func L1ERC721Bridge(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L100-L102
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return fmt.Errorf("setBytes32: %w", err)
		}
		args := []any{
			proxyAddresses["L1ERC721BridgeProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	l1ERC721BridgeABI, err := bindings.L1ERC721BridgeMetaData.GetAbi()
	if err != nil {
		return err
	}

	l1ERC721Bridge, err := bindings.NewL1ERC721BridgeCaller(proxyAddresses["L1ERC721BridgeProxy"], backend)
	if err != nil {
		return err
	}
	messenger, err := l1ERC721Bridge.Messenger(&bind.CallOpts{})
	if err != nil {
		return err
	}

	calldata, err := l1ERC721BridgeABI.Pack("initialize", messenger, proxyAddresses["SuperChainConfigProxy"])
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["L1ERC721BridgeProxy"],
		implAddresses["L1ERC721Bridge"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}

// L1StandardBridge will add a call to the batch that upgrades the L1StandardBridge.
func L1StandardBridge(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L36-L37
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return err
		}
		args := []any{
			proxyAddresses["L1StandardBridgeProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	l1StandardBridgeABI, err := newBindings.L1StandardBridgeMetaData.GetAbi()
	if err != nil {
		return err
	}

	l1StandardBridge, err := oldBindings.NewL1StandardBridgeCaller(proxyAddresses["L1StandardBridgeProxy"], backend)
	if err != nil {
		return err
	}

	messenger, err := l1StandardBridge.MESSENGER(&bind.CallOpts{})
	if err != nil {
		return err
	}

	calldata, err := l1StandardBridgeABI.Pack("initialize", messenger, proxyAddresses["SuperChainConfigProxy"], proxyAddresses["SystemConfigProxy"])
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["L1StandardBridgeProxy"],
		implAddresses["L1StandardBridge"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}

// L2OutputOracle will add a call to the batch that upgrades the L2OutputOracle.
func L2OutputOracle(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L50-L51
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return err
		}
		args := []any{
			proxyAddresses["L2OutputOracleProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	l2OutputOracleABI, err := newBindings.L2OutputOracleMetaData.GetAbi()
	if err != nil {
		return err
	}

	l2OutputOracle, err := oldBindings.NewL2OutputOracleCaller(proxyAddresses["L2OutputOracleProxy"], backend)
	if err != nil {
		return err
	}

	l2OutputOracleSubmissionInterval, err := l2OutputOracle.SUBMISSIONINTERVAL(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2BlockTime, err := l2OutputOracle.L2BLOCKTIME(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2OutputOracleStartingBlockNumber, err := l2OutputOracle.StartingBlockNumber(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2OutputOracleStartingTimestamp, err := l2OutputOracle.StartingTimestamp(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2OutputOracleProposer, err := l2OutputOracle.PROPOSER(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2OutputOracleChallenger, err := l2OutputOracle.CHALLENGER(&bind.CallOpts{})
	if err != nil {
		return err
	}

	finalizationPeriodSeconds, err := l2OutputOracle.FINALIZATIONPERIODSECONDS(&bind.CallOpts{})
	if err != nil {
		return err
	}
	calldata, err := l2OutputOracleABI.Pack(
		"initialize",
		l2OutputOracleSubmissionInterval,
		l2BlockTime,
		l2OutputOracleStartingBlockNumber,
		l2OutputOracleStartingTimestamp,
		l2OutputOracleProposer,
		l2OutputOracleChallenger,
		finalizationPeriodSeconds,
	)
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["L2OutputOracleProxy"],
		implAddresses["L2OutputOracle"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}

// OptimismMintableERC20Factory will add a call to the batch that upgrades the OptimismMintableERC20Factory.
func OptimismMintableERC20Factory(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L287-L289
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return err
		}
		args := []any{
			proxyAddresses["OptimismMintableERC20FactoryProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	optimismMintableERC20FactoryABI, err := newBindings.OptimismMintableERC20FactoryMetaData.GetAbi()
	if err != nil {
		return err
	}

	optimismMintableERC20Factory, err := oldBindings.NewOptimismMintableERC20FactoryCaller(proxyAddresses["OptimismMintableERC20FactoryProxy"], backend)
	if err != nil {
		return err
	}

	bridge, err := optimismMintableERC20Factory.BRIDGE(&bind.CallOpts{})
	if err != nil {
		return err
	}

	calldata, err := optimismMintableERC20FactoryABI.Pack("initialize", bridge)
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["OptimismMintableERC20FactoryProxy"],
		implAddresses["OptimismMintableERC20Factory"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}

// OptimismPortal will add a call to the batch that upgrades the OptimismPortal.
func OptimismPortal(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L64-L65
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return err
		}
		args := []any{
			proxyAddresses["OptimismPortalProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	optimismPortalABI, err := newBindings.OptimismPortalMetaData.GetAbi()
	if err != nil {
		return err
	}

	optimismPortal, err := oldBindings.NewOptimismPortalCaller(proxyAddresses["OptimismPortalProxy"], backend)
	if err != nil {
		return err
	}
	l2OutputOracle, err := optimismPortal.L2ORACLE(&bind.CallOpts{})
	if err != nil {
		return err
	}
	systemConfig, err := optimismPortal.SYSTEMCONFIG(&bind.CallOpts{})
	if err != nil {
		return err
	}

	calldata, err := optimismPortalABI.Pack("initialize", l2OutputOracle, systemConfig, proxyAddresses["SuperChainConfigProxy"])
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["OptimismPortalProxy"],
		implAddresses["OptimismPortal"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}

// SystemConfig will add a call to the batch that upgrades the SystemConfig.
func SystemConfig(batch *safe.Batch, proxyAddresses map[string]common.Address, implAddresses map[string]common.Address, backend bind.ContractBackend, chainId *big.Int) error {
	proxyAdminABI, err := oldBindings.ProxyAdminMetaData.GetAbi()
	if err != nil {
		return err
	}

	// 2 Step Upgrade
	{
		storageSetterABI, err := newBindings.StorageSetterMetaData.GetAbi()
		if err != nil {
			return err
		}

		// set startBlock genesis block l1 origin
		startBlock := common.BigToHash(new(big.Int).SetUint64(BscQAnetStartBlock))
		if chainId.Uint64() == BscTestnet {
			startBlock = common.BigToHash(new(big.Int).SetUint64(BscTestnetStartBlock))
		} else if chainId.Uint64() == BscMainnet {
			startBlock = common.BigToHash(new(big.Int).SetUint64(BscMainnetStartBlock))
		}

		input := []bindings.StorageSetterSlot{
			// https://github.com/ethereum-optimism/optimism/blob/86a96023ffd04d119296dff095d02fff79fa15de/packages/contracts-bedrock/.storage-layout#L82-L83
			{
				Key:   common.Hash{},
				Value: common.Hash{},
			},
			// bytes32 public constant START_BLOCK_SLOT = bytes32(uint256(keccak256("systemconfig.startBlock")) - 1);
			{
				Key:   common.HexToHash("0xa11ee3ab75b40e88a0105e935d17cd36c8faee0138320d776c411291bdbbb19f"),
				Value: startBlock,
			},
		}

		calldata, err := storageSetterABI.Pack(method, input)
		if err != nil {
			return err
		}
		args := []any{
			proxyAddresses["SystemConfigProxy"],
			implAddresses["StorageSetter"],
			calldata,
		}
		if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
			return err
		}
	}

	systemConfigABI, err := newBindings.SystemConfigMetaData.GetAbi()
	if err != nil {
		return err
	}

	systemConfig, err := oldBindings.NewSystemConfigCaller(proxyAddresses["SystemConfigProxy"], backend)
	if err != nil {
		return err
	}

	batcherHash, err := systemConfig.BatcherHash(&bind.CallOpts{})
	if err != nil {
		return err
	}

	l2GenesisBlockGasLimit, err := systemConfig.GasLimit(&bind.CallOpts{})
	if err != nil {
		return err
	}

	p2pSequencerAddress, err := systemConfig.UnsafeBlockSigner(&bind.CallOpts{})
	if err != nil {
		return err
	}

	finalSystemOwner, err := systemConfig.Owner(&bind.CallOpts{})
	if err != nil {
		return err
	}

	resourceConfig, err := systemConfig.ResourceConfig(&bind.CallOpts{})
	if err != nil {
		return err
	}

	// set startBlock genesis block l1 origin
	batchInboxAddr := BscQAnetBatcherInbox
	if chainId.Uint64() == BscTestnet {
		batchInboxAddr = BscTestnetBatcherInbox
	} else if chainId.Uint64() == BscMainnet {
		batchInboxAddr = BscMainnetBatcherInbox
	}
	_basefeeScalar := BasefeeScalar
	_blobbasefeeScalar := Blobbasefeescala
	calldata, err := systemConfigABI.Pack(
		"initialize",
		finalSystemOwner,
		_basefeeScalar,
		_blobbasefeeScalar,
		batcherHash,
		l2GenesisBlockGasLimit,
		p2pSequencerAddress,
		resourceConfig,
		batchInboxAddr,
		bindings.SystemConfigAddresses{
			L1CrossDomainMessenger:       proxyAddresses["L1CrossDomainMessengerProxy"],
			L1ERC721Bridge:               proxyAddresses["L1ERC721BridgeProxy"],
			L1StandardBridge:             proxyAddresses["L1StandardBridgeProxy"],
			DisputeGameFactory:           common.Address{},
			OptimismPortal:               proxyAddresses["OptimismPortalProxy"],
			OptimismMintableERC20Factory: proxyAddresses["OptimismMintableERC20FactoryProxy"],
			GasPayingToken:               common.Address{},
		},
	)
	if err != nil {
		return err
	}

	args := []any{
		proxyAddresses["SystemConfigProxy"],
		implAddresses["SystemConfig"],
		calldata,
	}

	if err := batch.AddCall(implAddresses["ProxyAdmin"], common.Big0, upgradeAndCall, args, proxyAdminABI); err != nil {
		return err
	}

	return nil
}
