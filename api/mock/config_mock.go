package mock

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"whisper-api/config"
)

func ConfigMock(t *testing.T) *config.Config {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..", "..")

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		os.Setenv("CONFIG_PATH", filepath.Join(projectRoot, "localdev"))
	}

	cfg, err := config.GetConfig()
	if err != nil {
		t.Fatal(err)
		return cfg
	}

	return cfg
}
