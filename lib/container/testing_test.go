package container

import (
	"testing"
	"go.uber.org/dig"
)

// BuildTest creates a container and allows callers to override providers for tests.
func BuildTest(overrides ...func(*dig.Container) error) (*dig.Container, error) {
	c, err := Build()
	if err != nil {
		return nil, err
	}
	for _, ov := range overrides {
		if ov == nil {
			continue
		}
		if err := ov(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// TestModuleIsolation tests that modules can be registered independently
func TestModuleIsolation(t *testing.T) {
	c := dig.New()
	
	// Test config module in isolation
	if err := provideConfigModule(c); err != nil {
		t.Fatalf("ConfigModule failed: %v", err)
	}
	
	// Test infra module in isolation (requires config)
	if err := provideInfraModule(c); err != nil {
		t.Fatalf("InfraModule failed: %v", err)
	}
	
	// Test repo module in isolation (requires infra)
	if err := provideRepoModule(c); err != nil {
		t.Fatalf("RepoModule failed: %v", err)
	}
	
	// Test service module in isolation (requires repos)
	if err := provideServiceModule(c); err != nil {
		t.Fatalf("ServiceModule failed: %v", err)
	}
	
	// Test delivery module in isolation (requires services)
	if err := provideDeliveryModule(c); err != nil {
		t.Fatalf("DeliveryModule failed: %v", err)
	}
}
