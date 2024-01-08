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

type mockBalanceRequest struct {
	PublicKey string `json:"public_key"`
}

func TestWalletBalanceGetterHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title          string
		requestBody    *mockBalanceRequest
		beforeFunc     func(*mocks.WalletBalanceGetter)
		wantStatusCode int
	}{
		{
			title: "successful balance retrieval",
			requestBody: &mockBalanceRequest{
				PublicKey: "testPublicKey",
			},
			beforeFunc: func(balanceGetter *mocks.WalletBalanceGetter) {
				balanceGetter.On("GetBalance", mock.Anything, "testPublicKey").
					Return("10.10", nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			title: "internal server error on balance retrieval",
			requestBody: &mockBalanceRequest{
				PublicKey: "testPublicKey",
			},
			beforeFunc: func(balanceGetter *mocks.WalletBalanceGetter) {
				balanceGetter.On("GetBalance", mock.Anything, "testPublicKey").
					Return("10.10", errors.New("internal server error"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	cupaloy := cupaloy.New(
		cupaloy.SnapshotSubdirectory("./.snapshots/wallet-balance-get-test"))

	for _, test := range tests {
		test := test
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()

			balanceGetter := &mocks.WalletBalanceGetter{}
			test.beforeFunc(balanceGetter)

			handler := handlers.NewWalletBalanceGetterHandler(balanceGetter)

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

			assert.True(t, balanceGetter.AssertExpectations(t))
		})
	}
}
