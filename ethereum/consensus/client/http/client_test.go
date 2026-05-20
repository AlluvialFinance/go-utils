//go:build !integration

//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	eth2client "github.com/kilnfi/go-utils/ethereum/consensus/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientImplementsInterface(t *testing.T) {
	iClient := (*eth2client.Client)(nil)
	client := new(Client)
	assert.Implements(t, iClient, client)
}

func TestNewClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"genesis_time":"1606824023","genesis_validators_root":"0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95","genesis_fork_version":"0x00000000"}}`))
	}))
	defer srv.Close()

	cfg := (&Config{Address: srv.URL}).SetDefault()
	c, err := NewClient(cfg)
	require.NoError(t, err)

	_, err = c.GetGenesis(t.Context())
	require.NoError(t, err)
}

func TestNewClientWithHeaders(t *testing.T) {
	var receivedHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("X-Api-Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"genesis_time":"1606824023","genesis_validators_root":"0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95","genesis_fork_version":"0x00000000"}}`))
	}))
	defer srv.Close()

	cfg := (&Config{Address: srv.URL}).SetDefault().WithHeaders(map[string]string{
		"X-Api-Key": "secret",
	})
	c, err := NewClient(cfg)
	require.NoError(t, err)

	_, err = c.GetGenesis(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "secret", receivedHeader)
}
