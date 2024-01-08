package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// WalletBalanceGetter defines the interface for getting wallet balances.
type WalletBalanceGetter interface {
	GetBalance(ctx context.Context, pubKey string) (string, error)
}

// WalletBalanceGetterHandler define the dependencies to perform http requests
// for getting the balance of wallets.
type WalletBalanceGetterHandler struct {
	walletBalanceGetter WalletBalanceGetter
}

// NewWalletBalanceGetterHandler creates a new WalletBalanceGetterHandler.
func NewWalletBalanceGetterHandler(
	walletBalanceGetter WalletBalanceGetter) *WalletBalanceGetterHandler {
	return &WalletBalanceGetterHandler{
		walletBalanceGetter: walletBalanceGetter,
	}
}

// Handler is the http handler func  for getting the balance of wallets.
func (wih *WalletBalanceGetterHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			PublicKey string `json:"public_key"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		balance, err := wih.walletBalanceGetter.GetBalance(r.Context(), request.PublicKey)
		if err != nil {
			http.Error(w, "Error initializing wallet", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(struct {
			PublicKey string `json:"balance"`
		}{
			PublicKey: balance,
		})
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(response); err != nil {
			slog.Error("error writing response", err)
		}
	}
}
