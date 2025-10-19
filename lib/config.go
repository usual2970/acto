package lib

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/usual2970/acto/auth"
	"github.com/usual2970/acto/points"

	goRedis "github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

// LibraryConfig holds the configuration for the library

// hidden global container singleton (internal use only)
var (
	containerOnce      sync.Once
	containerSingleton *dig.Container
)

func getGlobalContainer() *dig.Container {
	containerOnce.Do(func() {
		containerSingleton = dig.New()
	})
	return containerSingleton
}

// Setup initializes library internals using optional container with provided infra.
// This keeps DI hidden while giving callers a one-shot setup.
func Setup(db *sql.DB, redis *goRedis.Client) error {
	c := getGlobalContainer()
	if err := provideConfigModule(c); err != nil {
		return err
	}

	// Provide infra
	if err := provideInfraModule(c); err != nil {
		return fmt.Errorf("provide infra module: %w", err)
	}

	// Provide modules
	if err := provideRepoModule(c); err != nil {
		return fmt.Errorf("provide repo module: %w", err)
	}
	if err := provideServiceModule(c); err != nil {
		return fmt.Errorf("provide service module: %w", err)
	}
	if err := provideDeliveryModule(c); err != nil {
		return fmt.Errorf("provide delivery module: %w", err)
	}

	return nil
}

// SetupOption allows callers to extend DI during Setup without referencing the container directly.
// Use helpers like WithProvide to register constructors.
type SetupOption func(*dig.Container) error

// WithProvide registers a constructor into the internal container during SetupWithOptions.
// Example: WithProvide(func(db *sql.DB) points.PointTypeRepository { return myRepo })
func WithProvide(constructor any) SetupOption {
	return func(c *dig.Container) error { return c.Provide(constructor) }
}

func WithRepositoryOverrides(overrides RepositoryOverrides) SetupOption {
	return func(c *dig.Container) error {
		if overrides.PointTypeRepo != nil {
			if err := c.Provide(func() points.PointTypeRepository {
				return overrides.PointTypeRepo
			}); err != nil {
				return err
			}
		}
		if overrides.BalanceRepo != nil {
			if err := c.Provide(func() points.BalanceRepository {
				return overrides.BalanceRepo
			}); err != nil {
				return err
			}
		}
		if overrides.RewardRepo != nil {
			if err := c.Provide(func() points.RewardRepository {
				return overrides.RewardRepo
			}); err != nil {
				return err
			}
		}
		if overrides.RedemptionRepo != nil {
			if err := c.Provide(func() points.RedemptionRepository {
				return overrides.RedemptionRepo
			}); err != nil {
				return err
			}
		}
		if overrides.RankingRepo != nil {
			if err := c.Provide(func() points.RankingRepository {
				return overrides.RankingRepo
			}); err != nil {
				return err
			}
		}
		return nil
	}
}

func SetupWithRepositories(overrides RepositoryOverrides) error {

	c := getGlobalContainer()

	if err := provideConfigModule(c); err != nil {
		return err
	}

	err := WithRepositoryOverrides(overrides)(c)
	if err != nil {
		return err
	}

	if err := provideServiceModule(c); err != nil {
		return err
	}

	if err := provideDeliveryModule(c); err != nil {
		return err
	}

	return nil
}

// Services holds references to all services
type Services struct {
	PointTypeService    *points.PointTypeService
	BalanceService      *points.BalanceService
	DistributionService *points.DistributionService
	RedemptionService   *points.RedemptionService
	RankingsService     points.RankingsService
	AuthService         *auth.AuthService
}

// RepositoryOverrides enables injecting custom repository implementations without exposing DI.
type RepositoryOverrides struct {
	PointTypeRepo  points.PointTypeRepository
	BalanceRepo    points.BalanceRepository
	RewardRepo     points.RewardRepository
	RedemptionRepo points.RedemptionRepository
	RankingRepo    points.RankingRepository
}

func GetServices() (*Services, error) {
	c := getGlobalContainer()
	var svc Services
	err := c.Invoke(func(
		pointTypeSvc *points.PointTypeService,
		balanceSvc *points.BalanceService,
		distributionSvc *points.DistributionService,
		redemptionSvc *points.RedemptionService,
		rankingsSvc points.RankingsService,

		authSvc *auth.AuthService,
	) {
		svc = Services{
			PointTypeService:    pointTypeSvc,
			BalanceService:      balanceSvc,
			DistributionService: distributionSvc,
			RedemptionService:   redemptionSvc,
			RankingsService:     rankingsSvc,
			AuthService:         authSvc,
		}
	})
	if err != nil {
		return nil, fmt.Errorf("get services: %w", err)
	}
	return &svc, nil
}
