// Package tracing provides request tracing utilities for correlating logs across HTTP requests.
package tracing

import (
	"context"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

// traceIDPattern allows alphanumeric, hyphen, underscore, and dot. Used to validate inbound X-Trace-ID.
var traceIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_\-.]*$`)

const maxTraceIDLen = 128

// traceIDKey is an unexported type for context keys to prevent collisions.
type traceIDKey struct{}

// parentTraceIDKey is the context key for a parent/upstream trace ID.
type parentTraceIDKey struct{}

// HeaderTraceID is the HTTP header name for trace IDs.
const HeaderTraceID = "X-Trace-ID"

// FieldTraceID is the log field name for trace IDs.
const FieldTraceID = "trace_id"

// FieldParentTraceID is the log field name for a parent/upstream trace ID.
const FieldParentTraceID = "parent_trace_id"

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

// WithParentTraceID stores a parent/upstream trace ID in the context (e.g. from an incoming request).
func WithParentTraceID(ctx context.Context, parentTraceID string) context.Context {
	return context.WithValue(ctx, parentTraceIDKey{}, parentTraceID)
}

// GetParentTraceID retrieves the parent/upstream trace ID from the context, if any.
func GetParentTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(parentTraceIDKey{}).(string); ok {
		return id
	}
	return ""
}

// StartSpan starts a new trace span: assigns a new trace ID for this operation and, if the context
// already had a trace ID (e.g. from an upstream request), stores it as parent_trace_id so both
// can be included in logs. Use at the start of each logical operation (e.g. loop tick, handler) so
// that each call has its own trace_id while still linking back to the upstream request when present.
func StartSpan(ctx context.Context) context.Context {
	if existing := GetTraceID(ctx); existing != "" {
		ctx = WithParentTraceID(ctx, existing)
	}
	return WithTraceID(ctx, NewTraceID())
}

// GetTraceIDFromRequest retrieves the trace ID from the request context.
// Returns an empty string if no trace ID is present.
func GetTraceIDFromRequest(r *http.Request) string {
	return GetTraceID(r.Context())
}

// GetParentTraceIDFromRequest retrieves the parent/upstream trace ID from the request context, if any.
func GetParentTraceIDFromRequest(r *http.Request) string {
	return GetParentTraceID(r.Context())
}

// LoggerWithTrace returns a logger with trace_id and, when present, parent_trace_id added.
// If no trace ID is present in the context, returns the original logger.
// If logger is nil, returns nil.
func LoggerWithTrace(ctx context.Context, logger logrus.FieldLogger) logrus.FieldLogger {
	if logger == nil {
		return nil
	}
	traceID := GetTraceID(ctx)
	if traceID == "" {
		return logger
	}
	entry := logger.WithField(FieldTraceID, traceID)
	if parentID := GetParentTraceID(ctx); parentID != "" {
		entry = entry.WithField(FieldParentTraceID, parentID)
	}
	return entry
}

// validateTraceID returns the given string if it is a valid trace ID (non-empty, max 128 chars,
// only alphanumeric/hyphen/underscore/dot, no CR/LF or other control characters). Otherwise returns "".
func validateTraceID(s string) string {
	if s == "" || len(s) > maxTraceIDLen {
		return ""
	}
	for _, c := range s {
		if c < 0x20 || c == 0x7f {
			return ""
		}
	}
	if !traceIDPattern.MatchString(s) {
		return ""
	}
	return s
}

// Middleware returns an HTTP middleware that generates a trace ID for each request.
// It checks for an existing X-Trace-ID header first; if present and valid, it uses that value.
// Otherwise, it generates a new trace ID. Inbound trace IDs are validated and sanitized to prevent
// log/response header injection (allowed: alphanumeric, hyphen, underscore, dot; max 128 chars; no control chars).
// The trace ID is stored in the request context and set as a response header.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := validateTraceID(r.Header.Get(HeaderTraceID))
		if traceID == "" {
			traceID = NewTraceID()
		}

		// Add trace ID to context
		ctx := WithTraceID(r.Context(), traceID)

		// Set response header for client correlation (validated value only)
		w.Header().Set(HeaderTraceID, traceID)

		// Call next handler with traced context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
