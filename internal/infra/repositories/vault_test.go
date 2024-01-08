package repositories_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jcleira/coding-challenge/internal/infra/repositories"
)

func TestNewVault(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "vault_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	vault, err := repositories.NewVault(tmpDir)
	assert.NoError(t, err)
	assert.NotNil(t, vault)
	assert.DirExists(t, tmpDir)
}

func TestCreateWallet(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "vault_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	vault, err := repositories.NewVault(tmpDir)
	require.NoError(t, err)

	wallet, err := vault.CreateWallet()
	require.NoError(t, err)
	assert.NotEmpty(t, wallet.PrivateKey)
	assert.NotEmpty(t, wallet.PublicKey)

	filename := filepath.Join(tmpDir, wallet.PublicKey)
	_, err = os.Stat(filename)
	assert.NoError(t, err)
}

func TestGetWallet(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "vault_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	vault, err := repositories.NewVault(tmpDir)
	require.NoError(t, err)

	createdWallet, err := vault.CreateWallet()
	require.NoError(t, err)

	tests := []struct {
		name string
		key  string
		// TODO I would change this bool to an error type.
		wantErr bool
	}{
		{
			name:    "wallet exists",
			key:     createdWallet.PublicKey,
			wantErr: false,
		},
		{
			name:    "wallet does not exist",
			key:     "non-existent-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := vault.GetWallet(tt.key)
			if tt.wantErr {
				assert.True(t, errors.Is(err, os.ErrNotExist))
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, createdWallet.PrivateKey, wallet.PrivateKey)
			assert.Equal(t, createdWallet.PublicKey, wallet.PublicKey)
		})
	}
}
