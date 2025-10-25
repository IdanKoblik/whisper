package main

import (
	"testing"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestInvalidConfig(t *testing.T) {
	router := SetupRouter(nil)
	assert.Nil(t, router)
}

func TestRouterInvalidMongo(t *testing.T) {
	cfg := mock.ConfigMock(t)
	cfg.Mongo.ConnectionString = "invalid"

	router := SetupRouter(cfg)
	assert.Nil(t, router)
}

func TestRouterInvalidRedis(t *testing.T) {
	cfg := mock.ConfigMock(t)
	cfg.Redis.Addr = "invalid"

	router := SetupRouter(cfg)
	assert.Nil(t, router)
}
