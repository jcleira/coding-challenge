package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// WalletInitializer defines the dependencies for initializing wallets.
type WalletInitializer interface {
	Initialize() (string, error)
}

// WalletInitializerHandler define the dependencies handling wallet init requests.
type WalletInitializerHandler struct {
	initializer WalletInitializer
}

// NewWalletInitializerHandler creates a new WalletInitializerHandler.
func NewWalletInitializerHandler(initializer WalletInitializer) *WalletInitializerHandler {
	return &WalletInitializerHandler{
		initializer: initializer,
	}
}

// Handler is the http handler func  for initializing wallets.
func (wih *WalletInitializerHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		publicKey, err := wih.initializer.Initialize()
		if err != nil {
			http.Error(w, "Error initializing wallet", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(struct {
			PublicKey string `json:"public_key"`
		}{
			PublicKey: publicKey,
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
