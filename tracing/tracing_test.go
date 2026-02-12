package tracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTraceID(t *testing.T) {
	traceID := NewTraceID()
	assert.NotEmpty(t, traceID)
	assert.Len(t, traceID, 26) // ULID format: 26 characters

	// Generate another one to ensure uniqueness
	traceID2 := NewTraceID()
	assert.NotEqual(t, traceID, traceID2)
}

func TestWithTraceID_and_GetTraceID(t *testing.T) {
	ctx := t.Context()

	// Initially no trace ID
	assert.Empty(t, GetTraceID(ctx))

	// Add trace ID
	traceID := "test-trace-id-123"
	ctx = WithTraceID(ctx, traceID)

	// Should retrieve the trace ID
	assert.Equal(t, traceID, GetTraceID(ctx))
}

func TestGetTraceID_empty_context(t *testing.T) {
	ctx := t.Context()
	assert.Empty(t, GetTraceID(ctx))
}

func TestGetTraceIDFromRequest(t *testing.T) {
	traceID := "request-trace-id-456"

	// Create request with trace ID in context
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	ctx := WithTraceID(req.Context(), traceID)
	req = req.WithContext(ctx)

	assert.Equal(t, traceID, GetTraceIDFromRequest(req))
}

func TestMiddleware_generates_new_traceID(t *testing.T) {
	var capturedTraceID string

	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedTraceID = GetTraceIDFromRequest(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check trace ID was generated and stored in context
	assert.NotEmpty(t, capturedTraceID)
	assert.Len(t, capturedTraceID, 26) // ULID format

	// Check response header was set
	assert.Equal(t, capturedTraceID, rr.Header().Get(HeaderTraceID))
}

func TestMiddleware_uses_existing_traceID_from_header(t *testing.T) {
	existingTraceID := "incoming-trace-id-789"
	var capturedTraceID string

	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedTraceID = GetTraceIDFromRequest(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set(HeaderTraceID, existingTraceID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check trace ID was preserved from header
	assert.Equal(t, existingTraceID, capturedTraceID)

	// Check response header was set to the same value
	assert.Equal(t, existingTraceID, rr.Header().Get(HeaderTraceID))
}

func TestMiddleware_rejects_invalid_traceID_from_header(t *testing.T) {
	invalidTraceID := "bad\x00trace\r\n"
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := GetTraceIDFromRequest(r)
		assert.NotEqual(t, invalidTraceID, traceID, "should not use invalid header value")
		assert.Len(t, traceID, 26, "should fall back to new ULID")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set(HeaderTraceID, invalidTraceID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Response header should be the generated trace ID, not the invalid one
	respTraceID := rr.Header().Get(HeaderTraceID)
	assert.NotEqual(t, invalidTraceID, respTraceID)
	assert.Len(t, respTraceID, 26)
}

func TestMiddleware_chains_correctly(t *testing.T) {
	var firstTraceID, secondTraceID string

	// Create a chain of middleware
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstTraceID = GetTraceIDFromRequest(r)

		// Simulate calling another handler that also reads trace ID
		secondTraceID = GetTraceID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Both should have the same trace ID
	assert.Equal(t, firstTraceID, secondTraceID)
	assert.NotEmpty(t, firstTraceID)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "X-Trace-ID", HeaderTraceID)
	assert.Equal(t, "trace_id", FieldTraceID)
	assert.Equal(t, "parent_trace_id", FieldParentTraceID)
}

func TestWithParentTraceID_and_GetParentTraceID(t *testing.T) {
	ctx := t.Context()
	assert.Empty(t, GetParentTraceID(ctx))

	ctx = WithParentTraceID(ctx, "parent-123")
	assert.Equal(t, "parent-123", GetParentTraceID(ctx))
}

func TestStartSpan_no_existing_trace_id(t *testing.T) {
	ctx := t.Context()
	ctx = StartSpan(ctx)

	assert.NotEmpty(t, GetTraceID(ctx))
	assert.Len(t, GetTraceID(ctx), 26)
	assert.Empty(t, GetParentTraceID(ctx))
}

func TestStartSpan_with_existing_trace_id(t *testing.T) {
	upstream := "upstream-trace-abc"
	ctx := WithTraceID(t.Context(), upstream)
	ctx = StartSpan(ctx)

	current := GetTraceID(ctx)
	parent := GetParentTraceID(ctx)

	assert.NotEmpty(t, current)
	assert.Len(t, current, 26)
	assert.NotEqual(t, upstream, current, "new trace ID should differ from upstream")
	assert.Equal(t, upstream, parent, "upstream should be stored as parent")
}

func TestMiddleware_response_writer_passthrough(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		require.NoError(t, err)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check response is passed through correctly
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, `{"status":"ok"}`, rr.Body.String())

	// And trace ID header should also be set
	assert.NotEmpty(t, rr.Header().Get(HeaderTraceID))
}
