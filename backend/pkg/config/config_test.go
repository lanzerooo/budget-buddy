package config

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

// Тест для проверки NewTestConfig
func TestNewTestConfig(t *testing.T) {
    cfg := NewTestConfig()
    assert.NotNil(t, cfg)
    assert.Equal(t, "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable", cfg.DBUrl)
    assert.Equal(t, "test-secret", cfg.JWTSecret)
    assert.Equal(t, ":8081", cfg.FinanceServicePort)
    assert.Equal(t, ":8080", cfg.UserServicePort)
}