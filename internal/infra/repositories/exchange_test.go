package repositories_test

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jcleira/coding-challenge/internal/domain/aggregates"
	"github.com/jcleira/coding-challenge/internal/infra/repositories"
)

func TestNewExchange(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]interface{}
		wantErr  bool
	}{
		{
			name: "success",
			response: map[string]interface{}{
				"result": map[string]struct {
					P []string `json:"p"`
				}{
					"SOLEUR": {
						P: []string{"", "1.2"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			response: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    123,
					"message": "error message",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			exchange, err := repositories.NewExchange(ctx, server.URL)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, exchange)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, exchange)
		})
	}
}

func TestGetRate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": map[string]struct {
				P []string `json:"p"`
			}{
				"SOLEUR": {
					P: []string{"", "1.2"},
				},
			},
		})
	}))
	defer server.Close()

	exchange, err := repositories.NewExchange(context.Background(), server.URL)
	assert.NoError(t, err)
	assert.NotNil(t, exchange)

	tests := []struct {
		name    string
		rate    aggregates.Rate
		wantErr error
	}{
		{
			name:    "success",
			rate:    aggregates.Rate{Currency: "SOLEUR", Value: big.NewRat(12, 10)},
			wantErr: nil,
		},
		// TODO: Add test about timeouts
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, err := exchange.GetRate()
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NotNil(t, rate)
		})
	}
}
