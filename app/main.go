package main

import (
	"log"
	"net/http"

	repoMysql "acto/internal/repository/mysql"
	restHandlers "acto/internal/rest/handlers"
	usecases "acto/points"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	goRedis "github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	// Minimal DI stubs: real DSN from config in later tasks
	db, _ := sql.Open("mysql", "acto:acto@tcp(127.0.0.1:3306)/acto?parseTime=true&charset=utf8mb4&loc=Local")
	ptRepo := repoMysql.NewPointTypeRepository(db)
	ptSvc := usecases.NewPointTypeService(ptRepo)
	ptHandler := restHandlers.NewPointTypesHandler(ptSvc)

	r.HandleFunc("/api/v1/point-types", ptHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/point-types", ptHandler.List).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/point-types", ptHandler.Update).Methods(http.MethodPatch)

	// US2 wiring (ranking repo omitted for now)
	bRepo := repoMysql.NewBalanceTxRepository(db)

	// US3 wiring: Redis ranking repo
	_ = goRedis.NewClient(&goRedis.Options{Addr: "127.0.0.1:6379"})

	bSvc := usecases.NewBalanceService(bRepo, nil, ptRepo)
	bHandler := restHandlers.NewBalancesHandler(bSvc)
	r.HandleFunc("/api/v1/users/balance/credit", bHandler.Credit).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/balance/debit", bHandler.Debit).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{userId}/transactions", bHandler.ListTransactions).Methods(http.MethodGet)

	// Rankings handler (read-only for now)
	rankingsHandler := restHandlers.NewRankingsHandler(nil, ptRepo)
	r.HandleFunc("/api/v1/rankings", rankingsHandler.Get).Methods(http.MethodGet)

	// US4 wiring (distribution)
	rewardsRepo := repoMysql.NewRewardsRepository(db)
	distSvc := usecases.NewDistributionService(rewardsRepo, bRepo, nil, ptRepo)
	distHandler := restHandlers.NewDistributionsHandler(distSvc)
	r.HandleFunc("/api/v1/distributions", distHandler.Execute).Methods(http.MethodPost)

	// US5 wiring (redemptions)
	redRepo := repoMysql.NewRedemptionRepository(db)
	redSvc := usecases.NewRedemptionService(redRepo, bRepo)
	redHandler := restHandlers.NewRedemptionsHandler(redSvc)
	r.HandleFunc("/api/v1/redeem", redHandler.Redeem).Methods(http.MethodPost)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
