package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jcleira/coding-challenge/internal/infra/handlers"
	"github.com/jcleira/coding-challenge/mocks"
)

type mockSendRequest struct {
	PublicKey string `json:"public_key"`
	To        string `json:"to"`
	Amount    string `json:"amount"`
}

func TestTransactionsSenderHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title          string
		requestBody    *mockSendRequest
		beforeFunc     func(*mocks.TransactionsSender)
		wantStatusCode int
	}{
		{
			title: "successful transaction sending",
			requestBody: &mockSendRequest{
				PublicKey: "testPublicKey",
				To:        "testReceiver",
				Amount:    "100 EUR",
			},
			beforeFunc: func(sender *mocks.TransactionsSender) {
				sender.On("SendTransaction",
					mock.Anything, mock.AnythingOfType("aggregates.Transaction")).
					Return("testSignature", nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			title: "bad request with invalid body",
			requestBody: &mockSendRequest{
				PublicKey: "",
			},
			beforeFunc: func(sender *mocks.TransactionsSender) {
				sender.AssertNotCalled(t, "SendTransaction")
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			title: "bad request with invalid amount format",
			requestBody: &mockSendRequest{
				PublicKey: "testPublicKey",
				To:        "testReceiver",
				Amount:    "invalidAmount",
			},
			beforeFunc: func(sender *mocks.TransactionsSender) {
				sender.AssertNotCalled(t, "SendTransaction")
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			title: "internal server error on transaction sending",
			requestBody: &mockSendRequest{
				PublicKey: "testPublicKey",
				To:        "testReceiver",
				Amount:    "100 EUR",
			},
			beforeFunc: func(sender *mocks.TransactionsSender) {
				sender.On("SendTransaction",
					mock.Anything, mock.AnythingOfType("aggregates.Transaction")).
					Return("", errors.New("internal server error"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	cupaloy := cupaloy.New(
		cupaloy.SnapshotSubdirectory("./.snapshots/transactions-send-test"))

	for _, test := range tests {
		test := test
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()

			sender := &mocks.TransactionsSender{}
			test.beforeFunc(sender)

			handler := handlers.NewTransactionsSenderHandler(sender)

			mux := http.NewServeMux()
			mux.Handle("/", handler.Handler())

			server := httptest.NewServer(mux)
			defer server.Close()

			requestBody, _ := json.Marshal(test.requestBody)
			req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(requestBody))
			assert.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)

			assert.Equal(t, test.wantStatusCode, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			require.NoError(t, cupaloy.SnapshotMulti(
				getSnapshotFileName(test.title),
				string(body)))

			assert.True(t, sender.AssertExpectations(t))
		})
	}
}
