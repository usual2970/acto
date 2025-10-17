package lib

import (
	"database/sql"

	appcfg "github.com/usual2970/acto/internal/config"
	repoMysql "github.com/usual2970/acto/internal/repository/mysql"
	repoRedis "github.com/usual2970/acto/internal/repository/redis"
	restHandlers "github.com/usual2970/acto/internal/rest/handlers"
	"github.com/usual2970/acto/points"
	usecases "github.com/usual2970/acto/points"

	_ "github.com/go-sql-driver/mysql"
	goRedis "github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

// ConfigModule provides configuration
func provideConfigModule(c *dig.Container) error {
	return c.Provide(appcfg.Load)
}

// InfraModule provides infrastructure dependencies (DB, Redis)
func provideInfraModule(c *dig.Container) error {
	// DB
	if err := c.Provide(func(cfg appcfg.Config) (*sql.DB, error) {
		return sql.Open("mysql", cfg.MySQLDSN)
	}); err != nil {
		return err
	}
	// Redis
	if err := c.Provide(func(cfg appcfg.Config) *goRedis.Client {
		return goRedis.NewClient(&goRedis.Options{Addr: cfg.RedisAddr})
	}); err != nil {
		return err
	}
	return nil
}

// RepoModule provides repository implementations
func provideRepoModule(c *dig.Container) error {
	// Bind concrete repositories as their interface types using dig.As
	if err := c.Provide(repoMysql.NewPointTypeRepository, dig.As(new(points.PointTypeRepository))); err != nil {
		return err
	}
	if err := c.Provide(repoMysql.NewBalanceTxRepository, dig.As(new(points.BalanceRepository))); err != nil {
		return err
	}
	if err := c.Provide(repoMysql.NewRewardsRepository, dig.As(new(points.RewardRepository))); err != nil {
		return err
	}
	if err := c.Provide(repoMysql.NewRedemptionRepository, dig.As(new(points.RedemptionRepository))); err != nil {
		return err
	}
	if err := c.Provide(repoRedis.NewRankingRepository, dig.As(new(points.RankingRepository))); err != nil {
		return err
	}
	return nil
}

// ServiceModule provides business services
func provideServiceModule(c *dig.Container) error {
	providers := []func() error{
		func() error { return c.Provide(usecases.NewPointTypeService) },
		func() error { return c.Provide(usecases.NewBalanceService) },
		func() error { return c.Provide(usecases.NewDistributionService) },
		func() error { return c.Provide(usecases.NewRedemptionService) },
		func() error { return c.Provide(usecases.NewRankingsService) },
	}

	for _, provider := range providers {
		if err := provider(); err != nil {
			return err
		}
	}
	return nil
}

// DeliveryModule provides HTTP handlers
func provideDeliveryModule(c *dig.Container) error {
	providers := []func() error{
		func() error { return c.Provide(restHandlers.NewPointTypesHandler) },
		func() error { return c.Provide(restHandlers.NewBalancesHandler) },
		func() error { return c.Provide(restHandlers.NewDistributionsHandler) },
		func() error { return c.Provide(restHandlers.NewRedemptionsHandler) },
		func() error { return c.Provide(restHandlers.NewRankingsHandler) },
	}

	for _, provider := range providers {
		if err := provider(); err != nil {
			return err
		}
	}
	return nil
}
