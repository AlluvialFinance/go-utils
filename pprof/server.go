package pprof

import (
	"context"
	"errors"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/sirupsen/logrus"
)

// Server wraps an HTTP server for pprof endpoints.
type Server struct {
	server *http.Server
	log    logrus.FieldLogger
}

// NewServer creates a new pprof server.
// Returns nil if cfg is nil or disabled.
//
// Available endpoints:
//   - GET /debug/pprof/           - index page with links to all profiles
//   - GET /debug/pprof/cmdline    - command line arguments
//   - GET /debug/pprof/profile    - CPU profile (use ?seconds=N)
//   - GET /debug/pprof/symbol     - symbol lookup
//   - GET /debug/pprof/trace      - execution trace (use ?seconds=N)
//   - GET /debug/pprof/heap       - heap profile
//   - GET /debug/pprof/goroutine  - goroutine profile
//   - GET /debug/pprof/block      - block profile
//   - GET /debug/pprof/threadcreate - thread creation profile
//   - GET /debug/pprof/mutex      - mutex profile
//   - GET /debug/pprof/allocs     - allocation profile
func NewServer(cfg *Config, log logrus.FieldLogger) *Server {
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return &Server{
		server: &http.Server{
			Addr:              cfg.Address,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
		},
		log: log,
	}
}

// Start starts the pprof server in a goroutine.
func (s *Server) Start() {
	if s == nil {
		return
	}
	go func() {
		s.log.WithField("address", s.server.Addr).Info("starting pprof server")
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.WithError(err).Error("pprof server error")
		}
	}()
}

// Stop gracefully shuts down the pprof server.
func (s *Server) Stop(ctx context.Context) error {
	if s == nil {
		return nil
	}
	s.log.Info("stopping pprof server")
	return s.server.Shutdown(ctx)
}

// Addr returns the server's listen address.
func (s *Server) Addr() string {
	if s == nil {
		return ""
	}
	return s.server.Addr
}
