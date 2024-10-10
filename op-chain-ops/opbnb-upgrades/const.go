package opbnb_upgrades

import "github.com/ethereum/go-ethereum/common"

const (
	BscTestnet = 97
	BscMainnet = 56
)

const (
	bscTestnetStartBlock = 30727847
	bscMainnetStartBlock = 30758357
	bscQAnetStartBlock   = 44147099
)

const (
	// all networks are the same
	BasefeeScalar    uint32 = 68000
	Blobbasefeescala uint32 = 655000
)

var (
	bscTestnetBatcherInbox = common.HexToAddress("0xff00000000000000000000000000000000005611")
	bscMainnetBatcherInbox = common.HexToAddress("0xff00000000000000000000000000000000000204")
	bscQAnetBatcherInbox   = common.HexToAddress("0xff00000000000000000000000000000000008848")
)

var BscTestnetProxyContracts = map[string]common.Address{
	"SuperChainConfigProxy":             common.HexToAddress("0xb19bFAC32Aa9aADEc286E2c918FF949bc0e59218"),
	"L1CrossDomainMessengerProxy":       common.HexToAddress("0xD506952e78eeCd5d4424B1990a0c99B1568E7c2C"),
	"L1ERC721BridgeProxy":               common.HexToAddress("0x17e1454015bFb3377c75bE7b6d47B236fd2ddbE7"),
	"L1StandardBridgeProxy":             common.HexToAddress("0x677311Fd2cCc511Bbc0f581E8d9a07B033D5E840"),
	"L2OutputOracleProxy":               common.HexToAddress("0xFf2394Bb843012562f4349C6632a0EcB92fC8810"),
	"OptimismMintableERC20FactoryProxy": common.HexToAddress("0x182cE4305791744202BB4F802C155B94cb66163B"),
	"OptimismPortalProxy":               common.HexToAddress("0x4386C8ABf2009aC0c263462Da568DD9d46e52a31"),
	"SystemConfigProxy":                 common.HexToAddress("0x406aC857817708eAf4ca3A82317eF4ae3D1EA23B"),
}

var BscTestnetImplContracts = map[string]common.Address{
	"SuperChainConfig":             common.HexToAddress("0x80E480e226F38c11f33A4cf7744EE92C90224B83"),
	"L1CrossDomainMessenger":       common.HexToAddress("0xc25C60b8f38AF900F520E19aA9eCeaBAF6501906"),
	"L1ERC721Bridge":               common.HexToAddress("0x7d68a6F2B21B28Bc058837Cd1Fb67f36D717e554"),
	"L1StandardBridge":             common.HexToAddress("0xCC3E823575e77B2E47c50C86800bdCcDb73d4185"),
	"L2OutputOracle":               common.HexToAddress("0xe80A6945CD3b32eA2d05702772232A056a218D13"),
	"OptimismMintableERC20Factory": common.HexToAddress("0x21751BDb682A73036B3dB0Fa5faB1Ec629b29941"),
	"OptimismPortal":               common.HexToAddress("0x2Dc5fAB884EC606b196381feb5743C9390b899F0"),
	"SystemConfig":                 common.HexToAddress("0x218E3dCdf2E0f29E569b6934fef6Ad50bbDe785C"),
	"ProxyAdmin":                   common.HexToAddress("0xE4925bD8Ac30b2d4e2bD7b8Ba495a5c92d4c5156"),
	"StorageSetter":                common.HexToAddress("0x14C1079C7F0e436d76D5fbBe6a212da7244D8244"),
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

var BscQAnetProxyContracts = map[string]common.Address{
	"SuperChainConfigProxy":             common.HexToAddress("0x83Fa2e9cA2DF536fF3cdC80c33e33a625cE75C0f"),
	"L1CrossDomainMessengerProxy":       common.HexToAddress("0x192EB5D7C741A8Ab047Ee16065672E13b75fD778"),
	"L1ERC721BridgeProxy":               common.HexToAddress("0xc45a6EeAa0ed6D07B000A60ea8Ef720247aa79DC"),
	"L1StandardBridgeProxy":             common.HexToAddress("0x8614FBaE2c54a8149F3DEDEe89b5Cf6e848f4d9E"),
	"L2OutputOracleProxy":               common.HexToAddress("0x9906a258D0adb25de3a73d57f27Bff73Eb1078b8"),
	"OptimismMintableERC20FactoryProxy": common.HexToAddress("0xe5e4733d55305D7266148781C7534AC02C6F2861"),
	"OptimismPortalProxy":               common.HexToAddress("0xDe279cb3237b7b322449E5cbc141BaE0EB450137"),
	"SystemConfigProxy":                 common.HexToAddress("0x38e24297458C0B6Aa4a44497086cbE7839dafb70"),
}

var BscQAnetImplContracts = map[string]common.Address{
	"SuperChainConfig":             common.HexToAddress("0x96eE40cEe4Db1A59AF5Fb62AA8110eE3909B3A24"),
	"L1CrossDomainMessenger":       common.HexToAddress("0xb811d5c3169df318cd7aaE4f518AF2B0591Be82b"),
	"L1ERC721Bridge":               common.HexToAddress("0x9b46a0729Db9D7db74a0a6447E56cd8eFBdbb7b2"),
	"L1StandardBridge":             common.HexToAddress("0xAd65f684432aA24166145b55492630797306888d"),
	"L2OutputOracle":               common.HexToAddress("0x96B8e03E9BE5eB60507e16b59B1A1Ea3c8199ba9"),
	"OptimismMintableERC20Factory": common.HexToAddress("0xcda36E872a271416cEaA62db667ef4334Db2C46c"),
	"OptimismPortal":               common.HexToAddress("0x9Afdf32859305A203d201e450b1860e92aeC15F3"),
	"SystemConfig":                 common.HexToAddress("0x448beed7e0a1D7ec5671A96Cb5703a4EE1282908"),
	"ProxyAdmin":                   common.HexToAddress("0x28Fa925aD1CC9C8F0c4daF601BfaB0b36e576a71"),
	"StorageSetter":                common.HexToAddress("0x954E3c51483BAb314e23F740bd9Ff9A4cDEA74dd"),
}
