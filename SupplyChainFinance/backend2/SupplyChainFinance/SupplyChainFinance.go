// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package SupplyChainFinance

import (
	"math/big"
	"strings"

	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
)

// SupplyChainFinanceEnterprise is an auto generated low-level Go binding around an user-defined struct.
type SupplyChainFinanceEnterprise struct {
	EnterpriseId string
	Name         string
	Wallet       common.Address
	IsCore       bool
	IsFinancial  bool
	RegisterTime *big.Int
}

// SupplyChainFinanceFinancingRecord is an auto generated low-level Go binding around an user-defined struct.
type SupplyChainFinanceFinancingRecord struct {
	FinancialInst common.Address
	Amount        *big.Int
	InterestRate  *big.Int
	FinancingTime *big.Int
}

// SupplyChainFinanceReceivable is an auto generated low-level Go binding around an user-defined struct.
type SupplyChainFinanceReceivable struct {
	ReceivableId string
	OrderId      string
	Issuer       common.Address
	Payer        common.Address
	Amount       *big.Int
	IssueTime    *big.Int
	DueTime      *big.Int
	Status       uint8
}

// SupplyChainFinanceABI is the input ABI used to generate the binding from.
const SupplyChainFinanceABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enterpriseId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isCore\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinancial\",\"type\":\"bool\"}],\"name\":\"EnterpriseRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receivableId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"ReceivableConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receivableId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"financialInst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"}],\"name\":\"ReceivableFinanced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receivableId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"orderId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ReceivableIssued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receivableId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"ReceivableSettled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_receivableId\",\"type\":\"string\"}],\"name\":\"confirmReceivable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"enterpriseIdToAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"enterpriseList\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"enterprises\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"enterpriseId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isCore\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isFinancial\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"registerTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"financingRecords\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"financialInst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"financingTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getEnterprise\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"enterpriseId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"wallet\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isCore\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isFinancial\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"registerTime\",\"type\":\"uint256\"}],\"internalType\":\"structSupplyChainFinance.Enterprise\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEnterpriseCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_receivableId\",\"type\":\"string\"}],\"name\":\"getFinancingRecords\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"financialInst\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"financingTime\",\"type\":\"uint256\"}],\"internalType\":\"structSupplyChainFinance.FinancingRecord[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_id\",\"type\":\"string\"}],\"name\":\"getReceivable\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"receivableId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"orderId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issueTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dueTime\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"}],\"internalType\":\"structSupplyChainFinance.Receivable\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getReceivableCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getReceivableIdByIndex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_receivableId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_orderId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_payer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_dueTime\",\"type\":\"uint256\"}],\"name\":\"issueReceivable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"receivableIdList\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"receivables\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"receivableId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"orderId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issueTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dueTime\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_receivableId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interestRate\",\"type\":\"uint256\"}],\"name\":\"recordFinancing\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_receivableId\",\"type\":\"string\"}],\"name\":\"recordSettlement\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enterpriseId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"_isCore\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"_isFinancial\",\"type\":\"bool\"}],\"name\":\"registerEnterprise\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SupplyChainFinance is an auto generated Go binding around a Solidity contract.
type SupplyChainFinance struct {
	SupplyChainFinanceCaller     // Read-only binding to the contract
	SupplyChainFinanceTransactor // Write-only binding to the contract
	SupplyChainFinanceFilterer   // Log filterer for contract events
}

// SupplyChainFinanceCaller is an auto generated read-only Go binding around a Solidity contract.
type SupplyChainFinanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SupplyChainFinanceTransactor is an auto generated write-only Go binding around a Solidity contract.
type SupplyChainFinanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SupplyChainFinanceFilterer is an auto generated log filtering Go binding around a Solidity contract events.
type SupplyChainFinanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SupplyChainFinanceSession is an auto generated Go binding around a Solidity contract,
// with pre-set call and transact options.
type SupplyChainFinanceSession struct {
	Contract     *SupplyChainFinance // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SupplyChainFinanceCallerSession is an auto generated read-only Go binding around a Solidity contract,
// with pre-set call options.
type SupplyChainFinanceCallerSession struct {
	Contract *SupplyChainFinanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// SupplyChainFinanceTransactorSession is an auto generated write-only Go binding around a Solidity contract,
// with pre-set transact options.
type SupplyChainFinanceTransactorSession struct {
	Contract     *SupplyChainFinanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// SupplyChainFinanceRaw is an auto generated low-level Go binding around a Solidity contract.
type SupplyChainFinanceRaw struct {
	Contract *SupplyChainFinance // Generic contract binding to access the raw methods on
}

// SupplyChainFinanceCallerRaw is an auto generated low-level read-only Go binding around a Solidity contract.
type SupplyChainFinanceCallerRaw struct {
	Contract *SupplyChainFinanceCaller // Generic read-only contract binding to access the raw methods on
}

// SupplyChainFinanceTransactorRaw is an auto generated low-level write-only Go binding around a Solidity contract.
type SupplyChainFinanceTransactorRaw struct {
	Contract *SupplyChainFinanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSupplyChainFinance creates a new instance of SupplyChainFinance, bound to a specific deployed contract.
func NewSupplyChainFinance(address common.Address, backend bind.ContractBackend) (*SupplyChainFinance, error) {
	contract, err := bindSupplyChainFinance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SupplyChainFinance{SupplyChainFinanceCaller: SupplyChainFinanceCaller{contract: contract}, SupplyChainFinanceTransactor: SupplyChainFinanceTransactor{contract: contract}, SupplyChainFinanceFilterer: SupplyChainFinanceFilterer{contract: contract}}, nil
}

// NewSupplyChainFinanceCaller creates a new read-only instance of SupplyChainFinance, bound to a specific deployed contract.
func NewSupplyChainFinanceCaller(address common.Address, caller bind.ContractCaller) (*SupplyChainFinanceCaller, error) {
	contract, err := bindSupplyChainFinance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SupplyChainFinanceCaller{contract: contract}, nil
}

// NewSupplyChainFinanceTransactor creates a new write-only instance of SupplyChainFinance, bound to a specific deployed contract.
func NewSupplyChainFinanceTransactor(address common.Address, transactor bind.ContractTransactor) (*SupplyChainFinanceTransactor, error) {
	contract, err := bindSupplyChainFinance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SupplyChainFinanceTransactor{contract: contract}, nil
}

// NewSupplyChainFinanceFilterer creates a new log filterer instance of SupplyChainFinance, bound to a specific deployed contract.
func NewSupplyChainFinanceFilterer(address common.Address, filterer bind.ContractFilterer) (*SupplyChainFinanceFilterer, error) {
	contract, err := bindSupplyChainFinance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SupplyChainFinanceFilterer{contract: contract}, nil
}

// bindSupplyChainFinance binds a generic wrapper to an already deployed contract.
func bindSupplyChainFinance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SupplyChainFinanceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SupplyChainFinance *SupplyChainFinanceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SupplyChainFinance.Contract.SupplyChainFinanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SupplyChainFinance *SupplyChainFinanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.SupplyChainFinanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SupplyChainFinance *SupplyChainFinanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.SupplyChainFinanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SupplyChainFinance *SupplyChainFinanceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SupplyChainFinance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SupplyChainFinance *SupplyChainFinanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SupplyChainFinance *SupplyChainFinanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.contract.Transact(opts, method, params...)
}

// EnterpriseIdToAddr is a free data retrieval call binding the contract method 0x68a5d3f3.
//
// Solidity: function enterpriseIdToAddr(string ) constant returns(address)
func (_SupplyChainFinance *SupplyChainFinanceCaller) EnterpriseIdToAddr(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "enterpriseIdToAddr", arg0)
	return *ret0, err
}

// EnterpriseIdToAddr is a free data retrieval call binding the contract method 0x68a5d3f3.
//
// Solidity: function enterpriseIdToAddr(string ) constant returns(address)
func (_SupplyChainFinance *SupplyChainFinanceSession) EnterpriseIdToAddr(arg0 string) (common.Address, error) {
	return _SupplyChainFinance.Contract.EnterpriseIdToAddr(&_SupplyChainFinance.CallOpts, arg0)
}

// EnterpriseIdToAddr is a free data retrieval call binding the contract method 0x68a5d3f3.
//
// Solidity: function enterpriseIdToAddr(string ) constant returns(address)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) EnterpriseIdToAddr(arg0 string) (common.Address, error) {
	return _SupplyChainFinance.Contract.EnterpriseIdToAddr(&_SupplyChainFinance.CallOpts, arg0)
}

// EnterpriseList is a free data retrieval call binding the contract method 0x1348990b.
//
// Solidity: function enterpriseList(uint256 ) constant returns(address)
func (_SupplyChainFinance *SupplyChainFinanceCaller) EnterpriseList(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "enterpriseList", arg0)
	return *ret0, err
}

// EnterpriseList is a free data retrieval call binding the contract method 0x1348990b.
//
// Solidity: function enterpriseList(uint256 ) constant returns(address)
func (_SupplyChainFinance *SupplyChainFinanceSession) EnterpriseList(arg0 *big.Int) (common.Address, error) {
	return _SupplyChainFinance.Contract.EnterpriseList(&_SupplyChainFinance.CallOpts, arg0)
}

// EnterpriseList is a free data retrieval call binding the contract method 0x1348990b.
//
// Solidity: function enterpriseList(uint256 ) constant returns(address)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) EnterpriseList(arg0 *big.Int) (common.Address, error) {
	return _SupplyChainFinance.Contract.EnterpriseList(&_SupplyChainFinance.CallOpts, arg0)
}

// Enterprises is a free data retrieval call binding the contract method 0x902e11b6.
//
// Solidity: function enterprises(address ) constant returns(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial, uint256 registerTime)
func (_SupplyChainFinance *SupplyChainFinanceCaller) Enterprises(opts *bind.CallOpts, arg0 common.Address) (struct {
	EnterpriseId string
	Name         string
	Wallet       common.Address
	IsCore       bool
	IsFinancial  bool
	RegisterTime *big.Int
}, error) {
	ret := new(struct {
		EnterpriseId string
		Name         string
		Wallet       common.Address
		IsCore       bool
		IsFinancial  bool
		RegisterTime *big.Int
	})
	out := ret
	err := _SupplyChainFinance.contract.Call(opts, out, "enterprises", arg0)
	return *ret, err
}

// Enterprises is a free data retrieval call binding the contract method 0x902e11b6.
//
// Solidity: function enterprises(address ) constant returns(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial, uint256 registerTime)
func (_SupplyChainFinance *SupplyChainFinanceSession) Enterprises(arg0 common.Address) (struct {
	EnterpriseId string
	Name         string
	Wallet       common.Address
	IsCore       bool
	IsFinancial  bool
	RegisterTime *big.Int
}, error) {
	return _SupplyChainFinance.Contract.Enterprises(&_SupplyChainFinance.CallOpts, arg0)
}

// Enterprises is a free data retrieval call binding the contract method 0x902e11b6.
//
// Solidity: function enterprises(address ) constant returns(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial, uint256 registerTime)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) Enterprises(arg0 common.Address) (struct {
	EnterpriseId string
	Name         string
	Wallet       common.Address
	IsCore       bool
	IsFinancial  bool
	RegisterTime *big.Int
}, error) {
	return _SupplyChainFinance.Contract.Enterprises(&_SupplyChainFinance.CallOpts, arg0)
}

// FinancingRecords is a free data retrieval call binding the contract method 0x0f874ebb.
//
// Solidity: function financingRecords(string , uint256 ) constant returns(address financialInst, uint256 amount, uint256 interestRate, uint256 financingTime)
func (_SupplyChainFinance *SupplyChainFinanceCaller) FinancingRecords(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	FinancialInst common.Address
	Amount        *big.Int
	InterestRate  *big.Int
	FinancingTime *big.Int
}, error) {
	ret := new(struct {
		FinancialInst common.Address
		Amount        *big.Int
		InterestRate  *big.Int
		FinancingTime *big.Int
	})
	out := ret
	err := _SupplyChainFinance.contract.Call(opts, out, "financingRecords", arg0, arg1)
	return *ret, err
}

// FinancingRecords is a free data retrieval call binding the contract method 0x0f874ebb.
//
// Solidity: function financingRecords(string , uint256 ) constant returns(address financialInst, uint256 amount, uint256 interestRate, uint256 financingTime)
func (_SupplyChainFinance *SupplyChainFinanceSession) FinancingRecords(arg0 string, arg1 *big.Int) (struct {
	FinancialInst common.Address
	Amount        *big.Int
	InterestRate  *big.Int
	FinancingTime *big.Int
}, error) {
	return _SupplyChainFinance.Contract.FinancingRecords(&_SupplyChainFinance.CallOpts, arg0, arg1)
}

// FinancingRecords is a free data retrieval call binding the contract method 0x0f874ebb.
//
// Solidity: function financingRecords(string , uint256 ) constant returns(address financialInst, uint256 amount, uint256 interestRate, uint256 financingTime)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) FinancingRecords(arg0 string, arg1 *big.Int) (struct {
	FinancialInst common.Address
	Amount        *big.Int
	InterestRate  *big.Int
	FinancingTime *big.Int
}, error) {
	return _SupplyChainFinance.Contract.FinancingRecords(&_SupplyChainFinance.CallOpts, arg0, arg1)
}

// GetEnterprise is a free data retrieval call binding the contract method 0x889cd253.
//
// Solidity: function getEnterprise(address _addr) constant returns(SupplyChainFinanceEnterprise)
func (_SupplyChainFinance *SupplyChainFinanceCaller) GetEnterprise(opts *bind.CallOpts, _addr common.Address) (SupplyChainFinanceEnterprise, error) {
	var (
		ret0 = new(SupplyChainFinanceEnterprise)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "getEnterprise", _addr)
	return *ret0, err
}

// GetEnterprise is a free data retrieval call binding the contract method 0x889cd253.
//
// Solidity: function getEnterprise(address _addr) constant returns(SupplyChainFinanceEnterprise)
func (_SupplyChainFinance *SupplyChainFinanceSession) GetEnterprise(_addr common.Address) (SupplyChainFinanceEnterprise, error) {
	return _SupplyChainFinance.Contract.GetEnterprise(&_SupplyChainFinance.CallOpts, _addr)
}

// GetEnterprise is a free data retrieval call binding the contract method 0x889cd253.
//
// Solidity: function getEnterprise(address _addr) constant returns(SupplyChainFinanceEnterprise)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) GetEnterprise(_addr common.Address) (SupplyChainFinanceEnterprise, error) {
	return _SupplyChainFinance.Contract.GetEnterprise(&_SupplyChainFinance.CallOpts, _addr)
}

// GetEnterpriseCount is a free data retrieval call binding the contract method 0x52d07cd0.
//
// Solidity: function getEnterpriseCount() constant returns(uint256)
func (_SupplyChainFinance *SupplyChainFinanceCaller) GetEnterpriseCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "getEnterpriseCount")
	return *ret0, err
}

// GetEnterpriseCount is a free data retrieval call binding the contract method 0x52d07cd0.
//
// Solidity: function getEnterpriseCount() constant returns(uint256)
func (_SupplyChainFinance *SupplyChainFinanceSession) GetEnterpriseCount() (*big.Int, error) {
	return _SupplyChainFinance.Contract.GetEnterpriseCount(&_SupplyChainFinance.CallOpts)
}

// GetEnterpriseCount is a free data retrieval call binding the contract method 0x52d07cd0.
//
// Solidity: function getEnterpriseCount() constant returns(uint256)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) GetEnterpriseCount() (*big.Int, error) {
	return _SupplyChainFinance.Contract.GetEnterpriseCount(&_SupplyChainFinance.CallOpts)
}

// GetFinancingRecords is a free data retrieval call binding the contract method 0x308b1481.
//
// Solidity: function getFinancingRecords(string _receivableId) constant returns([]SupplyChainFinanceFinancingRecord)
func (_SupplyChainFinance *SupplyChainFinanceCaller) GetFinancingRecords(opts *bind.CallOpts, _receivableId string) ([]SupplyChainFinanceFinancingRecord, error) {
	var (
		ret0 = new([]SupplyChainFinanceFinancingRecord)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "getFinancingRecords", _receivableId)
	return *ret0, err
}

// GetFinancingRecords is a free data retrieval call binding the contract method 0x308b1481.
//
// Solidity: function getFinancingRecords(string _receivableId) constant returns([]SupplyChainFinanceFinancingRecord)
func (_SupplyChainFinance *SupplyChainFinanceSession) GetFinancingRecords(_receivableId string) ([]SupplyChainFinanceFinancingRecord, error) {
	return _SupplyChainFinance.Contract.GetFinancingRecords(&_SupplyChainFinance.CallOpts, _receivableId)
}

// GetFinancingRecords is a free data retrieval call binding the contract method 0x308b1481.
//
// Solidity: function getFinancingRecords(string _receivableId) constant returns([]SupplyChainFinanceFinancingRecord)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) GetFinancingRecords(_receivableId string) ([]SupplyChainFinanceFinancingRecord, error) {
	return _SupplyChainFinance.Contract.GetFinancingRecords(&_SupplyChainFinance.CallOpts, _receivableId)
}

// GetReceivable is a free data retrieval call binding the contract method 0xc133fb3a.
//
// Solidity: function getReceivable(string _id) constant returns(SupplyChainFinanceReceivable)
func (_SupplyChainFinance *SupplyChainFinanceCaller) GetReceivable(opts *bind.CallOpts, _id string) (SupplyChainFinanceReceivable, error) {
	var (
		ret0 = new(SupplyChainFinanceReceivable)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "getReceivable", _id)
	return *ret0, err
}

// GetReceivable is a free data retrieval call binding the contract method 0xc133fb3a.
//
// Solidity: function getReceivable(string _id) constant returns(SupplyChainFinanceReceivable)
func (_SupplyChainFinance *SupplyChainFinanceSession) GetReceivable(_id string) (SupplyChainFinanceReceivable, error) {
	return _SupplyChainFinance.Contract.GetReceivable(&_SupplyChainFinance.CallOpts, _id)
}

// GetReceivable is a free data retrieval call binding the contract method 0xc133fb3a.
//
// Solidity: function getReceivable(string _id) constant returns(SupplyChainFinanceReceivable)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) GetReceivable(_id string) (SupplyChainFinanceReceivable, error) {
	return _SupplyChainFinance.Contract.GetReceivable(&_SupplyChainFinance.CallOpts, _id)
}

// GetReceivableCount is a free data retrieval call binding the contract method 0xd4223d55.
//
// Solidity: function getReceivableCount() constant returns(uint256)
func (_SupplyChainFinance *SupplyChainFinanceCaller) GetReceivableCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "getReceivableCount")
	return *ret0, err
}

// GetReceivableCount is a free data retrieval call binding the contract method 0xd4223d55.
//
// Solidity: function getReceivableCount() constant returns(uint256)
func (_SupplyChainFinance *SupplyChainFinanceSession) GetReceivableCount() (*big.Int, error) {
	return _SupplyChainFinance.Contract.GetReceivableCount(&_SupplyChainFinance.CallOpts)
}

// GetReceivableCount is a free data retrieval call binding the contract method 0xd4223d55.
//
// Solidity: function getReceivableCount() constant returns(uint256)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) GetReceivableCount() (*big.Int, error) {
	return _SupplyChainFinance.Contract.GetReceivableCount(&_SupplyChainFinance.CallOpts)
}

// GetReceivableIdByIndex is a free data retrieval call binding the contract method 0xfb277a9f.
//
// Solidity: function getReceivableIdByIndex(uint256 _index) constant returns(string)
func (_SupplyChainFinance *SupplyChainFinanceCaller) GetReceivableIdByIndex(opts *bind.CallOpts, _index *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "getReceivableIdByIndex", _index)
	return *ret0, err
}

// GetReceivableIdByIndex is a free data retrieval call binding the contract method 0xfb277a9f.
//
// Solidity: function getReceivableIdByIndex(uint256 _index) constant returns(string)
func (_SupplyChainFinance *SupplyChainFinanceSession) GetReceivableIdByIndex(_index *big.Int) (string, error) {
	return _SupplyChainFinance.Contract.GetReceivableIdByIndex(&_SupplyChainFinance.CallOpts, _index)
}

// GetReceivableIdByIndex is a free data retrieval call binding the contract method 0xfb277a9f.
//
// Solidity: function getReceivableIdByIndex(uint256 _index) constant returns(string)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) GetReceivableIdByIndex(_index *big.Int) (string, error) {
	return _SupplyChainFinance.Contract.GetReceivableIdByIndex(&_SupplyChainFinance.CallOpts, _index)
}

// ReceivableIdList is a free data retrieval call binding the contract method 0x00c2d786.
//
// Solidity: function receivableIdList(uint256 ) constant returns(string)
func (_SupplyChainFinance *SupplyChainFinanceCaller) ReceivableIdList(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SupplyChainFinance.contract.Call(opts, out, "receivableIdList", arg0)
	return *ret0, err
}

// ReceivableIdList is a free data retrieval call binding the contract method 0x00c2d786.
//
// Solidity: function receivableIdList(uint256 ) constant returns(string)
func (_SupplyChainFinance *SupplyChainFinanceSession) ReceivableIdList(arg0 *big.Int) (string, error) {
	return _SupplyChainFinance.Contract.ReceivableIdList(&_SupplyChainFinance.CallOpts, arg0)
}

// ReceivableIdList is a free data retrieval call binding the contract method 0x00c2d786.
//
// Solidity: function receivableIdList(uint256 ) constant returns(string)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) ReceivableIdList(arg0 *big.Int) (string, error) {
	return _SupplyChainFinance.Contract.ReceivableIdList(&_SupplyChainFinance.CallOpts, arg0)
}

// Receivables is a free data retrieval call binding the contract method 0x97f36de6.
//
// Solidity: function receivables(string ) constant returns(string receivableId, string orderId, address issuer, address payer, uint256 amount, uint256 issueTime, uint256 dueTime, uint8 status)
func (_SupplyChainFinance *SupplyChainFinanceCaller) Receivables(opts *bind.CallOpts, arg0 string) (struct {
	ReceivableId string
	OrderId      string
	Issuer       common.Address
	Payer        common.Address
	Amount       *big.Int
	IssueTime    *big.Int
	DueTime      *big.Int
	Status       uint8
}, error) {
	ret := new(struct {
		ReceivableId string
		OrderId      string
		Issuer       common.Address
		Payer        common.Address
		Amount       *big.Int
		IssueTime    *big.Int
		DueTime      *big.Int
		Status       uint8
	})
	out := ret
	err := _SupplyChainFinance.contract.Call(opts, out, "receivables", arg0)
	return *ret, err
}

// Receivables is a free data retrieval call binding the contract method 0x97f36de6.
//
// Solidity: function receivables(string ) constant returns(string receivableId, string orderId, address issuer, address payer, uint256 amount, uint256 issueTime, uint256 dueTime, uint8 status)
func (_SupplyChainFinance *SupplyChainFinanceSession) Receivables(arg0 string) (struct {
	ReceivableId string
	OrderId      string
	Issuer       common.Address
	Payer        common.Address
	Amount       *big.Int
	IssueTime    *big.Int
	DueTime      *big.Int
	Status       uint8
}, error) {
	return _SupplyChainFinance.Contract.Receivables(&_SupplyChainFinance.CallOpts, arg0)
}

// Receivables is a free data retrieval call binding the contract method 0x97f36de6.
//
// Solidity: function receivables(string ) constant returns(string receivableId, string orderId, address issuer, address payer, uint256 amount, uint256 issueTime, uint256 dueTime, uint8 status)
func (_SupplyChainFinance *SupplyChainFinanceCallerSession) Receivables(arg0 string) (struct {
	ReceivableId string
	OrderId      string
	Issuer       common.Address
	Payer        common.Address
	Amount       *big.Int
	IssueTime    *big.Int
	DueTime      *big.Int
	Status       uint8
}, error) {
	return _SupplyChainFinance.Contract.Receivables(&_SupplyChainFinance.CallOpts, arg0)
}

// ConfirmReceivable is a paid mutator transaction binding the contract method 0xdfe7afa4.
//
// Solidity: function confirmReceivable(string _receivableId) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactor) ConfirmReceivable(opts *bind.TransactOpts, _receivableId string) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.contract.Transact(opts, "confirmReceivable", _receivableId)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactor) AsyncConfirmReceivable(handler func(*types.Receipt, error), opts *bind.TransactOpts, _receivableId string) (*types.Transaction, error) {
	return _SupplyChainFinance.contract.AsyncTransact(opts, handler, "confirmReceivable", _receivableId)
}

// ConfirmReceivable is a paid mutator transaction binding the contract method 0xdfe7afa4.
//
// Solidity: function confirmReceivable(string _receivableId) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceSession) ConfirmReceivable(_receivableId string) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.ConfirmReceivable(&_SupplyChainFinance.TransactOpts, _receivableId)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) AsyncConfirmReceivable(handler func(*types.Receipt, error), _receivableId string) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncConfirmReceivable(handler, &_SupplyChainFinance.TransactOpts, _receivableId)
}

// ConfirmReceivable is a paid mutator transaction binding the contract method 0xdfe7afa4.
//
// Solidity: function confirmReceivable(string _receivableId) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) ConfirmReceivable(_receivableId string) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.ConfirmReceivable(&_SupplyChainFinance.TransactOpts, _receivableId)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) AsyncConfirmReceivable(handler func(*types.Receipt, error), _receivableId string) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncConfirmReceivable(handler, &_SupplyChainFinance.TransactOpts, _receivableId)
}

// IssueReceivable is a paid mutator transaction binding the contract method 0x9ad3f9f2.
//
// Solidity: function issueReceivable(string _receivableId, string _orderId, address _payer, uint256 _amount, uint256 _dueTime) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactor) IssueReceivable(opts *bind.TransactOpts, _receivableId string, _orderId string, _payer common.Address, _amount *big.Int, _dueTime *big.Int) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.contract.Transact(opts, "issueReceivable", _receivableId, _orderId, _payer, _amount, _dueTime)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactor) AsyncIssueReceivable(handler func(*types.Receipt, error), opts *bind.TransactOpts, _receivableId string, _orderId string, _payer common.Address, _amount *big.Int, _dueTime *big.Int) (*types.Transaction, error) {
	return _SupplyChainFinance.contract.AsyncTransact(opts, handler, "issueReceivable", _receivableId, _orderId, _payer, _amount, _dueTime)
}

// IssueReceivable is a paid mutator transaction binding the contract method 0x9ad3f9f2.
//
// Solidity: function issueReceivable(string _receivableId, string _orderId, address _payer, uint256 _amount, uint256 _dueTime) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceSession) IssueReceivable(_receivableId string, _orderId string, _payer common.Address, _amount *big.Int, _dueTime *big.Int) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.IssueReceivable(&_SupplyChainFinance.TransactOpts, _receivableId, _orderId, _payer, _amount, _dueTime)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) AsyncIssueReceivable(handler func(*types.Receipt, error), _receivableId string, _orderId string, _payer common.Address, _amount *big.Int, _dueTime *big.Int) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncIssueReceivable(handler, &_SupplyChainFinance.TransactOpts, _receivableId, _orderId, _payer, _amount, _dueTime)
}

// IssueReceivable is a paid mutator transaction binding the contract method 0x9ad3f9f2.
//
// Solidity: function issueReceivable(string _receivableId, string _orderId, address _payer, uint256 _amount, uint256 _dueTime) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) IssueReceivable(_receivableId string, _orderId string, _payer common.Address, _amount *big.Int, _dueTime *big.Int) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.IssueReceivable(&_SupplyChainFinance.TransactOpts, _receivableId, _orderId, _payer, _amount, _dueTime)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) AsyncIssueReceivable(handler func(*types.Receipt, error), _receivableId string, _orderId string, _payer common.Address, _amount *big.Int, _dueTime *big.Int) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncIssueReceivable(handler, &_SupplyChainFinance.TransactOpts, _receivableId, _orderId, _payer, _amount, _dueTime)
}

// RecordFinancing is a paid mutator transaction binding the contract method 0x3c3a3ebe.
//
// Solidity: function recordFinancing(string _receivableId, uint256 _amount, uint256 _interestRate) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactor) RecordFinancing(opts *bind.TransactOpts, _receivableId string, _amount *big.Int, _interestRate *big.Int) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.contract.Transact(opts, "recordFinancing", _receivableId, _amount, _interestRate)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactor) AsyncRecordFinancing(handler func(*types.Receipt, error), opts *bind.TransactOpts, _receivableId string, _amount *big.Int, _interestRate *big.Int) (*types.Transaction, error) {
	return _SupplyChainFinance.contract.AsyncTransact(opts, handler, "recordFinancing", _receivableId, _amount, _interestRate)
}

// RecordFinancing is a paid mutator transaction binding the contract method 0x3c3a3ebe.
//
// Solidity: function recordFinancing(string _receivableId, uint256 _amount, uint256 _interestRate) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceSession) RecordFinancing(_receivableId string, _amount *big.Int, _interestRate *big.Int) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.RecordFinancing(&_SupplyChainFinance.TransactOpts, _receivableId, _amount, _interestRate)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) AsyncRecordFinancing(handler func(*types.Receipt, error), _receivableId string, _amount *big.Int, _interestRate *big.Int) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncRecordFinancing(handler, &_SupplyChainFinance.TransactOpts, _receivableId, _amount, _interestRate)
}

// RecordFinancing is a paid mutator transaction binding the contract method 0x3c3a3ebe.
//
// Solidity: function recordFinancing(string _receivableId, uint256 _amount, uint256 _interestRate) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) RecordFinancing(_receivableId string, _amount *big.Int, _interestRate *big.Int) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.RecordFinancing(&_SupplyChainFinance.TransactOpts, _receivableId, _amount, _interestRate)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) AsyncRecordFinancing(handler func(*types.Receipt, error), _receivableId string, _amount *big.Int, _interestRate *big.Int) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncRecordFinancing(handler, &_SupplyChainFinance.TransactOpts, _receivableId, _amount, _interestRate)
}

// RecordSettlement is a paid mutator transaction binding the contract method 0xcfff7343.
//
// Solidity: function recordSettlement(string _receivableId) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactor) RecordSettlement(opts *bind.TransactOpts, _receivableId string) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.contract.Transact(opts, "recordSettlement", _receivableId)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactor) AsyncRecordSettlement(handler func(*types.Receipt, error), opts *bind.TransactOpts, _receivableId string) (*types.Transaction, error) {
	return _SupplyChainFinance.contract.AsyncTransact(opts, handler, "recordSettlement", _receivableId)
}

// RecordSettlement is a paid mutator transaction binding the contract method 0xcfff7343.
//
// Solidity: function recordSettlement(string _receivableId) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceSession) RecordSettlement(_receivableId string) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.RecordSettlement(&_SupplyChainFinance.TransactOpts, _receivableId)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) AsyncRecordSettlement(handler func(*types.Receipt, error), _receivableId string) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncRecordSettlement(handler, &_SupplyChainFinance.TransactOpts, _receivableId)
}

// RecordSettlement is a paid mutator transaction binding the contract method 0xcfff7343.
//
// Solidity: function recordSettlement(string _receivableId) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) RecordSettlement(_receivableId string) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.RecordSettlement(&_SupplyChainFinance.TransactOpts, _receivableId)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) AsyncRecordSettlement(handler func(*types.Receipt, error), _receivableId string) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncRecordSettlement(handler, &_SupplyChainFinance.TransactOpts, _receivableId)
}

// RegisterEnterprise is a paid mutator transaction binding the contract method 0x35238a92.
//
// Solidity: function registerEnterprise(string _enterpriseId, string _name, bool _isCore, bool _isFinancial) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactor) RegisterEnterprise(opts *bind.TransactOpts, _enterpriseId string, _name string, _isCore bool, _isFinancial bool) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.contract.Transact(opts, "registerEnterprise", _enterpriseId, _name, _isCore, _isFinancial)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactor) AsyncRegisterEnterprise(handler func(*types.Receipt, error), opts *bind.TransactOpts, _enterpriseId string, _name string, _isCore bool, _isFinancial bool) (*types.Transaction, error) {
	return _SupplyChainFinance.contract.AsyncTransact(opts, handler, "registerEnterprise", _enterpriseId, _name, _isCore, _isFinancial)
}

// RegisterEnterprise is a paid mutator transaction binding the contract method 0x35238a92.
//
// Solidity: function registerEnterprise(string _enterpriseId, string _name, bool _isCore, bool _isFinancial) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceSession) RegisterEnterprise(_enterpriseId string, _name string, _isCore bool, _isFinancial bool) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.RegisterEnterprise(&_SupplyChainFinance.TransactOpts, _enterpriseId, _name, _isCore, _isFinancial)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) AsyncRegisterEnterprise(handler func(*types.Receipt, error), _enterpriseId string, _name string, _isCore bool, _isFinancial bool) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncRegisterEnterprise(handler, &_SupplyChainFinance.TransactOpts, _enterpriseId, _name, _isCore, _isFinancial)
}

// RegisterEnterprise is a paid mutator transaction binding the contract method 0x35238a92.
//
// Solidity: function registerEnterprise(string _enterpriseId, string _name, bool _isCore, bool _isFinancial) returns(bool)
func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) RegisterEnterprise(_enterpriseId string, _name string, _isCore bool, _isFinancial bool) (*types.Transaction, *types.Receipt, error) {
	return _SupplyChainFinance.Contract.RegisterEnterprise(&_SupplyChainFinance.TransactOpts, _enterpriseId, _name, _isCore, _isFinancial)
}

func (_SupplyChainFinance *SupplyChainFinanceTransactorSession) AsyncRegisterEnterprise(handler func(*types.Receipt, error), _enterpriseId string, _name string, _isCore bool, _isFinancial bool) (*types.Transaction, error) {
	return _SupplyChainFinance.Contract.AsyncRegisterEnterprise(handler, &_SupplyChainFinance.TransactOpts, _enterpriseId, _name, _isCore, _isFinancial)
}

// SupplyChainFinanceEnterpriseRegistered represents a EnterpriseRegistered event raised by the SupplyChainFinance contract.
type SupplyChainFinanceEnterpriseRegistered struct {
	EnterpriseId string
	Name         string
	Wallet       common.Address
	IsCore       bool
	IsFinancial  bool
	Raw          types.Log // Blockchain specific contextual infos
}

// WatchEnterpriseRegistered is a free log subscription operation binding the contract event 0xc7eb38f697e467ff0b23344b8266c87010410b034a0568f86ea0ea562164953e.
//
// Solidity: event EnterpriseRegistered(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchEnterpriseRegistered(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "EnterpriseRegistered")
	return err
}

func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchAllEnterpriseRegistered(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "EnterpriseRegistered")
	return err
}

// ParseEnterpriseRegistered is a log parse operation binding the contract event 0xc7eb38f697e467ff0b23344b8266c87010410b034a0568f86ea0ea562164953e.
//
// Solidity: event EnterpriseRegistered(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) ParseEnterpriseRegistered(log types.Log) (*SupplyChainFinanceEnterpriseRegistered, error) {
	event := new(SupplyChainFinanceEnterpriseRegistered)
	if err := _SupplyChainFinance.contract.UnpackLog(event, "EnterpriseRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchEnterpriseRegistered is a free log subscription operation binding the contract event 0xc7eb38f697e467ff0b23344b8266c87010410b034a0568f86ea0ea562164953e.
//
// Solidity: event EnterpriseRegistered(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial)
func (_SupplyChainFinance *SupplyChainFinanceSession) WatchEnterpriseRegistered(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchEnterpriseRegistered(fromBlock, handler)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) WatchAllEnterpriseRegistered(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchAllEnterpriseRegistered(fromBlock, handler)
}

// ParseEnterpriseRegistered is a log parse operation binding the contract event 0xc7eb38f697e467ff0b23344b8266c87010410b034a0568f86ea0ea562164953e.
//
// Solidity: event EnterpriseRegistered(string enterpriseId, string name, address wallet, bool isCore, bool isFinancial)
func (_SupplyChainFinance *SupplyChainFinanceSession) ParseEnterpriseRegistered(log types.Log) (*SupplyChainFinanceEnterpriseRegistered, error) {
	return _SupplyChainFinance.Contract.ParseEnterpriseRegistered(log)
}

// SupplyChainFinanceReceivableConfirmed represents a ReceivableConfirmed event raised by the SupplyChainFinance contract.
type SupplyChainFinanceReceivableConfirmed struct {
	ReceivableId string
	Payer        common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// WatchReceivableConfirmed is a free log subscription operation binding the contract event 0x5b9df97ac7c8eebb0af305dbebe607ec1004a4f8ed24e14de349aeb2fa38fa09.
//
// Solidity: event ReceivableConfirmed(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchReceivableConfirmed(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableConfirmed")
	return err
}

func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchAllReceivableConfirmed(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableConfirmed")
	return err
}

// ParseReceivableConfirmed is a log parse operation binding the contract event 0x5b9df97ac7c8eebb0af305dbebe607ec1004a4f8ed24e14de349aeb2fa38fa09.
//
// Solidity: event ReceivableConfirmed(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) ParseReceivableConfirmed(log types.Log) (*SupplyChainFinanceReceivableConfirmed, error) {
	event := new(SupplyChainFinanceReceivableConfirmed)
	if err := _SupplyChainFinance.contract.UnpackLog(event, "ReceivableConfirmed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchReceivableConfirmed is a free log subscription operation binding the contract event 0x5b9df97ac7c8eebb0af305dbebe607ec1004a4f8ed24e14de349aeb2fa38fa09.
//
// Solidity: event ReceivableConfirmed(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceSession) WatchReceivableConfirmed(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchReceivableConfirmed(fromBlock, handler)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) WatchAllReceivableConfirmed(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchAllReceivableConfirmed(fromBlock, handler)
}

// ParseReceivableConfirmed is a log parse operation binding the contract event 0x5b9df97ac7c8eebb0af305dbebe607ec1004a4f8ed24e14de349aeb2fa38fa09.
//
// Solidity: event ReceivableConfirmed(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceSession) ParseReceivableConfirmed(log types.Log) (*SupplyChainFinanceReceivableConfirmed, error) {
	return _SupplyChainFinance.Contract.ParseReceivableConfirmed(log)
}

// SupplyChainFinanceReceivableFinanced represents a ReceivableFinanced event raised by the SupplyChainFinance contract.
type SupplyChainFinanceReceivableFinanced struct {
	ReceivableId  string
	FinancialInst common.Address
	Amount        *big.Int
	InterestRate  *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// WatchReceivableFinanced is a free log subscription operation binding the contract event 0xfc2e84e8b0698bc821cf46054ab7207d8e257345c13f85b68fc23ca931493f2b.
//
// Solidity: event ReceivableFinanced(string receivableId, address financialInst, uint256 amount, uint256 interestRate)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchReceivableFinanced(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableFinanced")
	return err
}

func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchAllReceivableFinanced(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableFinanced")
	return err
}

// ParseReceivableFinanced is a log parse operation binding the contract event 0xfc2e84e8b0698bc821cf46054ab7207d8e257345c13f85b68fc23ca931493f2b.
//
// Solidity: event ReceivableFinanced(string receivableId, address financialInst, uint256 amount, uint256 interestRate)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) ParseReceivableFinanced(log types.Log) (*SupplyChainFinanceReceivableFinanced, error) {
	event := new(SupplyChainFinanceReceivableFinanced)
	if err := _SupplyChainFinance.contract.UnpackLog(event, "ReceivableFinanced", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchReceivableFinanced is a free log subscription operation binding the contract event 0xfc2e84e8b0698bc821cf46054ab7207d8e257345c13f85b68fc23ca931493f2b.
//
// Solidity: event ReceivableFinanced(string receivableId, address financialInst, uint256 amount, uint256 interestRate)
func (_SupplyChainFinance *SupplyChainFinanceSession) WatchReceivableFinanced(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchReceivableFinanced(fromBlock, handler)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) WatchAllReceivableFinanced(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchAllReceivableFinanced(fromBlock, handler)
}

// ParseReceivableFinanced is a log parse operation binding the contract event 0xfc2e84e8b0698bc821cf46054ab7207d8e257345c13f85b68fc23ca931493f2b.
//
// Solidity: event ReceivableFinanced(string receivableId, address financialInst, uint256 amount, uint256 interestRate)
func (_SupplyChainFinance *SupplyChainFinanceSession) ParseReceivableFinanced(log types.Log) (*SupplyChainFinanceReceivableFinanced, error) {
	return _SupplyChainFinance.Contract.ParseReceivableFinanced(log)
}

// SupplyChainFinanceReceivableIssued represents a ReceivableIssued event raised by the SupplyChainFinance contract.
type SupplyChainFinanceReceivableIssued struct {
	ReceivableId string
	OrderId      string
	Issuer       common.Address
	Payer        common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// WatchReceivableIssued is a free log subscription operation binding the contract event 0x1c363599b5e17eb034389258d850faabbdefa9ae45f8f19e1e08a7c80339ac5c.
//
// Solidity: event ReceivableIssued(string receivableId, string orderId, address issuer, address payer, uint256 amount)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchReceivableIssued(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableIssued")
	return err
}

func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchAllReceivableIssued(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableIssued")
	return err
}

// ParseReceivableIssued is a log parse operation binding the contract event 0x1c363599b5e17eb034389258d850faabbdefa9ae45f8f19e1e08a7c80339ac5c.
//
// Solidity: event ReceivableIssued(string receivableId, string orderId, address issuer, address payer, uint256 amount)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) ParseReceivableIssued(log types.Log) (*SupplyChainFinanceReceivableIssued, error) {
	event := new(SupplyChainFinanceReceivableIssued)
	if err := _SupplyChainFinance.contract.UnpackLog(event, "ReceivableIssued", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchReceivableIssued is a free log subscription operation binding the contract event 0x1c363599b5e17eb034389258d850faabbdefa9ae45f8f19e1e08a7c80339ac5c.
//
// Solidity: event ReceivableIssued(string receivableId, string orderId, address issuer, address payer, uint256 amount)
func (_SupplyChainFinance *SupplyChainFinanceSession) WatchReceivableIssued(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchReceivableIssued(fromBlock, handler)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) WatchAllReceivableIssued(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchAllReceivableIssued(fromBlock, handler)
}

// ParseReceivableIssued is a log parse operation binding the contract event 0x1c363599b5e17eb034389258d850faabbdefa9ae45f8f19e1e08a7c80339ac5c.
//
// Solidity: event ReceivableIssued(string receivableId, string orderId, address issuer, address payer, uint256 amount)
func (_SupplyChainFinance *SupplyChainFinanceSession) ParseReceivableIssued(log types.Log) (*SupplyChainFinanceReceivableIssued, error) {
	return _SupplyChainFinance.Contract.ParseReceivableIssued(log)
}

// SupplyChainFinanceReceivableSettled represents a ReceivableSettled event raised by the SupplyChainFinance contract.
type SupplyChainFinanceReceivableSettled struct {
	ReceivableId string
	Payer        common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// WatchReceivableSettled is a free log subscription operation binding the contract event 0x8fde5a98eae57b20ae26cd4c7add1359fd043ce2a5a75c0240ac846c9dfd587e.
//
// Solidity: event ReceivableSettled(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchReceivableSettled(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableSettled")
	return err
}

func (_SupplyChainFinance *SupplyChainFinanceFilterer) WatchAllReceivableSettled(fromBlock *uint64, handler func(int, []types.Log)) error {
	_, err := _SupplyChainFinance.contract.WatchLogs(fromBlock, handler, "ReceivableSettled")
	return err
}

// ParseReceivableSettled is a log parse operation binding the contract event 0x8fde5a98eae57b20ae26cd4c7add1359fd043ce2a5a75c0240ac846c9dfd587e.
//
// Solidity: event ReceivableSettled(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceFilterer) ParseReceivableSettled(log types.Log) (*SupplyChainFinanceReceivableSettled, error) {
	event := new(SupplyChainFinanceReceivableSettled)
	if err := _SupplyChainFinance.contract.UnpackLog(event, "ReceivableSettled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchReceivableSettled is a free log subscription operation binding the contract event 0x8fde5a98eae57b20ae26cd4c7add1359fd043ce2a5a75c0240ac846c9dfd587e.
//
// Solidity: event ReceivableSettled(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceSession) WatchReceivableSettled(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchReceivableSettled(fromBlock, handler)
}

func (_SupplyChainFinance *SupplyChainFinanceSession) WatchAllReceivableSettled(fromBlock *uint64, handler func(int, []types.Log)) error {
	return _SupplyChainFinance.Contract.WatchAllReceivableSettled(fromBlock, handler)
}

// ParseReceivableSettled is a log parse operation binding the contract event 0x8fde5a98eae57b20ae26cd4c7add1359fd043ce2a5a75c0240ac846c9dfd587e.
//
// Solidity: event ReceivableSettled(string receivableId, address payer)
func (_SupplyChainFinance *SupplyChainFinanceSession) ParseReceivableSettled(log types.Log) (*SupplyChainFinanceReceivableSettled, error) {
	return _SupplyChainFinance.Contract.ParseReceivableSettled(log)
}
