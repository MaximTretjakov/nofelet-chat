package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout       = 5 * time.Second
	_defaultReadHeaderTimeout = 500 * time.Millisecond
	_defaultWriteTimeout      = 15 * time.Second
	_defaultShutdownTimeout   = 3 * time.Second
	_defaultAddress           = ":8080"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	serverCrt       string
	serverKey       string
}

func New(handler http.Handler, options ...Option) *Server {
	httpServer := &http.Server{
		Handler:           handler,
		ReadTimeout:       _defaultReadTimeout,
		WriteTimeout:      _defaultWriteTimeout,
		ReadHeaderTimeout: _defaultReadHeaderTimeout,
		Addr:              _defaultAddress,
	}

	s := Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, option := range options {
		option(&s)
	}
	s.start()
	return &s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServeTLS(s.serverCrt, s.serverKey)
		close(s.notify)
	}()
}

// Notify - provide notify channel for errors
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown - graceful shutdown server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
