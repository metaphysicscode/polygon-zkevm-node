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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"depositAmounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"depositOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"punish\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"punishAmounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"expect\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"real\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISlotAdapter\",\"name\":\"_slotAdapter\",\"type\":\"address\"}],\"name\":\"setSlotAdapter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slotAdapter\",\"outputs\":[{\"internalType\":\"contractISlotAdapter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDeposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611356806100206000396000f3fe6080604052600436106100d25760003560e01c80638da5cb5b1161007f578063d0e30db011610059578063d0e30db014610259578063d8e423b214610261578063e476f9ff14610281578063f2fde38b146102ca57600080fd5b80638da5cb5b146101e1578063a62125c91461020c578063bf0294d01461022c57600080fd5b8063715018a6116100b0578063715018a6146101a15780637d882097146101b65780638129fc1c146101cc57600080fd5b80631c52e346146100d757806323e3fbd51461012e5780632e1a7d4d1461017f575b600080fd5b3480156100e357600080fd5b506066546101049073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b34801561013a57600080fd5b506101716101493660046111f7565b73ffffffffffffffffffffffffffffffffffffffff1660009081526067602052604090205490565b604051908152602001610125565b34801561018b57600080fd5b5061019f61019a36600461121b565b6102ea565b005b3480156101ad57600080fd5b5061019f6107ee565b3480156101c257600080fd5b5061017160655481565b3480156101d857600080fd5b5061019f610802565b3480156101ed57600080fd5b5060335473ffffffffffffffffffffffffffffffffffffffff16610104565b34801561021857600080fd5b5061019f6102273660046111f7565b610994565b34801561023857600080fd5b506101716102473660046111f7565b60676020526000908152604090205481565b61019f610a61565b34801561026d57600080fd5b5061019f61027c366004611234565b610dc5565b34801561028d57600080fd5b506102b561029c3660046111f7565b6068602052600090815260409020805460019091015482565b60408051928352602083019190915201610125565b3480156102d657600080fd5b5061019f6102e53660046111f7565b610f89565b60695415610359576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600660248201527f4c4f434b4544000000000000000000000000000000000000000000000000000060448201526064015b60405180910390fd5b6001606955333b156103c7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f4163636f756e74206e6f7420454f4100000000000000000000000000000000006044820152606401610350565b606654604080517feb428505000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163eb428505916004808301926020929190829003018187875af1158015610438573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061045c9190611260565b6040517f6a256b2900000000000000000000000000000000000000000000000000000000815233600482015290915073ffffffffffffffffffffffffffffffffffffffff821690636a256b2990602401600060405180830381600087803b1580156104c657600080fd5b505af11580156104da573d6000803e3d6000fd5b505033600090815260676020526040902054915050828110156105025750600091508161050f565b61050c83826112ac565b90505b33600090815260676020526040902081905560655461052f9084906112ac565b606555604080516000808252602082019092523390859060405161055391906112c5565b60006040518083038185875af1925050503d8060008114610590576040519150601f19603f3d011682016040523d82523d6000602084013e610595565b606091505b5050905080610600576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f7769746864726177616c3a207472616e73666572206661696c656400000000006044820152606401610350565b6066546040517f550f283100000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091169063550f283190602401602060405180830381865afa15801561066e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061069291906112f4565b606660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166365868d606040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610701573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061072591906112f4565b036107ae576066546040517f5748a7e200000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911690635748a7e290602401600060405180830381600087803b15801561079557600080fd5b505af11580156107a9573d6000803e3d6000fd5b505050505b60405184815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a2505060006069555050565b6107f661103d565b61080060006110be565b565b600054610100900460ff16158080156108225750600054600160ff909116105b8061083c5750303b15801561083c575060005460ff166001145b6108c8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610350565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561092657600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b61092e611135565b801561099157600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b61099c61103d565b73ffffffffffffffffffffffffffffffffffffffff81163b610a1a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600b60248201527f4163636f756e7420454f410000000000000000000000000000000000000000006044820152606401610350565b606680547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60695415610acb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600660248201527f4c4f434b454400000000000000000000000000000000000000000000000000006044820152606401610350565b600160695534610b39576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103509060208082526004908201527f7a65726f00000000000000000000000000000000000000000000000000000000604082015260600190565b6066546040517f550f283100000000000000000000000000000000000000000000000000000000815233600482015260009173ffffffffffffffffffffffffffffffffffffffff169063550f283190602401602060405180830381865afa158015610ba8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bcc91906112f4565b90508015610cd057606654604080517f65868d600000000000000000000000000000000000000000000000000000000081529051839273ffffffffffffffffffffffffffffffffffffffff16916365868d6091600480830192602092919082900301816000875af1158015610c45573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c6991906112f4565b14610cd0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f696e76616c6964206465706f73697400000000000000000000000000000000006044820152606401610350565b34606554610cde919061130d565b6065553360009081526067602052604081208054349290610d0090849061130d565b90915550506066546040517f04b99a8900000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff909116906304b99a8990602401600060405180830381600087803b158015610d7057600080fd5b505af1158015610d84573d6000803e3d6000fd5b50506040513481523392507fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c915060200160405180910390a2506000606955565b60695415610e2f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600660248201527f4c4f434b454400000000000000000000000000000000000000000000000000006044820152606401610350565b600160695560665473ffffffffffffffffffffffffffffffffffffffff163314610eb5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600f60248201527f4e6f7420736c6f744164617074657200000000000000000000000000000000006044820152606401610350565b73ffffffffffffffffffffffffffffffffffffffff82166000908152606760205260409020548180821015610eed5750600090610efa565b610ef783836112ac565b91505b73ffffffffffffffffffffffffffffffffffffffff84166000908152606760209081526040808320859055606890915281206001018054839290610f3f90849061130d565b909155505073ffffffffffffffffffffffffffffffffffffffff841660009081526068602052604081208054859290610f7990849061130d565b9091555050600060695550505050565b610f9161103d565b73ffffffffffffffffffffffffffffffffffffffff8116611034576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610350565b610991816110be565b60335473ffffffffffffffffffffffffffffffffffffffff163314610800576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610350565b6033805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166111cc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610350565b610800336110be565b73ffffffffffffffffffffffffffffffffffffffff8116811461099157600080fd5b60006020828403121561120957600080fd5b8135611214816111d5565b9392505050565b60006020828403121561122d57600080fd5b5035919050565b6000806040838503121561124757600080fd5b8235611252816111d5565b946020939093013593505050565b60006020828403121561127257600080fd5b8151611214816111d5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156112bf576112bf61127d565b92915050565b6000825160005b818110156112e657602081860181015185830152016112cc565b506000920191825250919050565b60006020828403121561130657600080fd5b5051919050565b808201808211156112bf576112bf61127d56fea2646970667358221220882e32192d145f58d132a339416155b790dc491deceed8d3fe654b53a837782a64736f6c63430008110033",
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
// Solidity: function punishAmounts(address ) view returns(uint256 expect, uint256 real)
func (_Deposit *DepositCaller) PunishAmounts(opts *bind.CallOpts, arg0 common.Address) (struct {
	Expect *big.Int
	Real   *big.Int
}, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "punishAmounts", arg0)

	outstruct := new(struct {
		Expect *big.Int
		Real   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Expect = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Real = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PunishAmounts is a free data retrieval call binding the contract method 0xe476f9ff.
//
// Solidity: function punishAmounts(address ) view returns(uint256 expect, uint256 real)
func (_Deposit *DepositSession) PunishAmounts(arg0 common.Address) (struct {
	Expect *big.Int
	Real   *big.Int
}, error) {
	return _Deposit.Contract.PunishAmounts(&_Deposit.CallOpts, arg0)
}

// PunishAmounts is a free data retrieval call binding the contract method 0xe476f9ff.
//
// Solidity: function punishAmounts(address ) view returns(uint256 expect, uint256 real)
func (_Deposit *DepositCallerSession) PunishAmounts(arg0 common.Address) (struct {
	Expect *big.Int
	Real   *big.Int
}, error) {
	return _Deposit.Contract.PunishAmounts(&_Deposit.CallOpts, arg0)
}

// SlotAdapter is a free data retrieval call binding the contract method 0x1c52e346.
//
// Solidity: function slotAdapter() view returns(address)
func (_Deposit *DepositCaller) SlotAdapter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "slotAdapter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SlotAdapter is a free data retrieval call binding the contract method 0x1c52e346.
//
// Solidity: function slotAdapter() view returns(address)
func (_Deposit *DepositSession) SlotAdapter() (common.Address, error) {
	return _Deposit.Contract.SlotAdapter(&_Deposit.CallOpts)
}

// SlotAdapter is a free data retrieval call binding the contract method 0x1c52e346.
//
// Solidity: function slotAdapter() view returns(address)
func (_Deposit *DepositCallerSession) SlotAdapter() (common.Address, error) {
	return _Deposit.Contract.SlotAdapter(&_Deposit.CallOpts)
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
