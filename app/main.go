package main

import (
	"log"
	"net/http"

	"github.com/usual2970/acto/internal/config"
	"github.com/usual2970/acto/lib"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	goRedis "github.com/redis/go-redis/v9"
)

// muxAdapter adapts *mux.Router to lib.RouteRegistrar without coupling lib to mux
type muxAdapter struct{ r *mux.Router }

func (m muxAdapter) Handle(method string, path string, h http.Handler) {
	m.r.Handle(path, h).Methods(method)
}

func main() {
	// Load config
	cfg := config.Load()

	// Create router
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

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

	// Create registrar adapter (mux-specific, kept out of lib to remain framework-agnostic)
	adapter := muxAdapter{r: r}

	// Register routes using DI container
	if err := lib.RegisterRoutes(adapter, "/api/v1"); err != nil {
		log.Fatalf("failed to register library routes: %v", err)
	}

	// Framework-specific helpers for path params
	getParams := func(req *http.Request) map[string]string { return mux.Vars(req) }
	setVars := func(req *http.Request, vars map[string]string) *http.Request { return mux.SetURLVars(req, vars) }

	// Register business routes using DI container
	if err := lib.RegisterBusinessRoutes(adapter, "/api/v1", getParams, setVars); err != nil {
		log.Fatalf("failed to register business routes: %v", err)
	}

	log.Println("listening on " + cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, r); err != nil {
		log.Fatal(err)
	}
}
