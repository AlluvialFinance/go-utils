// Package tracing provides request tracing utilities for correlating logs across HTTP requests.
package tracing

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

// traceIDKey is an unexported type for context keys to prevent collisions.
type traceIDKey struct{}

// HeaderTraceID is the HTTP header name for trace IDs.
const HeaderTraceID = "X-Trace-ID"

// FieldTraceID is the log field name for trace IDs.
const FieldTraceID = "trace_id"

// NewTraceID generates a new ULID-based trace ID.
func NewTraceID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// WithTraceID stores a trace ID in the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// GetTraceID retrieves the trace ID from the context.
// Returns an empty string if no trace ID is present.
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey{}).(string); ok {
		return id
	}
	return ""
}

// GetTraceIDFromRequest retrieves the trace ID from the request context.
// Returns an empty string if no trace ID is present.
func GetTraceIDFromRequest(r *http.Request) string {
	return GetTraceID(r.Context())
}

// LoggerWithTrace returns a logger with the trace ID field added.
// If no trace ID is present in the context, returns the original logger.
// If logger is nil, returns nil.
func LoggerWithTrace(ctx context.Context, logger logrus.FieldLogger) logrus.FieldLogger {
	if logger == nil {
		return nil
	}
	if traceID := GetTraceID(ctx); traceID != "" {
		return logger.WithField(FieldTraceID, traceID)
	}
	return logger
}

// Middleware returns an HTTP middleware that generates a trace ID for each request.
// It checks for an existing X-Trace-ID header first; if present, it uses that value.
// Otherwise, it generates a new trace ID.
// The trace ID is stored in the request context and set as a response header.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if trace ID already exists in request header (for distributed tracing)
		traceID := r.Header.Get(HeaderTraceID)
		if traceID == "" {
			traceID = NewTraceID()
		}

		// Add trace ID to context
		ctx := WithTraceID(r.Context(), traceID)

		// Set response header for client correlation
		w.Header().Set(HeaderTraceID, traceID)

		// Call next handler with traced context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
