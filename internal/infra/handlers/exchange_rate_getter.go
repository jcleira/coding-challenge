package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

// ExchangeRateGetter defines the methods to get exchange rates.
type ExchangeRateGetter interface {
	GetRate() (float64, error)
}

// ExchangeRateGetterHandler handles the exchange rate getter.
type ExchangeRateGetterHandler struct {
	getter ExchangeRateGetter
}

// NewExchangeRateGetterHandler creates a new ExchangeRateGetterHandler.
func NewExchangeRateGetterHandler(getter ExchangeRateGetter) *ExchangeRateGetterHandler {
	return &ExchangeRateGetterHandler{
		getter: getter,
	}
}

// Handler handles the exchange rate getter.
func (h *ExchangeRateGetterHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rate, err := h.getter.GetRate()
		if err != nil {
			switch err {
			case aggregates.ErrCurrencyNotSupported:
				http.Error(w, "Currency not supported", http.StatusBadRequest)
				return
			case aggregates.ErrRateExpired:
				http.Error(w, "Currency Rate Expired", http.StatusUnprocessableEntity)
				return
			default:
				http.Error(w, "Error getting exchange rate", http.StatusInternalServerError)
				return
			}
		}

		response, err := json.Marshal(
			struct {
				Rate float64 `json:"sol_eur"`
			}{
				Rate: rate,
			},
		)
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		if _, err = w.Write(response); err != nil {
			slog.Error("Error writing response", err)
		}
	}
}
