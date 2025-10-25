package endpoints

import (
	"testing"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestInvalidConfig(t *testing.T) {
	router := SetupRouter(nil)
	assert.Nil(t, router)
}

func TestValidConfig(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)
	assert.NotNil(t, router)
}
