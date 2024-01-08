package services

import (
	"context"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

// SolanaGetter defines the methods for getting transactions from the Solana
// blockchain.
type SolanaGetter interface {
	GetTransactions(ctx context.Context, publicKey string) ([]aggregates.Transaction, error)
}

// SolanaSender is an interface that defines the methods for sending
// transactions to the Solana blockchain.
type SolanaSender interface {
	SendTransaction(context.Context,
		aggregates.Transaction,
		aggregates.Wallet,
	) (string, error)
}

// SolanaBalanceGetter defines the methods for getting the balance from the
// Solana blockchain from an address.
type SolanaBalanceGetter interface {
	GetBalance(ctx context.Context, publicKey string) (uint64, error)
}

// ExchangeGetter defines the methods for getting exchange rates.
type ExchangeGetter interface {
	GetRate() (aggregates.Rate, error)
}

// WalletGetter defines the methods for getting wallets.
type WalletGetter interface {
	GetWallet(publicKey string) (aggregates.Wallet, error)
}

// WalletCreator defines the methods for creating wallets.
type WalletCreator interface {
	CreateWallet() (aggregates.Wallet, error)
}

// RateGetter defines the methods for getting exchange rates.
type RateGetter interface {
	GetRate() (aggregates.Rate, error)
}
