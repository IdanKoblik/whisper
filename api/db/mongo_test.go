package db

import (
	"testing"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestMongoConnection_Success(t *testing.T) {
	cfg := mock.ConfigMock(t)

	client, err := MongoConnection(&cfg)
	assert.NoError(t, err, "expected no error connecting to MongoDB")
	assert.NotNil(t, client, "expected a non-nil Mongo client")
}

func TestMongoConnection_Failure(t *testing.T) {
	cfg := mock.ConfigMock(t)
	cfg.Mongo.ConnectionURL = "mongodb://invalidhost:27017"

	client, err := MongoConnection(&cfg)
	assert.Error(t, err, "expected error with invalid URI")
	assert.Nil(t, client, "expected nil client on failure")
}

