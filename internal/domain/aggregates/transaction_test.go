package aggregates_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

func TestTransaction_SetEURAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		amountLAM uint64
		rate      *big.Rat
		wantEUR   string
	}{
		{
			name:      "Set EUR amount correctly",
			amountLAM: 2000000000,
			rate:      big.NewRat(12345, 10000),
			wantEUR:   "1.62",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := aggregates.Transaction{
				AmountLAM: tt.amountLAM,
			}
			rate := aggregates.Rate{
				Value: tt.rate,
			}

			transaction.SetEURAmount(rate)

			assert.Equal(t, tt.wantEUR, transaction.AmountEUR)
		})
	}
}

func TestTransaction_SetLamportsAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		amountEUR string
		rate      *big.Rat
		wantLAM   uint64
		wantErr   bool
	}{
		{
			name:      "Set Lamports amount correctly",
			amountEUR: "1.62",
			rate:      big.NewRat(12345, 10000),
			wantLAM:   1312272174,
			wantErr:   false,
		},
		{
			name:      "Error with invalid EUR amount",
			amountEUR: "invalid",
			rate:      big.NewRat(12345, 10000),
			wantLAM:   0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := aggregates.Transaction{
				AmountEUR: tt.amountEUR,
			}
			rate := aggregates.Rate{
				Value: tt.rate,
			}

			err := transaction.SetLamportsAmount(rate)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLAM, transaction.AmountLAM)
			}
		})
	}
}
