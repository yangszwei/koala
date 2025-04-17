// Package http provides the HTTP server setup using the Gin framework, including routing,
// graceful shutdown support, and hooks for modular route and shutdown configuration.

package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/yangszwei/koala/config"
)

// Server wraps a Gin engine and provides methods for running an HTTP server,
// registering routes, and handling graceful shutdown.
type Server struct {
	cfg           config.HttpConfig
	engine        *gin.Engine
	httpServer    *http.Server
	shutdownHooks []func(context.Context) error
	mu            sync.Mutex
}

// NewServer initializes a new Server instance with default middleware and routes.
func NewServer(cfg config.HttpConfig) (*Server, error) {
	engine := gin.Default()

	server := Server{
		cfg:        cfg,
		engine:     engine,
		httpServer: &http.Server{Addr: cfg.Addr, Handler: engine},
	}

	return &server, nil
}

// Run starts the HTTP server on the specified address.
func (s *Server) Run() error {
	fmt.Printf("Starting server on %s\n", s.cfg.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server and calls any registered shutdown hooks.
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, hook := range s.shutdownHooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}

	fmt.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

// RegisterShutdownHook adds a function to be called during server shutdown.
func (s *Server) RegisterShutdownHook(hook func(context.Context) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shutdownHooks = append(s.shutdownHooks, hook)
}
