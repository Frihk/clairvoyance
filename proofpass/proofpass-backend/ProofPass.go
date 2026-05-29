// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

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

// ProofPassMetaData contains all meta data concerning the ProofPass contract.
var ProofPassMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"credentialId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"dataHash\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"}],\"name\":\"CredentialIssued\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"}],\"name\":\"addIssuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"authorizedIssuers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"credentialId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dataHash\",\"type\":\"string\"}],\"name\":\"issueCredential\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"credentialId\",\"type\":\"string\"}],\"name\":\"verifyCredential\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ProofPassABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofPassMetaData.ABI instead.
var ProofPassABI = ProofPassMetaData.ABI

// ProofPass is an auto generated Go binding around an Ethereum contract.
type ProofPass struct {
	ProofPassCaller     // Read-only binding to the contract
	ProofPassTransactor // Write-only binding to the contract
	ProofPassFilterer   // Log filterer for contract events
}

// ProofPassCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProofPassCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofPassTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProofPassTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofPassFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProofPassFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofPassSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProofPassSession struct {
	Contract     *ProofPass        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofPassCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProofPassCallerSession struct {
	Contract *ProofPassCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ProofPassTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProofPassTransactorSession struct {
	Contract     *ProofPassTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ProofPassRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProofPassRaw struct {
	Contract *ProofPass // Generic contract binding to access the raw methods on
}

// ProofPassCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProofPassCallerRaw struct {
	Contract *ProofPassCaller // Generic read-only contract binding to access the raw methods on
}

// ProofPassTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProofPassTransactorRaw struct {
	Contract *ProofPassTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProofPass creates a new instance of ProofPass, bound to a specific deployed contract.
func NewProofPass(address common.Address, backend bind.ContractBackend) (*ProofPass, error) {
	contract, err := bindProofPass(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProofPass{ProofPassCaller: ProofPassCaller{contract: contract}, ProofPassTransactor: ProofPassTransactor{contract: contract}, ProofPassFilterer: ProofPassFilterer{contract: contract}}, nil
}

// NewProofPassCaller creates a new read-only instance of ProofPass, bound to a specific deployed contract.
func NewProofPassCaller(address common.Address, caller bind.ContractCaller) (*ProofPassCaller, error) {
	contract, err := bindProofPass(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProofPassCaller{contract: contract}, nil
}

// NewProofPassTransactor creates a new write-only instance of ProofPass, bound to a specific deployed contract.
func NewProofPassTransactor(address common.Address, transactor bind.ContractTransactor) (*ProofPassTransactor, error) {
	contract, err := bindProofPass(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProofPassTransactor{contract: contract}, nil
}

// NewProofPassFilterer creates a new log filterer instance of ProofPass, bound to a specific deployed contract.
func NewProofPassFilterer(address common.Address, filterer bind.ContractFilterer) (*ProofPassFilterer, error) {
	contract, err := bindProofPass(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProofPassFilterer{contract: contract}, nil
}

// bindProofPass binds a generic wrapper to an already deployed contract.
func bindProofPass(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProofPassMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofPass *ProofPassRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofPass.Contract.ProofPassCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofPass *ProofPassRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofPass.Contract.ProofPassTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofPass *ProofPassRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofPass.Contract.ProofPassTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofPass *ProofPassCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofPass.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofPass *ProofPassTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofPass.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofPass *ProofPassTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofPass.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedIssuers is a free data retrieval call binding the contract method 0xf731fa0f.
//
// Solidity: function authorizedIssuers(address ) view returns(bool)
func (_ProofPass *ProofPassCaller) AuthorizedIssuers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ProofPass.contract.Call(opts, &out, "authorizedIssuers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedIssuers is a free data retrieval call binding the contract method 0xf731fa0f.
//
// Solidity: function authorizedIssuers(address ) view returns(bool)
func (_ProofPass *ProofPassSession) AuthorizedIssuers(arg0 common.Address) (bool, error) {
	return _ProofPass.Contract.AuthorizedIssuers(&_ProofPass.CallOpts, arg0)
}

// AuthorizedIssuers is a free data retrieval call binding the contract method 0xf731fa0f.
//
// Solidity: function authorizedIssuers(address ) view returns(bool)
func (_ProofPass *ProofPassCallerSession) AuthorizedIssuers(arg0 common.Address) (bool, error) {
	return _ProofPass.Contract.AuthorizedIssuers(&_ProofPass.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProofPass *ProofPassCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofPass.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProofPass *ProofPassSession) Owner() (common.Address, error) {
	return _ProofPass.Contract.Owner(&_ProofPass.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ProofPass *ProofPassCallerSession) Owner() (common.Address, error) {
	return _ProofPass.Contract.Owner(&_ProofPass.CallOpts)
}

// VerifyCredential is a free data retrieval call binding the contract method 0xbb3670ab.
//
// Solidity: function verifyCredential(string credentialId) view returns(string, address, uint256)
func (_ProofPass *ProofPassCaller) VerifyCredential(opts *bind.CallOpts, credentialId string) (string, common.Address, *big.Int, error) {
	var out []interface{}
	err := _ProofPass.contract.Call(opts, &out, "verifyCredential", credentialId)

	if err != nil {
		return *new(string), *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// VerifyCredential is a free data retrieval call binding the contract method 0xbb3670ab.
//
// Solidity: function verifyCredential(string credentialId) view returns(string, address, uint256)
func (_ProofPass *ProofPassSession) VerifyCredential(credentialId string) (string, common.Address, *big.Int, error) {
	return _ProofPass.Contract.VerifyCredential(&_ProofPass.CallOpts, credentialId)
}

// VerifyCredential is a free data retrieval call binding the contract method 0xbb3670ab.
//
// Solidity: function verifyCredential(string credentialId) view returns(string, address, uint256)
func (_ProofPass *ProofPassCallerSession) VerifyCredential(credentialId string) (string, common.Address, *big.Int, error) {
	return _ProofPass.Contract.VerifyCredential(&_ProofPass.CallOpts, credentialId)
}

// AddIssuer is a paid mutator transaction binding the contract method 0x20694db0.
//
// Solidity: function addIssuer(address issuer) returns()
func (_ProofPass *ProofPassTransactor) AddIssuer(opts *bind.TransactOpts, issuer common.Address) (*types.Transaction, error) {
	return _ProofPass.contract.Transact(opts, "addIssuer", issuer)
}

// AddIssuer is a paid mutator transaction binding the contract method 0x20694db0.
//
// Solidity: function addIssuer(address issuer) returns()
func (_ProofPass *ProofPassSession) AddIssuer(issuer common.Address) (*types.Transaction, error) {
	return _ProofPass.Contract.AddIssuer(&_ProofPass.TransactOpts, issuer)
}

// AddIssuer is a paid mutator transaction binding the contract method 0x20694db0.
//
// Solidity: function addIssuer(address issuer) returns()
func (_ProofPass *ProofPassTransactorSession) AddIssuer(issuer common.Address) (*types.Transaction, error) {
	return _ProofPass.Contract.AddIssuer(&_ProofPass.TransactOpts, issuer)
}

// IssueCredential is a paid mutator transaction binding the contract method 0x51ef0326.
//
// Solidity: function issueCredential(string credentialId, string dataHash) returns()
func (_ProofPass *ProofPassTransactor) IssueCredential(opts *bind.TransactOpts, credentialId string, dataHash string) (*types.Transaction, error) {
	return _ProofPass.contract.Transact(opts, "issueCredential", credentialId, dataHash)
}

// IssueCredential is a paid mutator transaction binding the contract method 0x51ef0326.
//
// Solidity: function issueCredential(string credentialId, string dataHash) returns()
func (_ProofPass *ProofPassSession) IssueCredential(credentialId string, dataHash string) (*types.Transaction, error) {
	return _ProofPass.Contract.IssueCredential(&_ProofPass.TransactOpts, credentialId, dataHash)
}

// IssueCredential is a paid mutator transaction binding the contract method 0x51ef0326.
//
// Solidity: function issueCredential(string credentialId, string dataHash) returns()
func (_ProofPass *ProofPassTransactorSession) IssueCredential(credentialId string, dataHash string) (*types.Transaction, error) {
	return _ProofPass.Contract.IssueCredential(&_ProofPass.TransactOpts, credentialId, dataHash)
}

// ProofPassCredentialIssuedIterator is returned from FilterCredentialIssued and is used to iterate over the raw logs and unpacked data for CredentialIssued events raised by the ProofPass contract.
type ProofPassCredentialIssuedIterator struct {
	Event *ProofPassCredentialIssued // Event containing the contract specifics and raw log

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
func (it *ProofPassCredentialIssuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofPassCredentialIssued)
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
		it.Event = new(ProofPassCredentialIssued)
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
func (it *ProofPassCredentialIssuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofPassCredentialIssuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofPassCredentialIssued represents a CredentialIssued event raised by the ProofPass contract.
type ProofPassCredentialIssued struct {
	CredentialId string
	DataHash     string
	Issuer       common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterCredentialIssued is a free log retrieval operation binding the contract event 0xbced6383e2dd56a440d8755db8d2520393bede4e079e12a544d29aa065174711.
//
// Solidity: event CredentialIssued(string credentialId, string dataHash, address issuer)
func (_ProofPass *ProofPassFilterer) FilterCredentialIssued(opts *bind.FilterOpts) (*ProofPassCredentialIssuedIterator, error) {

	logs, sub, err := _ProofPass.contract.FilterLogs(opts, "CredentialIssued")
	if err != nil {
		return nil, err
	}
	return &ProofPassCredentialIssuedIterator{contract: _ProofPass.contract, event: "CredentialIssued", logs: logs, sub: sub}, nil
}

// WatchCredentialIssued is a free log subscription operation binding the contract event 0xbced6383e2dd56a440d8755db8d2520393bede4e079e12a544d29aa065174711.
//
// Solidity: event CredentialIssued(string credentialId, string dataHash, address issuer)
func (_ProofPass *ProofPassFilterer) WatchCredentialIssued(opts *bind.WatchOpts, sink chan<- *ProofPassCredentialIssued) (event.Subscription, error) {

	logs, sub, err := _ProofPass.contract.WatchLogs(opts, "CredentialIssued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofPassCredentialIssued)
				if err := _ProofPass.contract.UnpackLog(event, "CredentialIssued", log); err != nil {
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

// ParseCredentialIssued is a log parse operation binding the contract event 0xbced6383e2dd56a440d8755db8d2520393bede4e079e12a544d29aa065174711.
//
// Solidity: event CredentialIssued(string credentialId, string dataHash, address issuer)
func (_ProofPass *ProofPassFilterer) ParseCredentialIssued(log types.Log) (*ProofPassCredentialIssued, error) {
	event := new(ProofPassCredentialIssued)
	if err := _ProofPass.contract.UnpackLog(event, "CredentialIssued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
