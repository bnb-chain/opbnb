package chaincfg

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

var Mainnet, Goerli, Sepolia *rollup.Config

func init() {
	mustCfg := func(name string) *rollup.Config {
		cfg, err := GetRollupConfig(name)
		if err != nil {
			panic(fmt.Errorf("failed to load rollup config %q: %w", name, err))
		}
		return cfg
	}
	Mainnet = mustCfg("op-mainnet")
	Goerli = mustCfg("op-goerli")
	Sepolia = mustCfg("op-sepolia")
}

var L2ChainIDToNetworkDisplayName = func() map[string]string {
	out := make(map[string]string)
	for _, netCfg := range superchain.OPChains {
		out[fmt.Sprintf("%d", netCfg.ChainID)] = netCfg.Name
	}
	return out
}()

// AvailableNetworks returns the selection of network configurations that is available by default.
func AvailableNetworks() []string {
	var networks []string
	for _, cfg := range superchain.OPChains {
		networks = append(networks, cfg.Chain+"-"+cfg.Superchain)
	}
	return networks
}

func handleLegacyName(name string) string {
	switch name {
	case "goerli":
		return "op-goerli"
	case "mainnet":
		return "op-mainnet"
	case "sepolia":
		return "op-sepolia"
	default:
		return name
	}
}

// ChainByName returns a chain, from known available configurations, by name.
// ChainByName returns nil when the chain name is unknown.
func ChainByName(name string) *superchain.ChainConfig {
	// Handle legacy name aliases
	name = handleLegacyName(name)
	for _, chainCfg := range superchain.OPChains {
		if strings.EqualFold(chainCfg.Chain+"-"+chainCfg.Superchain, name) {
			return chainCfg
		}
	}
	return nil
}

func GetRollupConfig(name string) (*rollup.Config, error) {
	chainCfg := ChainByName(name)
	if chainCfg == nil {
		return nil, fmt.Errorf("invalid network: %q", name)
	}
	rollupCfg, err := rollup.LoadOPStackRollupConfig(chainCfg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to load rollup config: %w", err)
	}
	return rollupCfg, nil
}

var NetworksByName = map[string]rollup.Config{
	"opBNBMainnet": OPBNBMainnet,
	"opBNBTestnet": OPBNBTestnet,
	"opBNBQANet":   OPBNBQANet,
}

var NetworksByChainId = map[string]rollup.Config{
	"204":  OPBNBMainnet,
	"5611": OPBNBTestnet,
	"2484": OPBNBQANet,
}

func GetRollupConfigByNetwork(name string) (rollup.Config, error) {
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
	Fermat:                 big.NewInt(9397477), // Nov-28-2023 06 AM +UTC
	SnowTime:               u64Ptr(1713160800),  // Apr-15-2024 06 AM +UTC
	CanyonTime:             u64Ptr(1718870400),  // Jun-20-2024 08:00 AM +UTC
	DeltaTime:              u64Ptr(1718871000),  // Jun-20-2024 08:10 AM +UTC
	EcotoneTime:            u64Ptr(1718871600),  // Jun-20-2024 08:20 AM +UTC
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
	Fermat:                 big.NewInt(12113000), // Nov-03-2023 06 AM +UTC
	SnowTime:               u64Ptr(1715752800),   // May-15-2024 06:00 AM +UTC
	CanyonTime:             u64Ptr(1715753400),   // May-15-2024 06:10 AM +UTC
	DeltaTime:              u64Ptr(1715754000),   // May-15-2024 06:20 AM +UTC
	EcotoneTime:            u64Ptr(1715754600),   // May-15-2024 06:30 AM +UTC
}

var OPBNBQANet = rollup.Config{
	Genesis: rollup.Genesis{
		L1: eth.BlockID{
			Hash:   common.HexToHash("0xd8b84c6811ad3eb68ad578e12312f797d84c59a97993a1f230409c1644fcb3d2"),
			Number: 373422,
		},
		L2: eth.BlockID{
			Hash:   common.HexToHash("0xe182e685b1ec05ca55f2374cb3a190d1ae8f3e196acb55a69efd61536fc3983f"),
			Number: 0,
		},
		L2Time: 1714291718,
		SystemConfig: eth.SystemConfig{
			BatcherAddr: common.HexToAddress("0xbd6353a2e43a0d8eaa370b2eceb80481bc5c4094"),
			Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
			Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
			GasLimit:    100000000,
		},
	},
	BlockTime:              1,
	MaxSequencerDrift:      600,
	SeqWindowSize:          14400,
	ChannelTimeout:         1200,
	L1ChainID:              big.NewInt(714),
	L2ChainID:              big.NewInt(2484),
	BatchInboxAddress:      common.HexToAddress("0xff00000000000000000000000000000000001484"),
	DepositContractAddress: common.HexToAddress("0xb22e158785dbfb055edddb24ad97b4e7c51a6624"),
	L1SystemConfigAddress:  common.HexToAddress("0xbf05c7e8ac1bd5ed042618762a7442f726ecae0b"),
	RegolithTime:           u64Ptr(0),
	Fermat:                 big.NewInt(0),
	SnowTime:               u64Ptr(1714993200), // May-06-2024 11:00 AM +UTC
	CanyonTime:             u64Ptr(1714993800), // May-06-2024 11:10 AM +UTC
	DeltaTime:              u64Ptr(1714994400), // May-06-2024 11:20 AM +UTC
	EcotoneTime:            u64Ptr(1714995000), // May-06-2024 11:30 AM +UTC
}

func u64Ptr(v uint64) *uint64 {
	return &v
}
