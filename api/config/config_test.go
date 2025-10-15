package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.dev.yaml")

	yamlContent := `
uri: "localhost:8080"
adminToken: "secret123"
mongo:
  connectionURL: "mongodb://localhost:27017"
  database: "testdb"
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
`

	err := os.WriteFile(cfgPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	os.Setenv("APP_ENV", "dev")
	os.Setenv("CONFIG_PATH", tmpDir)

	reader := ConfigReader{}
	cfg, err := reader.ReadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080", cfg.Addr)
	assert.Equal(t, "secret123", cfg.AdminToken)
	assert.Equal(t, "mongodb://localhost:27017", cfg.Mongo.ConnectionURL)
	assert.Equal(t, "testdb", cfg.Mongo.Database)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "", cfg.Redis.Password)
	assert.Equal(t, 0, cfg.Redis.DB)
}

func TestReadConfig_FileNotFound(t *testing.T) {
	os.Setenv("APP_ENV", "missing")
	os.Setenv("CONFIG_PATH", "/nonexistent/path")

	reader := ConfigReader{}
	_, err := reader.ReadConfig()

	assert.Error(t, err)
}

func TestReadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.dev.yaml")

	yamlContent := `
uri: "localhost:8080"
mongo: invalid_yaml_here
`

	err := os.WriteFile(cfgPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	os.Setenv("APP_ENV", "dev")
	os.Setenv("CONFIG_PATH", tmpDir)

	reader := ConfigReader{}
	_, err = reader.ReadConfig()

	assert.Error(t, err)
}

