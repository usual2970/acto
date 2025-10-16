package main

import (
	"log"
	"net/http"

	"acto/internal/config"
	repoMysql "acto/internal/repository/mysql"
	repoRedis "acto/internal/repository/redis"
	restHandlers "acto/internal/rest/handlers"
	usecases "acto/points"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	goRedis "github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
)

func main() {
	config := config.Load()

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	// Minimal DI stubs: real DSN from config in later tasks
	db, err := sql.Open("mysql", config.MySQLDSN)
	if err != nil {
		log.Fatalf("failed to open mysql connection: %v", err)
	}
	ptRepo := repoMysql.NewPointTypeRepository(db)
	ptSvc := usecases.NewPointTypeService(ptRepo)
	ptHandler := restHandlers.NewPointTypesHandler(ptSvc)

	r.HandleFunc("/api/v1/point-types", ptHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/point-types", ptHandler.List).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/point-types/{name}", ptHandler.Update).Methods(http.MethodPatch)
	r.HandleFunc("/api/v1/point-types/{name}", ptHandler.Delete).Methods(http.MethodDelete)

	// US2 wiring (ranking repo omitted for now)
	bRepo := repoMysql.NewBalanceTxRepository(db)

	// US3 wiring: Redis ranking repo
	redisClient := goRedis.NewClient(&goRedis.Options{Addr: config.RedisAddr})
	rankingRepo := repoRedis.NewRankingRepository(redisClient)

	bSvc := usecases.NewBalanceService(bRepo, rankingRepo, ptRepo)
	bHandler := restHandlers.NewBalancesHandler(bSvc)
	r.HandleFunc("/api/v1/users/balance/credit", bHandler.Credit).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/balance/debit", bHandler.Debit).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{userId}/transactions", bHandler.ListTransactions).Methods(http.MethodGet)

	// Rankings handler (read-only for now)
	rankingsHandler := restHandlers.NewRankingsHandler(rankingRepo, ptRepo)
	r.HandleFunc("/api/v1/rankings", rankingsHandler.Get).Methods(http.MethodGet)

	// US4 wiring (distribution)
	rewardsRepo := repoMysql.NewRewardsRepository(db)
	distSvc := usecases.NewDistributionService(rewardsRepo, bRepo, rankingRepo, ptRepo)
	distHandler := restHandlers.NewDistributionsHandler(distSvc)
	r.HandleFunc("/api/v1/distributions", distHandler.Execute).Methods(http.MethodPost)

	// US5 wiring (redemptions)
	redRepo := repoMysql.NewRedemptionRepository(db)
	redSvc := usecases.NewRedemptionService(redRepo, bRepo)
	redHandler := restHandlers.NewRedemptionsHandler(redSvc)
	r.HandleFunc("/api/v1/redeem", redHandler.Redeem).Methods(http.MethodPost)

	log.Println("listening on " + config.HTTPAddr)
	if err := http.ListenAndServe(config.HTTPAddr, r); err != nil {
		log.Fatal(err)
	}
}
