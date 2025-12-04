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

// AnchorStateRegistryStartingAnchorRoot is an auto generated low-level Go binding around an user-defined struct.
type AnchorStateRegistryStartingAnchorRoot struct {
	GameType   uint32
	OutputRoot OutputRoot
}

// OutputRoot is an auto generated low-level Go binding around an user-defined struct.
type OutputRoot struct {
	Root          [32]byte
	L2BlockNumber *big.Int
}

// AnchorStateRegistryMetaData contains all meta data concerning the AnchorStateRegistry contract.
var AnchorStateRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_disputeGameFactory\",\"type\":\"address\",\"internalType\":\"contractIDisputeGameFactory\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchors\",\"inputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"GameType\"}],\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"Hash\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"disputeGameFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIDisputeGameFactory\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_startingAnchorRoots\",\"type\":\"tuple[]\",\"internalType\":\"structAnchorStateRegistry.StartingAnchorRoot[]\",\"components\":[{\"name\":\"gameType\",\"type\":\"uint32\",\"internalType\":\"GameType\"},{\"name\":\"outputRoot\",\"type\":\"tuple\",\"internalType\":\"structOutputRoot\",\"components\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"Hash\"},{\"name\":\"l2BlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tryUpdateAnchorState\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162000fc938038062000fc983398101604081905262000034916200026b565b6001600160a01b0381166080526040805160008082526020820190925262000081916200007a565b620000666200022d565b8152602001906001900390816200005c5790505b5062000088565b50620002db565b600054610100900460ff1615808015620000a95750600054600160ff909116105b80620000d95750620000c6306200021e60201b620007c61760201c565b158015620000d9575060005460ff166001145b620001415760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b606482015260840160405180910390fd5b6000805460ff19166001179055801562000165576000805461ff0019166101001790555b60005b8251811015620001d25760008382815181106200018957620001896200029d565b60209081029190910181015180820151905163ffffffff166000908152600180845260409091208251815591909201519101555080620001c981620002b3565b91505062000168565b5080156200021a576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b6001600160a01b03163b151590565b6040518060400160405280600063ffffffff16815260200162000266604051806040016040528060008019168152602001600081525090565b905290565b6000602082840312156200027e57600080fd5b81516001600160a01b03811681146200029657600080fd5b9392505050565b634e487b7160e01b600052603260045260246000fd5b600060018201620002d457634e487b7160e01b600052601160045260246000fd5b5060010190565b608051610ccb620002fe6000396000818161013101526102000152610ccb6000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c8063838c2d1e11610050578063838c2d1e146100fa578063c303f0df14610104578063f2b4e6171461011757600080fd5b806354fd4d501461006c5780637258a807146100be575b600080fd5b6100a86040518060400160405280600581526020017f312e302e3000000000000000000000000000000000000000000000000000000081525081565b6040516100b5919061085c565b60405180910390f35b6100e56100cc36600461088b565b6001602081905260009182526040909120805491015482565b604080519283526020830191909152016100b5565b61010261015b565b005b61010261011236600461094f565b6105d4565b60405173ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001681526020016100b5565b600033905060008060008373ffffffffffffffffffffffffffffffffffffffff1663fa24f7436040518163ffffffff1660e01b8152600401600060405180830381865afa1580156101b0573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526101f69190810190610a68565b92509250925060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635f0150cb8585856040518463ffffffff1660e01b815260040161025b93929190610b39565b6040805180830381865afa158015610277573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061029b9190610b67565b5090508473ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610384576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152604360248201527f416e63686f72537461746552656769737472793a206661756c7420646973707560448201527f74652067616d65206e6f7420726567697374657265642077697468206661637460648201527f6f72790000000000000000000000000000000000000000000000000000000000608482015260a4015b60405180910390fd5b600160008563ffffffff1663ffffffff168152602001908152602001600020600101548573ffffffffffffffffffffffffffffffffffffffff16638b85902b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103f2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104169190610bc7565b11610422575050505050565b60028573ffffffffffffffffffffffffffffffffffffffff1663200d2ed26040518163ffffffff1660e01b8152600401602060405180830381865afa15801561046f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104939190610c0f565b60028111156104a4576104a4610be0565b146104b0575050505050565b60405180604001604052806105308773ffffffffffffffffffffffffffffffffffffffff1663bcef3b556040518163ffffffff1660e01b8152600401602060405180830381865afa158015610509573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061052d9190610bc7565b90565b81526020018673ffffffffffffffffffffffffffffffffffffffff16638b85902b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610580573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a49190610bc7565b905263ffffffff909416600090815260016020818152604090922086518155959091015194019390935550505050565b600054610100900460ff16158080156105f45750600054600160ff909116105b8061060e5750303b15801561060e575060005460ff166001145b61069a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840161037b565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156106f857600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b60005b825181101561075e57600083828151811061071857610718610c30565b60209081029190910181015180820151905163ffffffff16600090815260018084526040909120825181559190920151910155508061075681610c5f565b9150506106fb565b5080156107c257600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b73ffffffffffffffffffffffffffffffffffffffff163b151590565b60005b838110156107fd5781810151838201526020016107e5565b8381111561080c576000848401525b50505050565b6000815180845261082a8160208601602086016107e2565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60208152600061086f6020830184610812565b9392505050565b63ffffffff8116811461088857600080fd5b50565b60006020828403121561089d57600080fd5b813561086f81610876565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156108fa576108fa6108a8565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610947576109476108a8565b604052919050565b6000602080838503121561096257600080fd5b823567ffffffffffffffff8082111561097a57600080fd5b818501915085601f83011261098e57600080fd5b8135818111156109a0576109a06108a8565b6109ae848260051b01610900565b818152848101925060609182028401850191888311156109cd57600080fd5b938501935b82851015610a5c57848903818112156109eb5760008081fd5b6109f36108d7565b86356109fe81610876565b815260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08301811315610a325760008081fd5b610a3a6108d7565b888a0135815290880135898201528189015285525093840193928501926109d2565b50979650505050505050565b600080600060608486031215610a7d57600080fd5b8351610a8881610876565b60208501516040860151919450925067ffffffffffffffff80821115610aad57600080fd5b818601915086601f830112610ac157600080fd5b815181811115610ad357610ad36108a8565b610b0460207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610900565b9150808252876020828501011115610b1b57600080fd5b610b2c8160208401602086016107e2565b5080925050509250925092565b63ffffffff84168152826020820152606060408201526000610b5e6060830184610812565b95945050505050565b60008060408385031215610b7a57600080fd5b825173ffffffffffffffffffffffffffffffffffffffff81168114610b9e57600080fd5b602084015190925067ffffffffffffffff81168114610bbc57600080fd5b809150509250929050565b600060208284031215610bd957600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600060208284031215610c2157600080fd5b81516003811061086f57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610cb7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b506001019056fea164736f6c634300080f000a",
}

// AnchorStateRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AnchorStateRegistryMetaData.ABI instead.
var AnchorStateRegistryABI = AnchorStateRegistryMetaData.ABI

// AnchorStateRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AnchorStateRegistryMetaData.Bin instead.
var AnchorStateRegistryBin = AnchorStateRegistryMetaData.Bin

// DeployAnchorStateRegistry deploys a new Ethereum contract, binding an instance of AnchorStateRegistry to it.
func DeployAnchorStateRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _disputeGameFactory common.Address) (common.Address, *types.Transaction, *AnchorStateRegistry, error) {
	parsed, err := AnchorStateRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AnchorStateRegistryBin), backend, _disputeGameFactory)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AnchorStateRegistry{AnchorStateRegistryCaller: AnchorStateRegistryCaller{contract: contract}, AnchorStateRegistryTransactor: AnchorStateRegistryTransactor{contract: contract}, AnchorStateRegistryFilterer: AnchorStateRegistryFilterer{contract: contract}}, nil
}

// AnchorStateRegistry is an auto generated Go binding around an Ethereum contract.
type AnchorStateRegistry struct {
	AnchorStateRegistryCaller     // Read-only binding to the contract
	AnchorStateRegistryTransactor // Write-only binding to the contract
	AnchorStateRegistryFilterer   // Log filterer for contract events
}

// AnchorStateRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AnchorStateRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AnchorStateRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AnchorStateRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AnchorStateRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AnchorStateRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AnchorStateRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AnchorStateRegistrySession struct {
	Contract     *AnchorStateRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AnchorStateRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AnchorStateRegistryCallerSession struct {
	Contract *AnchorStateRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// AnchorStateRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AnchorStateRegistryTransactorSession struct {
	Contract     *AnchorStateRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// AnchorStateRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AnchorStateRegistryRaw struct {
	Contract *AnchorStateRegistry // Generic contract binding to access the raw methods on
}

// AnchorStateRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AnchorStateRegistryCallerRaw struct {
	Contract *AnchorStateRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AnchorStateRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AnchorStateRegistryTransactorRaw struct {
	Contract *AnchorStateRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAnchorStateRegistry creates a new instance of AnchorStateRegistry, bound to a specific deployed contract.
func NewAnchorStateRegistry(address common.Address, backend bind.ContractBackend) (*AnchorStateRegistry, error) {
	contract, err := bindAnchorStateRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AnchorStateRegistry{AnchorStateRegistryCaller: AnchorStateRegistryCaller{contract: contract}, AnchorStateRegistryTransactor: AnchorStateRegistryTransactor{contract: contract}, AnchorStateRegistryFilterer: AnchorStateRegistryFilterer{contract: contract}}, nil
}

// NewAnchorStateRegistryCaller creates a new read-only instance of AnchorStateRegistry, bound to a specific deployed contract.
func NewAnchorStateRegistryCaller(address common.Address, caller bind.ContractCaller) (*AnchorStateRegistryCaller, error) {
	contract, err := bindAnchorStateRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AnchorStateRegistryCaller{contract: contract}, nil
}

// NewAnchorStateRegistryTransactor creates a new write-only instance of AnchorStateRegistry, bound to a specific deployed contract.
func NewAnchorStateRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AnchorStateRegistryTransactor, error) {
	contract, err := bindAnchorStateRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AnchorStateRegistryTransactor{contract: contract}, nil
}

// NewAnchorStateRegistryFilterer creates a new log filterer instance of AnchorStateRegistry, bound to a specific deployed contract.
func NewAnchorStateRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AnchorStateRegistryFilterer, error) {
	contract, err := bindAnchorStateRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AnchorStateRegistryFilterer{contract: contract}, nil
}

// bindAnchorStateRegistry binds a generic wrapper to an already deployed contract.
func bindAnchorStateRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AnchorStateRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AnchorStateRegistry *AnchorStateRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AnchorStateRegistry.Contract.AnchorStateRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AnchorStateRegistry *AnchorStateRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.AnchorStateRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AnchorStateRegistry *AnchorStateRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.AnchorStateRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AnchorStateRegistry *AnchorStateRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AnchorStateRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AnchorStateRegistry *AnchorStateRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AnchorStateRegistry *AnchorStateRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.contract.Transact(opts, method, params...)
}

// Anchors is a free data retrieval call binding the contract method 0x7258a807.
//
// Solidity: function anchors(uint32 ) view returns(bytes32 root, uint256 l2BlockNumber)
func (_AnchorStateRegistry *AnchorStateRegistryCaller) Anchors(opts *bind.CallOpts, arg0 uint32) (struct {
	Root          [32]byte
	L2BlockNumber *big.Int
}, error) {
	var out []interface{}
	err := _AnchorStateRegistry.contract.Call(opts, &out, "anchors", arg0)

	outstruct := new(struct {
		Root          [32]byte
		L2BlockNumber *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Root = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.L2BlockNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Anchors is a free data retrieval call binding the contract method 0x7258a807.
//
// Solidity: function anchors(uint32 ) view returns(bytes32 root, uint256 l2BlockNumber)
func (_AnchorStateRegistry *AnchorStateRegistrySession) Anchors(arg0 uint32) (struct {
	Root          [32]byte
	L2BlockNumber *big.Int
}, error) {
	return _AnchorStateRegistry.Contract.Anchors(&_AnchorStateRegistry.CallOpts, arg0)
}

// Anchors is a free data retrieval call binding the contract method 0x7258a807.
//
// Solidity: function anchors(uint32 ) view returns(bytes32 root, uint256 l2BlockNumber)
func (_AnchorStateRegistry *AnchorStateRegistryCallerSession) Anchors(arg0 uint32) (struct {
	Root          [32]byte
	L2BlockNumber *big.Int
}, error) {
	return _AnchorStateRegistry.Contract.Anchors(&_AnchorStateRegistry.CallOpts, arg0)
}

// DisputeGameFactory is a free data retrieval call binding the contract method 0xf2b4e617.
//
// Solidity: function disputeGameFactory() view returns(address)
func (_AnchorStateRegistry *AnchorStateRegistryCaller) DisputeGameFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AnchorStateRegistry.contract.Call(opts, &out, "disputeGameFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DisputeGameFactory is a free data retrieval call binding the contract method 0xf2b4e617.
//
// Solidity: function disputeGameFactory() view returns(address)
func (_AnchorStateRegistry *AnchorStateRegistrySession) DisputeGameFactory() (common.Address, error) {
	return _AnchorStateRegistry.Contract.DisputeGameFactory(&_AnchorStateRegistry.CallOpts)
}

// DisputeGameFactory is a free data retrieval call binding the contract method 0xf2b4e617.
//
// Solidity: function disputeGameFactory() view returns(address)
func (_AnchorStateRegistry *AnchorStateRegistryCallerSession) DisputeGameFactory() (common.Address, error) {
	return _AnchorStateRegistry.Contract.DisputeGameFactory(&_AnchorStateRegistry.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_AnchorStateRegistry *AnchorStateRegistryCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AnchorStateRegistry.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_AnchorStateRegistry *AnchorStateRegistrySession) Version() (string, error) {
	return _AnchorStateRegistry.Contract.Version(&_AnchorStateRegistry.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_AnchorStateRegistry *AnchorStateRegistryCallerSession) Version() (string, error) {
	return _AnchorStateRegistry.Contract.Version(&_AnchorStateRegistry.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc303f0df.
//
// Solidity: function initialize((uint32,(bytes32,uint256))[] _startingAnchorRoots) returns()
func (_AnchorStateRegistry *AnchorStateRegistryTransactor) Initialize(opts *bind.TransactOpts, _startingAnchorRoots []AnchorStateRegistryStartingAnchorRoot) (*types.Transaction, error) {
	return _AnchorStateRegistry.contract.Transact(opts, "initialize", _startingAnchorRoots)
}

// Initialize is a paid mutator transaction binding the contract method 0xc303f0df.
//
// Solidity: function initialize((uint32,(bytes32,uint256))[] _startingAnchorRoots) returns()
func (_AnchorStateRegistry *AnchorStateRegistrySession) Initialize(_startingAnchorRoots []AnchorStateRegistryStartingAnchorRoot) (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.Initialize(&_AnchorStateRegistry.TransactOpts, _startingAnchorRoots)
}

// Initialize is a paid mutator transaction binding the contract method 0xc303f0df.
//
// Solidity: function initialize((uint32,(bytes32,uint256))[] _startingAnchorRoots) returns()
func (_AnchorStateRegistry *AnchorStateRegistryTransactorSession) Initialize(_startingAnchorRoots []AnchorStateRegistryStartingAnchorRoot) (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.Initialize(&_AnchorStateRegistry.TransactOpts, _startingAnchorRoots)
}

// TryUpdateAnchorState is a paid mutator transaction binding the contract method 0x838c2d1e.
//
// Solidity: function tryUpdateAnchorState() returns()
func (_AnchorStateRegistry *AnchorStateRegistryTransactor) TryUpdateAnchorState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AnchorStateRegistry.contract.Transact(opts, "tryUpdateAnchorState")
}

// TryUpdateAnchorState is a paid mutator transaction binding the contract method 0x838c2d1e.
//
// Solidity: function tryUpdateAnchorState() returns()
func (_AnchorStateRegistry *AnchorStateRegistrySession) TryUpdateAnchorState() (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.TryUpdateAnchorState(&_AnchorStateRegistry.TransactOpts)
}

// TryUpdateAnchorState is a paid mutator transaction binding the contract method 0x838c2d1e.
//
// Solidity: function tryUpdateAnchorState() returns()
func (_AnchorStateRegistry *AnchorStateRegistryTransactorSession) TryUpdateAnchorState() (*types.Transaction, error) {
	return _AnchorStateRegistry.Contract.TryUpdateAnchorState(&_AnchorStateRegistry.TransactOpts)
}

// AnchorStateRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AnchorStateRegistry contract.
type AnchorStateRegistryInitializedIterator struct {
	Event *AnchorStateRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *AnchorStateRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AnchorStateRegistryInitialized)
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
		it.Event = new(AnchorStateRegistryInitialized)
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
func (it *AnchorStateRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AnchorStateRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AnchorStateRegistryInitialized represents a Initialized event raised by the AnchorStateRegistry contract.
type AnchorStateRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_AnchorStateRegistry *AnchorStateRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*AnchorStateRegistryInitializedIterator, error) {

	logs, sub, err := _AnchorStateRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AnchorStateRegistryInitializedIterator{contract: _AnchorStateRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_AnchorStateRegistry *AnchorStateRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AnchorStateRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _AnchorStateRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AnchorStateRegistryInitialized)
				if err := _AnchorStateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AnchorStateRegistry *AnchorStateRegistryFilterer) ParseInitialized(log types.Log) (*AnchorStateRegistryInitialized, error) {
	event := new(AnchorStateRegistryInitialized)
	if err := _AnchorStateRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
