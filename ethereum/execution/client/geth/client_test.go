//go:build !integration

package geth

import (
	"testing"

	"github.com/kilnfi/go-utils/ethereum/execution/client"
	"github.com/stretchr/testify/assert"
)

// Compile-time interface check.
var _ client.Client = (*Client)(nil)

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8545")
	assert.NotNil(t, c)
	assert.Equal(t, "http://localhost:8545", c.address)
	assert.Nil(t, c.headers)
}

func TestNewClientWithHeaders(t *testing.T) {
	headers := map[string]string{"X-Api-Key": "secret"}
	c := NewClientWithHeaders("http://localhost:8545", headers)
	assert.NotNil(t, c)
	assert.Equal(t, "http://localhost:8545", c.address)
	assert.Equal(t, headers, c.headers)
}
