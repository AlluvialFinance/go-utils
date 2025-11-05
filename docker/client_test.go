package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientConfig_SetDefault_DoesNotSetVersion(t *testing.T) {
	cfg := &ClientConfig{}
	cfg.SetDefault()

	// Version should be empty to enable API version negotiation
	assert.Equal(t, "", cfg.Version, "Version should be empty by default to enable API negotiation")
}

func TestClientConfig_SetDefault_PreservesExplicitVersion(t *testing.T) {
	cfg := &ClientConfig{
		Version: "1.48",
	}
	cfg.SetDefault()

	// Explicitly set version should be preserved
	assert.Equal(t, "1.48", cfg.Version, "Explicitly set version should be preserved")
}

func TestNewClient_WithEmptyVersion(t *testing.T) {
	cfg := &ClientConfig{}
	cfg.SetDefault()

	client, err := NewClient(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Client should be created successfully with API version negotiation
	client.Close()
}

func TestNewClient_WithExplicitVersion(t *testing.T) {
	cfg := &ClientConfig{
		Version: "1.48",
	}
	cfg.SetDefault()

	client, err := NewClient(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Client should be created successfully with explicit version
	client.Close()
}
