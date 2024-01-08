package services_test

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
	"github.com/jcleira/coding-challenge/internal/domain/services"
	"github.com/jcleira/coding-challenge/mocks"
)

func TestExchangeRateGetter_GetRate(t *testing.T) {
	t.Parallel()

	rate := aggregates.Rate{
		Currency:  "USD",
		Value:     big.NewRat(12345, 10000),
		ExpiredAt: time.Now().Add(1 * time.Hour),
	}

	tests := []struct {
		name       string
		beforeFunc func(*mocks.RateGetter)
		wantRate   float64
		wantError  error
	}{
		{
			name: "successful rate retrieval",
			beforeFunc: func(erg *mocks.RateGetter) {
				erg.On("GetRate").Return(rate, nil)
			},
			wantRate: 1.234500,
		},
		{
			name: "error in rate retrieval",
			beforeFunc: func(erg *mocks.RateGetter) {
				erg.On("GetRate").
					Return(aggregates.Rate{}, errors.New("error"))
			},
			wantError: errors.New("error getting exchange rate: error"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rateGetter := mocks.NewRateGetter(t)
			test.beforeFunc(rateGetter)

			exchangeRateGetter := services.NewExchangeRateGetter(rateGetter)

			rate, err := exchangeRateGetter.GetRate()

			rateGetter.AssertExpectations(t)

			if test.wantError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.wantError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.wantRate, rate)

		})
	}
}
