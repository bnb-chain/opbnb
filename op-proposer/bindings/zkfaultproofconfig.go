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

// ZKFaultProofConfigMetaData contains all meta data concerning the ZKFaultProofConfig contract.
var ZKFaultProofConfigMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"aggregationVkey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockDistance\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"chainId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"historicBlockHashes\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_blockDistance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_aggregationVkey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_rangeVkeyCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_verifierGateway\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_rollupConfigHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rangeVkeyCommitment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollupConfigHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAggregationVKey\",\"inputs\":[{\"name\":\"_aggregationVKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRangeVkeyCommitment\",\"inputs\":[{\"name\":\"_rangeVkeyCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRollupConfigHash\",\"inputs\":[{\"name\":\"_rollupConfigHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateVerifierGateway\",\"inputs\":[{\"name\":\"_verifierGateway\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifierGateway\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractSP1VerifierGateway\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdatedAggregationVKey\",\"inputs\":[{\"name\":\"oldVkey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newVkey\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdatedRangeVkeyCommitment\",\"inputs\":[{\"name\":\"oldRangeVkeyCommitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newRangeVkeyCommitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdatedRollupConfigHash\",\"inputs\":[{\"name\":\"oldRollupConfigHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newRollupConfigHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdatedVerifierGateway\",\"inputs\":[{\"name\":\"oldVerifierGateway\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newVerifierGateway\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b506109fd806100206000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80638da5cb5b11610097578063bc91ce3311610066578063bc91ce331461023a578063c32e4e3e1461024d578063f2fde38b14610256578063fb3c491c1461026957600080fd5b80638da5cb5b146101bf5780639a8a0592146101fe578063a196b52514610207578063a660743b1461022757600080fd5b806354fd4d50116100d357806354fd4d50146101525780636d9a1c8b1461019b578063715018a6146101a457806380fdb3e1146101ac57600080fd5b80631bdd450c146101055780632b31841e1461011a57806337416d82146101365780634418db5e1461013f575b600080fd5b6101186101133660046108b6565b610289565b005b61012360695481565b6040519081526020015b60405180910390f35b61012360655481565b61011861014d3660046108f8565b6102c6565b61018e6040518060400160405280600581526020017f312e302e3000000000000000000000000000000000000000000000000000000081525081565b60405161012d919061091a565b61012360665481565b61011861033f565b6101186101ba3660046108b6565b610353565b60335473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012d565b61012360675481565b6101236102153660046108b6565b606b6020526000908152604090205481565b61011861023536600461098d565b610390565b6101186102483660046108b6565b610590565b61012360685481565b6101186102643660046108f8565b6105cd565b606a546101d99073ffffffffffffffffffffffffffffffffffffffff1681565b610291610684565b6066819055604051819081907fda2f5f014ada26cff39a0f2a9dc6fa4fca1581376fc91ec09506c8fb8657bc3590600090a350565b6102ce610684565b606a80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff831690811790915560405181907f1379941631ff0ed9178ab16ab67a2e5db3aeada7f87e518f761e79c8e38377e390600090a350565b610347610684565b6103516000610705565b565b61035b610684565b6068819055604051819081907fb81f9c41933b730a90fba96ab14541de7cab774f762ea0c183054947bc49aee790600090a350565b600054610100900460ff16158080156103b05750600054600160ff909116105b806103ca5750303b1580156103ca575060005460ff166001145b61045b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156104b957600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6104c161077c565b6104ca886105cd565b6065879055606786905560688590556069849055606a80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff85161790556066829055801561058657600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050505050565b610598610684565b6069819055604051819081907f1035606f0606905acdf851342466a5b64406fa798b7440235cd5811cea2850fd90600090a350565b6105d5610684565b73ffffffffffffffffffffffffffffffffffffffff8116610678576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610452565b61068181610705565b50565b60335473ffffffffffffffffffffffffffffffffffffffff163314610351576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610452565b6033805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff16610813576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610452565b610351600054610100900460ff166108ad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610452565b61035133610705565b6000602082840312156108c857600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146108f357600080fd5b919050565b60006020828403121561090a57600080fd5b610913826108cf565b9392505050565b600060208083528351808285015260005b818110156109475785810183015185820160400152820161092b565b81811115610959576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b600080600080600080600060e0888a0312156109a857600080fd5b6109b1886108cf565b9650602088013595506040880135945060608801359350608088013592506109db60a089016108cf565b915060c088013590509295989194975092955056fea164736f6c634300080f000a",
}

// ZKFaultProofConfigABI is the input ABI used to generate the binding from.
// Deprecated: Use ZKFaultProofConfigMetaData.ABI instead.
var ZKFaultProofConfigABI = ZKFaultProofConfigMetaData.ABI

// ZKFaultProofConfigBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZKFaultProofConfigMetaData.Bin instead.
var ZKFaultProofConfigBin = ZKFaultProofConfigMetaData.Bin

// DeployZKFaultProofConfig deploys a new Ethereum contract, binding an instance of ZKFaultProofConfig to it.
func DeployZKFaultProofConfig(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ZKFaultProofConfig, error) {
	parsed, err := ZKFaultProofConfigMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZKFaultProofConfigBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ZKFaultProofConfig{ZKFaultProofConfigCaller: ZKFaultProofConfigCaller{contract: contract}, ZKFaultProofConfigTransactor: ZKFaultProofConfigTransactor{contract: contract}, ZKFaultProofConfigFilterer: ZKFaultProofConfigFilterer{contract: contract}}, nil
}

// ZKFaultProofConfig is an auto generated Go binding around an Ethereum contract.
type ZKFaultProofConfig struct {
	ZKFaultProofConfigCaller     // Read-only binding to the contract
	ZKFaultProofConfigTransactor // Write-only binding to the contract
	ZKFaultProofConfigFilterer   // Log filterer for contract events
}

// ZKFaultProofConfigCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZKFaultProofConfigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKFaultProofConfigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZKFaultProofConfigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKFaultProofConfigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZKFaultProofConfigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKFaultProofConfigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZKFaultProofConfigSession struct {
	Contract     *ZKFaultProofConfig // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ZKFaultProofConfigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZKFaultProofConfigCallerSession struct {
	Contract *ZKFaultProofConfigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ZKFaultProofConfigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZKFaultProofConfigTransactorSession struct {
	Contract     *ZKFaultProofConfigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ZKFaultProofConfigRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZKFaultProofConfigRaw struct {
	Contract *ZKFaultProofConfig // Generic contract binding to access the raw methods on
}

// ZKFaultProofConfigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZKFaultProofConfigCallerRaw struct {
	Contract *ZKFaultProofConfigCaller // Generic read-only contract binding to access the raw methods on
}

// ZKFaultProofConfigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZKFaultProofConfigTransactorRaw struct {
	Contract *ZKFaultProofConfigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZKFaultProofConfig creates a new instance of ZKFaultProofConfig, bound to a specific deployed contract.
func NewZKFaultProofConfig(address common.Address, backend bind.ContractBackend) (*ZKFaultProofConfig, error) {
	contract, err := bindZKFaultProofConfig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfig{ZKFaultProofConfigCaller: ZKFaultProofConfigCaller{contract: contract}, ZKFaultProofConfigTransactor: ZKFaultProofConfigTransactor{contract: contract}, ZKFaultProofConfigFilterer: ZKFaultProofConfigFilterer{contract: contract}}, nil
}

// NewZKFaultProofConfigCaller creates a new read-only instance of ZKFaultProofConfig, bound to a specific deployed contract.
func NewZKFaultProofConfigCaller(address common.Address, caller bind.ContractCaller) (*ZKFaultProofConfigCaller, error) {
	contract, err := bindZKFaultProofConfig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigCaller{contract: contract}, nil
}

// NewZKFaultProofConfigTransactor creates a new write-only instance of ZKFaultProofConfig, bound to a specific deployed contract.
func NewZKFaultProofConfigTransactor(address common.Address, transactor bind.ContractTransactor) (*ZKFaultProofConfigTransactor, error) {
	contract, err := bindZKFaultProofConfig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigTransactor{contract: contract}, nil
}

// NewZKFaultProofConfigFilterer creates a new log filterer instance of ZKFaultProofConfig, bound to a specific deployed contract.
func NewZKFaultProofConfigFilterer(address common.Address, filterer bind.ContractFilterer) (*ZKFaultProofConfigFilterer, error) {
	contract, err := bindZKFaultProofConfig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigFilterer{contract: contract}, nil
}

// bindZKFaultProofConfig binds a generic wrapper to an already deployed contract.
func bindZKFaultProofConfig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZKFaultProofConfigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZKFaultProofConfig *ZKFaultProofConfigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZKFaultProofConfig.Contract.ZKFaultProofConfigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZKFaultProofConfig *ZKFaultProofConfigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.ZKFaultProofConfigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZKFaultProofConfig *ZKFaultProofConfigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.ZKFaultProofConfigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZKFaultProofConfig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.contract.Transact(opts, method, params...)
}

// AggregationVkey is a free data retrieval call binding the contract method 0xc32e4e3e.
//
// Solidity: function aggregationVkey() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) AggregationVkey(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "aggregationVkey")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AggregationVkey is a free data retrieval call binding the contract method 0xc32e4e3e.
//
// Solidity: function aggregationVkey() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) AggregationVkey() ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.AggregationVkey(&_ZKFaultProofConfig.CallOpts)
}

// AggregationVkey is a free data retrieval call binding the contract method 0xc32e4e3e.
//
// Solidity: function aggregationVkey() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) AggregationVkey() ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.AggregationVkey(&_ZKFaultProofConfig.CallOpts)
}

// BlockDistance is a free data retrieval call binding the contract method 0x37416d82.
//
// Solidity: function blockDistance() view returns(uint256)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) BlockDistance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "blockDistance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BlockDistance is a free data retrieval call binding the contract method 0x37416d82.
//
// Solidity: function blockDistance() view returns(uint256)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) BlockDistance() (*big.Int, error) {
	return _ZKFaultProofConfig.Contract.BlockDistance(&_ZKFaultProofConfig.CallOpts)
}

// BlockDistance is a free data retrieval call binding the contract method 0x37416d82.
//
// Solidity: function blockDistance() view returns(uint256)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) BlockDistance() (*big.Int, error) {
	return _ZKFaultProofConfig.Contract.BlockDistance(&_ZKFaultProofConfig.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "chainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) ChainId() (*big.Int, error) {
	return _ZKFaultProofConfig.Contract.ChainId(&_ZKFaultProofConfig.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() view returns(uint256)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) ChainId() (*big.Int, error) {
	return _ZKFaultProofConfig.Contract.ChainId(&_ZKFaultProofConfig.CallOpts)
}

// HistoricBlockHashes is a free data retrieval call binding the contract method 0xa196b525.
//
// Solidity: function historicBlockHashes(uint256 ) view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) HistoricBlockHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "historicBlockHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HistoricBlockHashes is a free data retrieval call binding the contract method 0xa196b525.
//
// Solidity: function historicBlockHashes(uint256 ) view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) HistoricBlockHashes(arg0 *big.Int) ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.HistoricBlockHashes(&_ZKFaultProofConfig.CallOpts, arg0)
}

// HistoricBlockHashes is a free data retrieval call binding the contract method 0xa196b525.
//
// Solidity: function historicBlockHashes(uint256 ) view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) HistoricBlockHashes(arg0 *big.Int) ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.HistoricBlockHashes(&_ZKFaultProofConfig.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) Owner() (common.Address, error) {
	return _ZKFaultProofConfig.Contract.Owner(&_ZKFaultProofConfig.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) Owner() (common.Address, error) {
	return _ZKFaultProofConfig.Contract.Owner(&_ZKFaultProofConfig.CallOpts)
}

// RangeVkeyCommitment is a free data retrieval call binding the contract method 0x2b31841e.
//
// Solidity: function rangeVkeyCommitment() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) RangeVkeyCommitment(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "rangeVkeyCommitment")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RangeVkeyCommitment is a free data retrieval call binding the contract method 0x2b31841e.
//
// Solidity: function rangeVkeyCommitment() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) RangeVkeyCommitment() ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.RangeVkeyCommitment(&_ZKFaultProofConfig.CallOpts)
}

// RangeVkeyCommitment is a free data retrieval call binding the contract method 0x2b31841e.
//
// Solidity: function rangeVkeyCommitment() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) RangeVkeyCommitment() ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.RangeVkeyCommitment(&_ZKFaultProofConfig.CallOpts)
}

// RollupConfigHash is a free data retrieval call binding the contract method 0x6d9a1c8b.
//
// Solidity: function rollupConfigHash() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) RollupConfigHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "rollupConfigHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RollupConfigHash is a free data retrieval call binding the contract method 0x6d9a1c8b.
//
// Solidity: function rollupConfigHash() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) RollupConfigHash() ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.RollupConfigHash(&_ZKFaultProofConfig.CallOpts)
}

// RollupConfigHash is a free data retrieval call binding the contract method 0x6d9a1c8b.
//
// Solidity: function rollupConfigHash() view returns(bytes32)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) RollupConfigHash() ([32]byte, error) {
	return _ZKFaultProofConfig.Contract.RollupConfigHash(&_ZKFaultProofConfig.CallOpts)
}

// VerifierGateway is a free data retrieval call binding the contract method 0xfb3c491c.
//
// Solidity: function verifierGateway() view returns(address)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) VerifierGateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "verifierGateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VerifierGateway is a free data retrieval call binding the contract method 0xfb3c491c.
//
// Solidity: function verifierGateway() view returns(address)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) VerifierGateway() (common.Address, error) {
	return _ZKFaultProofConfig.Contract.VerifierGateway(&_ZKFaultProofConfig.CallOpts)
}

// VerifierGateway is a free data retrieval call binding the contract method 0xfb3c491c.
//
// Solidity: function verifierGateway() view returns(address)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) VerifierGateway() (common.Address, error) {
	return _ZKFaultProofConfig.Contract.VerifierGateway(&_ZKFaultProofConfig.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ZKFaultProofConfig *ZKFaultProofConfigCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZKFaultProofConfig.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) Version() (string, error) {
	return _ZKFaultProofConfig.Contract.Version(&_ZKFaultProofConfig.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_ZKFaultProofConfig *ZKFaultProofConfigCallerSession) Version() (string, error) {
	return _ZKFaultProofConfig.Contract.Version(&_ZKFaultProofConfig.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xa660743b.
//
// Solidity: function initialize(address _owner, uint256 _blockDistance, uint256 _chainId, bytes32 _aggregationVkey, bytes32 _rangeVkeyCommitment, address _verifierGateway, bytes32 _rollupConfigHash) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _blockDistance *big.Int, _chainId *big.Int, _aggregationVkey [32]byte, _rangeVkeyCommitment [32]byte, _verifierGateway common.Address, _rollupConfigHash [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.contract.Transact(opts, "initialize", _owner, _blockDistance, _chainId, _aggregationVkey, _rangeVkeyCommitment, _verifierGateway, _rollupConfigHash)
}

// Initialize is a paid mutator transaction binding the contract method 0xa660743b.
//
// Solidity: function initialize(address _owner, uint256 _blockDistance, uint256 _chainId, bytes32 _aggregationVkey, bytes32 _rangeVkeyCommitment, address _verifierGateway, bytes32 _rollupConfigHash) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) Initialize(_owner common.Address, _blockDistance *big.Int, _chainId *big.Int, _aggregationVkey [32]byte, _rangeVkeyCommitment [32]byte, _verifierGateway common.Address, _rollupConfigHash [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.Initialize(&_ZKFaultProofConfig.TransactOpts, _owner, _blockDistance, _chainId, _aggregationVkey, _rangeVkeyCommitment, _verifierGateway, _rollupConfigHash)
}

// Initialize is a paid mutator transaction binding the contract method 0xa660743b.
//
// Solidity: function initialize(address _owner, uint256 _blockDistance, uint256 _chainId, bytes32 _aggregationVkey, bytes32 _rangeVkeyCommitment, address _verifierGateway, bytes32 _rollupConfigHash) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorSession) Initialize(_owner common.Address, _blockDistance *big.Int, _chainId *big.Int, _aggregationVkey [32]byte, _rangeVkeyCommitment [32]byte, _verifierGateway common.Address, _rollupConfigHash [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.Initialize(&_ZKFaultProofConfig.TransactOpts, _owner, _blockDistance, _chainId, _aggregationVkey, _rangeVkeyCommitment, _verifierGateway, _rollupConfigHash)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZKFaultProofConfig.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.RenounceOwnership(&_ZKFaultProofConfig.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.RenounceOwnership(&_ZKFaultProofConfig.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ZKFaultProofConfig.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.TransferOwnership(&_ZKFaultProofConfig.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.TransferOwnership(&_ZKFaultProofConfig.TransactOpts, newOwner)
}

// UpdateAggregationVKey is a paid mutator transaction binding the contract method 0x80fdb3e1.
//
// Solidity: function updateAggregationVKey(bytes32 _aggregationVKey) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactor) UpdateAggregationVKey(opts *bind.TransactOpts, _aggregationVKey [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.contract.Transact(opts, "updateAggregationVKey", _aggregationVKey)
}

// UpdateAggregationVKey is a paid mutator transaction binding the contract method 0x80fdb3e1.
//
// Solidity: function updateAggregationVKey(bytes32 _aggregationVKey) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) UpdateAggregationVKey(_aggregationVKey [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateAggregationVKey(&_ZKFaultProofConfig.TransactOpts, _aggregationVKey)
}

// UpdateAggregationVKey is a paid mutator transaction binding the contract method 0x80fdb3e1.
//
// Solidity: function updateAggregationVKey(bytes32 _aggregationVKey) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorSession) UpdateAggregationVKey(_aggregationVKey [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateAggregationVKey(&_ZKFaultProofConfig.TransactOpts, _aggregationVKey)
}

// UpdateRangeVkeyCommitment is a paid mutator transaction binding the contract method 0xbc91ce33.
//
// Solidity: function updateRangeVkeyCommitment(bytes32 _rangeVkeyCommitment) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactor) UpdateRangeVkeyCommitment(opts *bind.TransactOpts, _rangeVkeyCommitment [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.contract.Transact(opts, "updateRangeVkeyCommitment", _rangeVkeyCommitment)
}

// UpdateRangeVkeyCommitment is a paid mutator transaction binding the contract method 0xbc91ce33.
//
// Solidity: function updateRangeVkeyCommitment(bytes32 _rangeVkeyCommitment) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) UpdateRangeVkeyCommitment(_rangeVkeyCommitment [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateRangeVkeyCommitment(&_ZKFaultProofConfig.TransactOpts, _rangeVkeyCommitment)
}

// UpdateRangeVkeyCommitment is a paid mutator transaction binding the contract method 0xbc91ce33.
//
// Solidity: function updateRangeVkeyCommitment(bytes32 _rangeVkeyCommitment) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorSession) UpdateRangeVkeyCommitment(_rangeVkeyCommitment [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateRangeVkeyCommitment(&_ZKFaultProofConfig.TransactOpts, _rangeVkeyCommitment)
}

// UpdateRollupConfigHash is a paid mutator transaction binding the contract method 0x1bdd450c.
//
// Solidity: function updateRollupConfigHash(bytes32 _rollupConfigHash) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactor) UpdateRollupConfigHash(opts *bind.TransactOpts, _rollupConfigHash [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.contract.Transact(opts, "updateRollupConfigHash", _rollupConfigHash)
}

// UpdateRollupConfigHash is a paid mutator transaction binding the contract method 0x1bdd450c.
//
// Solidity: function updateRollupConfigHash(bytes32 _rollupConfigHash) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) UpdateRollupConfigHash(_rollupConfigHash [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateRollupConfigHash(&_ZKFaultProofConfig.TransactOpts, _rollupConfigHash)
}

// UpdateRollupConfigHash is a paid mutator transaction binding the contract method 0x1bdd450c.
//
// Solidity: function updateRollupConfigHash(bytes32 _rollupConfigHash) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorSession) UpdateRollupConfigHash(_rollupConfigHash [32]byte) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateRollupConfigHash(&_ZKFaultProofConfig.TransactOpts, _rollupConfigHash)
}

// UpdateVerifierGateway is a paid mutator transaction binding the contract method 0x4418db5e.
//
// Solidity: function updateVerifierGateway(address _verifierGateway) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactor) UpdateVerifierGateway(opts *bind.TransactOpts, _verifierGateway common.Address) (*types.Transaction, error) {
	return _ZKFaultProofConfig.contract.Transact(opts, "updateVerifierGateway", _verifierGateway)
}

// UpdateVerifierGateway is a paid mutator transaction binding the contract method 0x4418db5e.
//
// Solidity: function updateVerifierGateway(address _verifierGateway) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigSession) UpdateVerifierGateway(_verifierGateway common.Address) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateVerifierGateway(&_ZKFaultProofConfig.TransactOpts, _verifierGateway)
}

// UpdateVerifierGateway is a paid mutator transaction binding the contract method 0x4418db5e.
//
// Solidity: function updateVerifierGateway(address _verifierGateway) returns()
func (_ZKFaultProofConfig *ZKFaultProofConfigTransactorSession) UpdateVerifierGateway(_verifierGateway common.Address) (*types.Transaction, error) {
	return _ZKFaultProofConfig.Contract.UpdateVerifierGateway(&_ZKFaultProofConfig.TransactOpts, _verifierGateway)
}

// ZKFaultProofConfigInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigInitializedIterator struct {
	Event *ZKFaultProofConfigInitialized // Event containing the contract specifics and raw log

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
func (it *ZKFaultProofConfigInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKFaultProofConfigInitialized)
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
		it.Event = new(ZKFaultProofConfigInitialized)
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
func (it *ZKFaultProofConfigInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKFaultProofConfigInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKFaultProofConfigInitialized represents a Initialized event raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) FilterInitialized(opts *bind.FilterOpts) (*ZKFaultProofConfigInitializedIterator, error) {

	logs, sub, err := _ZKFaultProofConfig.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigInitializedIterator{contract: _ZKFaultProofConfig.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ZKFaultProofConfigInitialized) (event.Subscription, error) {

	logs, sub, err := _ZKFaultProofConfig.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKFaultProofConfigInitialized)
				if err := _ZKFaultProofConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) ParseInitialized(log types.Log) (*ZKFaultProofConfigInitialized, error) {
	event := new(ZKFaultProofConfigInitialized)
	if err := _ZKFaultProofConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKFaultProofConfigOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigOwnershipTransferredIterator struct {
	Event *ZKFaultProofConfigOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ZKFaultProofConfigOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKFaultProofConfigOwnershipTransferred)
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
		it.Event = new(ZKFaultProofConfigOwnershipTransferred)
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
func (it *ZKFaultProofConfigOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKFaultProofConfigOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKFaultProofConfigOwnershipTransferred represents a OwnershipTransferred event raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ZKFaultProofConfigOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigOwnershipTransferredIterator{contract: _ZKFaultProofConfig.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ZKFaultProofConfigOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKFaultProofConfigOwnershipTransferred)
				if err := _ZKFaultProofConfig.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) ParseOwnershipTransferred(log types.Log) (*ZKFaultProofConfigOwnershipTransferred, error) {
	event := new(ZKFaultProofConfigOwnershipTransferred)
	if err := _ZKFaultProofConfig.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKFaultProofConfigUpdatedAggregationVKeyIterator is returned from FilterUpdatedAggregationVKey and is used to iterate over the raw logs and unpacked data for UpdatedAggregationVKey events raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedAggregationVKeyIterator struct {
	Event *ZKFaultProofConfigUpdatedAggregationVKey // Event containing the contract specifics and raw log

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
func (it *ZKFaultProofConfigUpdatedAggregationVKeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKFaultProofConfigUpdatedAggregationVKey)
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
		it.Event = new(ZKFaultProofConfigUpdatedAggregationVKey)
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
func (it *ZKFaultProofConfigUpdatedAggregationVKeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKFaultProofConfigUpdatedAggregationVKeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKFaultProofConfigUpdatedAggregationVKey represents a UpdatedAggregationVKey event raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedAggregationVKey struct {
	OldVkey [32]byte
	NewVkey [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUpdatedAggregationVKey is a free log retrieval operation binding the contract event 0xb81f9c41933b730a90fba96ab14541de7cab774f762ea0c183054947bc49aee7.
//
// Solidity: event UpdatedAggregationVKey(bytes32 indexed oldVkey, bytes32 indexed newVkey)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) FilterUpdatedAggregationVKey(opts *bind.FilterOpts, oldVkey [][32]byte, newVkey [][32]byte) (*ZKFaultProofConfigUpdatedAggregationVKeyIterator, error) {

	var oldVkeyRule []interface{}
	for _, oldVkeyItem := range oldVkey {
		oldVkeyRule = append(oldVkeyRule, oldVkeyItem)
	}
	var newVkeyRule []interface{}
	for _, newVkeyItem := range newVkey {
		newVkeyRule = append(newVkeyRule, newVkeyItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.FilterLogs(opts, "UpdatedAggregationVKey", oldVkeyRule, newVkeyRule)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigUpdatedAggregationVKeyIterator{contract: _ZKFaultProofConfig.contract, event: "UpdatedAggregationVKey", logs: logs, sub: sub}, nil
}

// WatchUpdatedAggregationVKey is a free log subscription operation binding the contract event 0xb81f9c41933b730a90fba96ab14541de7cab774f762ea0c183054947bc49aee7.
//
// Solidity: event UpdatedAggregationVKey(bytes32 indexed oldVkey, bytes32 indexed newVkey)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) WatchUpdatedAggregationVKey(opts *bind.WatchOpts, sink chan<- *ZKFaultProofConfigUpdatedAggregationVKey, oldVkey [][32]byte, newVkey [][32]byte) (event.Subscription, error) {

	var oldVkeyRule []interface{}
	for _, oldVkeyItem := range oldVkey {
		oldVkeyRule = append(oldVkeyRule, oldVkeyItem)
	}
	var newVkeyRule []interface{}
	for _, newVkeyItem := range newVkey {
		newVkeyRule = append(newVkeyRule, newVkeyItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.WatchLogs(opts, "UpdatedAggregationVKey", oldVkeyRule, newVkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKFaultProofConfigUpdatedAggregationVKey)
				if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedAggregationVKey", log); err != nil {
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

// ParseUpdatedAggregationVKey is a log parse operation binding the contract event 0xb81f9c41933b730a90fba96ab14541de7cab774f762ea0c183054947bc49aee7.
//
// Solidity: event UpdatedAggregationVKey(bytes32 indexed oldVkey, bytes32 indexed newVkey)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) ParseUpdatedAggregationVKey(log types.Log) (*ZKFaultProofConfigUpdatedAggregationVKey, error) {
	event := new(ZKFaultProofConfigUpdatedAggregationVKey)
	if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedAggregationVKey", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKFaultProofConfigUpdatedRangeVkeyCommitmentIterator is returned from FilterUpdatedRangeVkeyCommitment and is used to iterate over the raw logs and unpacked data for UpdatedRangeVkeyCommitment events raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedRangeVkeyCommitmentIterator struct {
	Event *ZKFaultProofConfigUpdatedRangeVkeyCommitment // Event containing the contract specifics and raw log

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
func (it *ZKFaultProofConfigUpdatedRangeVkeyCommitmentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKFaultProofConfigUpdatedRangeVkeyCommitment)
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
		it.Event = new(ZKFaultProofConfigUpdatedRangeVkeyCommitment)
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
func (it *ZKFaultProofConfigUpdatedRangeVkeyCommitmentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKFaultProofConfigUpdatedRangeVkeyCommitmentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKFaultProofConfigUpdatedRangeVkeyCommitment represents a UpdatedRangeVkeyCommitment event raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedRangeVkeyCommitment struct {
	OldRangeVkeyCommitment [32]byte
	NewRangeVkeyCommitment [32]byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterUpdatedRangeVkeyCommitment is a free log retrieval operation binding the contract event 0x1035606f0606905acdf851342466a5b64406fa798b7440235cd5811cea2850fd.
//
// Solidity: event UpdatedRangeVkeyCommitment(bytes32 indexed oldRangeVkeyCommitment, bytes32 indexed newRangeVkeyCommitment)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) FilterUpdatedRangeVkeyCommitment(opts *bind.FilterOpts, oldRangeVkeyCommitment [][32]byte, newRangeVkeyCommitment [][32]byte) (*ZKFaultProofConfigUpdatedRangeVkeyCommitmentIterator, error) {

	var oldRangeVkeyCommitmentRule []interface{}
	for _, oldRangeVkeyCommitmentItem := range oldRangeVkeyCommitment {
		oldRangeVkeyCommitmentRule = append(oldRangeVkeyCommitmentRule, oldRangeVkeyCommitmentItem)
	}
	var newRangeVkeyCommitmentRule []interface{}
	for _, newRangeVkeyCommitmentItem := range newRangeVkeyCommitment {
		newRangeVkeyCommitmentRule = append(newRangeVkeyCommitmentRule, newRangeVkeyCommitmentItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.FilterLogs(opts, "UpdatedRangeVkeyCommitment", oldRangeVkeyCommitmentRule, newRangeVkeyCommitmentRule)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigUpdatedRangeVkeyCommitmentIterator{contract: _ZKFaultProofConfig.contract, event: "UpdatedRangeVkeyCommitment", logs: logs, sub: sub}, nil
}

// WatchUpdatedRangeVkeyCommitment is a free log subscription operation binding the contract event 0x1035606f0606905acdf851342466a5b64406fa798b7440235cd5811cea2850fd.
//
// Solidity: event UpdatedRangeVkeyCommitment(bytes32 indexed oldRangeVkeyCommitment, bytes32 indexed newRangeVkeyCommitment)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) WatchUpdatedRangeVkeyCommitment(opts *bind.WatchOpts, sink chan<- *ZKFaultProofConfigUpdatedRangeVkeyCommitment, oldRangeVkeyCommitment [][32]byte, newRangeVkeyCommitment [][32]byte) (event.Subscription, error) {

	var oldRangeVkeyCommitmentRule []interface{}
	for _, oldRangeVkeyCommitmentItem := range oldRangeVkeyCommitment {
		oldRangeVkeyCommitmentRule = append(oldRangeVkeyCommitmentRule, oldRangeVkeyCommitmentItem)
	}
	var newRangeVkeyCommitmentRule []interface{}
	for _, newRangeVkeyCommitmentItem := range newRangeVkeyCommitment {
		newRangeVkeyCommitmentRule = append(newRangeVkeyCommitmentRule, newRangeVkeyCommitmentItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.WatchLogs(opts, "UpdatedRangeVkeyCommitment", oldRangeVkeyCommitmentRule, newRangeVkeyCommitmentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKFaultProofConfigUpdatedRangeVkeyCommitment)
				if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedRangeVkeyCommitment", log); err != nil {
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

// ParseUpdatedRangeVkeyCommitment is a log parse operation binding the contract event 0x1035606f0606905acdf851342466a5b64406fa798b7440235cd5811cea2850fd.
//
// Solidity: event UpdatedRangeVkeyCommitment(bytes32 indexed oldRangeVkeyCommitment, bytes32 indexed newRangeVkeyCommitment)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) ParseUpdatedRangeVkeyCommitment(log types.Log) (*ZKFaultProofConfigUpdatedRangeVkeyCommitment, error) {
	event := new(ZKFaultProofConfigUpdatedRangeVkeyCommitment)
	if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedRangeVkeyCommitment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKFaultProofConfigUpdatedRollupConfigHashIterator is returned from FilterUpdatedRollupConfigHash and is used to iterate over the raw logs and unpacked data for UpdatedRollupConfigHash events raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedRollupConfigHashIterator struct {
	Event *ZKFaultProofConfigUpdatedRollupConfigHash // Event containing the contract specifics and raw log

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
func (it *ZKFaultProofConfigUpdatedRollupConfigHashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKFaultProofConfigUpdatedRollupConfigHash)
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
		it.Event = new(ZKFaultProofConfigUpdatedRollupConfigHash)
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
func (it *ZKFaultProofConfigUpdatedRollupConfigHashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKFaultProofConfigUpdatedRollupConfigHashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKFaultProofConfigUpdatedRollupConfigHash represents a UpdatedRollupConfigHash event raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedRollupConfigHash struct {
	OldRollupConfigHash [32]byte
	NewRollupConfigHash [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterUpdatedRollupConfigHash is a free log retrieval operation binding the contract event 0xda2f5f014ada26cff39a0f2a9dc6fa4fca1581376fc91ec09506c8fb8657bc35.
//
// Solidity: event UpdatedRollupConfigHash(bytes32 indexed oldRollupConfigHash, bytes32 indexed newRollupConfigHash)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) FilterUpdatedRollupConfigHash(opts *bind.FilterOpts, oldRollupConfigHash [][32]byte, newRollupConfigHash [][32]byte) (*ZKFaultProofConfigUpdatedRollupConfigHashIterator, error) {

	var oldRollupConfigHashRule []interface{}
	for _, oldRollupConfigHashItem := range oldRollupConfigHash {
		oldRollupConfigHashRule = append(oldRollupConfigHashRule, oldRollupConfigHashItem)
	}
	var newRollupConfigHashRule []interface{}
	for _, newRollupConfigHashItem := range newRollupConfigHash {
		newRollupConfigHashRule = append(newRollupConfigHashRule, newRollupConfigHashItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.FilterLogs(opts, "UpdatedRollupConfigHash", oldRollupConfigHashRule, newRollupConfigHashRule)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigUpdatedRollupConfigHashIterator{contract: _ZKFaultProofConfig.contract, event: "UpdatedRollupConfigHash", logs: logs, sub: sub}, nil
}

// WatchUpdatedRollupConfigHash is a free log subscription operation binding the contract event 0xda2f5f014ada26cff39a0f2a9dc6fa4fca1581376fc91ec09506c8fb8657bc35.
//
// Solidity: event UpdatedRollupConfigHash(bytes32 indexed oldRollupConfigHash, bytes32 indexed newRollupConfigHash)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) WatchUpdatedRollupConfigHash(opts *bind.WatchOpts, sink chan<- *ZKFaultProofConfigUpdatedRollupConfigHash, oldRollupConfigHash [][32]byte, newRollupConfigHash [][32]byte) (event.Subscription, error) {

	var oldRollupConfigHashRule []interface{}
	for _, oldRollupConfigHashItem := range oldRollupConfigHash {
		oldRollupConfigHashRule = append(oldRollupConfigHashRule, oldRollupConfigHashItem)
	}
	var newRollupConfigHashRule []interface{}
	for _, newRollupConfigHashItem := range newRollupConfigHash {
		newRollupConfigHashRule = append(newRollupConfigHashRule, newRollupConfigHashItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.WatchLogs(opts, "UpdatedRollupConfigHash", oldRollupConfigHashRule, newRollupConfigHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKFaultProofConfigUpdatedRollupConfigHash)
				if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedRollupConfigHash", log); err != nil {
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

// ParseUpdatedRollupConfigHash is a log parse operation binding the contract event 0xda2f5f014ada26cff39a0f2a9dc6fa4fca1581376fc91ec09506c8fb8657bc35.
//
// Solidity: event UpdatedRollupConfigHash(bytes32 indexed oldRollupConfigHash, bytes32 indexed newRollupConfigHash)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) ParseUpdatedRollupConfigHash(log types.Log) (*ZKFaultProofConfigUpdatedRollupConfigHash, error) {
	event := new(ZKFaultProofConfigUpdatedRollupConfigHash)
	if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedRollupConfigHash", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKFaultProofConfigUpdatedVerifierGatewayIterator is returned from FilterUpdatedVerifierGateway and is used to iterate over the raw logs and unpacked data for UpdatedVerifierGateway events raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedVerifierGatewayIterator struct {
	Event *ZKFaultProofConfigUpdatedVerifierGateway // Event containing the contract specifics and raw log

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
func (it *ZKFaultProofConfigUpdatedVerifierGatewayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKFaultProofConfigUpdatedVerifierGateway)
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
		it.Event = new(ZKFaultProofConfigUpdatedVerifierGateway)
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
func (it *ZKFaultProofConfigUpdatedVerifierGatewayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKFaultProofConfigUpdatedVerifierGatewayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKFaultProofConfigUpdatedVerifierGateway represents a UpdatedVerifierGateway event raised by the ZKFaultProofConfig contract.
type ZKFaultProofConfigUpdatedVerifierGateway struct {
	OldVerifierGateway common.Address
	NewVerifierGateway common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpdatedVerifierGateway is a free log retrieval operation binding the contract event 0x1379941631ff0ed9178ab16ab67a2e5db3aeada7f87e518f761e79c8e38377e3.
//
// Solidity: event UpdatedVerifierGateway(address indexed oldVerifierGateway, address indexed newVerifierGateway)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) FilterUpdatedVerifierGateway(opts *bind.FilterOpts, oldVerifierGateway []common.Address, newVerifierGateway []common.Address) (*ZKFaultProofConfigUpdatedVerifierGatewayIterator, error) {

	var oldVerifierGatewayRule []interface{}
	for _, oldVerifierGatewayItem := range oldVerifierGateway {
		oldVerifierGatewayRule = append(oldVerifierGatewayRule, oldVerifierGatewayItem)
	}
	var newVerifierGatewayRule []interface{}
	for _, newVerifierGatewayItem := range newVerifierGateway {
		newVerifierGatewayRule = append(newVerifierGatewayRule, newVerifierGatewayItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.FilterLogs(opts, "UpdatedVerifierGateway", oldVerifierGatewayRule, newVerifierGatewayRule)
	if err != nil {
		return nil, err
	}
	return &ZKFaultProofConfigUpdatedVerifierGatewayIterator{contract: _ZKFaultProofConfig.contract, event: "UpdatedVerifierGateway", logs: logs, sub: sub}, nil
}

// WatchUpdatedVerifierGateway is a free log subscription operation binding the contract event 0x1379941631ff0ed9178ab16ab67a2e5db3aeada7f87e518f761e79c8e38377e3.
//
// Solidity: event UpdatedVerifierGateway(address indexed oldVerifierGateway, address indexed newVerifierGateway)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) WatchUpdatedVerifierGateway(opts *bind.WatchOpts, sink chan<- *ZKFaultProofConfigUpdatedVerifierGateway, oldVerifierGateway []common.Address, newVerifierGateway []common.Address) (event.Subscription, error) {

	var oldVerifierGatewayRule []interface{}
	for _, oldVerifierGatewayItem := range oldVerifierGateway {
		oldVerifierGatewayRule = append(oldVerifierGatewayRule, oldVerifierGatewayItem)
	}
	var newVerifierGatewayRule []interface{}
	for _, newVerifierGatewayItem := range newVerifierGateway {
		newVerifierGatewayRule = append(newVerifierGatewayRule, newVerifierGatewayItem)
	}

	logs, sub, err := _ZKFaultProofConfig.contract.WatchLogs(opts, "UpdatedVerifierGateway", oldVerifierGatewayRule, newVerifierGatewayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKFaultProofConfigUpdatedVerifierGateway)
				if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedVerifierGateway", log); err != nil {
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

// ParseUpdatedVerifierGateway is a log parse operation binding the contract event 0x1379941631ff0ed9178ab16ab67a2e5db3aeada7f87e518f761e79c8e38377e3.
//
// Solidity: event UpdatedVerifierGateway(address indexed oldVerifierGateway, address indexed newVerifierGateway)
func (_ZKFaultProofConfig *ZKFaultProofConfigFilterer) ParseUpdatedVerifierGateway(log types.Log) (*ZKFaultProofConfigUpdatedVerifierGateway, error) {
	event := new(ZKFaultProofConfigUpdatedVerifierGateway)
	if err := _ZKFaultProofConfig.contract.UnpackLog(event, "UpdatedVerifierGateway", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
