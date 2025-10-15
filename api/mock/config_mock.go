package mock

import (
	"os"
	"testing"
	"runtime"
	"path/filepath"
	"whisper-api/config"
)

func ConfigMock(t *testing.T) (config.Config) {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..", "..")

	os.Setenv("APP_ENV", "test")
	os.Setenv("CONFIG_PATH", filepath.Join(projectRoot, "localdev"))

	cfg, err := config.ConfigReader{}.ReadConfig()
	if err != nil {
		t.Fatal(err)
		return cfg 
	}

	return cfg
}
