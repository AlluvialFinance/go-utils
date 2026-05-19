//go:build !integration

//revive:disable-next-line:package-directory-mismatch
package jsonrpchttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		cfg := new(Config).SetDefault()
		require.NotNil(t, cfg.HTTP)
	})

	t.Run("WithHeaders", func(t *testing.T) {
		headers := map[string]string{"X-Api-Key": "secret"}
		cfg := new(Config).SetDefault().WithHeaders(headers)
		assert.Equal(t, headers, cfg.Headers)
	})
}
