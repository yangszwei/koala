// Package registry wires together the application components and provides lifecycle hooks
// for starting and shutting down the app.

package registry

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yangszwei/go-micala/config"
	httpserver "github.com/yangszwei/go-micala/internal/interface/http"
)

// App defines the application lifecycle interface, exposing methods to start and shut down the
// application server.
type App interface {
	// Run starts the application server.
	Run() error
	// Shutdown gracefully shuts down the application.
	Shutdown() error
}

// app is a concrete implementation of the App interface.
type app struct {
	server *httpserver.Server
	cfg    *config.Config
}

// NewApp creates a new App instance with the given address.
func NewApp() App {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	server := httpserver.NewServer()

	return &app{
		server: server,
		cfg:    cfg,
	}
}

// Run starts the HTTP server on the specified address.
func (a *app) Run() error {
	go func() {
		err := a.server.Run(a.cfg.Http.Addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	return a.Shutdown()
}

// Shutdown gracefully shuts down the HTTP server.
func (a *app) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.server.Shutdown(ctx)
}
