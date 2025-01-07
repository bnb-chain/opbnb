// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IDisputeGameFactoryGameSearchResult is an auto generated low-level Go binding around an user-defined struct.
type IDisputeGameFactoryGameSearchResult struct {
	Index     *big.Int
	Metadata  [32]byte
	Timestamp uint64
	RootClaim [32]byte
	ExtraData []byte
}

// DisputeGameFactoryMetaData contains all meta data concerning the DisputeGameFactory contract.
var DisputeGameFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"create\",\"inputs\":[{\"name\":\"_gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"_rootClaim\",\"type\":\"bytes32\",\"internalType\":\"Claim\"},{\"name\":\"_extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"proxy_\",\"type\":\"address\",\"internalType\":\"contractIDisputeGame\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"createZkFaultDisputeGame\",\"inputs\":[{\"name\":\"_gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"_claims\",\"type\":\"bytes32[]\",\"internalType\":\"Claim[]\"},{\"name\":\"_parentGameIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_l2BlockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"proxy_\",\"type\":\"address\",\"internalType\":\"contractIDisputeGame\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"findLatestGames\",\"inputs\":[{\"name\":\"_gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"_start\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_n\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"games_\",\"type\":\"tuple[]\",\"internalType\":\"structIDisputeGameFactory.GameSearchResult[]\",\"components\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadata\",\"type\":\"bytes32\",\"internalType\":\"GameId\"},{\"name\":\"timestamp\",\"type\":\"uint64\",\"internalType\":\"Timestamp\"},{\"name\":\"rootClaim\",\"type\":\"bytes32\",\"internalType\":\"Claim\"},{\"name\":\"extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gameAtIndex\",\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"gameType_\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"timestamp_\",\"type\":\"uint64\",\"internalType\":\"Timestamp\"},{\"name\":\"proxy_\",\"type\":\"address\",\"internalType\":\"contractIDisputeGame\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gameCount\",\"inputs\":[],\"outputs\":[{\"name\":\"gameCount_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gameImpls\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"GameType\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDisputeGame\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"games\",\"inputs\":[{\"name\":\"_gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"_rootClaim\",\"type\":\"bytes32\",\"internalType\":\"Claim\"},{\"name\":\"_extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"proxy_\",\"type\":\"address\",\"internalType\":\"contractIDisputeGame\"},{\"name\":\"timestamp_\",\"type\":\"uint64\",\"internalType\":\"Timestamp\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameUUID\",\"inputs\":[{\"name\":\"_gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"_rootClaim\",\"type\":\"bytes32\",\"internalType\":\"Claim\"},{\"name\":\"_extraData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"uuid_\",\"type\":\"bytes32\",\"internalType\":\"Hash\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initBonds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"GameType\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setImplementation\",\"inputs\":[{\"name\":\"_gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"_impl\",\"type\":\"address\",\"internalType\":\"contractIDisputeGame\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setInitBond\",\"inputs\":[{\"name\":\"_gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"_initBond\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DisputeGameCreated\",\"inputs\":[{\"name\":\"disputeProxy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gameType\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"GameType\"},{\"name\":\"rootClaim\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"Claim\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ImplementationSet\",\"inputs\":[{\"name\":\"impl\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gameType\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"GameType\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InitBondUpdated\",\"inputs\":[{\"name\":\"gameType\",\"type\":\"uint32\",\"indexed\":true,\"internalType\":\"GameType\"},{\"name\":\"newBond\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"GameAlreadyExists\",\"inputs\":[{\"name\":\"uuid\",\"type\":\"bytes32\",\"internalType\":\"Hash\"}]},{\"type\":\"error\",\"name\":\"IncorrectBondAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidParentGameIndex\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidParentGameStatus\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidParentGameType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoClaims\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoImplementation\",\"inputs\":[{\"name\":\"gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"}]}]",
	Bin: "0x60806040523480156200001157600080fd5b506200001e600062000024565b62000292565b600054610100900460ff1615808015620000455750600054600160ff909116105b8062000075575062000062306200016260201b6200101d1760201c565b15801562000075575060005460ff166001145b620000de5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805460ff19166001179055801562000102576000805461ff0019166101001790555b6200010c62000171565b6200011782620001d9565b80156200015e576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b6001600160a01b03163b151590565b600054610100900460ff16620001cd5760405162461bcd60e51b815260206004820152602b6024820152600080516020620021ba83398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000d5565b620001d76200022b565b565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff16620002875760405162461bcd60e51b815260206004820152602b6024820152600080516020620021ba83398151915260448201526a6e697469616c697a696e6760a81b6064820152608401620000d5565b620001d733620001d9565b611f1880620002a26000396000f3fe6080604052600436106100f35760003560e01c8063715018a61161008a57806396cd97201161005957806396cd972014610331578063bb8aa1fc14610351578063c4d66de8146103b2578063f2fde38b146103d257600080fd5b8063715018a6146102cb57806377f63c84146102e057806382ecf2f6146102f35780638da5cb5b1461030657600080fd5b80634d1975b4116100c65780634d1975b4146101d457806354fd4d50146101f35780635f0150cb146102495780636593dc6e1461029e57600080fd5b806314f6b1a3146100f85780631b685b9e1461011a5780631e33424014610187578063254bd683146101a7575b600080fd5b34801561010457600080fd5b50610118610113366004611703565b6103f2565b005b34801561012657600080fd5b5061015d61013536600461173a565b60656020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b34801561019357600080fd5b506101186101a2366004611755565b61047c565b3480156101b357600080fd5b506101c76101c236600461177f565b6104c8565b60405161017e919061182c565b3480156101e057600080fd5b506068545b60405190815260200161017e565b3480156101ff57600080fd5b5061023c6040518060400160405280600581526020017f312e302e3000000000000000000000000000000000000000000000000000000081525081565b60405161017e91906118e9565b34801561025557600080fd5b50610269610264366004611945565b61070c565b6040805173ffffffffffffffffffffffffffffffffffffffff909316835267ffffffffffffffff90911660208301520161017e565b3480156102aa57600080fd5b506101e56102b936600461173a565b60666020526000908152604090205481565b3480156102d757600080fd5b50610118610794565b61015d6102ee3660046119b7565b6107a8565b61015d610301366004611945565b610a61565b34801561031257600080fd5b5060335473ffffffffffffffffffffffffffffffffffffffff1661015d565b34801561033d57600080fd5b506101e561034c366004611b4c565b610d2b565b34801561035d57600080fd5b5061037161036c366004611be6565b610d61565b6040805163ffffffff909416845267ffffffffffffffff909216602084015273ffffffffffffffffffffffffffffffffffffffff169082015260600161017e565b3480156103be57600080fd5b506101186103cd366004611bff565b610dca565b3480156103de57600080fd5b506101186103ed366004611bff565b610f66565b6103fa611039565b63ffffffff821660008181526065602052604080822080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8616908117909155905190917fff513d80e2c7fa487608f70a618dfbc0cf415699dc69588c747e8c71566c88de91a35050565b610484611039565b63ffffffff8216600081815260666020526040808220849055518392917f74d6665c4b26d5596a5aa13d3014e0c06af4d322075a797f87b03cd4c5bc91ca91a35050565b606854606090831015806104da575081155b610705575060408051600583901b8101602001909152825b8381116107035760006068828154811061050e5761050e611c1c565b600091825260209091200154905060e081901c60a082901c67ffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff831663ffffffff891683036106d4576001865101865260008173ffffffffffffffffffffffffffffffffffffffff1663609d33346040518163ffffffff1660e01b8152600401600060405180830381865afa1580156105a8573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526105ee9190810190611c4b565b905060008273ffffffffffffffffffffffffffffffffffffffff1663bcef3b556040518163ffffffff1660e01b8152600401602060405180830381865afa15801561063d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106619190611cb9565b90506040518060a001604052808881526020018781526020018567ffffffffffffffff168152602001828152602001838152508860018a516106a39190611cd2565b815181106106b3576106b3611c1c565b6020026020010181905250888851106106d157505050505050610703565b50505b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90920191506104f29050565b505b9392505050565b6000806000610752878787878080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610d2b92505050565b60009081526067602052604090205473ffffffffffffffffffffffffffffffffffffffff81169860a09190911c67ffffffffffffffff16975095505050505050565b61079c611039565b6107a660006110ba565b565b63ffffffff871660009081526065602052604081205473ffffffffffffffffffffffffffffffffffffffff16816107e28a8a8a8a86611131565b9050600089896107f3600182611cd2565b81811061080257610802611c1c565b90506020020135905060006108818b8b8080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020601f8d018190048102820181019092528b81528893508d9250908c908c90819084018382808284376000920191909152506113d492505050565b90506108f33383610893600143611cd2565b40846040516020016108a89493929190611d10565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905273ffffffffffffffffffffffffffffffffffffffff861690611435565b94508473ffffffffffffffffffffffffffffffffffffffff16638129fc1c346040518263ffffffff1660e01b81526004016000604051808303818588803b15801561093d57600080fd5b505af1158015610951573d6000803e3d6000fd5b505050505060006109638d8484610d2b565b600081815260676020526040902054909150156109b4576040517f014f6fe5000000000000000000000000000000000000000000000000000000008152600481018290526024015b60405180910390fd5b60004260a01b60e08f901b17871790508060676000848152602001908152602001600020819055506068819080600181540180825580915050600190039060005260206000200160009091909190915055838e63ffffffff168873ffffffffffffffffffffffffffffffffffffffff167f5b565efe82411da98814f356d0e7bcb8f0219b8d970307c5afb4a6903a8b2e3560405160405180910390a4505050505050979650505050505050565b63ffffffff841660009081526065602052604081205473ffffffffffffffffffffffffffffffffffffffff1680610acc576040517f031c6de400000000000000000000000000000000000000000000000000000000815263ffffffff871660048201526024016109ab565b63ffffffff86166000908152606660205260409020543414610b1a576040517f8620aa1900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610b27600143611cd2565b409050610b913387838888604051602001610b46959493929190611d69565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291905273ffffffffffffffffffffffffffffffffffffffff841690611435565b92508273ffffffffffffffffffffffffffffffffffffffff16638129fc1c346040518263ffffffff1660e01b81526004016000604051808303818588803b158015610bdb57600080fd5b505af1158015610bef573d6000803e3d6000fd5b50505050506000610c37888888888080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610d2b92505050565b60008181526067602052604090205490915015610c83576040517f014f6fe5000000000000000000000000000000000000000000000000000000008152600481018290526024016109ab565b60004260a01b60e08a901b178517600083815260676020526040808220839055606880546001810182559083527fa2153420d844928b4421650203c77babc8b33d7f2e7b450e2966db0c220977530183905551919250899163ffffffff8c169173ffffffffffffffffffffffffffffffffffffffff8916917f5b565efe82411da98814f356d0e7bcb8f0219b8d970307c5afb4a6903a8b2e359190a450505050949350505050565b6000838383604051602001610d4293929190611db6565b6040516020818303038152906040528051906020012090509392505050565b600080600080600080610dba60688881548110610d8057610d80611c1c565b906000526020600020015460e081901c9160a082901c67ffffffffffffffff169173ffffffffffffffffffffffffffffffffffffffff1690565b9199909850909650945050505050565b600054610100900460ff1615808015610dea5750600054600160ff909116105b80610e045750303b158015610e04575060005460ff166001145b610e90576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016109ab565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790558015610eee57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b610ef6611443565b610eff826110ba565b8015610f6257600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b610f6e611039565b73ffffffffffffffffffffffffffffffffffffffff8116611011576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016109ab565b61101a816110ba565b50565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b60335473ffffffffffffffffffffffffffffffffffffffff1633146107a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016109ab565b6033805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600073ffffffffffffffffffffffffffffffffffffffff8216611188576040517f031c6de400000000000000000000000000000000000000000000000000000000815263ffffffff871660048201526024016109ab565b63ffffffff861660009081526066602052604090205434146111d6576040517f8620aa1900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000849003611211576040517fe120557300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff83811614801590611237575060685467ffffffffffffffff841610155b1561126e576040517fd56dc5b400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff838116146113c7576000806112a260688667ffffffffffffffff1681548110610d8057610d80611c1c565b92505091508092506112b78863ffffffff1690565b63ffffffff166112ca8363ffffffff1690565b63ffffffff1614611307576040517f7af5628400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018373ffffffffffffffffffffffffffffffffffffffff1663200d2ed26040518163ffffffff1660e01b8152600401602060405180830381865afa158015611354573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113789190611e0a565b600281111561138957611389611ddb565b036113c0576040517f4a99819c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50506113cb565b5060005b95945050505050565b60606000856040516020016113e99190611e2b565b60405160208183030381529060405280519060200120905080865186868660405160200161141b959493929190611e61565b604051602081830303815290604052915050949350505050565b6000610705600084846114e2565b600054610100900460ff166114da576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e6700000000000000000000000000000000000000000060648201526084016109ab565b6107a6611628565b600060608203516040830351602084035184518060208701018051600283016c5af43d3d93803e606057fd5bf3895289600d8a035278593da1005b363d3d373d3d3d3d610000806062363936013d738160481b1760218a03527f9e4ac34f21c619cefc926c8bd93b54bf5a39c7ab2127a895af1cc0691d7e3dff603a8a035272fd6100003d81600a3d39f336602c57343d527f6062820160781b1761ff9e82106059018a03528060f01b8352606c8101604c8a038cf0975050866115ae5763301164256000526004601cfd5b905285527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08501527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08401527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa09092019190915292915050565b600054610100900460ff166116bf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e6700000000000000000000000000000000000000000060648201526084016109ab565b6107a6336110ba565b803563ffffffff811681146116dc57600080fd5b919050565b73ffffffffffffffffffffffffffffffffffffffff8116811461101a57600080fd5b6000806040838503121561171657600080fd5b61171f836116c8565b9150602083013561172f816116e1565b809150509250929050565b60006020828403121561174c57600080fd5b610705826116c8565b6000806040838503121561176857600080fd5b611771836116c8565b946020939093013593505050565b60008060006060848603121561179457600080fd5b61179d846116c8565b95602085013595506040909401359392505050565b60005b838110156117cd5781810151838201526020016117b5565b838111156117dc576000848401525b50505050565b600081518084526117fa8160208601602086016117b2565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60006020808301818452808551808352604092508286019150828160051b87010184880160005b838110156118db578883037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc001855281518051845287810151888501528681015167ffffffffffffffff16878501526060808201519085015260809081015160a0918501829052906118c7818601836117e2565b968901969450505090860190600101611853565b509098975050505050505050565b60208152600061070560208301846117e2565b60008083601f84011261190e57600080fd5b50813567ffffffffffffffff81111561192657600080fd5b60208301915083602082850101111561193e57600080fd5b9250929050565b6000806000806060858703121561195b57600080fd5b611964856116c8565b935060208501359250604085013567ffffffffffffffff81111561198757600080fd5b611993878288016118fc565b95989497509550505050565b803567ffffffffffffffff811681146116dc57600080fd5b600080600080600080600060a0888a0312156119d257600080fd5b6119db886116c8565b9650602088013567ffffffffffffffff808211156119f857600080fd5b818a0191508a601f830112611a0c57600080fd5b813581811115611a1b57600080fd5b8b60208260051b8501011115611a3057600080fd5b6020830198509650611a4460408b0161199f565b9550611a5260608b0161199f565b945060808a0135915080821115611a6857600080fd5b50611a758a828b016118fc565b989b979a50959850939692959293505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611afe57611afe611a88565b604052919050565b600067ffffffffffffffff821115611b2057611b20611a88565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600080600060608486031215611b6157600080fd5b611b6a846116c8565b925060208401359150604084013567ffffffffffffffff811115611b8d57600080fd5b8401601f81018613611b9e57600080fd5b8035611bb1611bac82611b06565b611ab7565b818152876020838501011115611bc657600080fd5b816020840160208301376000602083830101528093505050509250925092565b600060208284031215611bf857600080fd5b5035919050565b600060208284031215611c1157600080fd5b8135610705816116e1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060208284031215611c5d57600080fd5b815167ffffffffffffffff811115611c7457600080fd5b8201601f81018413611c8557600080fd5b8051611c93611bac82611b06565b818152856020838501011115611ca857600080fd5b6113cb8260208301602086016117b2565b600060208284031215611ccb57600080fd5b5051919050565b600082821015611d0b577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500390565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008560601b16815283601482015282603482015260008251611d598160548501602087016117b2565b9190910160540195945050505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008660601b1681528460148201528360348201528183605483013760009101605401908152949350505050565b63ffffffff841681528260208201526060604082015260006113cb60608301846117e2565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600060208284031215611e1c57600080fd5b81516003811061070557600080fd5b815160009082906020808601845b83811015611e5557815185529382019390820190600101611e39565b50929695505050505050565b8581527fffffffff000000000000000000000000000000000000000000000000000000008560e01b1660208201527fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008460601b1660248201527fffffffffffffffff0000000000000000000000000000000000000000000000008360c01b16603882015260008251611efa8160408501602087016117b2565b91909101604001969550505050505056fea164736f6c634300080f000a496e697469616c697a61626c653a20636f6e7472616374206973206e6f742069",
}

// DisputeGameFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use DisputeGameFactoryMetaData.ABI instead.
var DisputeGameFactoryABI = DisputeGameFactoryMetaData.ABI

// DisputeGameFactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DisputeGameFactoryMetaData.Bin instead.
var DisputeGameFactoryBin = DisputeGameFactoryMetaData.Bin

// DeployDisputeGameFactory deploys a new Ethereum contract, binding an instance of DisputeGameFactory to it.
func DeployDisputeGameFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DisputeGameFactory, error) {
	parsed, err := DisputeGameFactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DisputeGameFactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DisputeGameFactory{DisputeGameFactoryCaller: DisputeGameFactoryCaller{contract: contract}, DisputeGameFactoryTransactor: DisputeGameFactoryTransactor{contract: contract}, DisputeGameFactoryFilterer: DisputeGameFactoryFilterer{contract: contract}}, nil
}

// DisputeGameFactory is an auto generated Go binding around an Ethereum contract.
type DisputeGameFactory struct {
	DisputeGameFactoryCaller     // Read-only binding to the contract
	DisputeGameFactoryTransactor // Write-only binding to the contract
	DisputeGameFactoryFilterer   // Log filterer for contract events
}

// DisputeGameFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DisputeGameFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisputeGameFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DisputeGameFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisputeGameFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DisputeGameFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisputeGameFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DisputeGameFactorySession struct {
	Contract     *DisputeGameFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DisputeGameFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DisputeGameFactoryCallerSession struct {
	Contract *DisputeGameFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DisputeGameFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DisputeGameFactoryTransactorSession struct {
	Contract     *DisputeGameFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DisputeGameFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DisputeGameFactoryRaw struct {
	Contract *DisputeGameFactory // Generic contract binding to access the raw methods on
}

// DisputeGameFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DisputeGameFactoryCallerRaw struct {
	Contract *DisputeGameFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// DisputeGameFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DisputeGameFactoryTransactorRaw struct {
	Contract *DisputeGameFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDisputeGameFactory creates a new instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactory(address common.Address, backend bind.ContractBackend) (*DisputeGameFactory, error) {
	contract, err := bindDisputeGameFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactory{DisputeGameFactoryCaller: DisputeGameFactoryCaller{contract: contract}, DisputeGameFactoryTransactor: DisputeGameFactoryTransactor{contract: contract}, DisputeGameFactoryFilterer: DisputeGameFactoryFilterer{contract: contract}}, nil
}

// NewDisputeGameFactoryCaller creates a new read-only instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactoryCaller(address common.Address, caller bind.ContractCaller) (*DisputeGameFactoryCaller, error) {
	contract, err := bindDisputeGameFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryCaller{contract: contract}, nil
}

// NewDisputeGameFactoryTransactor creates a new write-only instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*DisputeGameFactoryTransactor, error) {
	contract, err := bindDisputeGameFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryTransactor{contract: contract}, nil
}

// NewDisputeGameFactoryFilterer creates a new log filterer instance of DisputeGameFactory, bound to a specific deployed contract.
func NewDisputeGameFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*DisputeGameFactoryFilterer, error) {
	contract, err := bindDisputeGameFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryFilterer{contract: contract}, nil
}

// bindDisputeGameFactory binds a generic wrapper to an already deployed contract.
func bindDisputeGameFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DisputeGameFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DisputeGameFactory *DisputeGameFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DisputeGameFactory.Contract.DisputeGameFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DisputeGameFactory *DisputeGameFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.DisputeGameFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DisputeGameFactory *DisputeGameFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.DisputeGameFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DisputeGameFactory *DisputeGameFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DisputeGameFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DisputeGameFactory *DisputeGameFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DisputeGameFactory *DisputeGameFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.contract.Transact(opts, method, params...)
}

// FindLatestGames is a free data retrieval call binding the contract method 0x254bd683.
//
// Solidity: function findLatestGames(uint32 _gameType, uint256 _start, uint256 _n) view returns((uint256,bytes32,uint64,bytes32,bytes)[] games_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) FindLatestGames(opts *bind.CallOpts, _gameType uint32, _start *big.Int, _n *big.Int) ([]IDisputeGameFactoryGameSearchResult, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "findLatestGames", _gameType, _start, _n)

	if err != nil {
		return *new([]IDisputeGameFactoryGameSearchResult), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDisputeGameFactoryGameSearchResult)).(*[]IDisputeGameFactoryGameSearchResult)

	return out0, err

}

// FindLatestGames is a free data retrieval call binding the contract method 0x254bd683.
//
// Solidity: function findLatestGames(uint32 _gameType, uint256 _start, uint256 _n) view returns((uint256,bytes32,uint64,bytes32,bytes)[] games_)
func (_DisputeGameFactory *DisputeGameFactorySession) FindLatestGames(_gameType uint32, _start *big.Int, _n *big.Int) ([]IDisputeGameFactoryGameSearchResult, error) {
	return _DisputeGameFactory.Contract.FindLatestGames(&_DisputeGameFactory.CallOpts, _gameType, _start, _n)
}

// FindLatestGames is a free data retrieval call binding the contract method 0x254bd683.
//
// Solidity: function findLatestGames(uint32 _gameType, uint256 _start, uint256 _n) view returns((uint256,bytes32,uint64,bytes32,bytes)[] games_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) FindLatestGames(_gameType uint32, _start *big.Int, _n *big.Int) ([]IDisputeGameFactoryGameSearchResult, error) {
	return _DisputeGameFactory.Contract.FindLatestGames(&_DisputeGameFactory.CallOpts, _gameType, _start, _n)
}

// GameAtIndex is a free data retrieval call binding the contract method 0xbb8aa1fc.
//
// Solidity: function gameAtIndex(uint256 _index) view returns(uint32 gameType_, uint64 timestamp_, address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GameAtIndex(opts *bind.CallOpts, _index *big.Int) (struct {
	GameType  uint32
	Timestamp uint64
	Proxy     common.Address
}, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "gameAtIndex", _index)

	outstruct := new(struct {
		GameType  uint32
		Timestamp uint64
		Proxy     common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameType = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Proxy = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GameAtIndex is a free data retrieval call binding the contract method 0xbb8aa1fc.
//
// Solidity: function gameAtIndex(uint256 _index) view returns(uint32 gameType_, uint64 timestamp_, address proxy_)
func (_DisputeGameFactory *DisputeGameFactorySession) GameAtIndex(_index *big.Int) (struct {
	GameType  uint32
	Timestamp uint64
	Proxy     common.Address
}, error) {
	return _DisputeGameFactory.Contract.GameAtIndex(&_DisputeGameFactory.CallOpts, _index)
}

// GameAtIndex is a free data retrieval call binding the contract method 0xbb8aa1fc.
//
// Solidity: function gameAtIndex(uint256 _index) view returns(uint32 gameType_, uint64 timestamp_, address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GameAtIndex(_index *big.Int) (struct {
	GameType  uint32
	Timestamp uint64
	Proxy     common.Address
}, error) {
	return _DisputeGameFactory.Contract.GameAtIndex(&_DisputeGameFactory.CallOpts, _index)
}

// GameCount is a free data retrieval call binding the contract method 0x4d1975b4.
//
// Solidity: function gameCount() view returns(uint256 gameCount_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GameCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "gameCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GameCount is a free data retrieval call binding the contract method 0x4d1975b4.
//
// Solidity: function gameCount() view returns(uint256 gameCount_)
func (_DisputeGameFactory *DisputeGameFactorySession) GameCount() (*big.Int, error) {
	return _DisputeGameFactory.Contract.GameCount(&_DisputeGameFactory.CallOpts)
}

// GameCount is a free data retrieval call binding the contract method 0x4d1975b4.
//
// Solidity: function gameCount() view returns(uint256 gameCount_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GameCount() (*big.Int, error) {
	return _DisputeGameFactory.Contract.GameCount(&_DisputeGameFactory.CallOpts)
}

// GameImpls is a free data retrieval call binding the contract method 0x1b685b9e.
//
// Solidity: function gameImpls(uint32 ) view returns(address)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GameImpls(opts *bind.CallOpts, arg0 uint32) (common.Address, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "gameImpls", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GameImpls is a free data retrieval call binding the contract method 0x1b685b9e.
//
// Solidity: function gameImpls(uint32 ) view returns(address)
func (_DisputeGameFactory *DisputeGameFactorySession) GameImpls(arg0 uint32) (common.Address, error) {
	return _DisputeGameFactory.Contract.GameImpls(&_DisputeGameFactory.CallOpts, arg0)
}

// GameImpls is a free data retrieval call binding the contract method 0x1b685b9e.
//
// Solidity: function gameImpls(uint32 ) view returns(address)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GameImpls(arg0 uint32) (common.Address, error) {
	return _DisputeGameFactory.Contract.GameImpls(&_DisputeGameFactory.CallOpts, arg0)
}

// Games is a free data retrieval call binding the contract method 0x5f0150cb.
//
// Solidity: function games(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) view returns(address proxy_, uint64 timestamp_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) Games(opts *bind.CallOpts, _gameType uint32, _rootClaim [32]byte, _extraData []byte) (struct {
	Proxy     common.Address
	Timestamp uint64
}, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "games", _gameType, _rootClaim, _extraData)

	outstruct := new(struct {
		Proxy     common.Address
		Timestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Proxy = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// Games is a free data retrieval call binding the contract method 0x5f0150cb.
//
// Solidity: function games(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) view returns(address proxy_, uint64 timestamp_)
func (_DisputeGameFactory *DisputeGameFactorySession) Games(_gameType uint32, _rootClaim [32]byte, _extraData []byte) (struct {
	Proxy     common.Address
	Timestamp uint64
}, error) {
	return _DisputeGameFactory.Contract.Games(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// Games is a free data retrieval call binding the contract method 0x5f0150cb.
//
// Solidity: function games(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) view returns(address proxy_, uint64 timestamp_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) Games(_gameType uint32, _rootClaim [32]byte, _extraData []byte) (struct {
	Proxy     common.Address
	Timestamp uint64
}, error) {
	return _DisputeGameFactory.Contract.Games(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// GetGameUUID is a free data retrieval call binding the contract method 0x96cd9720.
//
// Solidity: function getGameUUID(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) pure returns(bytes32 uuid_)
func (_DisputeGameFactory *DisputeGameFactoryCaller) GetGameUUID(opts *bind.CallOpts, _gameType uint32, _rootClaim [32]byte, _extraData []byte) ([32]byte, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "getGameUUID", _gameType, _rootClaim, _extraData)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetGameUUID is a free data retrieval call binding the contract method 0x96cd9720.
//
// Solidity: function getGameUUID(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) pure returns(bytes32 uuid_)
func (_DisputeGameFactory *DisputeGameFactorySession) GetGameUUID(_gameType uint32, _rootClaim [32]byte, _extraData []byte) ([32]byte, error) {
	return _DisputeGameFactory.Contract.GetGameUUID(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// GetGameUUID is a free data retrieval call binding the contract method 0x96cd9720.
//
// Solidity: function getGameUUID(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) pure returns(bytes32 uuid_)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) GetGameUUID(_gameType uint32, _rootClaim [32]byte, _extraData []byte) ([32]byte, error) {
	return _DisputeGameFactory.Contract.GetGameUUID(&_DisputeGameFactory.CallOpts, _gameType, _rootClaim, _extraData)
}

// InitBonds is a free data retrieval call binding the contract method 0x6593dc6e.
//
// Solidity: function initBonds(uint32 ) view returns(uint256)
func (_DisputeGameFactory *DisputeGameFactoryCaller) InitBonds(opts *bind.CallOpts, arg0 uint32) (*big.Int, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "initBonds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitBonds is a free data retrieval call binding the contract method 0x6593dc6e.
//
// Solidity: function initBonds(uint32 ) view returns(uint256)
func (_DisputeGameFactory *DisputeGameFactorySession) InitBonds(arg0 uint32) (*big.Int, error) {
	return _DisputeGameFactory.Contract.InitBonds(&_DisputeGameFactory.CallOpts, arg0)
}

// InitBonds is a free data retrieval call binding the contract method 0x6593dc6e.
//
// Solidity: function initBonds(uint32 ) view returns(uint256)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) InitBonds(arg0 uint32) (*big.Int, error) {
	return _DisputeGameFactory.Contract.InitBonds(&_DisputeGameFactory.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DisputeGameFactory *DisputeGameFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DisputeGameFactory *DisputeGameFactorySession) Owner() (common.Address, error) {
	return _DisputeGameFactory.Contract.Owner(&_DisputeGameFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) Owner() (common.Address, error) {
	return _DisputeGameFactory.Contract.Owner(&_DisputeGameFactory.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_DisputeGameFactory *DisputeGameFactoryCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DisputeGameFactory.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_DisputeGameFactory *DisputeGameFactorySession) Version() (string, error) {
	return _DisputeGameFactory.Contract.Version(&_DisputeGameFactory.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_DisputeGameFactory *DisputeGameFactoryCallerSession) Version() (string, error) {
	return _DisputeGameFactory.Contract.Version(&_DisputeGameFactory.CallOpts)
}

// Create is a paid mutator transaction binding the contract method 0x82ecf2f6.
//
// Solidity: function create(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) payable returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryTransactor) Create(opts *bind.TransactOpts, _gameType uint32, _rootClaim [32]byte, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "create", _gameType, _rootClaim, _extraData)
}

// Create is a paid mutator transaction binding the contract method 0x82ecf2f6.
//
// Solidity: function create(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) payable returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactorySession) Create(_gameType uint32, _rootClaim [32]byte, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.Create(&_DisputeGameFactory.TransactOpts, _gameType, _rootClaim, _extraData)
}

// Create is a paid mutator transaction binding the contract method 0x82ecf2f6.
//
// Solidity: function create(uint32 _gameType, bytes32 _rootClaim, bytes _extraData) payable returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) Create(_gameType uint32, _rootClaim [32]byte, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.Create(&_DisputeGameFactory.TransactOpts, _gameType, _rootClaim, _extraData)
}

// CreateZkFaultDisputeGame is a paid mutator transaction binding the contract method 0x77f63c84.
//
// Solidity: function createZkFaultDisputeGame(uint32 _gameType, bytes32[] _claims, uint64 _parentGameIndex, uint64 _l2BlockNumber, bytes _extraData) payable returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryTransactor) CreateZkFaultDisputeGame(opts *bind.TransactOpts, _gameType uint32, _claims [][32]byte, _parentGameIndex uint64, _l2BlockNumber uint64, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "createZkFaultDisputeGame", _gameType, _claims, _parentGameIndex, _l2BlockNumber, _extraData)
}

// CreateZkFaultDisputeGame is a paid mutator transaction binding the contract method 0x77f63c84.
//
// Solidity: function createZkFaultDisputeGame(uint32 _gameType, bytes32[] _claims, uint64 _parentGameIndex, uint64 _l2BlockNumber, bytes _extraData) payable returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactorySession) CreateZkFaultDisputeGame(_gameType uint32, _claims [][32]byte, _parentGameIndex uint64, _l2BlockNumber uint64, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.CreateZkFaultDisputeGame(&_DisputeGameFactory.TransactOpts, _gameType, _claims, _parentGameIndex, _l2BlockNumber, _extraData)
}

// CreateZkFaultDisputeGame is a paid mutator transaction binding the contract method 0x77f63c84.
//
// Solidity: function createZkFaultDisputeGame(uint32 _gameType, bytes32[] _claims, uint64 _parentGameIndex, uint64 _l2BlockNumber, bytes _extraData) payable returns(address proxy_)
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) CreateZkFaultDisputeGame(_gameType uint32, _claims [][32]byte, _parentGameIndex uint64, _l2BlockNumber uint64, _extraData []byte) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.CreateZkFaultDisputeGame(&_DisputeGameFactory.TransactOpts, _gameType, _claims, _parentGameIndex, _l2BlockNumber, _extraData)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "initialize", _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_DisputeGameFactory *DisputeGameFactorySession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.Initialize(&_DisputeGameFactory.TransactOpts, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.Initialize(&_DisputeGameFactory.TransactOpts, _owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DisputeGameFactory *DisputeGameFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.RenounceOwnership(&_DisputeGameFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.RenounceOwnership(&_DisputeGameFactory.TransactOpts)
}

// SetImplementation is a paid mutator transaction binding the contract method 0x14f6b1a3.
//
// Solidity: function setImplementation(uint32 _gameType, address _impl) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactor) SetImplementation(opts *bind.TransactOpts, _gameType uint32, _impl common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "setImplementation", _gameType, _impl)
}

// SetImplementation is a paid mutator transaction binding the contract method 0x14f6b1a3.
//
// Solidity: function setImplementation(uint32 _gameType, address _impl) returns()
func (_DisputeGameFactory *DisputeGameFactorySession) SetImplementation(_gameType uint32, _impl common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.SetImplementation(&_DisputeGameFactory.TransactOpts, _gameType, _impl)
}

// SetImplementation is a paid mutator transaction binding the contract method 0x14f6b1a3.
//
// Solidity: function setImplementation(uint32 _gameType, address _impl) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) SetImplementation(_gameType uint32, _impl common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.SetImplementation(&_DisputeGameFactory.TransactOpts, _gameType, _impl)
}

// SetInitBond is a paid mutator transaction binding the contract method 0x1e334240.
//
// Solidity: function setInitBond(uint32 _gameType, uint256 _initBond) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactor) SetInitBond(opts *bind.TransactOpts, _gameType uint32, _initBond *big.Int) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "setInitBond", _gameType, _initBond)
}

// SetInitBond is a paid mutator transaction binding the contract method 0x1e334240.
//
// Solidity: function setInitBond(uint32 _gameType, uint256 _initBond) returns()
func (_DisputeGameFactory *DisputeGameFactorySession) SetInitBond(_gameType uint32, _initBond *big.Int) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.SetInitBond(&_DisputeGameFactory.TransactOpts, _gameType, _initBond)
}

// SetInitBond is a paid mutator transaction binding the contract method 0x1e334240.
//
// Solidity: function setInitBond(uint32 _gameType, uint256 _initBond) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) SetInitBond(_gameType uint32, _initBond *big.Int) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.SetInitBond(&_DisputeGameFactory.TransactOpts, _gameType, _initBond)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DisputeGameFactory *DisputeGameFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.TransferOwnership(&_DisputeGameFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DisputeGameFactory *DisputeGameFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DisputeGameFactory.Contract.TransferOwnership(&_DisputeGameFactory.TransactOpts, newOwner)
}

// DisputeGameFactoryDisputeGameCreatedIterator is returned from FilterDisputeGameCreated and is used to iterate over the raw logs and unpacked data for DisputeGameCreated events raised by the DisputeGameFactory contract.
type DisputeGameFactoryDisputeGameCreatedIterator struct {
	Event *DisputeGameFactoryDisputeGameCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DisputeGameFactoryDisputeGameCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DisputeGameFactoryDisputeGameCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DisputeGameFactoryDisputeGameCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DisputeGameFactoryDisputeGameCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DisputeGameFactoryDisputeGameCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DisputeGameFactoryDisputeGameCreated represents a DisputeGameCreated event raised by the DisputeGameFactory contract.
type DisputeGameFactoryDisputeGameCreated struct {
	DisputeProxy common.Address
	GameType     uint32
	RootClaim    [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDisputeGameCreated is a free log retrieval operation binding the contract event 0x5b565efe82411da98814f356d0e7bcb8f0219b8d970307c5afb4a6903a8b2e35.
//
// Solidity: event DisputeGameCreated(address indexed disputeProxy, uint32 indexed gameType, bytes32 indexed rootClaim)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) FilterDisputeGameCreated(opts *bind.FilterOpts, disputeProxy []common.Address, gameType []uint32, rootClaim [][32]byte) (*DisputeGameFactoryDisputeGameCreatedIterator, error) {

	var disputeProxyRule []interface{}
	for _, disputeProxyItem := range disputeProxy {
		disputeProxyRule = append(disputeProxyRule, disputeProxyItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}
	var rootClaimRule []interface{}
	for _, rootClaimItem := range rootClaim {
		rootClaimRule = append(rootClaimRule, rootClaimItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.FilterLogs(opts, "DisputeGameCreated", disputeProxyRule, gameTypeRule, rootClaimRule)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryDisputeGameCreatedIterator{contract: _DisputeGameFactory.contract, event: "DisputeGameCreated", logs: logs, sub: sub}, nil
}

// WatchDisputeGameCreated is a free log subscription operation binding the contract event 0x5b565efe82411da98814f356d0e7bcb8f0219b8d970307c5afb4a6903a8b2e35.
//
// Solidity: event DisputeGameCreated(address indexed disputeProxy, uint32 indexed gameType, bytes32 indexed rootClaim)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) WatchDisputeGameCreated(opts *bind.WatchOpts, sink chan<- *DisputeGameFactoryDisputeGameCreated, disputeProxy []common.Address, gameType []uint32, rootClaim [][32]byte) (event.Subscription, error) {

	var disputeProxyRule []interface{}
	for _, disputeProxyItem := range disputeProxy {
		disputeProxyRule = append(disputeProxyRule, disputeProxyItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}
	var rootClaimRule []interface{}
	for _, rootClaimItem := range rootClaim {
		rootClaimRule = append(rootClaimRule, rootClaimItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.WatchLogs(opts, "DisputeGameCreated", disputeProxyRule, gameTypeRule, rootClaimRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DisputeGameFactoryDisputeGameCreated)
				if err := _DisputeGameFactory.contract.UnpackLog(event, "DisputeGameCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDisputeGameCreated is a log parse operation binding the contract event 0x5b565efe82411da98814f356d0e7bcb8f0219b8d970307c5afb4a6903a8b2e35.
//
// Solidity: event DisputeGameCreated(address indexed disputeProxy, uint32 indexed gameType, bytes32 indexed rootClaim)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) ParseDisputeGameCreated(log types.Log) (*DisputeGameFactoryDisputeGameCreated, error) {
	event := new(DisputeGameFactoryDisputeGameCreated)
	if err := _DisputeGameFactory.contract.UnpackLog(event, "DisputeGameCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DisputeGameFactoryImplementationSetIterator is returned from FilterImplementationSet and is used to iterate over the raw logs and unpacked data for ImplementationSet events raised by the DisputeGameFactory contract.
type DisputeGameFactoryImplementationSetIterator struct {
	Event *DisputeGameFactoryImplementationSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DisputeGameFactoryImplementationSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DisputeGameFactoryImplementationSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DisputeGameFactoryImplementationSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DisputeGameFactoryImplementationSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DisputeGameFactoryImplementationSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DisputeGameFactoryImplementationSet represents a ImplementationSet event raised by the DisputeGameFactory contract.
type DisputeGameFactoryImplementationSet struct {
	Impl     common.Address
	GameType uint32
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterImplementationSet is a free log retrieval operation binding the contract event 0xff513d80e2c7fa487608f70a618dfbc0cf415699dc69588c747e8c71566c88de.
//
// Solidity: event ImplementationSet(address indexed impl, uint32 indexed gameType)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) FilterImplementationSet(opts *bind.FilterOpts, impl []common.Address, gameType []uint32) (*DisputeGameFactoryImplementationSetIterator, error) {

	var implRule []interface{}
	for _, implItem := range impl {
		implRule = append(implRule, implItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.FilterLogs(opts, "ImplementationSet", implRule, gameTypeRule)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryImplementationSetIterator{contract: _DisputeGameFactory.contract, event: "ImplementationSet", logs: logs, sub: sub}, nil
}

// WatchImplementationSet is a free log subscription operation binding the contract event 0xff513d80e2c7fa487608f70a618dfbc0cf415699dc69588c747e8c71566c88de.
//
// Solidity: event ImplementationSet(address indexed impl, uint32 indexed gameType)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) WatchImplementationSet(opts *bind.WatchOpts, sink chan<- *DisputeGameFactoryImplementationSet, impl []common.Address, gameType []uint32) (event.Subscription, error) {

	var implRule []interface{}
	for _, implItem := range impl {
		implRule = append(implRule, implItem)
	}
	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.WatchLogs(opts, "ImplementationSet", implRule, gameTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DisputeGameFactoryImplementationSet)
				if err := _DisputeGameFactory.contract.UnpackLog(event, "ImplementationSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseImplementationSet is a log parse operation binding the contract event 0xff513d80e2c7fa487608f70a618dfbc0cf415699dc69588c747e8c71566c88de.
//
// Solidity: event ImplementationSet(address indexed impl, uint32 indexed gameType)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) ParseImplementationSet(log types.Log) (*DisputeGameFactoryImplementationSet, error) {
	event := new(DisputeGameFactoryImplementationSet)
	if err := _DisputeGameFactory.contract.UnpackLog(event, "ImplementationSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DisputeGameFactoryInitBondUpdatedIterator is returned from FilterInitBondUpdated and is used to iterate over the raw logs and unpacked data for InitBondUpdated events raised by the DisputeGameFactory contract.
type DisputeGameFactoryInitBondUpdatedIterator struct {
	Event *DisputeGameFactoryInitBondUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DisputeGameFactoryInitBondUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DisputeGameFactoryInitBondUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DisputeGameFactoryInitBondUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DisputeGameFactoryInitBondUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DisputeGameFactoryInitBondUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DisputeGameFactoryInitBondUpdated represents a InitBondUpdated event raised by the DisputeGameFactory contract.
type DisputeGameFactoryInitBondUpdated struct {
	GameType uint32
	NewBond  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInitBondUpdated is a free log retrieval operation binding the contract event 0x74d6665c4b26d5596a5aa13d3014e0c06af4d322075a797f87b03cd4c5bc91ca.
//
// Solidity: event InitBondUpdated(uint32 indexed gameType, uint256 indexed newBond)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) FilterInitBondUpdated(opts *bind.FilterOpts, gameType []uint32, newBond []*big.Int) (*DisputeGameFactoryInitBondUpdatedIterator, error) {

	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}
	var newBondRule []interface{}
	for _, newBondItem := range newBond {
		newBondRule = append(newBondRule, newBondItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.FilterLogs(opts, "InitBondUpdated", gameTypeRule, newBondRule)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryInitBondUpdatedIterator{contract: _DisputeGameFactory.contract, event: "InitBondUpdated", logs: logs, sub: sub}, nil
}

// WatchInitBondUpdated is a free log subscription operation binding the contract event 0x74d6665c4b26d5596a5aa13d3014e0c06af4d322075a797f87b03cd4c5bc91ca.
//
// Solidity: event InitBondUpdated(uint32 indexed gameType, uint256 indexed newBond)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) WatchInitBondUpdated(opts *bind.WatchOpts, sink chan<- *DisputeGameFactoryInitBondUpdated, gameType []uint32, newBond []*big.Int) (event.Subscription, error) {

	var gameTypeRule []interface{}
	for _, gameTypeItem := range gameType {
		gameTypeRule = append(gameTypeRule, gameTypeItem)
	}
	var newBondRule []interface{}
	for _, newBondItem := range newBond {
		newBondRule = append(newBondRule, newBondItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.WatchLogs(opts, "InitBondUpdated", gameTypeRule, newBondRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DisputeGameFactoryInitBondUpdated)
				if err := _DisputeGameFactory.contract.UnpackLog(event, "InitBondUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitBondUpdated is a log parse operation binding the contract event 0x74d6665c4b26d5596a5aa13d3014e0c06af4d322075a797f87b03cd4c5bc91ca.
//
// Solidity: event InitBondUpdated(uint32 indexed gameType, uint256 indexed newBond)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) ParseInitBondUpdated(log types.Log) (*DisputeGameFactoryInitBondUpdated, error) {
	event := new(DisputeGameFactoryInitBondUpdated)
	if err := _DisputeGameFactory.contract.UnpackLog(event, "InitBondUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DisputeGameFactoryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DisputeGameFactory contract.
type DisputeGameFactoryInitializedIterator struct {
	Event *DisputeGameFactoryInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DisputeGameFactoryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DisputeGameFactoryInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DisputeGameFactoryInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DisputeGameFactoryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DisputeGameFactoryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DisputeGameFactoryInitialized represents a Initialized event raised by the DisputeGameFactory contract.
type DisputeGameFactoryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) FilterInitialized(opts *bind.FilterOpts) (*DisputeGameFactoryInitializedIterator, error) {

	logs, sub, err := _DisputeGameFactory.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryInitializedIterator{contract: _DisputeGameFactory.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DisputeGameFactoryInitialized) (event.Subscription, error) {

	logs, sub, err := _DisputeGameFactory.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DisputeGameFactoryInitialized)
				if err := _DisputeGameFactory.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) ParseInitialized(log types.Log) (*DisputeGameFactoryInitialized, error) {
	event := new(DisputeGameFactoryInitialized)
	if err := _DisputeGameFactory.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DisputeGameFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DisputeGameFactory contract.
type DisputeGameFactoryOwnershipTransferredIterator struct {
	Event *DisputeGameFactoryOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DisputeGameFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DisputeGameFactoryOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DisputeGameFactoryOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DisputeGameFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DisputeGameFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DisputeGameFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the DisputeGameFactory contract.
type DisputeGameFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DisputeGameFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DisputeGameFactoryOwnershipTransferredIterator{contract: _DisputeGameFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DisputeGameFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DisputeGameFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DisputeGameFactoryOwnershipTransferred)
				if err := _DisputeGameFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DisputeGameFactory *DisputeGameFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*DisputeGameFactoryOwnershipTransferred, error) {
	event := new(DisputeGameFactoryOwnershipTransferred)
	if err := _DisputeGameFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
