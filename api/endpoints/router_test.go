package endpoints

import (
	"testing"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestRouterInvalidMongo(t *testing.T) {
	cfg := mock.ConfigMock(t)
	cfg.Mongo.ConnectionURL = "mongodb://invalidhost:27017"

	router := SetupRouter(&cfg)
	assert.Nil(t, router)
}
