package services

import (
	"testing"
	"whisper-api/db"
	"whisper-api/mock"
	"whisper-api/utils"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	cfg := mock.ConfigMock(t)

	resp, err := RegisterUser(cfg)
	assert.NoError(t, err, "RegisterUser should not return an error")
	assert.NotEmpty(t, resp.ApiToken, "API token should not be empty")

	hash := utils.HashToken(resp.ApiToken)

	collection, err := db.MongoCollection(cfg)
	assert.NoError(t, err, "MongoCollection should not return an error")

	found, err := db.FindData(collection, hash)
	assert.NoError(t, err, "FindData should not return an error")
	assert.True(t, found, "User hash should exist in DB after registration")

	err = db.DeleteData(cfg, hash)
	assert.NoError(t, err, "DeleteData cleanup should not return an error")
}

func TestRemoveUser(t *testing.T) {
	cfg := mock.ConfigMock(t)

	resp, err := RegisterUser(cfg)
	assert.NoError(t, err, "RegisterUser should not return an error")

	err = RemoveUser(cfg, resp.ApiToken)
	assert.NoError(t, err, "RemoveUser should not return an error")

	hash := utils.HashToken(resp.ApiToken)

	collection, err := db.MongoCollection(cfg)
	assert.NoError(t, err, "MongoCollection should not return an error")

	found, err := db.FindData(collection, hash)
	assert.NoError(t, err, "FindData should not return an error")
	assert.False(t, found, "User hash should be deleted from DB")
}
