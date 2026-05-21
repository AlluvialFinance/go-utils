// Package observability provides shared helpers for wiring OpenTelemetry
// tracing into Alluvial Go services.
//
// Typical use from a service main:
//
//	shutdown, err := observability.InitTracingFromEnv(ctx, observability.TracingConfig{
//	    ServiceName:    "my-service",
//	    ServiceVersion: version.Build.Version,
//	}, logger)
//	if err != nil { return err }
//	defer func() {
//	    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	    defer cancel()
//	    _ = shutdown(ctx)
//	}()
//
// When the resolved Endpoint is empty (e.g. OTEL_EXPORTER_OTLP_ENDPOINT is
// unset) tracing stays disabled: the global no-op TracerProvider is left in
// place and shutdown is a no-op. This makes tracing strictly opt-in per
// deployment.
package observability

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// TracingConfig is the minimal configuration for InitTracing. When Endpoint
// is empty, InitTracing leaves the global no-op TracerProvider in place and
// returns a no-op shutdown.
type TracingConfig struct {
	// Endpoint is the OTLP gRPC endpoint as host:port.
	Endpoint string
	// Insecure disables TLS for the OTLP gRPC connection. For in-cluster
	// collectors (e.g. an Alloy DaemonSet on plain HTTP/2) this should be
	// true.
	Insecure bool
	// ServiceName is the "service.name" resource attribute. Required when
	// Endpoint is non-empty.
	ServiceName string
	// ServiceVersion is the "service.version" resource attribute. Optional.
	ServiceVersion string
	// Environment is the "deployment.environment" resource attribute (e.g.
	// "dev", "staging", "prod"). Optional.
	Environment string
	// InstanceID is the "service.instance.id" resource attribute. Must be
	// unique per running process (e.g. the pod name). When empty, InitTracing
	// fills it from os.Hostname() and falls back to omitting the attribute
	// if that also fails.
	InstanceID string
	// SamplerRatio configures the ParentBased(TraceIDRatioBased) sampler. A
	// value <= 0 or >= 1 selects AlwaysSample. Parent-sampled spans always
	// propagate regardless of this ratio.
	SamplerRatio float64
}

// InitTracing registers a global OTel TracerProvider that exports spans over
// OTLP gRPC, along with the W3C TraceContext and Baggage propagators.
//
// The returned shutdown function MUST be called on process exit (typically
// via defer with a bounded context) to flush buffered spans.
func InitTracing(ctx context.Context, cfg TracingConfig, logger *logrus.Logger) (shutdown func(context.Context) error, err error) {
	if cfg.Endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}
	if cfg.ServiceName == "" {
		return nil, errors.New("observability: ServiceName is required when Endpoint is set")
	}

	opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.Endpoint)}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
	}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ServiceVersion))
	}
	if cfg.Environment != "" {
		attrs = append(attrs, semconv.DeploymentEnvironmentName(cfg.Environment))
	}
	if id := cfg.InstanceID; id != "" {
		attrs = append(attrs, semconv.ServiceInstanceID(id))
	} else if hostname, herr := os.Hostname(); herr == nil && hostname != "" {
		attrs = append(attrs, semconv.ServiceInstanceID(hostname))
	}

	// Compose the Resource with the SDK's standard detectors plus our explicit
	// attributes. resource.New does NOT auto-include detectors — each one must
	// be opted into. Order matters: detectors run first so our explicit
	// WithAttributes wins on conflicting keys (e.g. ServiceName overrides any
	// OTEL_SERVICE_NAME that WithFromEnv picked up).
	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithTelemetrySDK(), // telemetry.sdk.{name,language,version}
		resource.WithProcess(),      // process.{pid,executable,command,...}
		resource.WithOS(),           // os.{type,description,...}
		resource.WithHost(),         // host.{name,id,...}
		resource.WithContainer(),    // container.id (k8s)
		resource.WithFromEnv(),      // OTEL_RESOURCE_ATTRIBUTES (operator-supplied tags)
		resource.WithAttributes(attrs...),
	)
	// resource.New can return ErrPartialResource together with a usable
	// resource — that happens when a built-in detector is unavailable but
	// the explicit attributes were still applied. Treat that as a warning,
	// not a fatal init failure.
	switch {
	case err == nil:
		// nothing to do
	case errors.Is(err, resource.ErrPartialResource):
		if logger != nil {
			logger.WithError(err).Warn("observability: partial OTel resource; some auto-detected attributes are missing")
		}
	default:
		_ = exporter.Shutdown(ctx)
		return nil, fmt.Errorf("build OTel resource: %w", err)
	}

	sampler := sdktrace.AlwaysSample()
	if cfg.SamplerRatio > 0 && cfg.SamplerRatio < 1.0 {
		sampler = sdktrace.TraceIDRatioBased(cfg.SamplerRatio)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sampler)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// InitTracingFromEnv reads the standard OTEL_EXPORTER_OTLP_* environment
// variables, merges them on top of the provided defaults, and calls
// InitTracing. Returns a no-op shutdown when OTEL_EXPORTER_OTLP_ENDPOINT is
// unset, so callers can safely defer the shutdown unconditionally.
//
// Recognised env vars:
//   - OTEL_EXPORTER_OTLP_ENDPOINT   — host:port; supports an optional
//     http://https:// scheme prefix and trailing path/query that is stripped
//     for the OTLP gRPC dial. Required to enable tracing.
//   - OTEL_EXPORTER_OTLP_INSECURE   — "true" to disable TLS. When unset, the
//     scheme of the endpoint decides (https:// → secure, otherwise insecure).
//   - OTEL_SERVICE_NAME             — overrides defaults.ServiceName.
//   - OTEL_DEPLOYMENT_ENVIRONMENT   — overrides defaults.Environment.
//   - OTEL_TRACES_SAMPLER_ARG       — float in [0, 1] for sampling ratio.
//     Invalid values keep the default and emit a warning when logger is non-nil.
func InitTracingFromEnv(ctx context.Context, defaults TracingConfig, logger *logrus.Logger) (shutdown func(context.Context) error, err error) {
	cfg, enabled := tracingConfigFromEnv(defaults, logger)
	if !enabled {
		return func(context.Context) error { return nil }, nil
	}
	return InitTracing(ctx, cfg, logger)
}

// tracingConfigFromEnv applies the OTEL_EXPORTER_OTLP_* env vars on top of
// defaults and returns the resolved config plus whether tracing should be
// initialized at all (false when OTEL_EXPORTER_OTLP_ENDPOINT is unset).
//
// Split from InitTracingFromEnv so tests can assert the env-to-config
// transform in isolation without standing up a real exporter.
func tracingConfigFromEnv(defaults TracingConfig, logger *logrus.Logger) (TracingConfig, bool) {
	cfg := defaults

	endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if endpoint == "" {
		return cfg, false
	}
	// Tolerate a URL-form endpoint. The OTLP gRPC dial wants host:port, so
	// strip an optional scheme and any trailing path/query.
	hasHTTPSScheme := strings.HasPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
	if i := strings.IndexAny(endpoint, "/?"); i != -1 {
		endpoint = endpoint[:i]
	}
	cfg.Endpoint = endpoint

	// OTEL_EXPORTER_OTLP_INSECURE takes precedence over the scheme heuristic
	// when explicitly set (per the OTel env-var spec). Otherwise infer from
	// the scheme: https:// → secure, anything else → insecure.
	if v := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_INSECURE")); v != "" {
		if b, perr := strconv.ParseBool(v); perr == nil {
			cfg.Insecure = b
		} else if logger != nil {
			logger.WithField("value", v).Warn("observability: invalid OTEL_EXPORTER_OTLP_INSECURE; ignoring")
		}
	} else {
		cfg.Insecure = !hasHTTPSScheme
	}

	if v := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME")); v != "" {
		cfg.ServiceName = v
	}
	if v := strings.TrimSpace(os.Getenv("OTEL_DEPLOYMENT_ENVIRONMENT")); v != "" {
		cfg.Environment = v
	}
	if v := strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER_ARG")); v != "" {
		if ratio, perr := strconv.ParseFloat(v, 64); perr == nil {
			cfg.SamplerRatio = ratio
		} else if logger != nil {
			logger.WithField("value", v).Warn("observability: invalid OTEL_TRACES_SAMPLER_ARG; keeping default sampler")
		}
	}

	return cfg, true
}
