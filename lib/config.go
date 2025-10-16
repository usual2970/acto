package lib

import (
	uc "acto/points"
	"errors"
	"net/http"
)

// LibraryConfig holds configuration for initializing the library.
// All JSON tags use camelCase per project conventions.
type LibraryConfig struct {
	DB           DBExecutor        `json:"db,omitempty"`
	Redis        RedisExecutor     `json:"redis,omitempty"`
	Repositories *RepositoryConfig `json:"repositories,omitempty"`
	Routes       *RouteConfig      `json:"routes,omitempty"`
}

// RepositoryConfig allows optional overrides of repository implementations.
// Each field must implement the corresponding use case repository interface.
type RepositoryConfig struct {
	PointTypeRepo any `json:"pointTypeRepo,omitempty"`
	BalanceRepo   any `json:"balanceRepo,omitempty"`
	RankingRepo   any `json:"rankingRepo,omitempty"`
}

// RouteConfig controls HTTP route integration.
// PathPrefix defaults to "/api/v1" when empty.
type RouteConfig struct {
	PathPrefix string           `json:"pathPrefix"`
	Middleware []MiddlewareFunc `json:"middleware,omitempty"`
	CORS       *CORSConfig      `json:"cors,omitempty"`
	// IncludeBusiness controls whether to mount business endpoints in addition to built-in routes
	IncludeBusiness bool `json:"includeBusiness,omitempty"`
	// Optional path-params adapters for framework-agnostic business routing
	ParamsGetter PathParamsGetter `json:"-"`
	ParamsSetter PathParamsSetter `json:"-"`
}

// CORSConfig defines CORS behavior for mounted routes.
type CORSConfig struct {
	AllowedOrigins   []string `json:"allowedOrigins"`
	AllowedMethods   []string `json:"allowedMethods"`
	AllowedHeaders   []string `json:"allowedHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
}

// MiddlewareFunc represents a generic HTTP middleware for mux-based routers.
// It is defined in terms of net/http to avoid framework lock-in.
type MiddlewareFunc func(next Handler) Handler

// The following tiny interfaces abstract the required capabilities for
// database and redis clients without pulling concrete types here.

// DBExecutor is satisfied by *sql.DB.
type DBExecutor interface{}

// RedisExecutor is satisfied by go-redis compatible clients.
type RedisExecutor interface{}

// PathParamsGetter extracts path parameters from a request (framework-specific)
type PathParamsGetter func(*http.Request) map[string]string

// PathParamsSetter injects path parameters back to request context as expected by handlers
type PathParamsSetter func(*http.Request, map[string]string) *http.Request

func (c *LibraryConfig) validate() error {
	if c == nil {
		return errors.New("config is nil")
	}
	if c.DB == nil {
		return errors.New("database connection required")
	}
	if c.Redis == nil {
		return errors.New("redis client required")
	}
	if c.Routes != nil {
		if c.Routes.PathPrefix == "" || c.Routes.PathPrefix[0] != '/' {
			return errors.New("routes.pathPrefix must start with '/' ")
		}
	}
	// Validate optional repository overrides implement required interfaces
	if c.Repositories != nil {
		if c.Repositories.PointTypeRepo != nil {
			if _, ok := c.Repositories.PointTypeRepo.(uc.PointTypeRepository); !ok {
				return errors.New("repositories.pointTypeRepo must implement PointTypeRepository")
			}
		}
		if c.Repositories.BalanceRepo != nil {
			if _, ok := c.Repositories.BalanceRepo.(uc.BalanceRepository); !ok {
				return errors.New("repositories.balanceRepo must implement BalanceRepository")
			}
		}
		if c.Repositories.RankingRepo != nil {
			if _, ok := c.Repositories.RankingRepo.(uc.RankingRepository); !ok {
				return errors.New("repositories.rankingRepo must implement RankingRepository")
			}
		}
	}
	return nil
}
