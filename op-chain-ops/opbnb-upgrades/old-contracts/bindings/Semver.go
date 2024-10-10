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

// SemverMetaData contains all meta data concerning the Semver contract.
var SemverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_major\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minor\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_patch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"}]",
}

// SemverABI is the input ABI used to generate the binding from.
// Deprecated: Use SemverMetaData.ABI instead.
var SemverABI = SemverMetaData.ABI

// Semver is an auto generated Go binding around an Ethereum contract.
type Semver struct {
	SemverCaller     // Read-only binding to the contract
	SemverTransactor // Write-only binding to the contract
	SemverFilterer   // Log filterer for contract events
}

// SemverCaller is an auto generated read-only Go binding around an Ethereum contract.
type SemverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SemverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SemverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SemverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SemverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SemverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SemverSession struct {
	Contract     *Semver           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SemverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SemverCallerSession struct {
	Contract *SemverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SemverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SemverTransactorSession struct {
	Contract     *SemverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SemverRaw is an auto generated low-level Go binding around an Ethereum contract.
type SemverRaw struct {
	Contract *Semver // Generic contract binding to access the raw methods on
}

// SemverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SemverCallerRaw struct {
	Contract *SemverCaller // Generic read-only contract binding to access the raw methods on
}

// SemverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SemverTransactorRaw struct {
	Contract *SemverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSemver creates a new instance of Semver, bound to a specific deployed contract.
func NewSemver(address common.Address, backend bind.ContractBackend) (*Semver, error) {
	contract, err := bindSemver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Semver{SemverCaller: SemverCaller{contract: contract}, SemverTransactor: SemverTransactor{contract: contract}, SemverFilterer: SemverFilterer{contract: contract}}, nil
}

// NewSemverCaller creates a new read-only instance of Semver, bound to a specific deployed contract.
func NewSemverCaller(address common.Address, caller bind.ContractCaller) (*SemverCaller, error) {
	contract, err := bindSemver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SemverCaller{contract: contract}, nil
}

// NewSemverTransactor creates a new write-only instance of Semver, bound to a specific deployed contract.
func NewSemverTransactor(address common.Address, transactor bind.ContractTransactor) (*SemverTransactor, error) {
	contract, err := bindSemver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SemverTransactor{contract: contract}, nil
}

// NewSemverFilterer creates a new log filterer instance of Semver, bound to a specific deployed contract.
func NewSemverFilterer(address common.Address, filterer bind.ContractFilterer) (*SemverFilterer, error) {
	contract, err := bindSemver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SemverFilterer{contract: contract}, nil
}

// bindSemver binds a generic wrapper to an already deployed contract.
func bindSemver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SemverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Semver *SemverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Semver.Contract.SemverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Semver *SemverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Semver.Contract.SemverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Semver *SemverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Semver.Contract.SemverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Semver *SemverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Semver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Semver *SemverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Semver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Semver *SemverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Semver.Contract.contract.Transact(opts, method, params...)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Semver *SemverCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Semver.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Semver *SemverSession) Version() (string, error) {
	return _Semver.Contract.Version(&_Semver.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Semver *SemverCallerSession) Version() (string, error) {
	return _Semver.Contract.Version(&_Semver.CallOpts)
}
