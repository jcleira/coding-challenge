// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// WalletInitializer is an autogenerated mock type for the WalletInitializer type
type WalletInitializer struct {
	mock.Mock
}

// Initialize provides a mock function with given fields:
func (_m *WalletInitializer) Initialize() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Initialize")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWalletInitializer creates a new instance of WalletInitializer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWalletInitializer(t interface {
	mock.TestingT
	Cleanup(func())
}) *WalletInitializer {
	mock := &WalletInitializer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
