//go:build !integration

//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := &Config{}

	t.Run("default", func(t *testing.T) {
		cfg.SetDefault()
		assert.NotNil(t, cfg.HTTP)
		assert.True(t, cfg.DisableLog)
	})
}
