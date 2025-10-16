package lib

import (
	repoMysql "acto/internal/repository/mysql"
	repoRedis "acto/internal/repository/redis"
	"acto/points"
	"context"
	"database/sql"

	goRedis "github.com/redis/go-redis/v9"
)

// Handler is a minimal interface compatible with net/http.Handler.
// It allows plugging into different HTTP stacks without a hard dependency.
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

// ResponseWriter and Request are tiny adapter aliases to avoid importing net/http
// here; actual router file will import net/http and adapt accordingly.
type ResponseWriter interface{}
type Request struct{}

// Services groups initialized business services for consumers.
type Services struct {
	PointTypeService    *points.PointTypeService
	BalanceService      *points.BalanceService
	DistributionService *points.DistributionService
	RedemptionService   *points.RedemptionService
	RankingsService     points.RankingsService
	// repositories exposed for adapters
	RankingRepo   points.RankingRepository
	PointTypeRepo points.PointTypeRepository
}

// Library holds initialized services and configuration.
type Library struct {
	cfg      LibraryConfig
	services Services
}

// GetServices returns the initialized services container.
func (l *Library) GetServices() *Services { return &l.services }

// Health checks the health of initialized components by pinging DB and Redis if available.
func (l *Library) Health() error {
	if db, ok := l.cfg.DB.(*sql.DB); ok && db != nil {
		if err := db.Ping(); err != nil {
			return NewLibraryError(ErrTypeConnection, "database ping failed", err)
		}
	}
	if rc, ok := l.cfg.Redis.(goRedis.UniversalClient); ok && rc != nil {
		if err := rc.Ping(context.Background()).Err(); err != nil {
			return NewLibraryError(ErrTypeConnection, "redis ping failed", err)
		}
	}
	return nil
}

// Close releases resources if needed.
func (l *Library) Close() error { return nil }

// NewLibrary constructs a Library instance using given configuration.
func NewLibrary(cfg LibraryConfig) (*Library, error) {
	if err := cfg.validate(); err != nil {
		return nil, NewLibraryError(ErrTypeConfiguration, "invalid configuration", err)
	}

	lib := &Library{cfg: cfg}

	// Build repositories with optional overrides
	var (
		ptRepo   points.PointTypeRepository
		balRepo  points.BalanceRepository
		rankRepo points.RankingRepository
	)

	if cfg.Repositories != nil && cfg.Repositories.PointTypeRepo != nil {
		if v, ok := cfg.Repositories.PointTypeRepo.(points.PointTypeRepository); ok {
			ptRepo = v
		} else {
			return nil, NewLibraryError(ErrTypeConfiguration, "repositories.pointTypeRepo must implement PointTypeRepository", nil)
		}
	} else {
		db, ok := cfg.DB.(*sql.DB)
		if !ok || db == nil {
			return nil, NewLibraryError(ErrTypeConfiguration, "db must be *sql.DB", nil)
		}
		ptRepo = repoMysql.NewPointTypeRepository(db)
	}

	if cfg.Repositories != nil && cfg.Repositories.BalanceRepo != nil {
		if v, ok := cfg.Repositories.BalanceRepo.(points.BalanceRepository); ok {
			balRepo = v
		} else {
			return nil, NewLibraryError(ErrTypeConfiguration, "repositories.balanceRepo must implement BalanceRepository", nil)
		}
	} else {
		db, ok := cfg.DB.(*sql.DB)
		if !ok || db == nil {
			return nil, NewLibraryError(ErrTypeConfiguration, "db must be *sql.DB", nil)
		}
		balRepo = repoMysql.NewBalanceTxRepository(db)
	}

	if cfg.Repositories != nil && cfg.Repositories.RankingRepo != nil {
		if v, ok := cfg.Repositories.RankingRepo.(points.RankingRepository); ok {
			rankRepo = v
		} else {
			return nil, NewLibraryError(ErrTypeConfiguration, "repositories.rankingRepo must implement RankingRepository", nil)
		}
	} else {
		// Only support *redis.Client for default ranking repository
		if c, ok := cfg.Redis.(*goRedis.Client); ok && c != nil {
			rankRepo = repoRedis.NewRankingRepository(c)
		}
	}

	// Build services
	lib.services.PointTypeRepo = ptRepo
	lib.services.PointTypeService = points.NewPointTypeService(ptRepo)
	lib.services.BalanceService = points.NewBalanceService(balRepo, rankRepo, ptRepo)
	lib.services.RankingsService = points.NewRankingsService(rankRepo, ptRepo)
	// Optional repositories/services (exposed for adapters)
	lib.services.RankingRepo = rankRepo

	return lib, nil
}
