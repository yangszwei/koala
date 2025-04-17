// Package web provides an HTTP handler to serve embedded frontend assets.
//
// This setup allows the Go application to serve static files (e.g., HTML, CSS, JS)
// embedded from the 'dist' directory, eliminating the need for external file dependencies.
//
// Note: Ensure the 'dist' directory exists and contains the necessary files before building.
// If the directory is missing, the build will fail due to the embed directive.
package web

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// assets embeds the contents of the 'dist' directory.
//
//go:embed dist/*
var assets embed.FS

// preloadStaticFiles loads and processes embedded files into memory with placeholder replacement.
func preloadStaticFiles(staticFS fs.FS, basePath string) map[string][]byte {
	files := make(map[string][]byte)

	baseName := basePath
	if baseName == "/" {
		baseName = "" // prevent double slashes in URLs
	}

	_ = fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".html" && ext != ".js" {
			return nil
		}
		content, err := fs.ReadFile(staticFS, path)
		if err != nil {
			return err
		}
		// Replace placeholder with actual basePath
		modified := strings.ReplaceAll(string(content), "/{{.BaseName}}", baseName)
		files[path] = []byte(modified)
		return nil
	})

	return files
}

// Handler returns a Gin middleware that serves static files from the embedded 'dist' directory.
// It accepts a basePath parameter to correctly handle applications served under a subpath.
// If a requested file is not found, it serves 'index.html' to support client-side routing in single-page applications.
func Handler(basePath string) gin.HandlerFunc {
	staticFS, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}

	processedFiles := preloadStaticFiles(staticFS, basePath)

	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, basePath) {
			c.Next()
			return
		}

		path := strings.TrimPrefix(strings.TrimPrefix(c.Request.URL.Path, basePath), "/")
		ext := filepath.Ext(path)

		if ext == "" {
			path = "index.html"
			ext = ".html"
		}

		if ext == ".html" || ext == ".js" {
			if content, ok := processedFiles[path]; ok {
				c.Data(http.StatusOK, mime.TypeByExtension(ext), content)
				return
			}
		}

		file, err := staticFS.Open(path)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		_ = file.Close()

		c.FileFromFS(path, http.FS(staticFS))
	}
}
