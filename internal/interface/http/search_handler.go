package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yangszwei/go-micala/internal/usecase/search"
)

// SearchHandler handles HTTP requests related to search operations.
type SearchHandler struct {
	svc search.Service
}

// NewSearchHandler creates a new handler and registers routes.
func NewSearchHandler(r *gin.Engine, svc search.Service) {
	h := &SearchHandler{svc: svc}

	r.POST("/search/index", h.Index)
	r.GET("/search", h.Search)
	r.GET("/search/categories", h.ListCategories) // New route for category listing
}

// Index handles POST /search/index to add a document.
func (h *SearchHandler) Index(c *gin.Context) {
	var doc search.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document"})
		return
	}
	if err := h.svc.Index(c.Request.Context(), doc); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "document indexed"})
}

// Search handles GET /search to perform a search.
func (h *SearchHandler) Search(c *gin.Context) {
	var q search.Query
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	results, err := h.svc.Search(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// ListCategories handles GET /search/categories to return all categories and their counts.
func (h *SearchHandler) ListCategories(c *gin.Context) {
	categories, err := h.svc.ListCategories(c.Request.Context(), c.Query("prefix"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}
