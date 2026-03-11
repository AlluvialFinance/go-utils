package pprof

import (
	"errors"
	"net"
)

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
	// Must be a loopback address (127.0.0.1, ::1, or localhost).
	// Access via: kubectl port-forward pod/<pod> 6060:6060
	Address string
}

// Validate checks that the pprof configuration is valid.
// Enforces that the address binds to loopback only (127.0.0.1, ::1, localhost).
func (cfg *Config) Validate() error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.Address == "" {
		return errors.New("pprof: address is required when enabled")
	}
	if err := validateLoopbackAddress(cfg.Address); err != nil {
		return err
	}
	return nil
}

// validateLoopbackAddress ensures the address binds to loopback only.
func validateLoopbackAddress(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		// Maybe it's just a host without port
		host = addr
	}

	// Reject empty host (wildcard bind)
	if host == "" {
		return errors.New("pprof: wildcard bind not allowed, use loopback address (127.0.0.1 or ::1)")
	}

	// Allow "localhost"
	if host == "localhost" {
		return nil
	}

	// Parse as IP and check if loopback
	ip := net.ParseIP(host)
	if ip == nil {
		return errors.New("pprof: invalid host in address, use loopback address (127.0.0.1 or ::1)")
	}

	if !ip.IsLoopback() {
		return errors.New("pprof: non-loopback bind not allowed, use 127.0.0.1 or ::1 (security: pprof must not be network-accessible)")
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
