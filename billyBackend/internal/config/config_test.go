package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestConfig(t *testing.T) string {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create a test configuration file
	configData := []byte(`
server:
  port: 8080
  timeout: 30
session:
  store: "redis"
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: "secret123"
tracing:
  enabled: true
  endpoint: "http://jaeger:14268/api/traces"
logging:
  level: "debug"
  encoding: "json"
`)

	err := os.WriteFile(configPath, configData, 0644)
	require.NoError(t, err)

	return configPath
}

func setupInvalidConfig(t *testing.T) string {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid-config.yaml")

	// Create an invalid YAML file with syntax errors
	invalidData := []byte(`
server:
  port: 8080
  timeout: [30
redis: {
  enabled: true,
  host: "localhost"
  port: @invalid
}
logging:
  level: debug:
    - invalid
`)

	err := os.WriteFile(configPath, invalidData, 0644)
	require.NoError(t, err)

	return configPath
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		validate    func(t *testing.T, cfg *Config, err error)
		expectError bool
	}{
		{
			name: "valid config",
			setupFunc: func(t *testing.T) string {
				return setupTestConfig(t)
			},
			validate: func(t *testing.T, cfg *Config, err error) {
				require.NoError(t, err)
				require.NotNil(t, cfg)

				// Validate Server config
				assert.Equal(t, 8080, cfg.Server.Port)
				assert.Equal(t, 30, cfg.Server.Timeout)

				// Validate Session config
				assert.Equal(t, "redis", cfg.Session.Store)

				// Validate Redis config
				assert.True(t, cfg.Redis.Enabled)
				assert.Equal(t, "localhost", cfg.Redis.Host)
				assert.Equal(t, 6379, cfg.Redis.Port)
				assert.Equal(t, "secret123", cfg.Redis.Password)

				// Validate Tracing config
				assert.True(t, cfg.Tracing.Enabled)
				assert.Equal(t, "http://jaeger:14268/api/traces", cfg.Tracing.Endpoint)

				// Validate Logging config
				assert.Equal(t, "debug", cfg.Logging.Level)
				assert.Equal(t, "json", cfg.Logging.Encoding)
			},
			expectError: false,
		},
		{
			name: "invalid yaml",
			setupFunc: func(t *testing.T) string {
				return setupInvalidConfig(t)
			},
			validate: func(t *testing.T, cfg *Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			},
			expectError: true,
		},
		{
			name: "non-existent file",
			setupFunc: func(t *testing.T) string {
				return "non-existent-config.yaml"
			},
			validate: func(t *testing.T, cfg *Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := tt.setupFunc(t)
			cfg, err := LoadConfig(configPath)
			tt.validate(t, cfg, err)
		})
	}
}

func TestConfigValidation(t *testing.T) {
	configPath := setupTestConfig(t)
	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Test that all required fields are present and correctly typed
	t.Run("server config validation", func(t *testing.T) {
		assert.GreaterOrEqual(t, cfg.Server.Port, 0)
		assert.GreaterOrEqual(t, cfg.Server.Timeout, 0)
	})

	t.Run("redis config validation", func(t *testing.T) {
		if cfg.Redis.Enabled {
			assert.NotEmpty(t, cfg.Redis.Host)
			assert.GreaterOrEqual(t, cfg.Redis.Port, 0)
		}
	})

	t.Run("logging config validation", func(t *testing.T) {
		assert.Contains(t, []string{"json", "console"}, cfg.Logging.Encoding)
		assert.Contains(t, []string{"debug", "info", "warn", "error"}, cfg.Logging.Level)
	})
}
