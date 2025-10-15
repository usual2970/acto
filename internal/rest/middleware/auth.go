package middleware

import (
	"net/http"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// RequireRole enforces that a request has the specified role.
// For demo purposes, role is read from header X-Role. Replace with real auth.
func RequireRole(required Role, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-Role")
		if Role(role) != required {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("forbidden"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
