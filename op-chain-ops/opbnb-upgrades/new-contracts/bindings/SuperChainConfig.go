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

// SuperChainConfigMetaData contains all meta data concerning the SuperChainConfig contract.
var SuperChainConfigMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"GUARDIAN_SLOT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAUSED_SLOT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_guardian\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_paused\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"_identifier\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"paused_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ConfigUpdate\",\"inputs\":[{\"name\":\"updateType\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumSuperchainConfig.UpdateType\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"identifier\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[],\"anonymous\":false}]",
}

// SuperChainConfigABI is the input ABI used to generate the binding from.
// Deprecated: Use SuperChainConfigMetaData.ABI instead.
var SuperChainConfigABI = SuperChainConfigMetaData.ABI

// SuperChainConfig is an auto generated Go binding around an Ethereum contract.
type SuperChainConfig struct {
	SuperChainConfigCaller     // Read-only binding to the contract
	SuperChainConfigTransactor // Write-only binding to the contract
	SuperChainConfigFilterer   // Log filterer for contract events
}

// SuperChainConfigCaller is an auto generated read-only Go binding around an Ethereum contract.
type SuperChainConfigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SuperChainConfigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SuperChainConfigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SuperChainConfigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SuperChainConfigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SuperChainConfigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SuperChainConfigSession struct {
	Contract     *SuperChainConfig // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SuperChainConfigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SuperChainConfigCallerSession struct {
	Contract *SuperChainConfigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SuperChainConfigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SuperChainConfigTransactorSession struct {
	Contract     *SuperChainConfigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SuperChainConfigRaw is an auto generated low-level Go binding around an Ethereum contract.
type SuperChainConfigRaw struct {
	Contract *SuperChainConfig // Generic contract binding to access the raw methods on
}

// SuperChainConfigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SuperChainConfigCallerRaw struct {
	Contract *SuperChainConfigCaller // Generic read-only contract binding to access the raw methods on
}

// SuperChainConfigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SuperChainConfigTransactorRaw struct {
	Contract *SuperChainConfigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSuperChainConfig creates a new instance of SuperChainConfig, bound to a specific deployed contract.
func NewSuperChainConfig(address common.Address, backend bind.ContractBackend) (*SuperChainConfig, error) {
	contract, err := bindSuperChainConfig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SuperChainConfig{SuperChainConfigCaller: SuperChainConfigCaller{contract: contract}, SuperChainConfigTransactor: SuperChainConfigTransactor{contract: contract}, SuperChainConfigFilterer: SuperChainConfigFilterer{contract: contract}}, nil
}

// NewSuperChainConfigCaller creates a new read-only instance of SuperChainConfig, bound to a specific deployed contract.
func NewSuperChainConfigCaller(address common.Address, caller bind.ContractCaller) (*SuperChainConfigCaller, error) {
	contract, err := bindSuperChainConfig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SuperChainConfigCaller{contract: contract}, nil
}

// NewSuperChainConfigTransactor creates a new write-only instance of SuperChainConfig, bound to a specific deployed contract.
func NewSuperChainConfigTransactor(address common.Address, transactor bind.ContractTransactor) (*SuperChainConfigTransactor, error) {
	contract, err := bindSuperChainConfig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SuperChainConfigTransactor{contract: contract}, nil
}

// NewSuperChainConfigFilterer creates a new log filterer instance of SuperChainConfig, bound to a specific deployed contract.
func NewSuperChainConfigFilterer(address common.Address, filterer bind.ContractFilterer) (*SuperChainConfigFilterer, error) {
	contract, err := bindSuperChainConfig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SuperChainConfigFilterer{contract: contract}, nil
}

// bindSuperChainConfig binds a generic wrapper to an already deployed contract.
func bindSuperChainConfig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SuperChainConfigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SuperChainConfig *SuperChainConfigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SuperChainConfig.Contract.SuperChainConfigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SuperChainConfig *SuperChainConfigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.SuperChainConfigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SuperChainConfig *SuperChainConfigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.SuperChainConfigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SuperChainConfig *SuperChainConfigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SuperChainConfig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SuperChainConfig *SuperChainConfigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SuperChainConfig *SuperChainConfigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.contract.Transact(opts, method, params...)
}

// GUARDIANSLOT is a free data retrieval call binding the contract method 0xc23a451a.
//
// Solidity: function GUARDIAN_SLOT() view returns(bytes32)
func (_SuperChainConfig *SuperChainConfigCaller) GUARDIANSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SuperChainConfig.contract.Call(opts, &out, "GUARDIAN_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GUARDIANSLOT is a free data retrieval call binding the contract method 0xc23a451a.
//
// Solidity: function GUARDIAN_SLOT() view returns(bytes32)
func (_SuperChainConfig *SuperChainConfigSession) GUARDIANSLOT() ([32]byte, error) {
	return _SuperChainConfig.Contract.GUARDIANSLOT(&_SuperChainConfig.CallOpts)
}

// GUARDIANSLOT is a free data retrieval call binding the contract method 0xc23a451a.
//
// Solidity: function GUARDIAN_SLOT() view returns(bytes32)
func (_SuperChainConfig *SuperChainConfigCallerSession) GUARDIANSLOT() ([32]byte, error) {
	return _SuperChainConfig.Contract.GUARDIANSLOT(&_SuperChainConfig.CallOpts)
}

// PAUSEDSLOT is a free data retrieval call binding the contract method 0x7fbf7b6a.
//
// Solidity: function PAUSED_SLOT() view returns(bytes32)
func (_SuperChainConfig *SuperChainConfigCaller) PAUSEDSLOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SuperChainConfig.contract.Call(opts, &out, "PAUSED_SLOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSEDSLOT is a free data retrieval call binding the contract method 0x7fbf7b6a.
//
// Solidity: function PAUSED_SLOT() view returns(bytes32)
func (_SuperChainConfig *SuperChainConfigSession) PAUSEDSLOT() ([32]byte, error) {
	return _SuperChainConfig.Contract.PAUSEDSLOT(&_SuperChainConfig.CallOpts)
}

// PAUSEDSLOT is a free data retrieval call binding the contract method 0x7fbf7b6a.
//
// Solidity: function PAUSED_SLOT() view returns(bytes32)
func (_SuperChainConfig *SuperChainConfigCallerSession) PAUSEDSLOT() ([32]byte, error) {
	return _SuperChainConfig.Contract.PAUSEDSLOT(&_SuperChainConfig.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address guardian_)
func (_SuperChainConfig *SuperChainConfigCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SuperChainConfig.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address guardian_)
func (_SuperChainConfig *SuperChainConfigSession) Guardian() (common.Address, error) {
	return _SuperChainConfig.Contract.Guardian(&_SuperChainConfig.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address guardian_)
func (_SuperChainConfig *SuperChainConfigCallerSession) Guardian() (common.Address, error) {
	return _SuperChainConfig.Contract.Guardian(&_SuperChainConfig.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_SuperChainConfig *SuperChainConfigCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SuperChainConfig.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_SuperChainConfig *SuperChainConfigSession) Paused() (bool, error) {
	return _SuperChainConfig.Contract.Paused(&_SuperChainConfig.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool paused_)
func (_SuperChainConfig *SuperChainConfigCallerSession) Paused() (bool, error) {
	return _SuperChainConfig.Contract.Paused(&_SuperChainConfig.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_SuperChainConfig *SuperChainConfigCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SuperChainConfig.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_SuperChainConfig *SuperChainConfigSession) Version() (string, error) {
	return _SuperChainConfig.Contract.Version(&_SuperChainConfig.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_SuperChainConfig *SuperChainConfigCallerSession) Version() (string, error) {
	return _SuperChainConfig.Contract.Version(&_SuperChainConfig.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x400ada75.
//
// Solidity: function initialize(address _guardian, bool _paused) returns()
func (_SuperChainConfig *SuperChainConfigTransactor) Initialize(opts *bind.TransactOpts, _guardian common.Address, _paused bool) (*types.Transaction, error) {
	return _SuperChainConfig.contract.Transact(opts, "initialize", _guardian, _paused)
}

// Initialize is a paid mutator transaction binding the contract method 0x400ada75.
//
// Solidity: function initialize(address _guardian, bool _paused) returns()
func (_SuperChainConfig *SuperChainConfigSession) Initialize(_guardian common.Address, _paused bool) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.Initialize(&_SuperChainConfig.TransactOpts, _guardian, _paused)
}

// Initialize is a paid mutator transaction binding the contract method 0x400ada75.
//
// Solidity: function initialize(address _guardian, bool _paused) returns()
func (_SuperChainConfig *SuperChainConfigTransactorSession) Initialize(_guardian common.Address, _paused bool) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.Initialize(&_SuperChainConfig.TransactOpts, _guardian, _paused)
}

// Pause is a paid mutator transaction binding the contract method 0x6da66355.
//
// Solidity: function pause(string _identifier) returns()
func (_SuperChainConfig *SuperChainConfigTransactor) Pause(opts *bind.TransactOpts, _identifier string) (*types.Transaction, error) {
	return _SuperChainConfig.contract.Transact(opts, "pause", _identifier)
}

// Pause is a paid mutator transaction binding the contract method 0x6da66355.
//
// Solidity: function pause(string _identifier) returns()
func (_SuperChainConfig *SuperChainConfigSession) Pause(_identifier string) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.Pause(&_SuperChainConfig.TransactOpts, _identifier)
}

// Pause is a paid mutator transaction binding the contract method 0x6da66355.
//
// Solidity: function pause(string _identifier) returns()
func (_SuperChainConfig *SuperChainConfigTransactorSession) Pause(_identifier string) (*types.Transaction, error) {
	return _SuperChainConfig.Contract.Pause(&_SuperChainConfig.TransactOpts, _identifier)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SuperChainConfig *SuperChainConfigTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SuperChainConfig.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SuperChainConfig *SuperChainConfigSession) Unpause() (*types.Transaction, error) {
	return _SuperChainConfig.Contract.Unpause(&_SuperChainConfig.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SuperChainConfig *SuperChainConfigTransactorSession) Unpause() (*types.Transaction, error) {
	return _SuperChainConfig.Contract.Unpause(&_SuperChainConfig.TransactOpts)
}

// SuperChainConfigConfigUpdateIterator is returned from FilterConfigUpdate and is used to iterate over the raw logs and unpacked data for ConfigUpdate events raised by the SuperChainConfig contract.
type SuperChainConfigConfigUpdateIterator struct {
	Event *SuperChainConfigConfigUpdate // Event containing the contract specifics and raw log

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
func (it *SuperChainConfigConfigUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SuperChainConfigConfigUpdate)
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
		it.Event = new(SuperChainConfigConfigUpdate)
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
func (it *SuperChainConfigConfigUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SuperChainConfigConfigUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SuperChainConfigConfigUpdate represents a ConfigUpdate event raised by the SuperChainConfig contract.
type SuperChainConfigConfigUpdate struct {
	UpdateType uint8
	Data       []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdate is a free log retrieval operation binding the contract event 0x7b743789cff01dafdeae47739925425aab5dfd02d0c8229e4a508bcd2b9f42bb.
//
// Solidity: event ConfigUpdate(uint8 indexed updateType, bytes data)
func (_SuperChainConfig *SuperChainConfigFilterer) FilterConfigUpdate(opts *bind.FilterOpts, updateType []uint8) (*SuperChainConfigConfigUpdateIterator, error) {

	var updateTypeRule []interface{}
	for _, updateTypeItem := range updateType {
		updateTypeRule = append(updateTypeRule, updateTypeItem)
	}

	logs, sub, err := _SuperChainConfig.contract.FilterLogs(opts, "ConfigUpdate", updateTypeRule)
	if err != nil {
		return nil, err
	}
	return &SuperChainConfigConfigUpdateIterator{contract: _SuperChainConfig.contract, event: "ConfigUpdate", logs: logs, sub: sub}, nil
}

// WatchConfigUpdate is a free log subscription operation binding the contract event 0x7b743789cff01dafdeae47739925425aab5dfd02d0c8229e4a508bcd2b9f42bb.
//
// Solidity: event ConfigUpdate(uint8 indexed updateType, bytes data)
func (_SuperChainConfig *SuperChainConfigFilterer) WatchConfigUpdate(opts *bind.WatchOpts, sink chan<- *SuperChainConfigConfigUpdate, updateType []uint8) (event.Subscription, error) {

	var updateTypeRule []interface{}
	for _, updateTypeItem := range updateType {
		updateTypeRule = append(updateTypeRule, updateTypeItem)
	}

	logs, sub, err := _SuperChainConfig.contract.WatchLogs(opts, "ConfigUpdate", updateTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SuperChainConfigConfigUpdate)
				if err := _SuperChainConfig.contract.UnpackLog(event, "ConfigUpdate", log); err != nil {
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

// ParseConfigUpdate is a log parse operation binding the contract event 0x7b743789cff01dafdeae47739925425aab5dfd02d0c8229e4a508bcd2b9f42bb.
//
// Solidity: event ConfigUpdate(uint8 indexed updateType, bytes data)
func (_SuperChainConfig *SuperChainConfigFilterer) ParseConfigUpdate(log types.Log) (*SuperChainConfigConfigUpdate, error) {
	event := new(SuperChainConfigConfigUpdate)
	if err := _SuperChainConfig.contract.UnpackLog(event, "ConfigUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SuperChainConfigInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SuperChainConfig contract.
type SuperChainConfigInitializedIterator struct {
	Event *SuperChainConfigInitialized // Event containing the contract specifics and raw log

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
func (it *SuperChainConfigInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SuperChainConfigInitialized)
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
		it.Event = new(SuperChainConfigInitialized)
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
func (it *SuperChainConfigInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SuperChainConfigInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SuperChainConfigInitialized represents a Initialized event raised by the SuperChainConfig contract.
type SuperChainConfigInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SuperChainConfig *SuperChainConfigFilterer) FilterInitialized(opts *bind.FilterOpts) (*SuperChainConfigInitializedIterator, error) {

	logs, sub, err := _SuperChainConfig.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SuperChainConfigInitializedIterator{contract: _SuperChainConfig.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SuperChainConfig *SuperChainConfigFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SuperChainConfigInitialized) (event.Subscription, error) {

	logs, sub, err := _SuperChainConfig.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SuperChainConfigInitialized)
				if err := _SuperChainConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_SuperChainConfig *SuperChainConfigFilterer) ParseInitialized(log types.Log) (*SuperChainConfigInitialized, error) {
	event := new(SuperChainConfigInitialized)
	if err := _SuperChainConfig.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SuperChainConfigPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SuperChainConfig contract.
type SuperChainConfigPausedIterator struct {
	Event *SuperChainConfigPaused // Event containing the contract specifics and raw log

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
func (it *SuperChainConfigPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SuperChainConfigPaused)
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
		it.Event = new(SuperChainConfigPaused)
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
func (it *SuperChainConfigPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SuperChainConfigPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SuperChainConfigPaused represents a Paused event raised by the SuperChainConfig contract.
type SuperChainConfigPaused struct {
	Identifier string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0xc32e6d5d6d1de257f64eac19ddb1f700ba13527983849c9486b1ab007ea28381.
//
// Solidity: event Paused(string identifier)
func (_SuperChainConfig *SuperChainConfigFilterer) FilterPaused(opts *bind.FilterOpts) (*SuperChainConfigPausedIterator, error) {

	logs, sub, err := _SuperChainConfig.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SuperChainConfigPausedIterator{contract: _SuperChainConfig.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0xc32e6d5d6d1de257f64eac19ddb1f700ba13527983849c9486b1ab007ea28381.
//
// Solidity: event Paused(string identifier)
func (_SuperChainConfig *SuperChainConfigFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SuperChainConfigPaused) (event.Subscription, error) {

	logs, sub, err := _SuperChainConfig.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SuperChainConfigPaused)
				if err := _SuperChainConfig.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0xc32e6d5d6d1de257f64eac19ddb1f700ba13527983849c9486b1ab007ea28381.
//
// Solidity: event Paused(string identifier)
func (_SuperChainConfig *SuperChainConfigFilterer) ParsePaused(log types.Log) (*SuperChainConfigPaused, error) {
	event := new(SuperChainConfigPaused)
	if err := _SuperChainConfig.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SuperChainConfigUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SuperChainConfig contract.
type SuperChainConfigUnpausedIterator struct {
	Event *SuperChainConfigUnpaused // Event containing the contract specifics and raw log

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
func (it *SuperChainConfigUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SuperChainConfigUnpaused)
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
		it.Event = new(SuperChainConfigUnpaused)
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
func (it *SuperChainConfigUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SuperChainConfigUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SuperChainConfigUnpaused represents a Unpaused event raised by the SuperChainConfig contract.
type SuperChainConfigUnpaused struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0xa45f47fdea8a1efdd9029a5691c7f759c32b7c698632b563573e155625d16933.
//
// Solidity: event Unpaused()
func (_SuperChainConfig *SuperChainConfigFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SuperChainConfigUnpausedIterator, error) {

	logs, sub, err := _SuperChainConfig.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SuperChainConfigUnpausedIterator{contract: _SuperChainConfig.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0xa45f47fdea8a1efdd9029a5691c7f759c32b7c698632b563573e155625d16933.
//
// Solidity: event Unpaused()
func (_SuperChainConfig *SuperChainConfigFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SuperChainConfigUnpaused) (event.Subscription, error) {

	logs, sub, err := _SuperChainConfig.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SuperChainConfigUnpaused)
				if err := _SuperChainConfig.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0xa45f47fdea8a1efdd9029a5691c7f759c32b7c698632b563573e155625d16933.
//
// Solidity: event Unpaused()
func (_SuperChainConfig *SuperChainConfigFilterer) ParseUnpaused(log types.Log) (*SuperChainConfigUnpaused, error) {
	event := new(SuperChainConfigUnpaused)
	if err := _SuperChainConfig.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
