//nolint:revive // package name intentionally reflects domain, not directory name
package gethkeystore

type Config struct {
	Path     string `json:"path"`
	Password string `json:"-"`
}

func (cfg *Config) SetDefault() *Config {
	if cfg.Path == "" {
		cfg.Path = "keystore"
	}

	return cfg
}
