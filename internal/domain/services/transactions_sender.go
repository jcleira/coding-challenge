package services

import (
	"context"
	"fmt"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

// TransactionsSender defines the dependencies for sending transactions to the
// Solana blockchain.
type TransactionsSender struct {
	vault    WalletGetter
	solana   SolanaSender
	exchange ExchangeGetter
}

// NewTransactionsSender creates a new TransactionsSender.
func NewTransactionsSender(
	vault WalletGetter,
	solana SolanaSender,
	exchange ExchangeGetter,
) *TransactionsSender {
	return &TransactionsSender{
		vault:    vault,
		solana:   solana,
		exchange: exchange,
	}
}

// SendTransaction sends a transaction to the Solana blockchain.
func (ts *TransactionsSender) SendTransaction(
	ctx context.Context, transaction aggregates.Transaction) (string, error) {
	wallet, err := ts.vault.GetWallet(transaction.Signer)
	if err != nil {
		return "", fmt.Errorf("error getting wallet: %w", err)
	}

	rate, err := ts.exchange.GetRate()
	if err != nil {
		return "", fmt.Errorf("error getting exchange rate: %w", err)
	}

	if err := transaction.SetLamportsAmount(rate); err != nil {
		return "", fmt.Errorf("error setting lamports amount: %w", err)
	}

	signature, err := ts.solana.SendTransaction(ctx, transaction, wallet)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %w", err)
	}

	return signature, nil
}
