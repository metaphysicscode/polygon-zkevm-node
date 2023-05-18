// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package deposit

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

// DepositMetaData contains all meta data concerning the Deposit contract.
var DepositMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"depositAmounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"depositOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"punish\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"punishAmounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_slotAdapter\",\"type\":\"address\"}],\"name\":\"setSlotAdapter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"slotAdapters\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDeposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506116fb806100206000396000f3fe6080604052600436106100c25760003560e01c8063a62125c91161007f578063d2d6551111610059578063d2d6551114610221578063d8e423b21461025e578063e476f9ff14610287578063f2fde38b146102c4576100c2565b8063a62125c9146101b1578063bf0294d0146101da578063d0e30db014610217576100c2565b806323e3fbd5146100c75780632e1a7d4d14610104578063715018a61461012d5780637d882097146101445780638129fc1c1461016f5780638da5cb5b14610186575b600080fd5b3480156100d357600080fd5b506100ee60048036038101906100e99190610e89565b6102ed565b6040516100fb9190610ecf565b60405180910390f35b34801561011057600080fd5b5061012b60048036038101906101269190610f16565b610336565b005b34801561013957600080fd5b50610142610650565b005b34801561015057600080fd5b50610159610664565b6040516101669190610ecf565b60405180910390f35b34801561017b57600080fd5b5061018461066a565b005b34801561019257600080fd5b5061019b6107a8565b6040516101a89190610f52565b60405180910390f35b3480156101bd57600080fd5b506101d860048036038101906101d39190610e89565b6107d2565b005b3480156101e657600080fd5b5061020160048036038101906101fc9190610e89565b61087d565b60405161020e9190610ecf565b60405180910390f35b61021f610895565b005b34801561022d57600080fd5b5061024860048036038101906102439190610e89565b6109e7565b6040516102559190610f88565b60405180910390f35b34801561026a57600080fd5b5061028560048036038101906102809190610fa3565b610a07565b005b34801561029357600080fd5b506102ae60048036038101906102a99190610e89565b610b98565b6040516102bb9190610ecf565b60405180910390f35b3480156102d057600080fd5b506102eb60048036038101906102e69190610e89565b610bb0565b005b6000606760008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60006069541461037b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161037290611040565b60405180910390fd5b6001606981905550600081116103c6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103bd906110ac565b60405180910390fd5b606760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054811115610448576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161043f90611118565b60405180910390fd5b61045133610c33565b15610491576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161048890611184565b60405180910390fd5b8060655461049f91906111d3565b60658190555080606760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546104f491906111d3565b9250508190555060003373ffffffffffffffffffffffffffffffffffffffff1682600067ffffffffffffffff8111156105305761052f611207565b5b6040519080825280601f01601f1916602001820160405280156105625781602001600182028036833780820191505090505b5060405161057091906112a7565b60006040518083038185875af1925050503d80600081146105ad576040519150601f19603f3d011682016040523d82523d6000602084013e6105b2565b606091505b50509050806105f6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105ed9061130a565b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff167f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243648360405161063c9190610ecf565b60405180910390a250600060698190555050565b610658610c56565b6106626000610cd4565b565b60655481565b60008060019054906101000a900460ff1615905080801561069b5750600160008054906101000a900460ff1660ff16105b806106c857506106aa30610d9a565b1580156106c75750600160008054906101000a900460ff1660ff16145b5b610707576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106fe9061139c565b60405180910390fd5b60016000806101000a81548160ff021916908360ff1602179055508015610744576001600060016101000a81548160ff0219169083151502179055505b61074c610dbd565b80156107a55760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498600160405161079c919061140e565b60405180910390a15b50565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6107da610c56565b6107e381610c33565b610822576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161081990611475565b60405180910390fd5b6001606660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b60676020528060005260406000206000915090505481565b6000606954146108da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108d190611040565b60405180910390fd5b600160698190555060003411610925576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161091c906110ac565b60405180910390fd5b346065546109339190611495565b60658190555034606760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109889190611495565b925050819055503373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c346040516109d59190610ecf565b60405180910390a26000606981905550565b60666020528060005260406000206000915054906101000a900460ff1681565b600060695414610a4c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a4390611040565b60405180910390fd5b6001606981905550606660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610ae0576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ad790611515565b60405180910390fd5b80606760008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610b2f91906111d3565b9250508190555080606860008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610b859190611495565b9250508190555060006069819055505050565b60686020528060005260406000206000915090505481565b610bb8610c56565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610c27576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c1e906115a7565b60405180910390fd5b610c3081610cd4565b50565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b610c5e610e1e565b73ffffffffffffffffffffffffffffffffffffffff16610c7c6107a8565b73ffffffffffffffffffffffffffffffffffffffff1614610cd2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cc990611613565b60405180910390fd5b565b6000603360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905081603360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b600060019054906101000a900460ff16610e0c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e03906116a5565b60405180910390fd5b610e1c610e17610e1e565b610cd4565b565b600033905090565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610e5682610e2b565b9050919050565b610e6681610e4b565b8114610e7157600080fd5b50565b600081359050610e8381610e5d565b92915050565b600060208284031215610e9f57610e9e610e26565b5b6000610ead84828501610e74565b91505092915050565b6000819050919050565b610ec981610eb6565b82525050565b6000602082019050610ee46000830184610ec0565b92915050565b610ef381610eb6565b8114610efe57600080fd5b50565b600081359050610f1081610eea565b92915050565b600060208284031215610f2c57610f2b610e26565b5b6000610f3a84828501610f01565b91505092915050565b610f4c81610e4b565b82525050565b6000602082019050610f676000830184610f43565b92915050565b60008115159050919050565b610f8281610f6d565b82525050565b6000602082019050610f9d6000830184610f79565b92915050565b60008060408385031215610fba57610fb9610e26565b5b6000610fc885828601610e74565b9250506020610fd985828601610f01565b9150509250929050565b600082825260208201905092915050565b7f4c4f434b45440000000000000000000000000000000000000000000000000000600082015250565b600061102a600683610fe3565b915061103582610ff4565b602082019050919050565b600060208201905081810360008301526110598161101d565b9050919050565b7f7a65726f00000000000000000000000000000000000000000000000000000000600082015250565b6000611096600483610fe3565b91506110a182611060565b602082019050919050565b600060208201905081810360008301526110c581611089565b9050919050565b7f496e76616c696420616d6f756e74000000000000000000000000000000000000600082015250565b6000611102600e83610fe3565b915061110d826110cc565b602082019050919050565b60006020820190508181036000830152611131816110f5565b9050919050565b7f4163636f756e74206e6f7420454f410000000000000000000000000000000000600082015250565b600061116e600f83610fe3565b915061117982611138565b602082019050919050565b6000602082019050818103600083015261119d81611161565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006111de82610eb6565b91506111e983610eb6565b9250828203905081811115611201576112006111a4565b5b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600081519050919050565b600081905092915050565b60005b8381101561126a57808201518184015260208101905061124f565b60008484015250505050565b600061128182611236565b61128b8185611241565b935061129b81856020860161124c565b80840191505092915050565b60006112b38284611276565b915081905092915050565b7f7769746864726177616c3a207472616e73666572206661696c65640000000000600082015250565b60006112f4601b83610fe3565b91506112ff826112be565b602082019050919050565b60006020820190508181036000830152611323816112e7565b9050919050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b6000611386602e83610fe3565b91506113918261132a565b604082019050919050565b600060208201905081810360008301526113b581611379565b9050919050565b6000819050919050565b600060ff82169050919050565b6000819050919050565b60006113f86113f36113ee846113bc565b6113d3565b6113c6565b9050919050565b611408816113dd565b82525050565b600060208201905061142360008301846113ff565b92915050565b7f4163636f756e7420454f41000000000000000000000000000000000000000000600082015250565b600061145f600b83610fe3565b915061146a82611429565b602082019050919050565b6000602082019050818103600083015261148e81611452565b9050919050565b60006114a082610eb6565b91506114ab83610eb6565b92508282019050808211156114c3576114c26111a4565b5b92915050565b7f4e6f7420736c6f74416461707465720000000000000000000000000000000000600082015250565b60006114ff600f83610fe3565b915061150a826114c9565b602082019050919050565b6000602082019050818103600083015261152e816114f2565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000611591602683610fe3565b915061159c82611535565b604082019050919050565b600060208201905081810360008301526115c081611584565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006115fd602083610fe3565b9150611608826115c7565b602082019050919050565b6000602082019050818103600083015261162c816115f0565b9050919050565b7f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960008201527f6e697469616c697a696e67000000000000000000000000000000000000000000602082015250565b600061168f602b83610fe3565b915061169a82611633565b604082019050919050565b600060208201905081810360008301526116be81611682565b905091905056fea26469706673582212207cbc2b127e8cce29303cc628e63ac9f25920e13e9facdcff92d71c7be23c00a564736f6c63430008120033",
}

// DepositABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositMetaData.ABI instead.
var DepositABI = DepositMetaData.ABI

// DepositBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepositMetaData.Bin instead.
var DepositBin = DepositMetaData.Bin

// DeployDeposit deploys a new Ethereum contract, binding an instance of Deposit to it.
func DeployDeposit(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Deposit, error) {
	parsed, err := DepositMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepositBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Deposit{DepositCaller: DepositCaller{contract: contract}, DepositTransactor: DepositTransactor{contract: contract}, DepositFilterer: DepositFilterer{contract: contract}}, nil
}

// Deposit is an auto generated Go binding around an Ethereum contract.
type Deposit struct {
	DepositCaller     // Read-only binding to the contract
	DepositTransactor // Write-only binding to the contract
	DepositFilterer   // Log filterer for contract events
}

// DepositCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositSession struct {
	Contract     *Deposit          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositCallerSession struct {
	Contract *DepositCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// DepositTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositTransactorSession struct {
	Contract     *DepositTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DepositRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositRaw struct {
	Contract *Deposit // Generic contract binding to access the raw methods on
}

// DepositCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositCallerRaw struct {
	Contract *DepositCaller // Generic read-only contract binding to access the raw methods on
}

// DepositTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositTransactorRaw struct {
	Contract *DepositTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeposit creates a new instance of Deposit, bound to a specific deployed contract.
func NewDeposit(address common.Address, backend bind.ContractBackend) (*Deposit, error) {
	contract, err := bindDeposit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Deposit{DepositCaller: DepositCaller{contract: contract}, DepositTransactor: DepositTransactor{contract: contract}, DepositFilterer: DepositFilterer{contract: contract}}, nil
}

// NewDepositCaller creates a new read-only instance of Deposit, bound to a specific deployed contract.
func NewDepositCaller(address common.Address, caller bind.ContractCaller) (*DepositCaller, error) {
	contract, err := bindDeposit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCaller{contract: contract}, nil
}

// NewDepositTransactor creates a new write-only instance of Deposit, bound to a specific deployed contract.
func NewDepositTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositTransactor, error) {
	contract, err := bindDeposit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositTransactor{contract: contract}, nil
}

// NewDepositFilterer creates a new log filterer instance of Deposit, bound to a specific deployed contract.
func NewDepositFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositFilterer, error) {
	contract, err := bindDeposit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositFilterer{contract: contract}, nil
}

// bindDeposit binds a generic wrapper to an already deployed contract.
func bindDeposit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.DepositCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transact(opts, method, params...)
}

// DepositAmounts is a free data retrieval call binding the contract method 0xbf0294d0.
//
// Solidity: function depositAmounts(address ) view returns(uint256)
func (_Deposit *DepositCaller) DepositAmounts(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "depositAmounts", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositAmounts is a free data retrieval call binding the contract method 0xbf0294d0.
//
// Solidity: function depositAmounts(address ) view returns(uint256)
func (_Deposit *DepositSession) DepositAmounts(arg0 common.Address) (*big.Int, error) {
	return _Deposit.Contract.DepositAmounts(&_Deposit.CallOpts, arg0)
}

// DepositAmounts is a free data retrieval call binding the contract method 0xbf0294d0.
//
// Solidity: function depositAmounts(address ) view returns(uint256)
func (_Deposit *DepositCallerSession) DepositAmounts(arg0 common.Address) (*big.Int, error) {
	return _Deposit.Contract.DepositAmounts(&_Deposit.CallOpts, arg0)
}

// DepositOf is a free data retrieval call binding the contract method 0x23e3fbd5.
//
// Solidity: function depositOf(address account) view returns(uint256)
func (_Deposit *DepositCaller) DepositOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "depositOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositOf is a free data retrieval call binding the contract method 0x23e3fbd5.
//
// Solidity: function depositOf(address account) view returns(uint256)
func (_Deposit *DepositSession) DepositOf(account common.Address) (*big.Int, error) {
	return _Deposit.Contract.DepositOf(&_Deposit.CallOpts, account)
}

// DepositOf is a free data retrieval call binding the contract method 0x23e3fbd5.
//
// Solidity: function depositOf(address account) view returns(uint256)
func (_Deposit *DepositCallerSession) DepositOf(account common.Address) (*big.Int, error) {
	return _Deposit.Contract.DepositOf(&_Deposit.CallOpts, account)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositSession) Owner() (common.Address, error) {
	return _Deposit.Contract.Owner(&_Deposit.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositCallerSession) Owner() (common.Address, error) {
	return _Deposit.Contract.Owner(&_Deposit.CallOpts)
}

// PunishAmounts is a free data retrieval call binding the contract method 0xe476f9ff.
//
// Solidity: function punishAmounts(address ) view returns(uint256)
func (_Deposit *DepositCaller) PunishAmounts(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "punishAmounts", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PunishAmounts is a free data retrieval call binding the contract method 0xe476f9ff.
//
// Solidity: function punishAmounts(address ) view returns(uint256)
func (_Deposit *DepositSession) PunishAmounts(arg0 common.Address) (*big.Int, error) {
	return _Deposit.Contract.PunishAmounts(&_Deposit.CallOpts, arg0)
}

// PunishAmounts is a free data retrieval call binding the contract method 0xe476f9ff.
//
// Solidity: function punishAmounts(address ) view returns(uint256)
func (_Deposit *DepositCallerSession) PunishAmounts(arg0 common.Address) (*big.Int, error) {
	return _Deposit.Contract.PunishAmounts(&_Deposit.CallOpts, arg0)
}

// SlotAdapters is a free data retrieval call binding the contract method 0xd2d65511.
//
// Solidity: function slotAdapters(address ) view returns(bool)
func (_Deposit *DepositCaller) SlotAdapters(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "slotAdapters", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SlotAdapters is a free data retrieval call binding the contract method 0xd2d65511.
//
// Solidity: function slotAdapters(address ) view returns(bool)
func (_Deposit *DepositSession) SlotAdapters(arg0 common.Address) (bool, error) {
	return _Deposit.Contract.SlotAdapters(&_Deposit.CallOpts, arg0)
}

// SlotAdapters is a free data retrieval call binding the contract method 0xd2d65511.
//
// Solidity: function slotAdapters(address ) view returns(bool)
func (_Deposit *DepositCallerSession) SlotAdapters(arg0 common.Address) (bool, error) {
	return _Deposit.Contract.SlotAdapters(&_Deposit.CallOpts, arg0)
}

// TotalDeposits is a free data retrieval call binding the contract method 0x7d882097.
//
// Solidity: function totalDeposits() view returns(uint256)
func (_Deposit *DepositCaller) TotalDeposits(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "totalDeposits")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalDeposits is a free data retrieval call binding the contract method 0x7d882097.
//
// Solidity: function totalDeposits() view returns(uint256)
func (_Deposit *DepositSession) TotalDeposits() (*big.Int, error) {
	return _Deposit.Contract.TotalDeposits(&_Deposit.CallOpts)
}

// TotalDeposits is a free data retrieval call binding the contract method 0x7d882097.
//
// Solidity: function totalDeposits() view returns(uint256)
func (_Deposit *DepositCallerSession) TotalDeposits() (*big.Int, error) {
	return _Deposit.Contract.TotalDeposits(&_Deposit.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Deposit *DepositTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Deposit *DepositSession) Deposit() (*types.Transaction, error) {
	return _Deposit.Contract.Deposit(&_Deposit.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Deposit *DepositTransactorSession) Deposit() (*types.Transaction, error) {
	return _Deposit.Contract.Deposit(&_Deposit.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Deposit *DepositTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Deposit *DepositSession) Initialize() (*types.Transaction, error) {
	return _Deposit.Contract.Initialize(&_Deposit.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Deposit *DepositTransactorSession) Initialize() (*types.Transaction, error) {
	return _Deposit.Contract.Initialize(&_Deposit.TransactOpts)
}

// Punish is a paid mutator transaction binding the contract method 0xd8e423b2.
//
// Solidity: function punish(address account, uint256 amount) returns()
func (_Deposit *DepositTransactor) Punish(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "punish", account, amount)
}

// Punish is a paid mutator transaction binding the contract method 0xd8e423b2.
//
// Solidity: function punish(address account, uint256 amount) returns()
func (_Deposit *DepositSession) Punish(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Punish(&_Deposit.TransactOpts, account, amount)
}

// Punish is a paid mutator transaction binding the contract method 0xd8e423b2.
//
// Solidity: function punish(address account, uint256 amount) returns()
func (_Deposit *DepositTransactorSession) Punish(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Punish(&_Deposit.TransactOpts, account, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deposit *DepositTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deposit *DepositSession) RenounceOwnership() (*types.Transaction, error) {
	return _Deposit.Contract.RenounceOwnership(&_Deposit.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deposit *DepositTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Deposit.Contract.RenounceOwnership(&_Deposit.TransactOpts)
}

// SetSlotAdapter is a paid mutator transaction binding the contract method 0xa62125c9.
//
// Solidity: function setSlotAdapter(address _slotAdapter) returns()
func (_Deposit *DepositTransactor) SetSlotAdapter(opts *bind.TransactOpts, _slotAdapter common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "setSlotAdapter", _slotAdapter)
}

// SetSlotAdapter is a paid mutator transaction binding the contract method 0xa62125c9.
//
// Solidity: function setSlotAdapter(address _slotAdapter) returns()
func (_Deposit *DepositSession) SetSlotAdapter(_slotAdapter common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetSlotAdapter(&_Deposit.TransactOpts, _slotAdapter)
}

// SetSlotAdapter is a paid mutator transaction binding the contract method 0xa62125c9.
//
// Solidity: function setSlotAdapter(address _slotAdapter) returns()
func (_Deposit *DepositTransactorSession) SetSlotAdapter(_slotAdapter common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetSlotAdapter(&_Deposit.TransactOpts, _slotAdapter)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deposit *DepositTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deposit *DepositSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.TransferOwnership(&_Deposit.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deposit *DepositTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.TransferOwnership(&_Deposit.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Deposit *DepositTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Deposit *DepositSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Withdraw(&_Deposit.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Deposit *DepositTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Withdraw(&_Deposit.TransactOpts, amount)
}

// DepositDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Deposit contract.
type DepositDepositIterator struct {
	Event *DepositDeposit // Event containing the contract specifics and raw log

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
func (it *DepositDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositDeposit)
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
		it.Event = new(DepositDeposit)
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
func (it *DepositDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositDeposit represents a Deposit event raised by the Deposit contract.
type DepositDeposit struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) FilterDeposit(opts *bind.FilterOpts, user []common.Address) (*DepositDepositIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "Deposit", userRule)
	if err != nil {
		return nil, err
	}
	return &DepositDepositIterator{contract: _Deposit.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *DepositDeposit, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "Deposit", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositDeposit)
				if err := _Deposit.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) ParseDeposit(log types.Log) (*DepositDeposit, error) {
	event := new(DepositDeposit)
	if err := _Deposit.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Deposit contract.
type DepositInitializedIterator struct {
	Event *DepositInitialized // Event containing the contract specifics and raw log

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
func (it *DepositInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositInitialized)
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
		it.Event = new(DepositInitialized)
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
func (it *DepositInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositInitialized represents a Initialized event raised by the Deposit contract.
type DepositInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Deposit *DepositFilterer) FilterInitialized(opts *bind.FilterOpts) (*DepositInitializedIterator, error) {

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DepositInitializedIterator{contract: _Deposit.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Deposit *DepositFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DepositInitialized) (event.Subscription, error) {

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositInitialized)
				if err := _Deposit.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Deposit *DepositFilterer) ParseInitialized(log types.Log) (*DepositInitialized, error) {
	event := new(DepositInitialized)
	if err := _Deposit.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Deposit contract.
type DepositOwnershipTransferredIterator struct {
	Event *DepositOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DepositOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositOwnershipTransferred)
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
		it.Event = new(DepositOwnershipTransferred)
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
func (it *DepositOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositOwnershipTransferred represents a OwnershipTransferred event raised by the Deposit contract.
type DepositOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Deposit *DepositFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DepositOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DepositOwnershipTransferredIterator{contract: _Deposit.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Deposit *DepositFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DepositOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositOwnershipTransferred)
				if err := _Deposit.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Deposit *DepositFilterer) ParseOwnershipTransferred(log types.Log) (*DepositOwnershipTransferred, error) {
	event := new(DepositOwnershipTransferred)
	if err := _Deposit.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Deposit contract.
type DepositWithdrawIterator struct {
	Event *DepositWithdraw // Event containing the contract specifics and raw log

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
func (it *DepositWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositWithdraw)
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
		it.Event = new(DepositWithdraw)
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
func (it *DepositWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositWithdraw represents a Withdraw event raised by the Deposit contract.
type DepositWithdraw struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) FilterWithdraw(opts *bind.FilterOpts, user []common.Address) (*DepositWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "Withdraw", userRule)
	if err != nil {
		return nil, err
	}
	return &DepositWithdrawIterator{contract: _Deposit.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *DepositWithdraw, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "Withdraw", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositWithdraw)
				if err := _Deposit.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) ParseWithdraw(log types.Log) (*DepositWithdraw, error) {
	event := new(DepositWithdraw)
	if err := _Deposit.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
