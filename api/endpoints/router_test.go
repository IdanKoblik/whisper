package endpoints

import (
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestRouterInvalidMongo(t *testing.T) {
	os.Setenv("MONGO_CONNECTION", "mongodb://invalidhost:27017")

	router := SetupRouter()
	assert.Nil(t, router)
}
