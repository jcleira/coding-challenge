package repositories

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gagliardetto/solana-go"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

// Vault is an abstraction for storing and retrieving wallets.
//
// At the moment, it does has a pretty simple implementation, as it does store
// the private key in a file, but it could be extended to use a provider, such as
// AWS Secrets Manager, Hashicorp Vault, or Evervault.
type Vault struct {
	Path string
}

// NewVault creates a new Vault instance.
func NewVault(path string) (*Vault, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating vault directory: %w", err)
	}
	return &Vault{Path: path}, nil
}

// CreateWallet creates a new wallet and stores it in the vault.
func (v *Vault) CreateWallet() (aggregates.Wallet, error) {
	privateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return aggregates.Wallet{}, fmt.Errorf("error generating new private key: %w", err)
	}

	publicKey := privateKey.PublicKey().String()

	wallet := aggregates.Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	if err := v.store(wallet); err != nil {
		return aggregates.Wallet{}, fmt.Errorf("error storing wallet: %w", err)
	}

	return wallet, nil
}

// Getaggregates.Wallet retrieves a wallet from the vault by its public key.
func (v *Vault) GetWallet(publicKey string) (aggregates.Wallet, error) {
	filename := filepath.Join(v.Path, publicKey)

	privateKeyBytes, err := os.ReadFile(filename)
	if err != nil {
		return aggregates.Wallet{}, fmt.Errorf("error reading private key from file: %w", err)
	}

	return aggregates.Wallet{
		PrivateKey: solana.PrivateKey(privateKeyBytes),
		PublicKey:  publicKey,
	}, nil
}

// store stores a wallet in the vault by its public key.
func (v *Vault) store(wallet aggregates.Wallet) error {
	filename := filepath.Join(v.Path, wallet.PublicKey)

	if err := os.WriteFile(filename, wallet.PrivateKey, 0644); err != nil {
		return fmt.Errorf("error writing private key to file: %w", err)
	}

	return nil
}
