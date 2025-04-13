// Package http provides the HTTP server setup using the Gin framework, including routing,
// graceful shutdown support, and hooks for modular route and shutdown configuration.

package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Server wraps a Gin engine and provides methods for running an HTTP server,
// registering routes, and handling graceful shutdown.
type Server struct {
	engine        *gin.Engine
	addr          string
	httpServer    *http.Server
	shutdownHooks []func(context.Context) error
	mu            sync.Mutex
}

// NewServer initializes a new Server instance with default middleware and routes.
func NewServer() *Server {
	engine := gin.Default()

	return &Server{
		engine: engine,
	}
}

// Run starts the HTTP server on the specified address.
func (s *Server) Run(addr string) error {
	s.addr = addr
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	fmt.Printf("Starting server on %s\n", addr)
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
