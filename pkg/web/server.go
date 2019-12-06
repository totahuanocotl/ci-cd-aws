package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Server contains the configuration for the HTTP server and desired service
type Server interface {
	Start() error
	Stop(timeout time.Duration) error
	RegisterEndpoint(path string, handler http.Handler)
}

type server struct {
	addr       string
	mux        *http.ServeMux
	httpServer *http.Server
}

// NewServer creates a new Server definition with an empty ServeMux
func NewServer(port int32) Server {
	return &server{
		addr: fmt.Sprintf("0.0.0.0:%d", port),
		mux:  http.NewServeMux(),
	}
}

// Start starts a new HTTP server listening at the specified port.
// Start is a blocking method that listens for SIGINT or SIGTERM to start a graceful shutdown,
// waiting for shutdownGracePeriod.
func (s *server) Start() error {
	s.httpServer = s.setup()
	log.Infof("Listening at: http://%s", s.addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("unexpected failure: %v", err)
	}
	return nil
}

func (s *server) Stop(timeout time.Duration) error {
	if s.httpServer == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Infof("Shutting server down with a grace period of %s", timeout)
	err := s.httpServer.Shutdown(ctx)
	return err
}

// RegisterService binds the http health
func (s *server) RegisterEndpoint(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

func (s *server) setup() *http.Server {
	return &http.Server{
		Addr:    s.addr,
		Handler: s.mux,
	}
}
