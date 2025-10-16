package app

import (
	"context"

	"github.com/hellofresh/health-go/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Loggaable is a service that holds a logger
type Loggable interface {
	SetLogger(logger logrus.FieldLogger)
}

// API is a service that exposes API routes
type API interface {
	RegisterHandler(mux *httprouter.Router)
}

// Middleware is a service that exposes a middleware to be set on an App
type Middleware interface {
	RegisterMiddleware(chain alice.Chain) alice.Chain
}

// Initializable is a service that can initialize
type Initializable interface {
	Init(context.Context) error
}

// Runnable is a service that maintains long living task(s)
type Runnable interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Closable is a service that needs to clean its state at the end of its execution
type Closable interface {
	Close() error
}

// Checkable is a service that can expose its health status
type Checkable interface {
	RegisterCheck(h *health.Health) error
}

// Checkable is a service that can expose metrics
type Measurable interface {
	RegisterMetrics(prometheus.Registerer) error
}

// ShutdownRequest carries context for a shutdown request.
type ShutdownRequest struct {
	Reason string // human-readable reason
	Err    error  // optional: underlying error
	Source string // optional: which service triggered it
}

// ShutdownRequester is what services use to request app shutdown.
type ShutdownRequester interface {
	RequestShutdown(ShutdownRequest)
}

// ShutdownAware is implemented by services that want a ShutdownRequester injected.
type ShutdownAware interface {
	SetShutdownRequester(ShutdownRequester)
}

// ShutdownFunc is a function adapter that satisfies ShutdownRequester.
// Use this when injecting into services to avoid handing them *App.
type ShutdownFunc func(ShutdownRequest)

func (f ShutdownFunc) RequestShutdown(r ShutdownRequest) { f(r) }
