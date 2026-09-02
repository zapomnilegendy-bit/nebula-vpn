package integration

import (
	"testing"

	"github.com/nebula-vpn/backend/internal/config"
)

func TestConfigLoading(t *testing.T) {
	// Test with default configuration path
	cfg, err := config.Load("../../backend/config/development.yaml")
	if err != nil {
		t.Skipf("config file not found (expected in CI): %v", err)
	}

	if cfg == nil {
		t.Fatal("config should not be nil")
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Database.Host == "" {
		t.Error("database host should be configured")
	}

	if cfg.Auth.AccessSecret == "" || cfg.Auth.AccessSecret == "change-me" {
		t.Error("auth secrets should be configured")
	}
}

func TestServerConfigDefaults(t *testing.T) {
	dbConfig := config.DatabaseConfig{
		Host: "localhost",
		Port: 5432,
		User: "test",
		Password: "test",
		Name: "testdb",
	}

	dsn := dbConfig.DSN()
	expected := "postgres://test:test@localhost:5432/testdb?sslmode=disable"
	if dsn != expected {
		t.Errorf("DSN = %s, want %s", dsn, expected)
	}
}

func TestRedisAddr(t *testing.T) {
	redisConfig := config.RedisConfig{
		Host: "redis.example.com",
		Port: 6379,
	}

	addr := redisConfig.Addr()
	if addr != "redis.example.com:6379" {
		t.Errorf("Addr = %s, want redis.example.com:6379", addr)
	}
}
