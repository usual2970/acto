package lib

import (
	"fmt"
	"net/http"

	handlers "github.com/usual2970/acto/internal/rest/handlers/admin"
	"github.com/usual2970/acto/internal/rest/middleware"
	actoHttp "github.com/usual2970/acto/pkg/http"
)

// RegisterAdminRoutes registers admin API endpoints using existing HTTP handlers.
// getParams: optional provider to extract path params from the request (framework-specific)
// setVars: optional setter to inject path params into request context as expected by handlers
func RegisterAdminRoutes(
	reg RouteRegistrar,
	basePath string,
) error {
	if basePath == "" {
		basePath = "/admin/v1"
	}
	svc, err := GetServices()
	if err != nil {
		return fmt.Errorf("failed to get services: %w", err)
	}

	// Use the framework-agnostic path-vars helpers by default. Routers/adapters
	// should inject path params into the request context (e.g. using
	// handlers.WithPathVars) so handlers can read them using GetPathVars.
	getParams := actoHttp.GetPathVars
	setVars := actoHttp.WithPathVars

	wrap := func(h http.HandlerFunc, needsVars bool) http.HandlerFunc {
		if !needsVars {
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

	// requireAdmin wraps a handler with admin role requirement
	requireAdmin := func(h http.Handler) http.Handler {
		return middleware.RequireAdmin(h)
	}

	if svc.AuthService != nil {
		authHandler := handlers.NewAuthHandler(svc.AuthService)
		reg.Handle(http.MethodPost, basePath+"/login", http.HandlerFunc(authHandler.Login))
	}

	if svc.PointTypeService != nil {
		// Note: handlers expect mux-style vars for {name}
		pt := handlers.NewPointTypesHandler(svc.PointTypeService)
		reg.Handle(http.MethodPost, basePath+"/point-types", requireAdmin(http.HandlerFunc(pt.Create)))
		reg.Handle(http.MethodGet, basePath+"/point-types", requireAdmin(http.HandlerFunc(pt.List)))
		reg.Handle(http.MethodPatch, basePath+"/point-types/{name}", requireAdmin(wrap(pt.Update, true)))
		reg.Handle(http.MethodDelete, basePath+"/point-types/{name}", requireAdmin(wrap(pt.Delete, true)))
	}
	if svc.BalanceService != nil {
		b := handlers.NewBalancesHandler(svc.BalanceService)
		reg.Handle(http.MethodPost, basePath+"/users/balance/credit", requireAdmin(http.HandlerFunc(b.Credit)))
		reg.Handle(http.MethodPost, basePath+"/users/balance/debit", requireAdmin(http.HandlerFunc(b.Debit)))
		reg.Handle(http.MethodGet, basePath+"/users/{userId}/transactions", requireAdmin(wrap(b.ListTransactions, true)))
	}

	if svc.RankingsService != nil {
		rk := handlers.NewRankingsHandler(svc.RankingsService)
		reg.Handle(http.MethodGet, basePath+"/rankings", requireAdmin(http.HandlerFunc(rk.Get)))
	}

	return nil
}
