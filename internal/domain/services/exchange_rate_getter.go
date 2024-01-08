package services

import (
	"fmt"
)

// ExchangeRateGetter define the dependencies to get exchange rates.
type ExchangeRateGetter struct {
	exchange RateGetter
}

// NewExchangeRateGetter creates a new ExchangeRateGetter.
func NewExchangeRateGetter(exchange RateGetter) *ExchangeRateGetter {
	return &ExchangeRateGetter{
		exchange: exchange,
	}
}

// GetRate gets the exchange rate.
func (e *ExchangeRateGetter) GetRate() (float64, error) {
	rate, err := e.exchange.GetRate()
	if err != nil {
		return 0, fmt.Errorf("error getting exchange rate: %w", err)
	}

	// Ignoring exact checks here
	value, _ := rate.Value.Float64()

	return value, nil
}
