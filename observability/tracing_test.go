package observability

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func derefRatio(p *float64) string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", *p)
}

func TestInitTracing_NoOpWhenEndpointEmpty(t *testing.T) {
	shutdown, err := InitTracing(t.Context(), TracingConfig{}, nil)
	if err != nil {
		t.Fatalf("expected no error with empty endpoint, got %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected non-nil shutdown function")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Fatalf("no-op shutdown returned error: %v", err)
	}
}

func TestInitTracing_RequiresServiceName(t *testing.T) {
	_, err := InitTracing(t.Context(), TracingConfig{Endpoint: "localhost:4317"}, nil)
	if err == nil {
		t.Fatal("expected error when ServiceName is missing")
	}
}

func TestInitTracing_HappyPath_RegistersTracerProvider(t *testing.T) {
	// Save and restore the global TracerProvider — InitTracing mutates global
	// state and we don't want test leakage.
	prev := otel.GetTracerProvider()
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	// Use a bogus port so we don't actually try to push spans anywhere. The
	// OTLP gRPC client uses lazy-dial, so construction succeeds even if no
	// collector is listening.
	shutdown, err := InitTracing(t.Context(), TracingConfig{
		Endpoint:    "localhost:14317",
		Insecure:    true,
		ServiceName: "happy-path-test",
	}, nil)
	if err != nil {
		t.Fatalf("InitTracing returned error: %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected non-nil shutdown")
	}

	// The global TracerProvider should now be an SDK TracerProvider, not the
	// SDK's default no-op.
	if _, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); !ok {
		t.Fatalf("expected *sdktrace.TracerProvider, got %T", otel.GetTracerProvider())
	}

	// Shutdown is idempotent and bounded; allow it to error if the bogus dial
	// can't reach a collector — we only care that it returns.
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel()
	_ = shutdown(ctx)
}

func TestSelectSampler(t *testing.T) {
	cases := []struct {
		name  string
		ratio *float64
		want  string // sampler.Description()
	}{
		{"nil → AlwaysSample (safe default)", nil, "AlwaysOnSampler"},
		{"explicit 0 → NeverSample (OTel spec)", Ratio(0), "AlwaysOffSampler"},
		{"explicit 1 → AlwaysSample", Ratio(1.0), "AlwaysOnSampler"},
		{"fractional 0.5 → TraceIDRatioBased", Ratio(0.5), "TraceIDRatioBased{0.5}"},
		{"negative → AlwaysSample (treated as invalid)", Ratio(-1.0), "AlwaysOnSampler"},
		{">1 → AlwaysSample (saturated)", Ratio(1.5), "AlwaysOnSampler"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := selectSampler(tc.ratio).Description()
			if got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestInitTracingFromEnv_NoOpWhenEndpointEnvUnset(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	shutdown, err := InitTracingFromEnv(t.Context(), TracingConfig{ServiceName: "test"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected non-nil shutdown function")
	}
}

func TestTracingConfigFromEnv(t *testing.T) {
	cases := []struct {
		name     string
		env      map[string]string
		defaults TracingConfig
		want     TracingConfig
		enabled  bool
	}{
		{
			name:     "endpoint unset → disabled, defaults preserved",
			env:      map[string]string{},
			defaults: TracingConfig{ServiceName: "svc", SamplerRatio: Ratio(0.5)},
			want:     TracingConfig{ServiceName: "svc", SamplerRatio: Ratio(0.5)},
			enabled:  false,
		},
		{
			name:     "bare host:port endpoint → enabled, insecure inferred",
			env:      map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "alloy:4317"},
			defaults: TracingConfig{ServiceName: "svc"},
			want:     TracingConfig{Endpoint: "alloy:4317", Insecure: true, ServiceName: "svc"},
			enabled:  true,
		},
		{
			name:     "https URL with path → scheme + path stripped, insecure=false",
			env:      map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "https://collector:4317/v1/traces?foo=bar"},
			defaults: TracingConfig{ServiceName: "svc"},
			want:     TracingConfig{Endpoint: "collector:4317", Insecure: false, ServiceName: "svc"},
			enabled:  true,
		},
		{
			name: "OTEL_EXPORTER_OTLP_INSECURE overrides scheme heuristic",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://collector:4317",
				"OTEL_EXPORTER_OTLP_INSECURE": "true",
			},
			defaults: TracingConfig{ServiceName: "svc"},
			want:     TracingConfig{Endpoint: "collector:4317", Insecure: true, ServiceName: "svc"},
			enabled:  true,
		},
		{
			name: "OTEL_SERVICE_NAME + OTEL_DEPLOYMENT_ENVIRONMENT override defaults",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "alloy:4317",
				"OTEL_SERVICE_NAME":           "from-env",
				"OTEL_DEPLOYMENT_ENVIRONMENT": "staging",
			},
			defaults: TracingConfig{ServiceName: "default", Environment: "dev"},
			want:     TracingConfig{Endpoint: "alloy:4317", Insecure: true, ServiceName: "from-env", Environment: "staging"},
			enabled:  true,
		},
		{
			name: "valid OTEL_TRACES_SAMPLER_ARG parses into ratio",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "alloy:4317",
				"OTEL_TRACES_SAMPLER_ARG":     "0.25",
			},
			defaults: TracingConfig{ServiceName: "svc"},
			want:     TracingConfig{Endpoint: "alloy:4317", Insecure: true, ServiceName: "svc", SamplerRatio: Ratio(0.25)},
			enabled:  true,
		},
		{
			name: "OTEL_TRACES_SAMPLER_ARG=0 → explicit NeverSample (pointer set to 0)",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "alloy:4317",
				"OTEL_TRACES_SAMPLER_ARG":     "0",
			},
			defaults: TracingConfig{ServiceName: "svc"},
			want:     TracingConfig{Endpoint: "alloy:4317", Insecure: true, ServiceName: "svc", SamplerRatio: Ratio(0)},
			enabled:  true,
		},
		{
			name: "invalid OTEL_TRACES_SAMPLER_ARG keeps default ratio",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "alloy:4317",
				"OTEL_TRACES_SAMPLER_ARG":     "0,25", // comma typo
			},
			defaults: TracingConfig{ServiceName: "svc", SamplerRatio: Ratio(1.0)},
			want:     TracingConfig{Endpoint: "alloy:4317", Insecure: true, ServiceName: "svc", SamplerRatio: Ratio(1.0)},
			enabled:  true,
		},
		{
			name: "invalid OTEL_EXPORTER_OTLP_INSECURE falls back to scheme heuristic",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "alloy:4317",
				"OTEL_EXPORTER_OTLP_INSECURE": "not-a-bool",
			},
			defaults: TracingConfig{ServiceName: "svc"},
			want:     TracingConfig{Endpoint: "alloy:4317", Insecure: true, ServiceName: "svc"}, // alloy:4317 has no https:// scheme → insecure=true
			enabled:  true,
		},
		{
			name: "invalid OTEL_EXPORTER_OTLP_INSECURE with https endpoint stays secure",
			env: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://collector:4317",
				"OTEL_EXPORTER_OTLP_INSECURE": "not-a-bool",
			},
			defaults: TracingConfig{ServiceName: "svc"},
			want:     TracingConfig{Endpoint: "collector:4317", Insecure: false, ServiceName: "svc"},
			enabled:  true,
		},
	}

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear all known env vars to isolate the case.
			for _, k := range []string{
				"OTEL_EXPORTER_OTLP_ENDPOINT",
				"OTEL_EXPORTER_OTLP_INSECURE",
				"OTEL_SERVICE_NAME",
				"OTEL_DEPLOYMENT_ENVIRONMENT",
				"OTEL_TRACES_SAMPLER_ARG",
			} {
				t.Setenv(k, "")
			}
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			got, enabled := tracingConfigFromEnv(tc.defaults, logger)
			if enabled != tc.enabled {
				t.Fatalf("enabled: want %v, got %v", tc.enabled, enabled)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("config mismatch:\n  want: %+v (ratio=%v)\n  got:  %+v (ratio=%v)",
					tc.want, derefRatio(tc.want.SamplerRatio),
					got, derefRatio(got.SamplerRatio))
			}
		})
	}
}
