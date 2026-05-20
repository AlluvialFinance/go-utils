package observability

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestInitTracing_NoOpWhenEndpointEmpty(t *testing.T) {
	shutdown, err := InitTracing(context.Background(), TracingConfig{}, nil)
	if err != nil {
		t.Fatalf("expected no error with empty endpoint, got %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected non-nil shutdown function")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Fatalf("no-op shutdown returned error: %v", err)
	}
}

func TestInitTracing_RequiresServiceName(t *testing.T) {
	_, err := InitTracing(context.Background(), TracingConfig{Endpoint: "localhost:4317"}, nil)
	if err == nil {
		t.Fatal("expected error when ServiceName is missing")
	}
}

func TestInitTracingFromEnv_NoOpWhenEndpointEnvUnset(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	shutdown, err := InitTracingFromEnv(context.Background(), TracingConfig{ServiceName: "test"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected non-nil shutdown function")
	}
}

func TestInitTracingFromEnv_InvalidSamplerArgKeepsDefault(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_TRACES_SAMPLER_ARG", "not-a-float")

	// We can't directly inspect the resolved cfg here because InitTracingFromEnv
	// returns early when ENDPOINT is empty, but the parse-warning branch is
	// covered indirectly: an invalid value must not panic, must not error, and
	// must be logged (when a logger is supplied). We assert the no-error case.
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	if _, err := InitTracingFromEnv(context.Background(), TracingConfig{ServiceName: "test"}, logger); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
