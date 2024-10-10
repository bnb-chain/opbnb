package opbnb_upgrades

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	newBindings "github.com/ethereum-optimism/optimism/op-chain-ops/opbnb-upgrades/new-contracts/bindings"
	oldBindings "github.com/ethereum-optimism/optimism/op-chain-ops/opbnb-upgrades/old-contracts/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"io"
	"os"
	"strconv"
)

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

type L1CrossDomainMessengerInfo struct {
	L1CrossDomainMessengerInfo string `yaml:"l1_cross_domain_messenger_info"`
	Name                       string `yaml:"name"`
	Address                    string `yaml:"address"`
	Version                    string `yaml:"version"`
	OtherMessenger             string `yaml:"other_messenger"`
	Portal                     string `yaml:"portal"`
	SuperChainConfig           string `yaml:"super_chain_config"`
	SystemConfig               string `yaml:"system_config"`
}

type L1ERC721BridgeInfo struct {
	L1ERC721BridgeInfo string `yaml:"l1_erc721_bridge_info"`
	Name               string `yaml:"name"`
	Address            string `yaml:"address"`
	Version            string `yaml:"version"`
	OtherMessenger     string `yaml:"other_messenger"`
	OtherBridge        string `yaml:"other_bridge"`
	SuperChainConfig   string `yaml:"super_chain_config"`
}

type L1StandardBridgeInfo struct {
	L1StandardBridgeInfo string `yaml:"l1_standard_bridge_info"`
	Name                 string `yaml:"name"`
	Address              string `yaml:"address"`
	Version              string `yaml:"version"`
	OtherMessenger       string `yaml:"other_messenger"`
	OtherBridge          string `yaml:"other_bridge"`
	SuperChainConfig     string `yaml:"super_chain_config"`
	SystemConfig         string `yaml:"system_config"`
}

type L2OutputOracleInfo struct {
	L2OutputOracleInfo        string `yaml:"l2_output_oracle_info"`
	Name                      string `yaml:"name"`
	Address                   string `yaml:"address"`
	Version                   string `yaml:"version"`
	SubmissionInterval        string `yaml:"submission_interval"`
	L2BlockTime               string `yaml:"l2_block_time"`
	Challenger                string `yaml:"challenger"`
	Proposer                  string `yaml:"proposer"`
	FinalizationPeriodSeconds string `yaml:"finalization_period_seconds"`
	StartingBlockNumber       string `yaml:"starting_block_number"`
	StartingTimestamp         string `yaml:"starting_time_stamp"`
}

type OptimismMintableERC20FactoryInfo struct {
	OptimismMintableERC20FactoryInfo string `yaml:"optimism_mintable_erc20_factory_info"`
	Name                             string `yaml:"name"`
	Address                          string `yaml:"address"`
	Version                          string `yaml:"version"`
	Bridge                           string `yaml:"bridge"`
}

type OptimismPortalInfo struct {
	OptimismPortalInfo string `yaml:"optimism_portal_info"`
	Name               string `yaml:"name"`
	Address            string `yaml:"address"`
	Version            string `yaml:"version"`
	L2Oracle           string `yaml:"l2_oracle"`
	L2Sender           string `yaml:"l2_sender"`
	SystemConfig       string `yaml:"system_config"`
	Guardian           string `yaml:"guardian"`
	SuperChainConfig   string `yaml:"super_chain_config"`
	Balance            string `yaml:"balance"`
}

type SystemConfigInfo struct {
	SystemConfigInfo  string                         `yaml:"system_config_info"`
	Name              string                         `yaml:"name"`
	Address           string                         `yaml:"address"`
	Version           string                         `yaml:"version"`
	Overhead          string                         `yaml:"overhead"`
	Scalar            string                         `yaml:"scalar"`
	BatcherHash       string                         `yaml:"batcher_hash"`
	GasLimit          string                         `yaml:"gas_limit"`
	ResourceConfig    ResourceMeteringResourceConfig `yaml:"resource_config"`
	StartBlock        string                         `yaml:"start_block"`
	BaseFeeScalar     string                         `yaml:"base_fee_scalar"`
	BlobBaseFeeScalar string                         `yaml:"blob_base_fee_scalar"`
}

// ResourceMeteringResourceConfig is an auto generated low-level Go binding around an user-defined struct.
type ResourceMeteringResourceConfig struct {
	MaxResourceLimit            uint32
	ElasticityMultiplier        uint8
	BaseFeeMaxChangeDenominator uint8
	MinimumBaseFee              uint32
	SystemTxMaxGas              uint32
	MaximumBaseFee              string
}

func CheckOldContracts(proxyAddresses map[string]common.Address, backend *ethclient.Client) ([]interface{}, error) {
	l1CrossDomainMessenger, err := oldBindings.NewL1CrossDomainMessengerCaller(proxyAddresses["L1CrossDomainMessengerProxy"], backend)
	version, err := l1CrossDomainMessenger.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherMessenger, err := l1CrossDomainMessenger.OTHERMESSENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	portal, err := l1CrossDomainMessenger.PORTAL(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	oldL1CrossDomainMessenger := L1CrossDomainMessengerInfo{
		Name:           "L1CrossDomainMessenger",
		Address:        proxyAddresses["L1CrossDomainMessengerProxy"].String(),
		Version:        version,
		OtherMessenger: otherMessenger.String(),
		Portal:         portal.String(),
	}
	l1ERC721Bridge, err := oldBindings.NewL1ERC721BridgeCaller(proxyAddresses["L1ERC721BridgeProxy"], backend)
	version, err = l1ERC721Bridge.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherMessenger, err = l1ERC721Bridge.MESSENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherBridge, err := l1ERC721Bridge.OtherBridge(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	oldL1ERC721Bridge := L1ERC721BridgeInfo{
		Name:           "L1ERC721Bridge",
		Address:        proxyAddresses["L1ERC721BridgeProxy"].String(),
		Version:        version,
		OtherMessenger: otherMessenger.String(),
		OtherBridge:    otherBridge.String(),
	}
	l1StandardBridge, err := oldBindings.NewL1StandardBridgeCaller(proxyAddresses["L1StandardBridgeProxy"], backend)
	version, err = l1StandardBridge.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherMessenger, err = l1StandardBridge.MESSENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherBridge, err = l1StandardBridge.OTHERBRIDGE(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	oldL1StandardBridge := L1StandardBridgeInfo{
		Name:           "L1StandardBridge",
		Address:        proxyAddresses["L1StandardBridgeProxy"].String(),
		Version:        version,
		OtherMessenger: otherMessenger.String(),
		OtherBridge:    otherBridge.String(),
	}
	l2OutputOracle, err := oldBindings.NewL2OutputOracleCaller(proxyAddresses["L2OutputOracleProxy"], backend)
	version, err = l2OutputOracle.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	submissionInterval, err := l2OutputOracle.SUBMISSIONINTERVAL(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	l2blockTime, err := l2OutputOracle.L2BLOCKTIME(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	challenger, err := l2OutputOracle.CHALLENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	proposer, err := l2OutputOracle.PROPOSER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	finalizationPeriodSeconds, err := l2OutputOracle.FINALIZATIONPERIODSECONDS(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	startingBlockNumber, err := l2OutputOracle.StartingBlockNumber(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	startingTimestamp, err := l2OutputOracle.StartingTimestamp(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	oldL2OutputOracle := L2OutputOracleInfo{
		Name:                      "L2OutputOracle",
		Address:                   proxyAddresses["L2OutputOracleProxy"].String(),
		Version:                   version,
		SubmissionInterval:        submissionInterval.String(),
		L2BlockTime:               l2blockTime.String(),
		Challenger:                challenger.String(),
		Proposer:                  proposer.String(),
		FinalizationPeriodSeconds: finalizationPeriodSeconds.String(),
		StartingBlockNumber:       startingBlockNumber.String(),
		StartingTimestamp:         startingTimestamp.String(),
	}
	optimismMintableERC20Factory, err := oldBindings.NewOptimismMintableERC20FactoryCaller(proxyAddresses["OptimismMintableERC20FactoryProxy"], backend)
	version, err = optimismMintableERC20Factory.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	bridge, err := optimismMintableERC20Factory.BRIDGE(&bind.CallOpts{})
	oldOptimismMintableERC20Factory := OptimismMintableERC20FactoryInfo{
		Name:    "OptimismMintableERC20Factory",
		Version: version,
		Bridge:  bridge.String(),
	}
	optimismPortal, err := oldBindings.NewOptimismPortalCaller(proxyAddresses["OptimismPortalProxy"], backend)
	version, err = optimismPortal.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	l2oracle, err := optimismPortal.L2ORACLE(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	l2Sender, err := optimismPortal.L2Sender(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	systemConfigAddr, err := optimismPortal.SYSTEMCONFIG(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	guardian, err := optimismPortal.GUARDIAN(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	balance, err := backend.BalanceAt(context.Background(), proxyAddresses["OptimismPortalProxy"], nil)
	if err != nil {
		return nil, err
	}
	oldOptimismPortal := OptimismPortalInfo{
		Name:         "OptimismPortal",
		Address:      proxyAddresses["OptimismPortalProxy"].String(),
		Version:      version,
		L2Oracle:     l2oracle.String(),
		L2Sender:     l2Sender.String(),
		SystemConfig: systemConfigAddr.String(),
		Guardian:     guardian.String(),
		Balance:      balance.String(),
	}
	systemConfig, err := oldBindings.NewSystemConfigCaller(proxyAddresses["SystemConfigProxy"], backend)
	version, err = systemConfig.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	overhead, err := systemConfig.Overhead(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	scalar, err := systemConfig.Scalar(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	batcherHash, err := systemConfig.BatcherHash(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	gasLimit, err := systemConfig.GasLimit(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	resourceConfig, err := systemConfig.ResourceConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newResourceConfig := ResourceMeteringResourceConfig{
		MaxResourceLimit:            resourceConfig.MaxResourceLimit,
		ElasticityMultiplier:        resourceConfig.ElasticityMultiplier,
		BaseFeeMaxChangeDenominator: resourceConfig.BaseFeeMaxChangeDenominator,
		MinimumBaseFee:              resourceConfig.MinimumBaseFee,
		SystemTxMaxGas:              resourceConfig.SystemTxMaxGas,
		MaximumBaseFee:              resourceConfig.MaximumBaseFee.String(),
	}
	oldSystemConfig := SystemConfigInfo{
		Name:           "SystemConfig",
		Address:        proxyAddresses["SystemConfigProxy"].String(),
		Version:        version,
		Overhead:       overhead.String(),
		Scalar:         scalar.String(),
		BatcherHash:    hex.EncodeToString(batcherHash[:]),
		GasLimit:       strconv.FormatUint(gasLimit, 10),
		ResourceConfig: newResourceConfig,
	}
	data := []interface{}{oldL1CrossDomainMessenger, oldL1ERC721Bridge,
		oldL1StandardBridge, oldL2OutputOracle, oldOptimismMintableERC20Factory, oldOptimismPortal, oldSystemConfig}
	return data, nil
}

func CheckNewContracts(proxyAddresses map[string]common.Address, backend *ethclient.Client) ([]interface{}, error) {
	l1CrossDomainMessenger, err := newBindings.NewL1CrossDomainMessengerCaller(proxyAddresses["L1CrossDomainMessengerProxy"], backend)
	version, err := l1CrossDomainMessenger.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherMessenger, err := l1CrossDomainMessenger.OTHERMESSENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	portal, err := l1CrossDomainMessenger.PORTAL(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	superChainConfig, err := l1CrossDomainMessenger.SuperchainConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	systemConfigAddr, err := l1CrossDomainMessenger.SystemConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newL1CrossDomainMessenger := L1CrossDomainMessengerInfo{
		Name:             "L1CrossDomainMessenger",
		Address:          proxyAddresses["L1CrossDomainMessengerProxy"].String(),
		Version:          version,
		OtherMessenger:   otherMessenger.String(),
		Portal:           portal.String(),
		SuperChainConfig: superChainConfig.String(),
		SystemConfig:     systemConfigAddr.String(),
	}
	l1ERC721Bridge, err := newBindings.NewL1ERC721BridgeCaller(proxyAddresses["L1ERC721BridgeProxy"], backend)
	version, err = l1ERC721Bridge.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherMessenger, err = l1ERC721Bridge.MESSENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherBridge, err := l1ERC721Bridge.OtherBridge(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	superChainConfig, err = l1ERC721Bridge.SuperchainConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newL1ERC721Bridge := L1ERC721BridgeInfo{
		Name:             "L1ERC721Bridge",
		Address:          proxyAddresses["L1ERC721BridgeProxy"].String(),
		Version:          version,
		OtherMessenger:   otherMessenger.String(),
		OtherBridge:      otherBridge.String(),
		SuperChainConfig: superChainConfig.String(),
	}
	l1StandardBridge, err := newBindings.NewL1StandardBridgeCaller(proxyAddresses["L1StandardBridgeProxy"], backend)
	version, err = l1StandardBridge.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherMessenger, err = l1StandardBridge.MESSENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	otherBridge, err = l1StandardBridge.OTHERBRIDGE(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	superChainConfig, err = l1StandardBridge.SuperchainConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	systemConfigAddr, err = l1StandardBridge.SystemConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newL1StandardBridge := L1StandardBridgeInfo{
		Name:             "L1StandardBridge",
		Address:          proxyAddresses["L1StandardBridgeProxy"].String(),
		Version:          version,
		OtherMessenger:   otherMessenger.String(),
		OtherBridge:      otherBridge.String(),
		SuperChainConfig: superChainConfig.String(),
		SystemConfig:     systemConfigAddr.String(),
	}
	l2OutputOracle, err := newBindings.NewL2OutputOracleCaller(proxyAddresses["L2OutputOracleProxy"], backend)
	version, err = l2OutputOracle.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	submissionInterval, err := l2OutputOracle.SUBMISSIONINTERVAL(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	l2blockTime, err := l2OutputOracle.L2BLOCKTIME(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	challenger, err := l2OutputOracle.CHALLENGER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	proposer, err := l2OutputOracle.PROPOSER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	finalizationPeriodSeconds, err := l2OutputOracle.FINALIZATIONPERIODSECONDS(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	startingBlockNumber, err := l2OutputOracle.StartingBlockNumber(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	startingTimestamp, err := l2OutputOracle.StartingTimestamp(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newL2OutputOracle := L2OutputOracleInfo{
		Name:                      "L2OutputOracle",
		Address:                   proxyAddresses["L2OutputOracleProxy"].String(),
		Version:                   version,
		SubmissionInterval:        submissionInterval.String(),
		L2BlockTime:               l2blockTime.String(),
		Challenger:                challenger.String(),
		Proposer:                  proposer.String(),
		FinalizationPeriodSeconds: finalizationPeriodSeconds.String(),
		StartingBlockNumber:       startingBlockNumber.String(),
		StartingTimestamp:         startingTimestamp.String(),
	}
	optimismMintableERC20Factory, err := newBindings.NewOptimismMintableERC20FactoryCaller(proxyAddresses["OptimismMintableERC20FactoryProxy"], backend)
	version, err = optimismMintableERC20Factory.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	bridge, err := optimismMintableERC20Factory.BRIDGE(&bind.CallOpts{})
	newOptimismMintableERC20Factory := OptimismMintableERC20FactoryInfo{
		Name:    "OptimismMintableERC20Factory",
		Version: version,
		Bridge:  bridge.String(),
	}
	optimismPortal, err := newBindings.NewOptimismPortalCaller(proxyAddresses["OptimismPortalProxy"], backend)
	version, err = optimismPortal.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	l2oracle, err := optimismPortal.L2Oracle(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	l2Sender, err := optimismPortal.L2Sender(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	systemConfigAddr, err = optimismPortal.SystemConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	guardian, err := optimismPortal.Guardian(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	balance, err := backend.BalanceAt(context.Background(), proxyAddresses["OptimismPortalProxy"], nil)
	if err != nil {
		return nil, err
	}
	superChainConfig, err = optimismPortal.SuperchainConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newOptimismPortal := OptimismPortalInfo{
		Name:             "OptimismPortal",
		Address:          proxyAddresses["OptimismPortalProxy"].String(),
		Version:          version,
		L2Oracle:         l2oracle.String(),
		L2Sender:         l2Sender.String(),
		SystemConfig:     systemConfigAddr.String(),
		Guardian:         guardian.String(),
		Balance:          balance.String(),
		SuperChainConfig: superChainConfig.String(),
	}
	systemConfig, err := newBindings.NewSystemConfigCaller(proxyAddresses["SystemConfigProxy"], backend)
	version, err = systemConfig.Version(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	overhead, err := systemConfig.Overhead(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	scalar, err := systemConfig.Scalar(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	batcherHash, err := systemConfig.BatcherHash(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	gasLimit, err := systemConfig.GasLimit(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	resourceConfig, err := systemConfig.ResourceConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newResourceConfig := ResourceMeteringResourceConfig{
		MaxResourceLimit:            resourceConfig.MaxResourceLimit,
		ElasticityMultiplier:        resourceConfig.ElasticityMultiplier,
		BaseFeeMaxChangeDenominator: resourceConfig.BaseFeeMaxChangeDenominator,
		MinimumBaseFee:              resourceConfig.MinimumBaseFee,
		SystemTxMaxGas:              resourceConfig.SystemTxMaxGas,
		MaximumBaseFee:              resourceConfig.MaximumBaseFee.String(),
	}
	startBlock, err := systemConfig.StartBlock(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	baseFeeScalar, err := systemConfig.BasefeeScalar(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	blobBaseFeeScalar, err := systemConfig.BlobbasefeeScalar(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	newSystemConfig := SystemConfigInfo{
		Name:              "SystemConfig",
		Address:           proxyAddresses["SystemConfigProxy"].String(),
		Version:           version,
		Overhead:          overhead.String(),
		Scalar:            scalar.String(),
		BatcherHash:       hex.EncodeToString(batcherHash[:]),
		GasLimit:          strconv.FormatUint(gasLimit, 10),
		ResourceConfig:    newResourceConfig,
		StartBlock:        startBlock.String(),
		BaseFeeScalar:     strconv.FormatUint(uint64(baseFeeScalar), 10),
		BlobBaseFeeScalar: strconv.FormatUint(uint64(blobBaseFeeScalar), 10),
	}
	data := []interface{}{newL1CrossDomainMessenger, newL1ERC721Bridge,
		newL1StandardBridge, newL2OutputOracle, newOptimismMintableERC20Factory, newOptimismPortal, newSystemConfig}
	return data, nil
}

func CompareContracts(oldContractsFilePath string, newContractsFilePath string) error {
	oldContractsFile, err := os.Open(oldContractsFilePath)
	if err != nil {
		return err
	}
	defer oldContractsFile.Close()
	newContractsFile, err := os.Open(newContractsFilePath)
	if err != nil {
		return err
	}
	defer newContractsFile.Close()

	oldContractByte, err := io.ReadAll(oldContractsFile)
	if err != nil {
		return err
	}

	var oldContract []map[string]interface{}
	err = json.Unmarshal(oldContractByte, &oldContract)
	if err != nil {
		return err
	}

	newContractByte, err := io.ReadAll(newContractsFile)
	if err != nil {
		return err
	}

	var newContract []map[string]interface{}
	err = json.Unmarshal(newContractByte, &newContract)
	if err != nil {
		return err
	}

	var oldL1CrossDomainMessengerInfo L1CrossDomainMessengerInfo
	var oldL1ERC721BridgeInfo L1ERC721BridgeInfo
	var oldL1StandardBridgeInfo L1StandardBridgeInfo
	var oldL2OutputOracleInfo L2OutputOracleInfo
	var oldOptimismMintableERC20FactoryInfo OptimismMintableERC20FactoryInfo
	var oldOptimismPortalInfo OptimismPortalInfo
	var oldSystemConfigInfo SystemConfigInfo
	for _, contract := range oldContract {
		if contract["L1CrossDomainMessengerInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &oldL1CrossDomainMessengerInfo)
			if err != nil {
				return err
			}
		} else if contract["L1ERC721BridgeInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &oldL1ERC721BridgeInfo)
			if err != nil {
				return err
			}
		} else if contract["L1StandardBridgeInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &oldL1StandardBridgeInfo)
			if err != nil {
				return err
			}
		} else if contract["L2OutputOracleInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &oldL2OutputOracleInfo)
			if err != nil {
				return err
			}
		} else if contract["OptimismMintableERC20FactoryInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &oldOptimismMintableERC20FactoryInfo)
			if err != nil {
				return err
			}
		} else if contract["OptimismPortalInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &oldOptimismPortalInfo)
			if err != nil {
				return err
			}
		} else if contract["SystemConfigInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &oldSystemConfigInfo)
			if err != nil {
				return err
			}
		}
	}

	var newL1CrossDomainMessengerInfo L1CrossDomainMessengerInfo
	var newL1ERC721BridgeInfo L1ERC721BridgeInfo
	var newL1StandardBridgeInfo L1StandardBridgeInfo
	var newL2OutputOracleInfo L2OutputOracleInfo
	var newOptimismMintableERC20FactoryInfo OptimismMintableERC20FactoryInfo
	var newOptimismPortalInfo OptimismPortalInfo
	var newSystemConfigInfo SystemConfigInfo
	for _, contract := range newContract {
		if contract["L1CrossDomainMessengerInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &newL1CrossDomainMessengerInfo)
			if err != nil {
				return err
			}
		} else if contract["L1ERC721BridgeInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &newL1ERC721BridgeInfo)
			if err != nil {
				return err
			}
		} else if contract["L1StandardBridgeInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &newL1StandardBridgeInfo)
			if err != nil {
				return err
			}
		} else if contract["L2OutputOracleInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &newL2OutputOracleInfo)
			if err != nil {
				return err
			}
		} else if contract["OptimismMintableERC20FactoryInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &newOptimismMintableERC20FactoryInfo)
			if err != nil {
				return err
			}
		} else if contract["OptimismPortalInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &newOptimismPortalInfo)
			if err != nil {
				return err
			}
		} else if contract["SystemConfigInfo"] != nil {
			jsonData, _ := json.Marshal(contract)
			err := json.Unmarshal(jsonData, &newSystemConfigInfo)
			if err != nil {
				return err
			}
		}
	}

	// compare L1CrossDomainMessenger
	if newL1CrossDomainMessengerInfo.Portal != oldL1CrossDomainMessengerInfo.Portal {
		return fmt.Errorf("L1CrossDomainMessenger Portal var check diff, expect %s actual %s", oldL1CrossDomainMessengerInfo.Portal, newL1CrossDomainMessengerInfo.Portal)
	}
	if newL1CrossDomainMessengerInfo.OtherMessenger != oldL1CrossDomainMessengerInfo.OtherMessenger {
		return fmt.Errorf("L1CrossDomainMessenger OtherMessenger var check diff, expect %s actual %s", oldL1CrossDomainMessengerInfo.OtherMessenger, newL1CrossDomainMessengerInfo.OtherMessenger)
	}

	// compare L1ERC721Bridge
	if newL1ERC721BridgeInfo.OtherMessenger != oldL1ERC721BridgeInfo.OtherMessenger {
		return fmt.Errorf("L1ERC721Bridge OtherMessenger var check diff, expect %s actual %s", oldL1ERC721BridgeInfo.OtherMessenger, newL1ERC721BridgeInfo.OtherMessenger)
	}
	if newL1ERC721BridgeInfo.OtherBridge != oldL1ERC721BridgeInfo.OtherBridge {
		return fmt.Errorf("L1ERC721Bridge OtherBridge var check diff, expect %s actual %s", oldL1ERC721BridgeInfo.OtherBridge, newL1ERC721BridgeInfo.OtherBridge)
	}

	// compare L1StandardBridge
	if newL1StandardBridgeInfo.OtherMessenger != oldL1StandardBridgeInfo.OtherMessenger {
		return fmt.Errorf("L1StandardBridge OtherMessenger var check diff, expect %s actual %s", oldL1StandardBridgeInfo.OtherMessenger, newL1StandardBridgeInfo.OtherMessenger)
	}
	if newL1StandardBridgeInfo.OtherBridge != oldL1StandardBridgeInfo.OtherBridge {
		return fmt.Errorf("L1StandardBridge OtherBridge var check diff, expect %s actual %s", oldL1StandardBridgeInfo.OtherBridge, newL1StandardBridgeInfo.OtherBridge)
	}

	// compare L2OutputOracle
	if newL2OutputOracleInfo.SubmissionInterval != oldL2OutputOracleInfo.SubmissionInterval {
		return fmt.Errorf("L2OutputOracle SubmissionInterval var check diff, expect %s actual %s", oldL2OutputOracleInfo.SubmissionInterval, newL2OutputOracleInfo.SubmissionInterval)
	}
	if newL2OutputOracleInfo.L2BlockTime != oldL2OutputOracleInfo.L2BlockTime {
		return fmt.Errorf("L2OutputOracle L2BlockTime var check diff, expect %s actual %s", oldL2OutputOracleInfo.L2BlockTime, newL2OutputOracleInfo.L2BlockTime)
	}
	if newL2OutputOracleInfo.Challenger != oldL2OutputOracleInfo.Challenger {
		return fmt.Errorf("L2OutputOracle Challenger var check diff, expect %s actual %s", oldL2OutputOracleInfo.Challenger, newL2OutputOracleInfo.Challenger)
	}
	if newL2OutputOracleInfo.Proposer != oldL2OutputOracleInfo.Proposer {
		return fmt.Errorf("L2OutputOracle Proposer var check diff, expect %s actual %s", oldL2OutputOracleInfo.Proposer, newL2OutputOracleInfo.Proposer)
	}
	if newL2OutputOracleInfo.FinalizationPeriodSeconds != oldL2OutputOracleInfo.FinalizationPeriodSeconds {
		return fmt.Errorf("L2OutputOracle FinalizationPeriodSeconds var check diff, expect %s actual %s", oldL2OutputOracleInfo.FinalizationPeriodSeconds, newL2OutputOracleInfo.FinalizationPeriodSeconds)
	}
	if newL2OutputOracleInfo.StartingBlockNumber != oldL2OutputOracleInfo.StartingBlockNumber {
		return fmt.Errorf("L2OutputOracle StartingBlockNumber var check diff, expect %s actual %s", oldL2OutputOracleInfo.StartingBlockNumber, newL2OutputOracleInfo.StartingBlockNumber)
	}
	if newL2OutputOracleInfo.StartingTimestamp != oldL2OutputOracleInfo.StartingTimestamp {
		return fmt.Errorf("L2OutputOracle StartingTimestamp var check diff, expect %s actual %s", oldL2OutputOracleInfo.StartingTimestamp, newL2OutputOracleInfo.StartingTimestamp)
	}

	// compare OptimismMintableERC20Factory
	if newOptimismMintableERC20FactoryInfo.Bridge != oldOptimismMintableERC20FactoryInfo.Bridge {
		return fmt.Errorf("OptimismMintableERC20Factory Bridge var check diff, expect %s actual %s", oldOptimismMintableERC20FactoryInfo.Bridge, newOptimismMintableERC20FactoryInfo.Bridge)
	}

	// compare OptimismPortal
	if newOptimismPortalInfo.L2Oracle != oldOptimismPortalInfo.L2Oracle {
		return fmt.Errorf("OptimismPortal L2Oracle var check diff, expect %s actual %s", oldOptimismPortalInfo.L2Oracle, newOptimismPortalInfo.L2Oracle)
	}
	if newOptimismPortalInfo.L2Sender != oldOptimismPortalInfo.L2Sender {
		return fmt.Errorf("OptimismPortal L2Sender var check diff, expect %s actual %s", oldOptimismPortalInfo.L2Sender, newOptimismPortalInfo.L2Sender)
	}
	if newOptimismPortalInfo.SystemConfig != oldOptimismPortalInfo.SystemConfig {
		return fmt.Errorf("OptimismPortal L2Oracle var check diff, expect %s actual %s", oldOptimismPortalInfo.L2Oracle, newOptimismPortalInfo.L2Oracle)
	}
	if newOptimismPortalInfo.Guardian != oldOptimismPortalInfo.Guardian {
		return fmt.Errorf("OptimismPortal L2Oracle var check diff, expect %s actual %s", oldOptimismPortalInfo.L2Oracle, newOptimismPortalInfo.L2Oracle)
	}

	// compare SystemConfig
	if newSystemConfigInfo.Overhead != oldSystemConfigInfo.Overhead {
		return fmt.Errorf("SystemConfig Overhead var check diff, expect %s actual %s", oldSystemConfigInfo.Overhead, newSystemConfigInfo.Overhead)
	}
	//if newSystemConfigInfo.Scalar != oldSystemConfigInfo.Scalar {
	//	return fmt.Errorf("SystemConfig Scalar var check diff, expect %s actual %s", oldSystemConfigInfo.Scalar, newSystemConfigInfo.Scalar)
	//}
	if newSystemConfigInfo.BatcherHash != oldSystemConfigInfo.BatcherHash {
		return fmt.Errorf("SystemConfig BatcherHash var check diff, expect %s actual %s", oldSystemConfigInfo.BatcherHash, newSystemConfigInfo.BatcherHash)
	}
	if newSystemConfigInfo.GasLimit != oldSystemConfigInfo.GasLimit {
		return fmt.Errorf("SystemConfig GasLimit var check diff, expect %s actual %s", oldSystemConfigInfo.GasLimit, newSystemConfigInfo.GasLimit)
	}
	if newSystemConfigInfo.BaseFeeScalar != strconv.FormatUint(uint64(BasefeeScalar), 10) {
		return fmt.Errorf("SystemConfig BaseFeeScalar var check diff, expect %s actual %s", strconv.FormatUint(uint64(BasefeeScalar), 10), newSystemConfigInfo.BaseFeeScalar)
	}
	if newSystemConfigInfo.BlobBaseFeeScalar != strconv.FormatUint(uint64(Blobbasefeescala), 10) {
		return fmt.Errorf("SystemConfig BlobBaseFeeScalar var check diff, expect %s actual %s", strconv.FormatUint(uint64(Blobbasefeescala), 10), newSystemConfigInfo.BlobBaseFeeScalar)
	}
	if newSystemConfigInfo.ResourceConfig.MaxResourceLimit != oldSystemConfigInfo.ResourceConfig.MaxResourceLimit {
		return fmt.Errorf("SystemConfig ResourceConfig.MaxResourceLimit var check diff, expect %d actual %d", oldSystemConfigInfo.ResourceConfig.MaxResourceLimit, newSystemConfigInfo.ResourceConfig.MaxResourceLimit)
	}
	if newSystemConfigInfo.ResourceConfig.ElasticityMultiplier != oldSystemConfigInfo.ResourceConfig.ElasticityMultiplier {
		return fmt.Errorf("SystemConfig ResourceConfig.ElasticityMultiplier var check diff, expect %d actual %d", oldSystemConfigInfo.ResourceConfig.ElasticityMultiplier, newSystemConfigInfo.ResourceConfig.ElasticityMultiplier)
	}
	if newSystemConfigInfo.ResourceConfig.BaseFeeMaxChangeDenominator != oldSystemConfigInfo.ResourceConfig.BaseFeeMaxChangeDenominator {
		return fmt.Errorf("SystemConfig ResourceConfig.MaxResourceLimit var check diff, expect %d actual %d", oldSystemConfigInfo.ResourceConfig.BaseFeeMaxChangeDenominator, newSystemConfigInfo.ResourceConfig.BaseFeeMaxChangeDenominator)
	}
	if newSystemConfigInfo.ResourceConfig.MinimumBaseFee != oldSystemConfigInfo.ResourceConfig.MinimumBaseFee {
		return fmt.Errorf("SystemConfig ResourceConfig.MinimumBaseFee var check diff, expect %d actual %d", oldSystemConfigInfo.ResourceConfig.MaxResourceLimit, newSystemConfigInfo.ResourceConfig.MaxResourceLimit)
	}
	if newSystemConfigInfo.ResourceConfig.SystemTxMaxGas != oldSystemConfigInfo.ResourceConfig.SystemTxMaxGas {
		return fmt.Errorf("SystemConfig ResourceConfig.SystemTxMaxGas var check diff, expect %d actual %d", oldSystemConfigInfo.ResourceConfig.SystemTxMaxGas, newSystemConfigInfo.ResourceConfig.SystemTxMaxGas)
	}
	if newSystemConfigInfo.ResourceConfig.MaximumBaseFee != oldSystemConfigInfo.ResourceConfig.MaximumBaseFee {
		return fmt.Errorf("SystemConfig ResourceConfig.MaximumBaseFee var check diff, expect %s actual %s", oldSystemConfigInfo.ResourceConfig.MaximumBaseFee, newSystemConfigInfo.ResourceConfig.MaximumBaseFee)
	}
	return nil
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
