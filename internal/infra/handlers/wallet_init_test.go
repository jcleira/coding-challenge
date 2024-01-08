package handlers_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jcleira/coding-challenge/internal/infra/handlers"
	"github.com/jcleira/coding-challenge/mocks"
)

func TestWalletInitializerHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title          string
		beforeFunc     func(*mocks.WalletInitializer)
		wantStatusCode int
	}{
		{
			title: "successful wallet initialization",
			beforeFunc: func(initializer *mocks.WalletInitializer) {
				initializer.On("Initialize").Return("testPublicKey", nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			title: "error during wallet initialization",
			beforeFunc: func(initializer *mocks.WalletInitializer) {
				initializer.On("Initialize").Return("", errors.New("initialization error"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	cupaloy := cupaloy.New(
		cupaloy.SnapshotSubdirectory("./.snapshots/wallet-init-test"))

	for _, test := range tests {
		test := test
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()

			initializer := &mocks.WalletInitializer{}
			test.beforeFunc(initializer)

			handler := handlers.NewWalletInitializerHandler(initializer)

			mux := http.NewServeMux()
			mux.Handle("/", handler.Handler())

			server := httptest.NewServer(mux)
			defer server.Close()

			resp, err := http.Get(server.URL)
			assert.NoError(t, err)

			assert.Equal(t, test.wantStatusCode, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			resp.Body.Close()

			require.NoError(t, cupaloy.SnapshotMulti(
				getSnapshotFileName(test.title),
				string(body)))

			assert.True(t, initializer.AssertExpectations(t))
		})
	}
}
