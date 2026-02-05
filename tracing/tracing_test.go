package tracing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTraceID(t *testing.T) {
	traceID := NewTraceID()
	assert.NotEmpty(t, traceID)
	assert.Len(t, traceID, 36) // UUID format: 8-4-4-4-12

	// Generate another one to ensure uniqueness
	traceID2 := NewTraceID()
	assert.NotEqual(t, traceID, traceID2)
}

func TestWithTraceID_and_GetTraceID(t *testing.T) {
	ctx := context.Background()

	// Initially no trace ID
	assert.Empty(t, GetTraceID(ctx))

	// Add trace ID
	traceID := "test-trace-id-123"
	ctx = WithTraceID(ctx, traceID)

	// Should retrieve the trace ID
	assert.Equal(t, traceID, GetTraceID(ctx))
}

func TestGetTraceID_empty_context(t *testing.T) {
	ctx := context.Background()
	assert.Empty(t, GetTraceID(ctx))
}

func TestGetTraceIDFromRequest(t *testing.T) {
	traceID := "request-trace-id-456"

	// Create request with trace ID in context
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
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

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check trace ID was generated and stored in context
	assert.NotEmpty(t, capturedTraceID)
	assert.Len(t, capturedTraceID, 36) // UUID format

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

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(HeaderTraceID, existingTraceID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check trace ID was preserved from header
	assert.Equal(t, existingTraceID, capturedTraceID)

	// Check response header was set to the same value
	assert.Equal(t, existingTraceID, rr.Header().Get(HeaderTraceID))
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

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Both should have the same trace ID
	assert.Equal(t, firstTraceID, secondTraceID)
	assert.NotEmpty(t, firstTraceID)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "X-Trace-ID", HeaderTraceID)
	assert.Equal(t, "trace_id", FieldTraceID)
}

func TestMiddleware_response_writer_passthrough(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		require.NoError(t, err)
	}))

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check response is passed through correctly
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, `{"status":"ok"}`, rr.Body.String())

	// And trace ID header should also be set
	assert.NotEmpty(t, rr.Header().Get(HeaderTraceID))
}
