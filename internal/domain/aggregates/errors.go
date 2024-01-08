package aggregates

import "errors"

var (
	// ErrCurrencyNotSupported is returned when the currency is not supported.
	ErrCurrencyNotSupported = errors.New("currency not supported")

	// ErrRateExpired is returned when the rate is expired.
	ErrRateExpired = errors.New("rate expired")

	// ErrTransactionConfirmationTimeout is returned when the transaction
	// confirmation times out.
	ErrTransactionConfirmationTimeout = errors.New("transaction confirmation timeout")

	// ErrNoRecentBlockHashValue is returned when the get recent block hash returns
	// an empty value.
	ErrNoRecentBlockHashValue = errors.New("not recent block hash value")

	// ErrNoCounterParty is returned when the counter party is not found
	ErrNoCounterParty = errors.New("no counter party found")

	// ErrNoInstructions is returned when the transaction has no instructions
	ErrNoInstructions = errors.New("no instructions found")

	// ErrMultipleInstructions is returned when the transaction has multiple instructions
	ErrMultipleInstructions = errors.New("multiple instructions found")

	// ErrInvalidInstruction is returned when the instruction is invalid
	ErrInvalidInstruction = errors.New("invalid instruction")
)
