package main

import (
	"log"

	"github.com/usual2970/acto/internal/config"
	restHandlers "github.com/usual2970/acto/internal/rest/handlers"
	"github.com/usual2970/acto/lib"

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
		reqWithVars := restHandlers.WithPathVars(c.Request, params)
		h.ServeHTTP(c.Writer, reqWithVars)
	}
	g.r.Handle(method, ginPath, gin.HandlerFunc(ginHandler))
}

func main() {
	// Load config
	cfg := config.Load()

	// Create Gin router
	r := gin.New()
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

	// Register routes using DI container
	if err := lib.RegisterRoutes(adapter, "/api/v1"); err != nil {
		log.Fatalf("failed to register library routes: %v", err)
	}

	// With Gin adapter we inject path params into the request context,
	// so call the simplified RegisterBusinessRoutes signature.
	if err := lib.RegisterBusinessRoutes(adapter, "/api/v1"); err != nil {
		log.Fatalf("failed to register business routes: %v", err)
	}

	log.Println("listening on " + cfg.HTTPAddr)
	if err := r.Run(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}
