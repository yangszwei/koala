// Package registry wires together the application components and provides lifecycle hooks
// for starting and shutting down the app.

package registry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yangszwei/go-micala/config"
	"github.com/yangszwei/go-micala/internal/infrastructure/elasticsearch"
	httpserver "github.com/yangszwei/go-micala/internal/interface/http"
	"github.com/yangszwei/go-micala/internal/usecase/completion"
)

// App defines the application lifecycle interface, exposing methods to start and shut down the
// application server.
type App interface {
	// Init initializes the application components.
	Init() error
	// Run starts the application server.
	Run() error
	// Shutdown gracefully shuts down the application.
	Shutdown() error
}

// app is a concrete implementation of the App interface.
type app struct {
	server *httpserver.Server
	cfg    *config.Config
	es     *elasticsearch.Client
}

// NewApp creates a new App instance.
func NewApp() App {
	return &app{
		server: httpserver.NewServer(),
	}
}

// Init initializes the application components.
func (a *app) Init() (err error) {
	// Load the configuration
	a.cfg, err = config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize the Elasticsearch client
	a.es, err = elasticsearch.NewClient(a.cfg.Elastic.Address)
	if err != nil {
		return fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// Create the Elasticsearch indices
	if err := a.es.EnsureIndices(); err != nil {
		panic(fmt.Sprintf("failed to initialize elasticsearch indices: %v", err))
	}

	// Register the HTTP server routes
	a.server.RegisterRoutes(httpserver.RoutesDeps{
		CompletionService: completion.NewService(a.es.Client),
	})

	return
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
