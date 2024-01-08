package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
)

// TransactionsGetter defines the methods for getting transactions from the Solana
// blockchain.
type TransactionsGetter interface {
	GetTransactions(ctx context.Context, publicKey string) ([]aggregates.Transaction, error)
}

// TransactionsGetterHandler define the dependencies handling transactions get requests.
type TransactionsGetterHandler struct {
	getter TransactionsGetter
}

// NewTransactionsGetterHandler creates a new TransactionsGetterHandler.
func NewTransactionsGetterHandler(getter TransactionsGetter) *TransactionsGetterHandler {
	return &TransactionsGetterHandler{
		getter: getter,
	}
}

// Handler is the http handler func  for getting transactions.
func (h *TransactionsGetterHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			PublicKey string `json:"public_key"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		transactions, err := h.getter.GetTransactions(context.Background(), request.PublicKey)
		if err != nil {
			http.Error(w, "Error getting transactions", http.StatusInternalServerError)
			return
		}

		httpTransactions := make([]httpTransaction, len(transactions))
		for i, transaction := range transactions {
			httpTransactions[i] = httpTransactionFromDomainTransaction(transaction)
		}

		response, err := json.Marshal(
			httpTransactionsResponse{
				HTTPTransactions: httpTransactions,
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

// httpTransactionsResponse is the http version for a list of domain transactions.
type httpTransactionsResponse struct {
	HTTPTransactions []httpTransaction `json:"transactions"`
}

// httpTransaction is the http version for a domain transaction.
type httpTransaction struct {
	Created      time.Time `json:"created"`
	Amount       string    `json:"amount"`
	CounterParty string    `json:"counter_party"`
	Signature    string    `json:"signature"`
}

// httpTransactionFromDomainTransaction converts a domain transaction to an http transaction.
func httpTransactionFromDomainTransaction(transaction aggregates.Transaction) httpTransaction {
	return httpTransaction{
		Created:      transaction.BlockTime,
		Amount:       transaction.AmountEUR,
		CounterParty: transaction.CounterParty,
		Signature:    transaction.Signature,
	}
}
