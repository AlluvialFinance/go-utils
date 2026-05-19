//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

type Config struct {
	Address string

	DisableLog bool

	Headers map[string]string

	HTTP *kilnhttp.ClientConfig
}

func (cfg *Config) SetDefault() *Config {
	if cfg.HTTP == nil {
		cfg.HTTP = new(kilnhttp.ClientConfig)
	}

	cfg.HTTP.SetDefault()

	cfg.DisableLog = true // Log disabled by default

	return cfg
}

func (cfg *Config) WithHeaders(headers map[string]string) *Config {
	snapshot := make(map[string]string, len(headers))
	for k, v := range headers {
		snapshot[k] = v
	}
	cfg.Headers = snapshot
	return cfg
}
