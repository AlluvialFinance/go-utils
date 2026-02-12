//nolint:revive // package name intentionally reflects domain, not directory name
package jsonrpchttp

import (
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

type Config struct {
	Address string

	HTTP *kilnhttp.ClientConfig
}

func (cfg *Config) SetDefault() *Config {
	if cfg.HTTP == nil {
		cfg.HTTP = new(kilnhttp.ClientConfig)
	}

	cfg.HTTP.SetDefault()

	return cfg
}
