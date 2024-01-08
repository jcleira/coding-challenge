// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	aggregates "github.com/jcleira/coding-challenge/internal/domain/aggregates"

	mock "github.com/stretchr/testify/mock"
)

// TransactionsGetter is an autogenerated mock type for the TransactionsGetter type
type TransactionsGetter struct {
	mock.Mock
}

// GetTransactions provides a mock function with given fields: ctx, publicKey
func (_m *TransactionsGetter) GetTransactions(ctx context.Context, publicKey string) ([]aggregates.Transaction, error) {
	ret := _m.Called(ctx, publicKey)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactions")
	}

	var r0 []aggregates.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]aggregates.Transaction, error)); ok {
		return rf(ctx, publicKey)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []aggregates.Transaction); ok {
		r0 = rf(ctx, publicKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]aggregates.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, publicKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransactionsGetter creates a new instance of TransactionsGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionsGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionsGetter {
	mock := &TransactionsGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}