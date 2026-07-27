package observability

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/grafana/pyroscope-go"
)

func TestInitProfiling_NoOpWhenServerAddressEmpty(t *testing.T) {
	stop, err := InitProfiling(t.Context(), ProfilingConfig{}, nil)
	if err != nil {
		t.Fatalf("expected no error with empty server address, got %v", err)
	}
	if stop == nil {
		t.Fatal("expected non-nil stop function")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	if err := stop(ctx); err != nil {
		t.Fatalf("no-op stop returned error: %v", err)
	}
}

func TestInitProfiling_RequiresServiceName(t *testing.T) {
	_, err := InitProfiling(t.Context(), ProfilingConfig{ServerAddress: "http://localhost:9998"}, nil)
	if err == nil {
		t.Fatal("expected error when ServiceName is missing")
	}
}

func TestInitProfiling_ErrorsWhenContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := InitProfiling(ctx, ProfilingConfig{
		ServerAddress: "http://localhost:9998",
		ServiceName:   "svc",
	}, nil)
	if err == nil {
		t.Fatal("expected error when context is already cancelled")
	}
}

func TestInitProfiling_HappyPath_StartsAndStops(t *testing.T) {
	// Use a bogus address and a non-CPU profile type: the SDK uploads
	// asynchronously so Start succeeds without a reachable server, and
	// avoiding ProfileCPU keeps this from contending for the process-global
	// CPU profiler under `go test`.
	stop, err := InitProfiling(t.Context(), ProfilingConfig{
		ServerAddress: "http://127.0.0.1:1",
		ServiceName:   "happy-path-test",
		ProfileTypes:  []pyroscope.ProfileType{pyroscope.ProfileAllocSpace},
	}, nil)
	if err != nil {
		t.Fatalf("InitProfiling returned error: %v", err)
	}
	if stop == nil {
		t.Fatal("expected non-nil stop")
	}
	// Stop flushes a final profile to the bogus address; allow it to error on
	// the unreachable target — we only care that it returns.
	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
	defer cancel()
	_ = stop(ctx)
}

func TestInitProfiling_StopIsIdempotent(t *testing.T) {
	// The SDK's uploader Stop closes its done channel unguarded, so a second
	// call panics unless we guard it. Callers reasonably defer stop and may
	// also call it explicitly on a shutdown path.
	stop, err := InitProfiling(t.Context(), ProfilingConfig{
		ServerAddress: "http://127.0.0.1:1",
		ServiceName:   "idempotent-stop-test",
		ProfileTypes:  []pyroscope.ProfileType{pyroscope.ProfileAllocSpace},
	}, nil)
	if err != nil {
		t.Fatalf("InitProfiling returned error: %v", err)
	}
	_ = stop(t.Context())
	_ = stop(t.Context()) // must not panic
}

func TestInitProfiling_DoesNotMutateCallerProfileTypes(t *testing.T) {
	// A caller-supplied slice with spare capacity must not be appended into:
	// the profiler would share its backing array and the caller could
	// overwrite the live profile set.
	callerTypes := make([]pyroscope.ProfileType, 1, 8)
	callerTypes[0] = pyroscope.ProfileAllocSpace

	stop, err := InitProfiling(t.Context(), ProfilingConfig{
		ServerAddress:        "http://127.0.0.1:1",
		ServiceName:          "no-mutate-test",
		ProfileTypes:         callerTypes,
		MutexProfileFraction: 5,
	}, nil)
	if err != nil {
		t.Fatalf("InitProfiling returned error: %v", err)
	}
	defer func() { _ = stop(t.Context()) }()

	// Extend over the caller's spare capacity: it must still be zero-valued.
	// (A single re-slice expression trips gosec's G602 bounds analysis.)
	full := callerTypes[:cap(callerTypes)]
	for i, pt := range full {
		if i == 0 {
			continue // the caller's own element
		}
		if pt != "" {
			t.Fatalf("InitProfiling wrote into the caller's backing array at +%d: %q", i, pt)
		}
	}
}

func TestProfilingTags(t *testing.T) {
	got := profilingTags(ProfilingConfig{
		ServiceVersion: "1.2.3",
		Environment:    "dev",
		InstanceID:     "pod-abc",
		Tags:           map[string]string{"environment": "override", "team": "staking"},
	})
	want := map[string]string{
		"service_version": "1.2.3",
		"environment":     "override", // explicit Tags win over derived
		"instance":        "pod-abc",
		"team":            "staking",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tags mismatch:\n  want: %+v\n  got:  %+v", want, got)
	}
}

func TestProfilingTags_InstanceFallsBackToHostname(t *testing.T) {
	got := profilingTags(ProfilingConfig{})
	// os.Hostname() is essentially always available in CI/dev; assert the tag
	// is populated rather than pinning a specific value.
	if got["instance"] == "" {
		t.Fatal("expected instance tag to fall back to a non-empty hostname")
	}
}

func TestInitProfilingFromEnv_NoOpWhenServerAddressEnvUnset(t *testing.T) {
	t.Setenv("PYROSCOPE_SERVER_ADDRESS", "")
	stop, err := InitProfilingFromEnv(t.Context(), ProfilingConfig{ServiceName: "test"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stop == nil {
		t.Fatal("expected non-nil stop function")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	if err := stop(ctx); err != nil {
		t.Fatalf("no-op stop returned error: %v", err)
	}
}

func TestProfilingConfigFromEnv(t *testing.T) {
	cases := []struct {
		name     string
		env      map[string]string
		defaults ProfilingConfig
		want     ProfilingConfig
		enabled  bool
	}{
		{
			name:     "server address unset → disabled, defaults preserved",
			env:      map[string]string{},
			defaults: ProfilingConfig{ServiceName: "svc", Environment: "dev"},
			want:     ProfilingConfig{ServiceName: "svc", Environment: "dev"},
			enabled:  false,
		},
		{
			name:     "server address set → enabled",
			env:      map[string]string{"PYROSCOPE_SERVER_ADDRESS": "http://alloy.telemetry.svc.cluster.local:9998"},
			defaults: ProfilingConfig{ServiceName: "svc"},
			want:     ProfilingConfig{ServerAddress: "http://alloy.telemetry.svc.cluster.local:9998", ServiceName: "svc"},
			enabled:  true,
		},
		{
			name: "PYROSCOPE_APPLICATION_NAME overrides ServiceName",
			env: map[string]string{
				"PYROSCOPE_SERVER_ADDRESS":   "http://alloy:9998",
				"PYROSCOPE_APPLICATION_NAME": "from-env",
			},
			defaults: ProfilingConfig{ServiceName: "default"},
			want:     ProfilingConfig{ServerAddress: "http://alloy:9998", ServiceName: "from-env"},
			enabled:  true,
		},
		{
			name:     "server address is trimmed",
			env:      map[string]string{"PYROSCOPE_SERVER_ADDRESS": "  http://alloy:9998  "},
			defaults: ProfilingConfig{ServiceName: "svc"},
			want:     ProfilingConfig{ServerAddress: "http://alloy:9998", ServiceName: "svc"},
			enabled:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, k := range []string{"PYROSCOPE_SERVER_ADDRESS", "PYROSCOPE_APPLICATION_NAME"} {
				t.Setenv(k, "")
			}
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			got, enabled := profilingConfigFromEnv(tc.defaults)
			if enabled != tc.enabled {
				t.Fatalf("enabled: want %v, got %v", tc.enabled, enabled)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("config mismatch:\n  want: %+v\n  got:  %+v", tc.want, got)
			}
		})
	}
}
