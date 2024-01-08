package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

// TransactionsGetter defines the dependencies for getting transactions from the Solana
// blockchain.
type TransactionsGetter struct {
	solana   SolanaGetter
	exchange ExchangeGetter
}

// NewTransactionsGetter creates a new TransactionsGetter.
func NewTransactionsGetter(solana SolanaGetter, exchange ExchangeGetter) *TransactionsGetter {
	return &TransactionsGetter{
		solana:   solana,
		exchange: exchange,
	}
}

func (t *TransactionsGetter) GetTransactions(ctx context.Context, publicKey string) ([]aggregates.Transaction, error) {
	transactions, err := t.solana.GetTransactions(ctx, publicKey)
	if err != nil {
		slog.Error("error getting transactions", err)
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}

	rate, err := t.exchange.GetRate()
	if err != nil {
		slog.Error("error getting rate", err)
		return nil, fmt.Errorf("error getting rate: %w", err)
	}

	for i := range transactions {
		transactions[i].SetEURAmount(rate)
	}

	return transactions, nil
}
