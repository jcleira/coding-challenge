package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

const (
	// exchangeTickerTime is the time between each exchange rate fetch.
	exchangeTickerTime = 5 * time.Second

	// exchangeRateExpirationTime is the time after which the exchange rate
	// expires. The value is set to 20 seconds for testing purposes. In a real
	// world scenario, this value should be set to X? (1 minute?) I'm not sure.
	//
	// 20 seconds is a good value for testing, as it allows up to 4 failures.
	exchangeRateExpirationTime = 20 * time.Second
)

var (
	// supportedCurrencies is a list of currencies supported by the exchange.
	// In this example, we only support SOLEUR, but eventually this list will
	// grow.
	//
	// Using a global variable here is not ideal, but it's good enough for this
	// example, as we don't have any database or other storage mechanism.
	//
	// TODO: Use a database to store the supported currencies.
	supportedCurrencies = []string{"SOLEUR"}
)

type Exchange struct {
	APIURL string

	// rates is a map of currency pairs to exchange rates.
	// For example, rates["SOLEUR"] = 1.2 means that 1 SOL is worth 1.2 EUR.
	// Using a map here is not idead as the map is not thread-safe, but it's
	// good enough for this example as it's syncronized with a mutex.
	rates map[string]aggregates.Rate

	rateMutex sync.RWMutex
}

// NewExchange creates a new exchange with the given API URL.
//
// The exchange will perform an initial fetch of the exchange rates, to prevent
// the service to operate with non initialized rates.
//
// Then the exchange will start a goroutine that will fetch the exchange rates
// every 5 seconds, I'm using a ticker here, but a better approach would be to
// use a cron job, with an external database to store the rates.
//
// The thing that I don't like about this approach is that there is no proper
// error handling in the goroutine, so if the exchange rate API fails, the
// request will start failing.
func NewExchange(ctx context.Context, apiURL string) (*Exchange, error) {
	e := &Exchange{
		APIURL: apiURL,
	}

	rates := make(map[string]aggregates.Rate)
	for _, currency := range supportedCurrencies {
		rate, err := e.fetchRate(currency)
		if err != nil {
			return nil, fmt.Errorf("error fetching rate: %w", err)
		}

		rates[currency] = rate
	}

	e.rates = rates

	go e.start(ctx)

	return e, nil
}

func (e *Exchange) GetRate() (aggregates.Rate, error) {
	e.rateMutex.RLock()
	defer e.rateMutex.RUnlock()

	// Hardcoded to SOLEUR for now, but eventually this will be a parameter.
	rate, ok := e.rates[supportedCurrencies[0]]
	if !ok {
		return aggregates.Rate{}, aggregates.ErrCurrencyNotSupported
	}

	if rate.ExpiredAt.Before(time.Now()) {
		return aggregates.Rate{}, aggregates.ErrRateExpired
	}

	return rate, nil
}

func (e *Exchange) start(ctx context.Context) {
	ticker := time.NewTicker(exchangeTickerTime)
	for {
		select {
		case <-ticker.C:
			e.rateMutex.Lock()
			for _, currency := range supportedCurrencies {
				rate, err := e.fetchRate(currency)
				if err != nil {
					slog.Error("error fetching rate: ", err)
					continue
				}
				e.rates[currency] = rate
			}
			e.rateMutex.Unlock()

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (e *Exchange) fetchRate(currency string) (aggregates.Rate, error) {
	url := fmt.Sprintf("%s?pair=%s", e.APIURL, currency)
	resp, err := http.Get(url)
	if err != nil {
		return aggregates.Rate{}, fmt.Errorf("error making request to exchange rate api: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Error  []string `json:"error"`
		Result map[string]struct {
			P []string `json:"p"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return aggregates.Rate{}, fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Error) != 0 {
		return aggregates.Rate{}, fmt.Errorf("exchange rate api returned an error: %v", result.Error)
	}

	pairData, ok := result.Result[currency]
	if !ok || len(pairData.P) < 2 {
		return aggregates.Rate{}, fmt.Errorf("invalid data received from exchange rate api")
	}

	// I've used big.Rat to do the calculations, as it does provide a better
	// precision than float64, but it's not ideal, as it's not as easy to work
	// with ,
	rate := &big.Rat{}
	rate, ok = rate.SetString(pairData.P[1])
	if !ok {
		return aggregates.Rate{}, fmt.Errorf("error parsing exchange rate: %s", pairData.P[1])
	}

	return aggregates.Rate{
		Currency:  currency,
		Value:     rate,
		ExpiredAt: time.Now().Add(exchangeRateExpirationTime),
	}, nil
}
