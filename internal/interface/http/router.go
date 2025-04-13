package http

import (
	"github.com/yangszwei/go-micala/internal/usecase/completion"
)

// RoutesDeps bundles all service dependencies for HTTP handlers.
type RoutesDeps struct {
	CompletionService completion.Service
}

// RegisterRoutes configures all HTTP routes for the app.
func (s *Server) RegisterRoutes(deps RoutesDeps) {
	NewCompletionHandler(s.engine, deps.CompletionService)
}
