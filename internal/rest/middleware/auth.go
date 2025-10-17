package middleware

import (
	"net/http"

	"github.com/usual2970/acto/internal/rest/handlers"
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
			handlers.WriteError(w, 1004, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	})
}
