package pprof

import "errors"

// DefaultAddress is the default pprof server address (localhost-only for security).
const DefaultAddress = "127.0.0.1:6060"

// Config configures the pprof debug server.
type Config struct {
	// Enabled controls whether the pprof server is started.
	// When enabled, a separate HTTP server is started for pprof endpoints.
	//
	// Security: The presence of port 6060 in GitOps configs is a clear signal
	// that debug endpoints are enabled - easy to spot during PR reviews.
	Enabled bool

	// Address is the listen address for the pprof server.
	// Defaults to "127.0.0.1:6060" (localhost-only for security).
	// Access via: kubectl port-forward pod/<pod> 6060:6060
	Address string
}

// Validate checks that the pprof configuration is valid.
func (cfg *Config) Validate() error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.Address == "" {
		return errors.New("pprof: address is required when enabled")
	}
	return nil
}

// SetDefaults sets default values for unset fields.
func (cfg *Config) SetDefaults() {
	if cfg == nil {
		return
	}
	if cfg.Address == "" {
		cfg.Address = DefaultAddress
	}
}
