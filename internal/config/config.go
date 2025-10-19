package config

import (
	"os"
	"sync"
)

type Config struct {
	MySQLDSN  string
	RedisAddr string
	HTTPAddr  string
	// Auth/JWT settings
	AuthUsername string
	AuthPassword string
	JWTSecret    string
	JWTIssuer    string
	JWTTTL       string // duration string, e.g. "1h", "30m"
}

var (
	loadedOnce sync.Once
	cachedCfg  Config
)

// Load returns the process configuration. It is loaded from environment
// variables only once and cached for subsequent calls.
func Load() Config {
	loadedOnce.Do(func() {
		cachedCfg = Config{
			MySQLDSN:  getenv("MYSQL_DSN", "acto:acto@tcp(127.0.0.1:33061)/acto?parseTime=true&charset=utf8mb4&loc=Local"),
			RedisAddr: getenv("REDIS_ADDR", "127.0.0.1:63791"),
			HTTPAddr:  getenv("HTTP_ADDR", ":1314"),
			// Defaults are intended for development only
			AuthUsername: getenv("AUTH_USERNAME", "admin@example.com"),
			AuthPassword: getenv("AUTH_PASSWORD", "admin123"),
			JWTSecret:    getenv("JWT_SECRET", "dev-secret"),
			JWTIssuer:    getenv("JWT_ISSUER", "acto-auth"),
			JWTTTL:       getenv("JWT_TTL", "1h"),
		}
	})
	return cachedCfg
}

// Current is an alias of Load for clarity at call sites; it returns the
// cached configuration and will trigger a one-time load if not initialized.
func Current() Config { return Load() }

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
