//nolint:revive // package name intentionally reflects domain, not directory name
package jsonrpchttp

import (
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

type Config struct {
	Address string
	Headers map[string]string
	HTTP    *kilnhttp.ClientConfig
}

func (cfg *Config) SetDefault() *Config {
	if cfg.HTTP == nil {
		cfg.HTTP = new(kilnhttp.ClientConfig)
	}

	cfg.HTTP.SetDefault()

	return cfg
}

func (cfg *Config) WithHeaders(headers map[string]string) *Config {
	cfg.Headers = headers
	return cfg
}
