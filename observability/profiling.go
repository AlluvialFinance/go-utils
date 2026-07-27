package observability

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/grafana/pyroscope-go"
	"github.com/sirupsen/logrus"
)

// ProfilingConfig is the minimal configuration for InitProfiling. When
// ServerAddress is empty, InitProfiling does nothing and returns a no-op stop
// function, making continuous profiling strictly opt-in per deployment
// (mirrors the empty-Endpoint behaviour of InitTracing).
//
// Profiles are pushed with the Pyroscope Go SDK to ServerAddress. In the
// Alluvial setup this is the in-cluster Alloy receiver
// (pyroscope.receive_http), which relays to the central Pyroscope adding the
// gateway auth headers and X-Scope-OrgID — so callers do NOT set BasicAuth or
// a tenant here; the app only ever talks to the local Alloy.
//
// Typical use from a service main, alongside InitTracingFromEnv:
//
//	stop, err := observability.InitProfilingFromEnv(ctx, observability.ProfilingConfig{
//	    ServiceName:    "my-service",
//	    ServiceVersion: version.Build.Version,
//	    Environment:    cfg.Environment,
//	}, logger)
//	if err != nil { return err }
//	defer func() {
//	    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	    defer cancel()
//	    _ = stop(ctx)
//	}()
type ProfilingConfig struct {
	// ServerAddress is the Pyroscope push endpoint as a full URL, e.g.
	// "http://alloy.telemetry.svc.cluster.local:9998". Required to enable
	// profiling; empty disables it.
	ServerAddress string
	// ServiceName maps to Pyroscope's ApplicationName (the profiled service,
	// e.g. "concorde"). Required when ServerAddress is non-empty.
	ServiceName string
	// ServiceVersion is attached as the "service_version" tag. Optional.
	ServiceVersion string
	// Environment is attached as the "environment" tag (e.g. "dev", "staging",
	// "prod"). Optional.
	Environment string
	// InstanceID is attached as the "instance" tag. Should be unique per
	// running process (e.g. the pod name). When empty, InitProfiling fills it
	// from os.Hostname() and omits the tag if that also fails.
	InstanceID string
	// Tags are extra tags merged onto the derived ones. Explicit keys here win
	// over the derived service_version/environment/instance tags.
	Tags map[string]string
	// ProfileTypes overrides the set of profiles collected. When empty, the
	// Pyroscope SDK's own defaults apply: CPU, alloc_objects, alloc_space,
	// inuse_objects and inuse_space.
	//
	// Mutex and block profiles are not collected by default; they additionally
	// require runtime.SetMutexProfileFraction / runtime.SetBlockProfileRate to
	// be set by the service, so enabling them is left to the caller.
	ProfileTypes []pyroscope.ProfileType
}

// InitProfiling starts the Pyroscope continuous profiler pushing to
// cfg.ServerAddress.
//
// When cfg.ServerAddress is empty it is a no-op: no profiler is started and a
// no-op stop is returned, so callers can defer the stop unconditionally.
//
// The returned stop function SHOULD be called on process exit to flush the
// final profile. It is safe to call more than once (subsequent calls are
// no-ops returning the first result), matching the idempotency of
// InitTracing's shutdown. It takes a context for signature symmetry with
// that shutdown; the Pyroscope SDK's own Stop does not observe it.
//
// logger is optional (pass nil to silence the SDK). It is typed as
// logrus.FieldLogger, which structurally satisfies the Pyroscope SDK's Logger
// interface, so it is handed straight through.
func InitProfiling(ctx context.Context, cfg ProfilingConfig, logger logrus.FieldLogger) (stop func(context.Context) error, err error) {
	if cfg.ServerAddress == "" {
		return func(context.Context) error { return nil }, nil
	}
	if cfg.ServiceName == "" {
		return nil, errors.New("observability: ServiceName is required when ServerAddress is set")
	}
	// Don't start a background profiler if the caller's context is already done.
	if cerr := ctx.Err(); cerr != nil {
		return nil, cerr
	}

	// ProfileTypes is passed through as-is: the SDK substitutes its own
	// defaults when the slice is empty.
	pc := pyroscope.Config{
		ApplicationName: cfg.ServiceName,
		ServerAddress:   cfg.ServerAddress,
		Tags:            profilingTags(cfg),
		ProfileTypes:    cfg.ProfileTypes,
	}
	if logger != nil {
		pc.Logger = logger
	}

	profiler, err := pyroscope.Start(pc)
	if err != nil {
		return nil, fmt.Errorf("start pyroscope profiler: %w", err)
	}

	// Guard with a Once: the SDK's uploader Stop closes its done channel
	// without a guard of its own, so a second call panics with "close of
	// closed channel". OTel's TracerProvider.Shutdown is idempotent, so
	// callers reasonably expect the same of this stop (e.g. a deferred stop
	// plus an explicit one on a shutdown path).
	var (
		once    sync.Once
		stopErr error
	)
	return func(context.Context) error {
		once.Do(func() { stopErr = profiler.Stop() })
		return stopErr
	}, nil
}

// profilingTags derives the Pyroscope tag set from the config: service_version,
// environment and instance (falling back to the hostname), with any explicit
// cfg.Tags merged last so callers can override.
func profilingTags(cfg ProfilingConfig) map[string]string {
	tags := map[string]string{}
	if cfg.ServiceVersion != "" {
		tags["service_version"] = cfg.ServiceVersion
	}
	if cfg.Environment != "" {
		tags["environment"] = cfg.Environment
	}
	if cfg.InstanceID != "" {
		tags["instance"] = cfg.InstanceID
	} else if hostname, herr := os.Hostname(); herr == nil && hostname != "" {
		tags["instance"] = hostname
	}
	for k, v := range cfg.Tags {
		tags[k] = v
	}
	return tags
}

// InitProfilingFromEnv reads the PYROSCOPE_* environment variables, merges
// them on top of the provided defaults, and calls InitProfiling. Returns a
// no-op stop when PYROSCOPE_SERVER_ADDRESS is unset, so callers can defer the
// stop unconditionally.
//
// Recognised env vars:
//   - PYROSCOPE_SERVER_ADDRESS    — full URL of the push endpoint (typically
//     the in-cluster Alloy receiver). Required to enable profiling.
//   - PYROSCOPE_APPLICATION_NAME  — overrides defaults.ServiceName.
func InitProfilingFromEnv(ctx context.Context, defaults ProfilingConfig, logger logrus.FieldLogger) (stop func(context.Context) error, err error) {
	cfg, enabled := profilingConfigFromEnv(defaults)
	if !enabled {
		return func(context.Context) error { return nil }, nil
	}
	return InitProfiling(ctx, cfg, logger)
}

// profilingConfigFromEnv applies the PYROSCOPE_* env vars on top of defaults
// and reports whether profiling should be initialized at all (false when
// PYROSCOPE_SERVER_ADDRESS is unset).
//
// Split from InitProfilingFromEnv so tests can assert the env-to-config
// transform without standing up a real profiler.
func profilingConfigFromEnv(defaults ProfilingConfig) (ProfilingConfig, bool) {
	cfg := defaults

	addr := strings.TrimSpace(os.Getenv("PYROSCOPE_SERVER_ADDRESS"))
	if addr == "" {
		return cfg, false
	}
	cfg.ServerAddress = addr

	if v := strings.TrimSpace(os.Getenv("PYROSCOPE_APPLICATION_NAME")); v != "" {
		cfg.ServiceName = v
	}

	return cfg, true
}
