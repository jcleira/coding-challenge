// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	aggregates "github.com/jcleira/coding-challenge/internal/domain/aggregates"
	mock "github.com/stretchr/testify/mock"
)

// ExchangeGetter is an autogenerated mock type for the ExchangeGetter type
type ExchangeGetter struct {
	mock.Mock
}

// GetRate provides a mock function with given fields:
func (_m *ExchangeGetter) GetRate() (aggregates.Rate, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRate")
	}

	var r0 aggregates.Rate
	var r1 error
	if rf, ok := ret.Get(0).(func() (aggregates.Rate, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() aggregates.Rate); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(aggregates.Rate)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExchangeGetter creates a new instance of ExchangeGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExchangeGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExchangeGetter {
	mock := &ExchangeGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}