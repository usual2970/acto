package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/usual2970/acto/ui"
)

// http adapters
type httpResponseWriter = http.ResponseWriter
type httpRequest = *http.Request

// RouteRegistrar is a minimal, framework-agnostic registration surface.
// Any router that can bind a path to an http.Handler can implement this.
type RouteRegistrar interface {
	Handle(method string, path string, h http.Handler)
	// NoRoute registers a handler for unmatched routes (SPA fallback / 404 handler)
	NoRoute(h http.Handler)
}

// RegisterUIRoutes registers UI static file serving routes
func RegisterUIRoutes(registrar RouteRegistrar) error {
	// Get the embedded filesystem
	distFS, err := ui.DistFS()
	if err != nil {
		return fmt.Errorf("failed to get embedded UI filesystem: %w", err)
	}

	// Helper function to serve index.html directly
	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			http.Error(w, "failed to stat index.html", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	}

	// Helper function to try serving a file from embedded FS
	tryServeFile := func(w http.ResponseWriter, r *http.Request, path string) bool {
		// Remove leading slash
		path = strings.TrimPrefix(path, "/")
		if path == "" {
			return false
		}

		file, err := distFS.Open(path)
		if err != nil {
			return false
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			return false
		}

		http.ServeContent(w, r, path, stat.ModTime(), file.(io.ReadSeeker))
		return true
	}

	// Serve index.html for the root path
	registrar.Handle("GET", "/", http.HandlerFunc(serveIndex))

	// Use NoRoute as catch-all: try to serve file, fallback to index.html for SPA
	registrar.NoRoute(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve as static file first
		if tryServeFile(w, r, r.URL.Path) {
			return
		}
		// Not a static file, serve index.html for SPA routing
		serveIndex(w, r)
	}))

	return nil
}

// RegisterRoutes registers built-in routes under the provided basePath using a generic registrar.
func RegisterRoutes(reg RouteRegistrar, basePath string) error {

	svc, err := GetServices()
	if err != nil {
		return fmt.Errorf("failed to get services: %w", err)
	}

	if basePath == "" {
		basePath = "/api/v1"
	}
	join := func(suffix string) string {
		bp := basePath
		if strings.HasSuffix(bp, "/") {
			bp = strings.TrimRight(bp, "/")
		}
		return bp + suffix
	}

	// /health
	reg.Handle(http.MethodGet, join("/health"), http.HandlerFunc(func(w httpResponseWriter, r httpRequest) {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

	}))

	// /config (redacted; only safe details)
	reg.Handle(http.MethodGet, join("/config"), http.HandlerFunc(func(w httpResponseWriter, r httpRequest) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"routes": map[string]any{
				"pathPrefix": basePath,
			},
		})
	}))

	// /services
	reg.Handle(http.MethodGet, join("/services"), http.HandlerFunc(func(w httpResponseWriter, r httpRequest) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"services": map[string]any{
				"pointType": svc.PointTypeService != nil,
				"balance":   svc.BalanceService != nil,
			},
			"repositories": map[string]any{
				// repository presence inferred from services
				"pointType": svc.PointTypeService != nil,
				"balance":   svc.BalanceService != nil,
				// ranking is optional; infer via balance service rankRepo presence is not exposed; report boolean as false if balance is nil
				"ranking": svc.BalanceService != nil,
			},
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	}))

	return nil
}
