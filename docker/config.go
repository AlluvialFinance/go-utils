package docker

import (
	docker "github.com/docker/docker/client"
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

type ClientConfig struct {
	Host    string
	Version string
	Client  *kilnhttp.ClientConfig
}

func (cfg *ClientConfig) SetDefault() *ClientConfig {
	if cfg.Client == nil {
		cfg.Client = &kilnhttp.ClientConfig{}
	}
	cfg.Client.SetDefault()

	if cfg.Host == "" {
		cfg.Host = docker.DefaultDockerHost
	}

	// Leave Version empty by default to enable API version negotiation
	// This allows the Docker client to automatically negotiate with the daemon
	// and avoids version mismatch issues between different Docker environments

	return cfg
}

type ComposeConfig struct {
	Client    *ClientConfig
	Namespace string
}

func (cfg *ComposeConfig) SetDefault() *ComposeConfig {
	if cfg.Client == nil {
		cfg.Client = &ClientConfig{}
	}
	cfg.Client.SetDefault()

	return cfg
}
