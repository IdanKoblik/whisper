package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAuthModel_JSONSerialization(t *testing.T) {
	auth := AuthModel{
		ApiToken: "token123",
		Devices:  []string{"device1", "device2"},
	}

	data, err := json.Marshal(auth)
	assert.NoError(t, err)

	expectedJSON := `{"_id":"token123","devices":["device1","device2"]}`
	assert.JSONEq(t, expectedJSON, string(data))

	var auth2 AuthModel
	err = json.Unmarshal(data, &auth2)
	assert.NoError(t, err)
	assert.Equal(t, auth, auth2)
}

func TestAuthModel_BSONSerialization(t *testing.T) {
	auth := AuthModel{
		ApiToken: "token456",
		Devices:  []string{"deviceA", "deviceB"},
	}

	data, err := bson.Marshal(auth)
	assert.NoError(t, err)

	var auth2 AuthModel
	err = bson.Unmarshal(data, &auth2)
	assert.NoError(t, err)
	assert.Equal(t, auth, auth2)
}
