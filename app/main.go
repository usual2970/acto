package main

import (
	"log"

	"github.com/usual2970/acto/internal/config"
	"github.com/usual2970/acto/lib"
	actoHttp "github.com/usual2970/acto/pkg/http"

	"database/sql"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	goRedis "github.com/redis/go-redis/v9"
)

// muxAdapter adapts *mux.Router to lib.RouteRegistrar without coupling lib to mux
// ginAdapter adapts gin.Engine to lib.RouteRegistrar while injecting path
// variables into the *http.Request context so handlers can use GetPathVars.
type ginAdapter struct{ r *gin.Engine }

var braceRe = regexp.MustCompile(`\{([^/}]+)\}`)

func convertPath(path string) string {
	return braceRe.ReplaceAllString(path, ":$1")
}

func (g ginAdapter) Handle(method string, path string, h http.Handler) {
	ginPath := convertPath(path)
	ginHandler := func(c *gin.Context) {
		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}
		// inject path vars into request so handlers using GetPathVars work
		reqWithVars := actoHttp.WithPathVars(c.Request, params)
		h.ServeHTTP(c.Writer, reqWithVars)
	}
	g.r.Handle(method, ginPath, gin.HandlerFunc(ginHandler))
}

func (g ginAdapter) NoRoute(h http.Handler) {
	ginHandler := func(c *gin.Context) {
		// pass through request without modifying path vars
		h.ServeHTTP(c.Writer, c.Request)
	}
	g.r.NoRoute(gin.HandlerFunc(ginHandler))
}

func main() {
	// Load config
	cfg := config.Load()

	// Create Gin router
	r := gin.New()
	// Global CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Role")
		c.Header("Vary", "Origin")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// Initialize shared connections
	db, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("failed to open mysql connection: %v", err)
	}
	redisClient := goRedis.NewClient(&goRedis.Options{Addr: cfg.RedisAddr})

	// Create library via one-shot setup (hides internal DI)
	if err := lib.Setup(db, redisClient); err != nil {
		log.Fatalf("failed to init library: %v", err)
	}

	// Create registrar adapter for Gin
	adapter := ginAdapter{r: r}

	// Register API routes first (specific routes)
	if err := lib.RegisterRoutes(adapter, "/api/v1"); err != nil {
		log.Fatalf("failed to register library routes: %v", err)
	}

	// With Gin adapter we inject path params into the request context,
	// so call the simplified RegisterApiRoutes signature.
	if err := lib.RegisterApiRoutes(adapter, "/api/v1"); err != nil {
		log.Fatalf("failed to register api routes: %v", err)
	}

	// Register admin routes
	if err := lib.RegisterAdminRoutes(adapter, "/admin/v1"); err != nil {
		log.Fatalf("failed to register admin routes: %v", err)
	}

	// Register UI routes last (catch-all route must be last in Gin)
	if err := lib.RegisterUIRoutes(adapter); err != nil {
		log.Fatalf("failed to register ui routes: %v", err)
	}

	log.Println("listening on " + cfg.HTTPAddr)
	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}
