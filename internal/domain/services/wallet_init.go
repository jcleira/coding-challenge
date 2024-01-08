package services

import (
	"fmt"
)

// WalletInitializer defines the dependencies for initializing wallets.
type WalletInitializer struct {
	vault WalletCreator
}

// NewWalletInitializer creates a new WalletInitializer.
func NewWalletInitializer(vault WalletCreator) *WalletInitializer {
	return &WalletInitializer{
		vault: vault,
	}
}

// Initialize initializes a wallet.
func (wi *WalletInitializer) Initialize() (string, error) {
	wallet, err := wi.vault.CreateWallet()
	if err != nil {
		return "", fmt.Errorf("error creating wallet: %w", err)
	}

	return wallet.PublicKey, nil
}
