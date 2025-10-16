package lib

import (
	"database/sql"
	"testing"
)

// Dummy Redis client substitute to satisfy type in config for tests
type dummyRedis struct{}

func TestNewLibrary_InvalidConfig(t *testing.T) {
	_, err := NewLibrary(LibraryConfig{})
	if err == nil {
		t.Fatalf("expected error for missing DB/Redis")
	}
}

func TestNewLibrary_WithDBOnlyFails(t *testing.T) {
	db := &sql.DB{}
	_, err := NewLibrary(LibraryConfig{DB: db})
	if err == nil {
		t.Fatalf("expected error when Redis missing")
	}
}
