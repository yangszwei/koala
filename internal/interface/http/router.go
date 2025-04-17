package http

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yangszwei/koala/internal/usecase/completion"
	"github.com/yangszwei/koala/internal/usecase/search"
	"github.com/yangszwei/koala/web"
)

// RoutesDeps defines the dependencies required to register HTTP routes.
type RoutesDeps struct {
	CompletionService completion.Service
	SearchService     search.Service
}

// RegisterRoutes sets up all HTTP routes, including static file serving and API endpoints.
func (s *Server) RegisterRoutes(deps RoutesDeps) {
	group := s.engine.Group(s.cfg.BasePath)

	// Web handler
	s.engine.NoRoute(NewWebHandler(s.cfg.BasePath))

	// API routes
	api := group.Group("/api")
	RegisterCompletionHandler(api, deps.CompletionService)
	RegisterSearchHandler(api, deps.SearchService)
}

// NewWebHandler returns a handler that serves static web content, excluding API routes.
func NewWebHandler(basePath string) gin.HandlerFunc {
	basePath = "/" + strings.Trim(basePath, "/")
	handler := web.Handler(basePath)
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, basePath+"/api") {
			c.Next()
			return
		}
		handler(c)
	}
}
