package container

import (
	"database/sql"
	"fmt"

	goRedis "github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

// Build builds and returns a DI container with all modules registered.
func Build() (*dig.Container, error) {
	c := dig.New()
	if err := provideModules(c); err != nil {
		return nil, fmt.Errorf("failed to provide modules: %w", err)
	}

	// Validate container can resolve core dependencies
	if err := validateContainer(c); err != nil {
		return nil, fmt.Errorf("container validation failed: %w", err)
	}

	return c, nil
}

// validateContainer checks that core dependencies can be resolved
func validateContainer(c *dig.Container) error {
	// Test config resolution
	var config interface{}
	if err := c.Invoke(func(cfg interface{}) { config = cfg }); err != nil {
		return fmt.Errorf("config resolution failed: %w", err)
	}
	_ = config // avoid unused variable

	// Test DB resolution
	var db interface{}
	if err := c.Invoke(func(d interface{}) { db = d }); err != nil {
		return fmt.Errorf("database resolution failed: %w", err)
	}
	_ = db // avoid unused variable

	return nil
}

// BuildWithInfra builds a container but binds provided DB and Redis clients
// instead of the defaults from configuration. This allows callers to supply
// their own infrastructure while keeping the container hidden behind setup APIs.
func BuildWithInfra(db *sql.DB, redis *goRedis.Client) (*dig.Container, error) {
	c := dig.New()
	// Provide fixed infra
	if err := c.Provide(func() *sql.DB { return db }); err != nil {
		return nil, fmt.Errorf("provide db: %w", err)
	}
	if err := c.Provide(func() *goRedis.Client { return redis }); err != nil {
		return nil, fmt.Errorf("provide redis: %w", err)
	}
	// Provide remaining modules (skip provideConfigModule/provideInfraModule)
	if err := provideRepoModule(c); err != nil {
		return nil, fmt.Errorf("provide repos: %w", err)
	}
	if err := provideServiceModule(c); err != nil {
		return nil, fmt.Errorf("provide services: %w", err)
	}
	if err := provideDeliveryModule(c); err != nil {
		return nil, fmt.Errorf("provide delivery: %w", err)
	}
	return c, nil
}

func BuildWithoutInfra() (*dig.Container, error) {
	c := dig.New()
	if err := provideModules(c); err != nil {
		return nil, fmt.Errorf("failed to provide modules: %w", err)
	}

	// Provide remaining modules (skip provideConfigModule/provideInfraModule)
	if err := provideRepoModule(c); err != nil {
		return nil, fmt.Errorf("provide repos: %w", err)
	}
	if err := provideServiceModule(c); err != nil {
		return nil, fmt.Errorf("provide services: %w", err)
	}
	if err := provideDeliveryModule(c); err != nil {
		return nil, fmt.Errorf("provide delivery: %w", err)
	}

	return c, nil
}
