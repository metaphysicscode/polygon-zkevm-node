// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"

	types "github.com/0xPolygonHermez/zkevm-node/etherman/types"
)

// Etherman is an autogenerated mock type for the etherman type
type Etherman struct {
	mock.Mock
}

func (_m *Etherman)BuildProofHashTxData(lastVerifiedBatch, newVerifiedBatch uint64, proofHash common.Hash) (to *common.Address, data []byte, err error) {
	ret := _m.Called(lastVerifiedBatch, newVerifiedBatch, proofHash)

	var r0 *common.Address
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func(uint64, uint64, common.Hash) (*common.Address, []byte, error)); ok {
		return rf(lastVerifiedBatch, newVerifiedBatch, proofHash)
	}
	if rf, ok := ret.Get(0).(func(uint64, uint64, common.Hash) *common.Address); ok {
		r0 = rf(lastVerifiedBatch, newVerifiedBatch, proofHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, uint64, common.Hash) []byte); ok {
		r1 = rf(lastVerifiedBatch, newVerifiedBatch, proofHash)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func(uint64, uint64, common.Hash) error); ok {
		r2 = rf(lastVerifiedBatch, newVerifiedBatch, proofHash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// BuildTrustedVerifyBatchesTxData provides a mock function with given fields: lastVerifiedBatch, newVerifiedBatch, inputs
func (_m *Etherman) BuildTrustedVerifyBatchesTxData(lastVerifiedBatch uint64, newVerifiedBatch uint64, inputs *types.FinalProofInputs) (*common.Address, []byte, error) {
	ret := _m.Called(lastVerifiedBatch, newVerifiedBatch, inputs)

	var r0 *common.Address
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func(uint64, uint64, *types.FinalProofInputs) (*common.Address, []byte, error)); ok {
		return rf(lastVerifiedBatch, newVerifiedBatch, inputs)
	}
	if rf, ok := ret.Get(0).(func(uint64, uint64, *types.FinalProofInputs) *common.Address); ok {
		r0 = rf(lastVerifiedBatch, newVerifiedBatch, inputs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, uint64, *types.FinalProofInputs) []byte); ok {
		r1 = rf(lastVerifiedBatch, newVerifiedBatch, inputs)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func(uint64, uint64, *types.FinalProofInputs) error); ok {
		r2 = rf(lastVerifiedBatch, newVerifiedBatch, inputs)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetLatestVerifiedBatchNum provides a mock function with given fields:
func (_m *Etherman) GetLatestVerifiedBatchNum() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func() (uint64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEtherman interface {
	mock.TestingT
	Cleanup(func())
}

// NewEtherman creates a new instance of Etherman. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEtherman(t mockConstructorTestingTNewEtherman) *Etherman {
	mock := &Etherman{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
