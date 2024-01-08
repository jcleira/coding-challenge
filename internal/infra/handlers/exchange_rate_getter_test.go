package handlers_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/require"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
	"github.com/jcleira/coding-challenge/internal/infra/handlers"
	"github.com/jcleira/coding-challenge/mocks"
)

type settingsTestExchangeRateHandler struct {
	exchangeRateGetter *mocks.ExchangeRateGetter
}

func TestExchangeRateGetterHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title          string
		beforeFunc     func(*settingsTestExchangeRateHandler)
		wantStatusCode int
	}{
		{
			title: "success",
			beforeFunc: func(s *settingsTestExchangeRateHandler) {
				s.exchangeRateGetter.On("GetRate").
					Return(1.2345, nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			title: "currency not supported",
			beforeFunc: func(s *settingsTestExchangeRateHandler) {
				s.exchangeRateGetter.On("GetRate").
					Return(0.0, aggregates.ErrCurrencyNotSupported)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			title: "rate expired error",
			beforeFunc: func(s *settingsTestExchangeRateHandler) {
				s.exchangeRateGetter.On("GetRate").
					Return(0.0, aggregates.ErrRateExpired)
			},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			title: "internal server error",
			beforeFunc: func(s *settingsTestExchangeRateHandler) {
				s.exchangeRateGetter.On("GetRate").
					Return(0.0, errors.New("internal error"))
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	cupaloy := cupaloy.New(
		cupaloy.SnapshotSubdirectory("./.snapshots/exchange-rate-test"))

	for _, test := range tests {
		test := test
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()

			settings := &settingsTestExchangeRateHandler{
				exchangeRateGetter: &mocks.ExchangeRateGetter{},
			}
			test.beforeFunc(settings)

			handler := handlers.NewExchangeRateGetterHandler(settings.exchangeRateGetter)

			mux := http.NewServeMux()
			mux.Handle("/", handler.Handler())

			server := httptest.NewServer(mux)
			defer server.Close()

			resp, err := http.Get(server.URL)
			require.NoError(t, err)

			require.True(t, settings.exchangeRateGetter.AssertExpectations(t))

			require.Equal(t, test.wantStatusCode, resp.StatusCode)

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.NoError(t, cupaloy.SnapshotMulti(
				getSnapshotFileName(test.title),
				string(body)))
		})
	}
}

func getSnapshotFileName(testName string) string {
	return fmt.Sprintf("%s_response", testName)
}
