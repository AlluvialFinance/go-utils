//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

type Config struct {
	Address string

	DisableLog bool

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
