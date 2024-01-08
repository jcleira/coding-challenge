package aggregates

import (
	"fmt"
	"math/big"
	"time"
)

const (
	lamportsPerSol = 1000000000
)

// Transaction is the domain representation of a transaction, it is used as an
// abstraction for the Solana transaction.
type Transaction struct {
	BlockTime    time.Time
	Signer       string
	CounterParty string
	Accounts     []string
	AmountLAM    uint64
	AmountEUR    string
	Signature    string
}

// SetEURAmount calculates the amount in EUR based on the exchange rate.
func (t *Transaction) SetEURAmount(rate Rate) {
	amountInSol := new(big.Rat).Quo(
		new(big.Rat).SetUint64(t.AmountLAM),
		new(big.Rat).SetInt64(lamportsPerSol),
	)

	amountInEUR := new(big.Rat).Quo(
		amountInSol,
		rate.Value,
	)

	t.AmountEUR = amountInEUR.FloatString(2)
}

// SetLamportsAmount sets the amount in lamports based on the exchange rate.
func (t *Transaction) SetLamportsAmount(rate Rate) error {
	amountRat := &big.Rat{}
	amountRat, ok := amountRat.SetString(t.AmountEUR)
	if !ok {
		return fmt.Errorf("error converting amount to big.Rat")
	}

	solAmount := new(big.Rat).Quo(amountRat, rate.Value)
	lamportsInSol := new(big.Rat).SetInt64(lamportsPerSol)

	lamports := new(big.Rat).Mul(solAmount, lamportsInSol)

	// For the purpose of this exercise, we will assume that the amount in
	// lamports is always precise.
	lamportsFloat, _ := lamports.Float64()

	t.AmountLAM = uint64(lamportsFloat)

	return nil
}
