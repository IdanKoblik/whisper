package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr string `yaml:"uri"`
	AdminToken string `yaml:"adminToken"`
	Redis RedisConfig `yaml:"redis"`
	Mongo MongoConfig `yaml:"mongo"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
	Password string `yaml:"password"`
	DB int `yaml:"db"` 
}

type MongoConfig struct {
	ConnectionURL string `yaml:"connectionURL"`
	Database string `yaml:"database"`
}

type Reader interface {
	ReadConfig() (Config, error)
}

type ConfigReader struct{}

func (ConfigReader) ReadConfig() (Config, error) {
	env := os.Getenv("APP_ENV")
	configPath := os.Getenv("CONFIG_PATH")
	cfgFile := fmt.Sprintf("%s/config.%s.yaml", configPath, env)

	var cfg Config
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg) 
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}


