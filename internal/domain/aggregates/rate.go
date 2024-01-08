package aggregates

import (
	"math/big"
	"time"
)

// Rate is a structure to store the exchange rate of a currency.
//
// We use big.Rat to store the rate value because it is more precise than
// float64, still my preference here would be to have a int64 with the value
// including the decimal part and keep track of the decimals, as it would be
// more precise and easier to work with.
type Rate struct {
	Currency  string
	Value     *big.Rat
	ExpiredAt time.Time
}
