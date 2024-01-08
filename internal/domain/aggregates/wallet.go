package aggregates

// Wallet is a structure to store the private and public keys of a wallet.
type Wallet struct {
	PrivateKey []byte
	PublicKey  string
}
