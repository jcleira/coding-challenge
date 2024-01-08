package services

import (
	"context"
	"fmt"
	"math/big"
)

const (
	lamportsPerSol = 1000000000
)

// WalletBalanceGetter defines the dependencies for getting the balance of a
// wallet.
type WalletBalanceGetter struct {
	solana   SolanaBalanceGetter
	exchange ExchangeGetter
}

// NewWalletBalanceGetter creates a new WalletBalanceGetter.
func NewWalletBalanceGetter(
	solana SolanaBalanceGetter,
	exchange ExchangeGetter,
) *WalletBalanceGetter {
	return &WalletBalanceGetter{
		solana:   solana,
		exchange: exchange,
	}
}

// GetBalance gets the balance of a wallet.
func (wbg *WalletBalanceGetter) GetBalance(
	ctx context.Context, publicKey string) (string, error) {
	balance, err := wbg.solana.GetBalance(ctx, publicKey)
	if err != nil {
		return "", fmt.Errorf("error getting balance: %w", err)
	}

	rate, err := wbg.exchange.GetRate()
	if err != nil {
		return "", fmt.Errorf("error getting rate: %w", err)
	}

	amountInSol := new(big.Rat).Quo(
		new(big.Rat).SetUint64(balance),
		new(big.Rat).SetInt64(lamportsPerSol),
	)

	amountInEUR := new(big.Rat).Mul(
		amountInSol,
		rate.Value,
	)

	return fmt.Sprintf("%s %s", "EUR", amountInEUR.FloatString(2)), nil
}
