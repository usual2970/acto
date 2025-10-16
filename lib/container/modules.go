package container

import (
	appcfg "acto/internal/config"
	repoMysql "acto/internal/repository/mysql"
	repoRedis "acto/internal/repository/redis"
	restHandlers "acto/internal/rest/handlers"
	"acto/points"
	usecases "acto/points"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	goRedis "github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

func provideModules(c *dig.Container) error {
	// Register all modules in order
	if err := provideConfigModule(c); err != nil {
		return err
	}
	if err := provideInfraModule(c); err != nil {
		return err
	}
	if err := provideRepoModule(c); err != nil {
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
