package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

// TransactionsSender defines the methods for sending transactions to the
// Solana blockchain.
type TransactionsSender interface {
	SendTransaction(context.Context, aggregates.Transaction) (string, error)
}

// TransactionsSenderHandler handles sending transactions to the Solana blockchain.
type TransactionsSenderHandler struct {
	TransactionsSender TransactionsSender
}

// NewTransactionsSenderHandler creates a new TransactionsSenderHandler.
func NewTransactionsSenderHandler(transactionsSender TransactionsSender) *TransactionsSenderHandler {
	return &TransactionsSenderHandler{
		TransactionsSender: transactionsSender,
	}
}

// Handler handles sending transactions to the Solana blockchain.
func (th *TransactionsSenderHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			PublicKey string `json:"public_key"`
			To        string `json:"to"`
			Amount    string `json:"amount"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		splitted := strings.Split(request.Amount, " ")
		if len(splitted) != 2 {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		amountEUR := splitted[1]

		transaction := aggregates.Transaction{
			Signer:       request.PublicKey,
			CounterParty: request.To,
			AmountEUR:    amountEUR,
		}

		signature, err := th.TransactionsSender.SendTransaction(r.Context(), transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(
			struct {
				Signature string `json:"signature"`
			}{Signature: signature},
		)
		if err != nil {
			http.Error(w, "Error marshalling response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(response); err != nil {
			slog.Error("Error writing response", err)
		}
	}
}
