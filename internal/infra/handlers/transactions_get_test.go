package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
	"github.com/jcleira/coding-challenge/internal/infra/handlers"
	"github.com/jcleira/coding-challenge/mocks"
)

type mockRequest struct {
	PublicKey string `json:"public_key"`
}

func TestTransactionsGetterHandler_Handle(t *testing.T) {
	t.Parallel()

	blockTime := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		title          string
		requestBody    *mockRequest
		beforeFunc     func(*mocks.TransactionsGetter)
		wantStatusCode int
	}{
		{
			title: "successful transactions retrieval",
			requestBody: &mockRequest{
				PublicKey: "testPublicKey",
			},
			beforeFunc: func(getter *mocks.TransactionsGetter) {
				getter.On("GetTransactions", context.Background(), "testPublicKey").
					Return([]aggregates.Transaction{
						{
							BlockTime:    blockTime,
							CounterParty: "testCounterParty",
							AmountEUR:    "100.00",
							Signature:    "testSignature",
						},
					}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			title: "internal server error on transaction retrieval",
			requestBody: &mockRequest{
				PublicKey: "testPublicKey",
			},
			beforeFunc: func(getter *mocks.TransactionsGetter) {
				getter.On("GetTransactions", context.Background(), "testPublicKey").
					Return(nil, errors.New("internal server error"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	cupaloy := cupaloy.New(
		cupaloy.SnapshotSubdirectory("./.snapshots/transactions-get-test"))

	for _, test := range tests {
		test := test
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()

			getter := &mocks.TransactionsGetter{}
			test.beforeFunc(getter)

			handler := handlers.NewTransactionsGetterHandler(getter)

			mux := http.NewServeMux()
			mux.Handle("/", handler.Handler())

			server := httptest.NewServer(mux)
			defer server.Close()

			requestBody, _ := json.Marshal(test.requestBody)
			req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(requestBody))
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			require.Equal(t, test.wantStatusCode, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.NoError(t, cupaloy.SnapshotMulti(
				getSnapshotFileName(test.title),
				string(body)))

			assert.True(t, getter.AssertExpectations(t))
		})
	}
}
