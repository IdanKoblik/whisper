package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr       string `yaml:"addr"`
	AdminToken string `yaml:"admin_token"`
	RateLimit  *int    `yaml:"rate_limit"`

	Redis RedisConfig `yaml:"redis"`
	Mongo MongoConfig `yaml:"mongo"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MongoConfig struct {
	ConnectionString string `yaml:"connection_string"`
	Database         string `yaml:"database"`
}

func GetConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	configPath := os.Getenv("CONFIG_PATH")
	cfgFile := fmt.Sprintf("%s/config.%s.yaml", configPath, env)

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
