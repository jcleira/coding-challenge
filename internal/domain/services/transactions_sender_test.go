package services_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
	"github.com/jcleira/coding-challenge/internal/domain/services"
	"github.com/jcleira/coding-challenge/mocks"
)

func TestTransactionsSender_SendTransaction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	transaction := aggregates.Transaction{
		Signer:       "Signer1",
		CounterParty: "CounterParty1",
		AmountEUR:    "10.12",
		AmountLAM:    8197650870,
	}

	rate := aggregates.Rate{
		Currency:  "USD",
		Value:     big.NewRat(12345, 10000),
		ExpiredAt: time.Now().Add(1 * time.Hour),
	}

	wallet := aggregates.Wallet{
		PublicKey: "testPublicKey",
	}

	tests := []struct {
		name       string
		beforeFunc func(*mocks.WalletGetter, *mocks.SolanaSender, *mocks.ExchangeGetter)
		want       string
		wantError  error
	}{
		{
			name: "successful transaction send",
			beforeFunc: func(vault *mocks.WalletGetter, solana *mocks.SolanaSender, exchange *mocks.ExchangeGetter) {
				vault.On("GetWallet", transaction.Signer).
					Return(wallet, nil)

				exchange.On("GetRate").Return(rate, nil)

				solana.On("SendTransaction", ctx, transaction, wallet).Return("signature", nil)
			},
			want: "signature",
		},
		{
			name: "error getting wallet",
			beforeFunc: func(vault *mocks.WalletGetter, solana *mocks.SolanaSender, exchange *mocks.ExchangeGetter) {
				vault.On("GetWallet", transaction.Signer).
					Return(aggregates.Wallet{}, errors.New("wallet error"))

				exchange.AssertNotCalled(t, "GetRate")
				solana.AssertNotCalled(t, "SendTransaction")
			},
			wantError: fmt.Errorf("error getting wallet: wallet error"),
		},
		{
			name: "error getting exchange rate",
			beforeFunc: func(vault *mocks.WalletGetter, solana *mocks.SolanaSender, exchange *mocks.ExchangeGetter) {
				vault.On("GetWallet", transaction.Signer).
					Return(wallet, nil)

				exchange.On("GetRate").
					Return(aggregates.Rate{}, errors.New("exchange rate error"))

				solana.AssertNotCalled(t, "SendTransaction")
			},
			wantError: fmt.Errorf("error getting exchange rate: exchange rate error"),
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				vault    = mocks.NewWalletGetter(t)
				solana   = mocks.NewSolanaSender(t)
				exchange = mocks.NewExchangeGetter(t)
			)

			tt.beforeFunc(vault, solana, exchange)

			service := services.NewTransactionsSender(vault, solana, exchange)

			result, err := service.SendTransaction(ctx, transaction)

			vault.AssertExpectations(t)
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
