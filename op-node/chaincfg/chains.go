package chaincfg

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
)

var Mainnet = rollup.Config{
	Genesis: rollup.Genesis{
		L1: eth.BlockID{
			Hash:   common.HexToHash("0x438335a20d98863a4c0c97999eb2481921ccd28553eac6f913af7c12aec04108"),
			Number: 17422590,
		},
		L2: eth.BlockID{
			Hash:   common.HexToHash("0xdbf6a80fef073de06add9b0d14026d6e5a86c85f6d102c36d3d8e9cf89c2afd3"),
			Number: 105235063,
		},
		L2Time: 1686068903,
		SystemConfig: eth.SystemConfig{
			BatcherAddr: common.HexToAddress("0x6887246668a3b87f54deb3b94ba47a6f63f32985"),
			Overhead:    eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc")),
			Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0")),
			GasLimit:    30_000_000,
		},
	},
	BlockTime:              2,
	MaxSequencerDrift:      600,
	SeqWindowSize:          3600,
	ChannelTimeout:         300,
	L1ChainID:              big.NewInt(1),
	L2ChainID:              big.NewInt(10),
	BatchInboxAddress:      common.HexToAddress("0xff00000000000000000000000000000000000010"),
	DepositContractAddress: common.HexToAddress("0xbEb5Fc579115071764c7423A4f12eDde41f106Ed"),
	L1SystemConfigAddress:  common.HexToAddress("0x229047fed2591dbec1eF1118d64F7aF3dB9EB290"),
	RegolithTime:           u64Ptr(0),
}

var Goerli = rollup.Config{
	Genesis: rollup.Genesis{
		L1: eth.BlockID{
			Hash:   common.HexToHash("0x6ffc1bf3754c01f6bb9fe057c1578b87a8571ce2e9be5ca14bace6eccfd336c7"),
			Number: 8300214,
		},
		L2: eth.BlockID{
			Hash:   common.HexToHash("0x0f783549ea4313b784eadd9b8e8a69913b368b7366363ea814d7707ac505175f"),
			Number: 4061224,
		},
		L2Time: 1673550516,
		SystemConfig: eth.SystemConfig{
			BatcherAddr: common.HexToAddress("0x7431310e026B69BFC676C0013E12A1A11411EEc9"),
			Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
			Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
			GasLimit:    25_000_000,
		},
	},
	BlockTime:              2,
	MaxSequencerDrift:      600,
	SeqWindowSize:          3600,
	ChannelTimeout:         300,
	L1ChainID:              big.NewInt(5),
	L2ChainID:              big.NewInt(420),
	BatchInboxAddress:      common.HexToAddress("0xff00000000000000000000000000000000000420"),
	DepositContractAddress: common.HexToAddress("0x5b47E1A08Ea6d985D6649300584e6722Ec4B1383"),
	L1SystemConfigAddress:  common.HexToAddress("0xAe851f927Ee40dE99aaBb7461C00f9622ab91d60"),
	RegolithTime:           u64Ptr(1679079600),
}

var OPBNBMainnet = rollup.Config{
	Genesis: rollup.Genesis{
		L1: eth.BlockID{
			Hash:   common.HexToHash("0x29443b21507894febe7700f7c5cd3569cc8bf1ba535df0489276d8004af81044"),
			Number: 30758357,
		},
		L2: eth.BlockID{
			Hash:   common.HexToHash("0x4dd61178c8b0f01670c231597e7bcb368e84545acd46d940a896d6a791dd6df4"),
			Number: 0,
		},
		L2Time: 1691753723,
		SystemConfig: eth.SystemConfig{
			BatcherAddr: common.HexToAddress("0xef8783382ef80ec23b66c43575a6103deca909c3"),
			Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
			Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
			GasLimit:    100000000,
		},
	},
	BlockTime:              1,
	MaxSequencerDrift:      600,
	SeqWindowSize:          14400,
	ChannelTimeout:         1200,
	L1ChainID:              big.NewInt(56),
	L2ChainID:              big.NewInt(204),
	BatchInboxAddress:      common.HexToAddress("0xff00000000000000000000000000000000000204"),
	DepositContractAddress: common.HexToAddress("0x1876ea7702c0ad0c6a2ae6036de7733edfbca519"),
	L1SystemConfigAddress:  common.HexToAddress("0x7ac836148c14c74086d57f7828f2d065672db3b8"),
	RegolithTime:           u64Ptr(0),
	// TODO update block number
	Fermat: nil,
}

var OPBNBTestnet = rollup.Config{
	Genesis: rollup.Genesis{
		L1: eth.BlockID{
			Hash:   common.HexToHash("0xc01a09840419cd993cf4666309f36e6d38de39771af8dbffecfa0386321c19f7"),
			Number: 30727847,
		},
		L2: eth.BlockID{
			Hash:   common.HexToHash("0x51fa57729dfb1c27542c21b06cb72a0459c57440ceb43a465dae1307cd04fe80"),
			Number: 0,
		},
		L2Time: 1686878506,
		SystemConfig: eth.SystemConfig{
			BatcherAddr: common.HexToAddress("0x1fd6a75cc72f39147756a663f3ef1fc95ef89495"),
			Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
			Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
			GasLimit:    100000000,
		},
	},
	BlockTime:              1,
	MaxSequencerDrift:      600,
	SeqWindowSize:          14400,
	ChannelTimeout:         1200,
	L1ChainID:              big.NewInt(97),
	L2ChainID:              big.NewInt(5611),
	BatchInboxAddress:      common.HexToAddress("0xff00000000000000000000000000000000005611"),
	DepositContractAddress: common.HexToAddress("0x4386c8abf2009ac0c263462da568dd9d46e52a31"),
	L1SystemConfigAddress:  common.HexToAddress("0x406ac857817708eaf4ca3a82317ef4ae3d1ea23b"),
	RegolithTime:           u64Ptr(0),
	// TODO update block number
	Fermat: nil,
}

var OPBNBDevnet = rollup.Config{
	Genesis: rollup.Genesis{
		L1: eth.BlockID{
			Hash:   common.HexToHash("0x29aee50ab3edefa64219e5c9b9c07f7d1953a98f2f4003d2c6fd93abeee4b706"),
			Number: 2890195,
		},
		L2: eth.BlockID{
			Hash:   common.HexToHash("0x49d448b8dc98cc95e3968615ff3dbd904d9eec8252c5f52271f029896e6147ee"),
			Number: 0,
		},
		L2Time: 1694166483,
		SystemConfig: eth.SystemConfig{
			BatcherAddr: common.HexToAddress("0x425a3598cb5e2d37213936e187914ea2059957ba"),
			Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
			Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
			GasLimit:    100000000,
		},
	},
	BlockTime:              1,
	MaxSequencerDrift:      600,
	SeqWindowSize:          14400,
	ChannelTimeout:         1200,
	L1ChainID:              big.NewInt(797),
	L2ChainID:              big.NewInt(1320),
	BatchInboxAddress:      common.HexToAddress("0xff00000000000000000000000000000000000204"),
	DepositContractAddress: common.HexToAddress("0xd93160096c5b65bb036b3269eb02328ddadb9856"),
	L1SystemConfigAddress:  common.HexToAddress("0xf053067cec8d8990de2ba9e17ec2f16c63c7bec4"),
	RegolithTime:           u64Ptr(0),
	Fermat:                 big.NewInt(3615117),
}

var NetworksByName = map[string]rollup.Config{
	"goerli":       Goerli,
	"mainnet":      Mainnet,
	"opBNBMainnet": OPBNBMainnet,
	"opBNBTestnet": OPBNBTestnet,
	"opBNBDevnet":  OPBNBDevnet,
}

var NetworksByChainId = map[string]rollup.Config{
	"420":  Goerli,
	"10":   Mainnet,
	"204":  OPBNBMainnet,
	"5611": OPBNBTestnet,
	"1320": OPBNBDevnet,
}

var L2ChainIDToNetworkName = func() map[string]string {
	out := make(map[string]string)
	for name, netCfg := range NetworksByName {
		out[netCfg.L2ChainID.String()] = name
	}
	return out
}()

func AvailableNetworks() []string {
	var networks []string
	for name := range NetworksByName {
		networks = append(networks, name)
	}
	return networks
}

func GetRollupConfig(name string) (rollup.Config, error) {
	network, ok := NetworksByName[name]
	if !ok {
		return rollup.Config{}, fmt.Errorf("invalid network %s", name)
	}

	return network, nil
}

func GetRollupConfigByChainId(chainId string) (rollup.Config, error) {
	network, ok := NetworksByChainId[chainId]
	if !ok {
		return rollup.Config{}, fmt.Errorf("no match pre-setting network chainId %s, use file config", chainId)
	}

	return network, nil
}

func u64Ptr(v uint64) *uint64 {
	return &v
}
