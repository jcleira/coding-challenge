// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ExchangeRateGetter is an autogenerated mock type for the ExchangeRateGetter type
type ExchangeRateGetter struct {
	mock.Mock
}

// GetRate provides a mock function with given fields:
func (_m *ExchangeRateGetter) GetRate() (float64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRate")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func() (float64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExchangeRateGetter creates a new instance of ExchangeRateGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExchangeRateGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExchangeRateGetter {
	mock := &ExchangeRateGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
