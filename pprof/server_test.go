package pprof

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.PanicLevel) // suppress logs in tests

	t.Run("nil config returns nil server", func(t *testing.T) {
		s := NewServer(nil, log)
		assert.Nil(t, s)
	})

	t.Run("disabled config returns nil server", func(t *testing.T) {
		s := NewServer(&Config{Enabled: false}, log)
		assert.Nil(t, s)
	})

	t.Run("enabled config creates server", func(t *testing.T) {
		cfg := &Config{Enabled: true, Address: "127.0.0.1:6060"}
		s := NewServer(cfg, log)
		require.NotNil(t, s)
		assert.Equal(t, "127.0.0.1:6060", s.Addr())
	})
}

func TestConfigValidate(t *testing.T) {
	t.Run("nil config is valid", func(t *testing.T) {
		var cfg *Config
		assert.NoError(t, cfg.Validate())
	})

	t.Run("disabled config is valid", func(t *testing.T) {
		assert.NoError(t, (&Config{Enabled: false}).Validate())
	})

	t.Run("enabled config without address fails", func(t *testing.T) {
		err := (&Config{Enabled: true, Address: ""}).Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "address is required")
	})

	t.Run("enabled config with address is valid", func(t *testing.T) {
		assert.NoError(t, (&Config{Enabled: true, Address: ":6060"}).Validate())
	})
}

func TestConfigSetDefaults(t *testing.T) {
	t.Run("sets default address", func(t *testing.T) {
		cfg := &Config{Enabled: true}
		cfg.SetDefaults()
		assert.Equal(t, DefaultAddress, cfg.Address)
	})

	t.Run("does not override explicit address", func(t *testing.T) {
		cfg := &Config{Enabled: true, Address: ":9999"}
		cfg.SetDefaults()
		assert.Equal(t, ":9999", cfg.Address)
	})
}

func TestServerStartStop(t *testing.T) {
	t.Run("nil server start/stop is safe", func(t *testing.T) {
		var s *Server
		s.Start() // should not panic
		assert.NoError(t, s.Stop(t.Context()))
	})
}
