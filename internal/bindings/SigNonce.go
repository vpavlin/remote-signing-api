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
)

// SigNonceMetaData contains all meta data concerning the SigNonce contract.
var SigNonceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sigNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SigNonceABI is the input ABI used to generate the binding from.
// Deprecated: Use SigNonceMetaData.ABI instead.
var SigNonceABI = SigNonceMetaData.ABI

// SigNonce is an auto generated Go binding around an Ethereum contract.
type SigNonce struct {
	SigNonceCaller     // Read-only binding to the contract
	SigNonceTransactor // Write-only binding to the contract
	SigNonceFilterer   // Log filterer for contract events
}

// SigNonceCaller is an auto generated read-only Go binding around an Ethereum contract.
type SigNonceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigNonceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SigNonceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigNonceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SigNonceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigNonceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SigNonceSession struct {
	Contract     *SigNonce         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SigNonceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SigNonceCallerSession struct {
	Contract *SigNonceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SigNonceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SigNonceTransactorSession struct {
	Contract     *SigNonceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SigNonceRaw is an auto generated low-level Go binding around an Ethereum contract.
type SigNonceRaw struct {
	Contract *SigNonce // Generic contract binding to access the raw methods on
}

// SigNonceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SigNonceCallerRaw struct {
	Contract *SigNonceCaller // Generic read-only contract binding to access the raw methods on
}

// SigNonceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SigNonceTransactorRaw struct {
	Contract *SigNonceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSigNonce creates a new instance of SigNonce, bound to a specific deployed contract.
func NewSigNonce(address common.Address, backend bind.ContractBackend) (*SigNonce, error) {
	contract, err := bindSigNonce(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SigNonce{SigNonceCaller: SigNonceCaller{contract: contract}, SigNonceTransactor: SigNonceTransactor{contract: contract}, SigNonceFilterer: SigNonceFilterer{contract: contract}}, nil
}

// NewSigNonceCaller creates a new read-only instance of SigNonce, bound to a specific deployed contract.
func NewSigNonceCaller(address common.Address, caller bind.ContractCaller) (*SigNonceCaller, error) {
	contract, err := bindSigNonce(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SigNonceCaller{contract: contract}, nil
}

// NewSigNonceTransactor creates a new write-only instance of SigNonce, bound to a specific deployed contract.
func NewSigNonceTransactor(address common.Address, transactor bind.ContractTransactor) (*SigNonceTransactor, error) {
	contract, err := bindSigNonce(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SigNonceTransactor{contract: contract}, nil
}

// NewSigNonceFilterer creates a new log filterer instance of SigNonce, bound to a specific deployed contract.
func NewSigNonceFilterer(address common.Address, filterer bind.ContractFilterer) (*SigNonceFilterer, error) {
	contract, err := bindSigNonce(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SigNonceFilterer{contract: contract}, nil
}

// bindSigNonce binds a generic wrapper to an already deployed contract.
func bindSigNonce(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SigNonceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SigNonce *SigNonceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SigNonce.Contract.SigNonceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SigNonce *SigNonceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SigNonce.Contract.SigNonceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SigNonce *SigNonceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SigNonce.Contract.SigNonceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SigNonce *SigNonceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SigNonce.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SigNonce *SigNonceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SigNonce.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SigNonce *SigNonceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SigNonce.Contract.contract.Transact(opts, method, params...)
}

// SigNonce is a free data retrieval call binding the contract method 0xb4cbf1cb.
//
// Solidity: function sigNonce(address ) view returns(uint256)
func (_SigNonce *SigNonceCaller) SigNonce(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SigNonce.contract.Call(opts, &out, "sigNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SigNonce is a free data retrieval call binding the contract method 0xb4cbf1cb.
//
// Solidity: function sigNonce(address ) view returns(uint256)
func (_SigNonce *SigNonceSession) SigNonce(arg0 common.Address) (*big.Int, error) {
	return _SigNonce.Contract.SigNonce(&_SigNonce.CallOpts, arg0)
}

// SigNonce is a free data retrieval call binding the contract method 0xb4cbf1cb.
//
// Solidity: function sigNonce(address ) view returns(uint256)
func (_SigNonce *SigNonceCallerSession) SigNonce(arg0 common.Address) (*big.Int, error) {
	return _SigNonce.Contract.SigNonce(&_SigNonce.CallOpts, arg0)
}

