package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/usual2970/acto/internal/rest/handlers"
)

// http adapters
type httpResponseWriter = http.ResponseWriter
type httpRequest = *http.Request

// RouteRegistrar is a minimal, framework-agnostic registration surface.
// Any router that can bind a path to an http.Handler can implement this.
type RouteRegistrar interface {
	Handle(method string, path string, h http.Handler)
}

// RegisterRoutes registers built-in routes under the provided basePath using a generic registrar.
func RegisterRoutes(reg RouteRegistrar, basePath string, library *Library) error {

	svc, err := library.GetServices()
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
		if err := library.Health(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": "unhealthy",
				"services": map[string]string{
					"database":     "unhealthy",
					"redis":        "unhealthy",
					"repositories": "unhealthy",
				},
				"timestamp": time.Now().Unix(),
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "healthy",
			"services": map[string]string{
				"database":     "healthy",
				"redis":        "healthy",
				"repositories": "healthy",
			},
			"timestamp": time.Now().Unix(),
		})
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

// RegisterBusinessRoutes registers business endpoints using existing HTTP handlers.
// getParams: optional provider to extract path params from the request (framework-specific)
// setVars: optional setter to inject path params into request context as expected by handlers
func RegisterBusinessRoutes(
	reg RouteRegistrar,
	basePath string,
	library *Library,
	getParams func(*http.Request) map[string]string,
	setVars func(*http.Request, map[string]string) *http.Request,
) error {
	if basePath == "" {
		basePath = "/api/v1"
	}
	svc, err := library.GetServices()
	if err != nil {
		return fmt.Errorf("failed to get services: %w", err)
	}

	wrap := func(h http.HandlerFunc, needsVars bool) http.HandlerFunc {
		if !needsVars || getParams == nil || setVars == nil {
			return h
		}
		return func(w http.ResponseWriter, r *http.Request) {
			params := getParams(r)
			if params != nil {
				r = setVars(r, params)
			}
			h(w, r)
		}
	}

	if svc.PointTypeService != nil {
		// Note: handlers expect mux-style vars for {name}
		pt := handlers.NewPointTypesHandler(svc.PointTypeService)
		reg.Handle(http.MethodPost, basePath+"/point-types", http.HandlerFunc(pt.Create))
		reg.Handle(http.MethodGet, basePath+"/point-types", http.HandlerFunc(pt.List))
		reg.Handle(http.MethodPatch, basePath+"/point-types/{name}", wrap(pt.Update, true))
		reg.Handle(http.MethodDelete, basePath+"/point-types/{name}", wrap(pt.Delete, true))
	}
	if svc.BalanceService != nil {
		b := handlers.NewBalancesHandler(svc.BalanceService)
		reg.Handle(http.MethodPost, basePath+"/users/balance/credit", http.HandlerFunc(b.Credit))
		reg.Handle(http.MethodPost, basePath+"/users/balance/debit", http.HandlerFunc(b.Debit))
		reg.Handle(http.MethodGet, basePath+"/users/{userId}/transactions", wrap(b.ListTransactions, true))
	}

	if svc.RankingsService != nil {
		rk := handlers.NewRankingsHandler(svc.RankingsService)
		reg.Handle(http.MethodGet, basePath+"/rankings", http.HandlerFunc(rk.Get))
	}

	return nil
}
