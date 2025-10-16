package lib

import (
	"database/sql"
	"fmt"
	"sync"

	"acto/lib/container"
	"acto/points"

	goRedis "github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

// LibraryConfig holds the configuration for the library
type LibraryConfig struct {
	DB    *sql.DB
	Redis *goRedis.Client
	// Optional repository overrides for library usage without DI
	Repositories *RepositoryOverrides
}

// Library represents the main library instance
type Library struct {
	config   LibraryConfig
	services Services
}

// hidden global container singleton (internal use only)
var (
	containerOnce      sync.Once
	containerSingleton *dig.Container
)

func initGlobalContainer(c *dig.Container) {
	containerOnce.Do(func() {
		containerSingleton = c
	})
}

func getGlobalContainer() *dig.Container { return containerSingleton }

// NewLibraryFromContainer creates a new library instance from DI container
func NewLibraryFromContainer(c *dig.Container) (*Library, error) {
	if c == nil {
		c = getGlobalContainer()
		if c == nil {
			return nil, fmt.Errorf("NewLibraryFromContainer: container is nil")
		}
	}
	var db *sql.DB
	var redis *goRedis.Client

	if err := c.Invoke(func(d *sql.DB, r *goRedis.Client) {
		db = d
		redis = r
	}); err != nil {
		return nil, err
	}

	return &Library{config: LibraryConfig{DB: db, Redis: redis}}, nil
}

// Setup initializes library internals using optional container with provided infra.
// This keeps DI hidden while giving callers a one-shot setup.
func Setup(db *sql.DB, redis *goRedis.Client) (*Library, error) {
	c, err := container.BuildWithInfra(db, redis)
	if err != nil {
		return nil, err
	}
	initGlobalContainer(c)
	lib, err := NewLibraryFromContainer(c)
	if err != nil {
		return nil, err
	}
	return lib, nil
}

// SetupOption allows callers to extend DI during Setup without referencing the container directly.
// Use helpers like WithProvide to register constructors.
type SetupOption func(*dig.Container) error

// WithProvide registers a constructor into the internal container during SetupWithOptions.
// Example: WithProvide(func(db *sql.DB) points.PointTypeRepository { return myRepo })
func WithProvide(constructor any) SetupOption {
	return func(c *dig.Container) error { return c.Provide(constructor) }
}

// SetupWithOptions builds the internal DI container with provided infra, then applies options
// (such as custom repositories/services) before returning the Library. The container remains hidden.
func SetupWithOptions(db *sql.DB, redis *goRedis.Client, options ...SetupOption) (*Library, error) {
	c, err := container.BuildWithInfra(db, redis)
	if err != nil {
		return nil, err
	}
	for _, opt := range options {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	initGlobalContainer(c)
	return NewLibraryFromContainer(c)
}

// Health checks the health of the library services
func (l *Library) Health() error {
	// Check database connection
	if err := l.config.DB.Ping(); err != nil {
		return err
	}

	// Check Redis connection
	if err := l.config.Redis.Ping(nil).Err(); err != nil {
		return err
	}

	return nil
}

// GetServices returns the services available in the library
func (l *Library) GetServices() *Services {
	c := getGlobalContainer()

	var (
		pointTypeService *points.PointTypeService
		balanceService   *points.BalanceService
		distributionSvc  *points.DistributionService
		redemptionSvc    *points.RedemptionService
		rankingsService  points.RankingsService
	)

	err := c.Invoke(func(
		resolvedPointTypeService *points.PointTypeService,
		resolvedBalanceService *points.BalanceService,
		resolvedDistributionService *points.DistributionService,
		resolvedRedemptionService *points.RedemptionService,
		resolvedRankingsService points.RankingsService,
	) {
		pointTypeService = resolvedPointTypeService
		balanceService = resolvedBalanceService
		distributionSvc = resolvedDistributionService
		redemptionSvc = resolvedRedemptionService
		rankingsService = resolvedRankingsService
	})
	if err != nil {
		panic(err)
	}

	return &Services{
		PointTypeService:    pointTypeService,
		BalanceService:      balanceService,
		DistributionService: distributionSvc,
		RedemptionService:   redemptionSvc,
		RankingsService:     rankingsService,
	}
}

// Services holds references to all services
type Services struct {
	PointTypeService    *points.PointTypeService
	BalanceService      *points.BalanceService
	DistributionService *points.DistributionService
	RedemptionService   *points.RedemptionService
	RankingsService     points.RankingsService
}

// RepositoryOverrides enables injecting custom repository implementations without exposing DI.
type RepositoryOverrides struct {
	PointTypeRepo  points.PointTypeRepository
	BalanceRepo    points.BalanceRepository
	RewardRepo     points.RewardRepository
	RedemptionRepo points.RedemptionRepository
	RankingRepo    points.RankingRepository
}
