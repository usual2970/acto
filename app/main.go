package main

import (
	"log"
	"net/http"

	"acto/internal/config"
	"acto/lib"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	goRedis "github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
)

// muxRegistrar adapts *mux.Router to lib.RouteRegistrar
type muxRegistrar struct{ r *mux.Router }

func (m muxRegistrar) Handle(method string, path string, h http.Handler) {
	m.r.Handle(path, h).Methods(method)
}

func main() {
	config := config.Load()

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	// Initialize shared connections
	db, err := sql.Open("mysql", config.MySQLDSN)
	if err != nil {
		log.Fatalf("failed to open mysql connection: %v", err)
	}
	redisClient := goRedis.NewClient(&goRedis.Options{Addr: config.RedisAddr})

	// Initialize library and mount its routes under /api/v1
	library, err := lib.NewLibrary(lib.LibraryConfig{DB: db, Redis: redisClient})
	if err != nil {
		log.Fatalf("failed to init library: %v", err)
	}

	// Framework-agnostic registration via small adapter
	if err := lib.RegisterRoutes(muxRegistrar{r: r}, "/api/v1", library); err != nil {
		log.Fatalf("failed to register library routes: %v", err)
	}

	// Register business routes (framework-agnostic with mux param adapters)
	getParams := func(req *http.Request) map[string]string { return mux.Vars(req) }
	if err := lib.RegisterBusinessRoutes(muxRegistrar{r: r}, "/api/v1", library, getParams, mux.SetURLVars); err != nil {
		log.Fatalf("failed to register business routes: %v", err)
	}

	log.Println("listening on " + config.HTTPAddr)
	if err := http.ListenAndServe(config.HTTPAddr, r); err != nil {
		log.Fatal(err)
	}
}
