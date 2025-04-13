package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yangszwei/go-micala/internal/usecase/completion"
)

// CompletionHandler handles HTTP requests related to completion terms.
type CompletionHandler struct {
	svc completion.Service
}

// NewCompletionHandler creates a new handler and registers routes.
func NewCompletionHandler(r *gin.Engine, svc completion.Service) {
	h := &CompletionHandler{svc: svc}

	// Public suggestion route
	r.GET("/terms/suggest", h.Suggest)

	// Management routes
	management := r.Group("/manage/completion-terms")
	{
		management.POST("", h.Add)
		management.DELETE("/:id", h.Remove)
	}
}

// Add handles POST /manage/completion-terms
func (h *CompletionHandler) Add(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	if err := h.svc.Upload(c.Request.Context(), f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload terms"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

// Remove handles DELETE /manage/completion-terms/:id
func (h *CompletionHandler) Remove(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	if err := h.svc.Remove(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete term"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Suggest handles GET /terms/suggest?q=
func (h *CompletionHandler) Suggest(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing query"})
		return
	}

	suggestions, err := h.svc.Suggest(c.Request.Context(), q, 10)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "suggestion failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": suggestions})
}
