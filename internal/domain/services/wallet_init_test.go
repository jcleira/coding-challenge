package services_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
	"github.com/jcleira/coding-challenge/internal/domain/services"
	"github.com/jcleira/coding-challenge/mocks"
)

func TestWalletInitializer_Initialize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		beforeFunc func(*mocks.WalletCreator)
		want       string
		wantError  error
	}{
		{
			name: "successful wallet initialization",
			beforeFunc: func(vault *mocks.WalletCreator) {
				vault.On("CreateWallet").
					Return(aggregates.Wallet{PublicKey: "testPublicKey"}, nil)
			},
			want: "testPublicKey",
		},
		{
			name: "error creating wallet",
			beforeFunc: func(vault *mocks.WalletCreator) {
				vault.On("CreateWallet").
					Return(aggregates.Wallet{}, errors.New("wallet creation error"))
			},
			wantError: fmt.Errorf("error creating wallet: wallet creation error"),
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vault := mocks.NewWalletCreator(t)

			tt.beforeFunc(vault)

			service := services.NewWalletInitializer(vault)

			result, err := service.Initialize()

			vault.AssertExpectations(t)

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
