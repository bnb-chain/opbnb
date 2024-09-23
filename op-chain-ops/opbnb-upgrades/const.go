package opbnb_upgrades

import "github.com/ethereum/go-ethereum/common"

const (
	BscTestnet = 97
	BscMainnet = 56
)

const (
	bscTestnetStartBlock = 30727847
	bscMainnetStartBlock = 30758357
	// TODO update for qa test
	bscQAnetStartBlock = 0
)

var (
	bscTestnetBatcherInbox = common.HexToAddress("0xff00000000000000000000000000000000005611")
	bscMainnetBatcherInbox = common.HexToAddress("0xff00000000000000000000000000000000000204")
	// TODO update for qa test
	bscQAnetBatcherInbox = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

var BscTestnetProxyContracts = map[string]common.Address{
	"SuperChainConfigProxy":             common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"L1CrossDomainMessengerProxy":       common.HexToAddress("0xD506952e78eeCd5d4424B1990a0c99B1568E7c2C"),
	"L1ERC721BridgeProxy":               common.HexToAddress("0x17e1454015bFb3377c75bE7b6d47B236fd2ddbE7"),
	"L1StandardBridgeProxy":             common.HexToAddress("0x677311Fd2cCc511Bbc0f581E8d9a07B033D5E840"),
	"L2OutputOracleProxy":               common.HexToAddress("0xFf2394Bb843012562f4349C6632a0EcB92fC8810"),
	"OptimismMintableERC20FactoryProxy": common.HexToAddress("0x182cE4305791744202BB4F802C155B94cb66163B"),
	"OptimismPortalProxy":               common.HexToAddress("0x4386C8ABf2009aC0c263462Da568DD9d46e52a31"),
	"SystemConfigProxy":                 common.HexToAddress("0x406aC857817708eAf4ca3A82317eF4ae3D1EA23B"),
}

var BscMainnetProxyContracts = map[string]common.Address{
	"SuperChainConfigProxy":             common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"L1CrossDomainMessengerProxy":       common.HexToAddress("0xd95D508f13f7029CCF0fb61984d5dfD11b879c4f"),
	"L1ERC721BridgeProxy":               common.HexToAddress("0xC7c796D3B712ad223Bc29Bf85E6cdD3045D998C4"),
	"L1StandardBridgeProxy":             common.HexToAddress("0xF05F0e4362859c3331Cb9395CBC201E3Fa6757Ea"),
	"L2OutputOracleProxy":               common.HexToAddress("0x153CAB79f4767E2ff862C94aa49573294B13D169"),
	"OptimismMintableERC20FactoryProxy": common.HexToAddress("0xAa53ddCDC64A53F65A5f570cc13eB13529d780f1"),
	"OptimismPortalProxy":               common.HexToAddress("0x1876EA7702C0ad0C6A2ae6036DE7733edfBca519"),
	"SystemConfigProxy":                 common.HexToAddress("0x7AC836148C14c74086D57F7828F2D065672Db3B8"),
}

var BscQAnetProxyContracts = map[string]common.Address{
	"SuperChainConfigProxy":             common.HexToAddress("0xcd5eA393ED6b7636837B7966c43084e59B4979A0"),
	"L1CrossDomainMessengerProxy":       common.HexToAddress("0xa606600E682e11233eC2aba2a05C0f317A963b58"),
	"L1ERC721BridgeProxy":               common.HexToAddress("0xE5367760eBC15d9B233600Ee267088b8311801ef"),
	"L1StandardBridgeProxy":             common.HexToAddress("0x3B64E6f473fc9c0B0e2095276d66160e4A52bd2A"),
	"L2OutputOracleProxy":               common.HexToAddress("0xab229f64359CB33A0c9f0FF82A8b3Ba4f816E1D5"),
	"OptimismMintableERC20FactoryProxy": common.HexToAddress("0xCDE19a8d97d3450e813fa0ad0286Ee53D5B1d2EF"),
	"OptimismPortalProxy":               common.HexToAddress("0xB0F2324bA4BF690dA650aECB088DC20826254a79"),
	"SystemConfigProxy":                 common.HexToAddress("0x58187c488e9708183e08da9d7cBA3A4B45F9Ea76"),
}

var BscTestnetImplContracts = map[string]common.Address{
	"SuperChainConfig":             common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"L1CrossDomainMessenger":       common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"L1ERC721Bridge":               common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"L1StandardBridge":             common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"L2OutputOracle":               common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"OptimismMintableERC20Factory": common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"OptimismPortal":               common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"SystemConfig":                 common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"ProxyAdmin":                   common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"StorageSetter":                common.HexToAddress("0x0000000000000000000000000000000000000000"),
}

var BscMainnetImplContracts = map[string]common.Address{
	"SuperChainConfig":             common.HexToAddress(""),
	"L1CrossDomainMessenger":       common.HexToAddress(""),
	"L1ERC721Bridge":               common.HexToAddress(""),
	"L1StandardBridge":             common.HexToAddress(""),
	"L2OutputOracle":               common.HexToAddress(""),
	"OptimismMintableERC20Factory": common.HexToAddress(""),
	"OptimismPortal":               common.HexToAddress(""),
	"SystemConfig":                 common.HexToAddress(""),
	"ProxyAdmin":                   common.HexToAddress(""),
	"StorageSetter":                common.HexToAddress(""),
}

var BscQAnetImplContracts = map[string]common.Address{
	"SuperChainConfig":             common.HexToAddress("0xc671C6c27A0138058E012317EBb1913e04eb5253"),
	"L1CrossDomainMessenger":       common.HexToAddress("0xEeD7BE74B1BFe8af300A303a75b9dCF71c75ada9"),
	"L1ERC721Bridge":               common.HexToAddress("0x2433c43608708A5b8Ca4D0f501F9D9166DCBb859"),
	"L1StandardBridge":             common.HexToAddress("0xC8b4bffaEFF2035AC9ad8DFc24D05f3148eF22A4"),
	"L2OutputOracle":               common.HexToAddress("0x539471d9F5AEB63fD144DEb33de81b6Ca90f219C"),
	"OptimismMintableERC20Factory": common.HexToAddress("0xD32CC03a5C7e2A7118aB311ABf0cc68aEB552833"),
	"OptimismPortal":               common.HexToAddress("0x84687764cfCc61D233F6e0247d3Ae36B3A537da6"),
	"SystemConfig":                 common.HexToAddress("0xCD747f007eF3B8E70642137B7DB9aeaa5089E1A0"),
	"ProxyAdmin":                   common.HexToAddress("0xea9d156a3F51dC608cb46934C290fD254fac2DeF"),
	"StorageSetter":                common.HexToAddress("0xD0554D38Fff88762Ee2425B180Cdf8efF996e357"),
}
