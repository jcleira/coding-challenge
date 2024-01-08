package services_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
	"github.com/jcleira/coding-challenge/internal/domain/services"
	"github.com/jcleira/coding-challenge/mocks"
)

func TestTransactionsGetter_GetTransactions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	publicKey := "testPublicKey"

	rate := aggregates.Rate{
		Currency:  "SOLEUR",
		Value:     big.NewRat(12345, 10000),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	transactions := []aggregates.Transaction{
		{
			BlockTime:    time.Now(),
			Signer:       "Signer1",
			CounterParty: "CounterParty1",
			Accounts:     []string{"Account1", "Account2"},
			AmountLAM:    5000000000,
			Signature:    "Signature1",
		},
		{
			BlockTime:    time.Now(),
			Signer:       "Signer2",
			CounterParty: "CounterParty2",
			Accounts:     []string{"Account3", "Account4"},
			AmountLAM:    10000000000, // 10 SOL
			Signature:    "Signature2",
		},
	}

	tests := []struct {
		name       string
		beforeFunc func(*mocks.SolanaGetter, *mocks.ExchangeGetter)
		want       []aggregates.Transaction
		wantError  error
	}{
		{
			name: "successful transaction retrieval",
			beforeFunc: func(solana *mocks.SolanaGetter, exchange *mocks.ExchangeGetter) {
				solana.On("GetTransactions", ctx, publicKey).
					Return(transactions, nil)

				exchange.On("GetRate").Return(rate, nil)
			},
			want: transactions,
		},
		{
			name: "error getting transactions from Solana",
			beforeFunc: func(solana *mocks.SolanaGetter, exchange *mocks.ExchangeGetter) {
				solana.On("GetTransactions", ctx, publicKey).
					Return(nil, errors.New("solana error"))

				exchange.AssertNotCalled(t, "GetRate")
			},
			wantError: errors.New("error getting transactions: solana error"),
		},
		{
			name: "error getting exchange rate",
			beforeFunc: func(solana *mocks.SolanaGetter, exchange *mocks.ExchangeGetter) {
				solana.On("GetTransactions", ctx, publicKey).
					Return(transactions, nil)

				exchange.On("GetRate").
					Return(aggregates.Rate{}, errors.New("exchange rate error"))
			},
			wantError: errors.New("error getting rate: exchange rate error"),
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				solana   = mocks.NewSolanaGetter(t)
				exchange = mocks.NewExchangeGetter(t)
			)

			tt.beforeFunc(solana, exchange)

			service := services.NewTransactionsGetter(solana, exchange)

			result, err := service.GetTransactions(ctx, publicKey)

			solana.AssertExpectations(t)
			exchange.AssertExpectations(t)

			if tt.wantError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}
