package http

import (
	"context"
	"net/http"
)

// pathVarsKey is the context key used to store path variables on *http.Request.
type pathVarsKeyType struct{}

var pathVarsKey = pathVarsKeyType{}

// WithPathVars returns a copy of r with path variables stored in context.
func WithPathVars(r *http.Request, vars map[string]string) *http.Request {
	if vars == nil {
		return r
	}
	return r.WithContext(context.WithValue(r.Context(), pathVarsKey, vars))
}

// GetPathVars extracts path variables from the request context. If none are present,
// an empty map is returned.
func GetPathVars(r *http.Request) map[string]string {
	if r == nil {
		return map[string]string{}
	}
	v := r.Context().Value(pathVarsKey)
	if v == nil {
		return map[string]string{}
	}
	if m, ok := v.(map[string]string); ok {
		return m
	}
	return map[string]string{}
}
