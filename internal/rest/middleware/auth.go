package middleware

import (
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/usual2970/acto/internal/config"
	"github.com/usual2970/acto/internal/rest/handlers"
	actoHttp "github.com/usual2970/acto/pkg/http"
)

// ...已迁移到 pkg/http/user.go ...

// RequireAdmin validates Authorization: Bearer {token} JWT and ensures role=admin.
// 认证成功后将用户信息写入 context，供后续 handler 使用。
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authz, prefix) {
			handlers.WriteError(w, 1004, "forbidden")
			return
		}
		tokenStr := strings.TrimSpace(authz[len(prefix):])
		if tokenStr == "" {
			handlers.WriteError(w, 1004, "forbidden")
			return
		}

		cfg := config.Load()
		keyFunc := func(t *jwt.Token) (any, error) { return []byte(cfg.JWTSecret), nil }

		tok, err := jwt.Parse(tokenStr, keyFunc, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !tok.Valid {
			handlers.WriteError(w, 1004, "forbidden")
			return
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			handlers.WriteError(w, 1004, "forbidden")
			return
		}
		if issCfg := cfg.JWTIssuer; issCfg != "" {
			if iss, _ := claims["iss"].(string); iss != issCfg {
				handlers.WriteError(w, 1004, "forbidden")
				return
			}
		}
		role, _ := claims["role"].(string)
		if role != "admin" {
			handlers.WriteError(w, 1004, "forbidden")
			return
		}
		username, _ := claims["sub"].(string)
		// 将用户信息写入 context
		user := &actoHttp.UserInfo{
			Username: username,
			Role:     role,
			Claims:   map[string]any{},
		}
		for k, v := range claims {
			user.Claims[k] = v
		}
		r = r.WithContext(actoHttp.WithUserInfo(r.Context(), user))
		next.ServeHTTP(w, r)
	})
}
