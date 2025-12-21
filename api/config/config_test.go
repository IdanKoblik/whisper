package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
addr: "localhost:8080"
admin_token: "secret123"
mongo:
  connection_string: "mongodb://localhost:27017"
  database: "testdb"
  collection: "users"
redis:
  addr: "localhost:6379"
  password: ""
`
	err := os.WriteFile(cfgPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	os.Setenv("CONFIG_PATH", tmpDir)
	cfg, err := GetConfig()

	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080", cfg.Addr)
	assert.Equal(t, "secret123", cfg.AdminToken)
	assert.Equal(t, 0, cfg.RateLimit)
	assert.Equal(t, "mongodb://localhost:27017", cfg.Mongo.ConnectionString)
	assert.Equal(t, "testdb", cfg.Mongo.Database)
	assert.Equal(t, "users", cfg.Mongo.Collection)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "", cfg.Redis.Password)
}

func TestReadConfig_FileNotFound(t *testing.T) {
	os.Setenv("APP_ENV", "missing")
	os.Setenv("CONFIG_PATH", "/nonexistent/path")

	_, err := GetConfig()

	assert.Error(t, err)
}

func TestReadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
uri: "localhost:8080"
mongo: invalid_yaml_here
`

	err := os.WriteFile(cfgPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	os.Setenv("APP_ENV", "dev")
	os.Setenv("CONFIG_PATH", tmpDir)

	_, err = GetConfig()

	assert.Error(t, err)
}
