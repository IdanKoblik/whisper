package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongoConnection_Success(t *testing.T) {
	os.Setenv("MONGO_CONNECTION", "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=test")

	client, err := MongoConnection()
	assert.NoError(t, err, "expected no error connecting to MongoDB")
	assert.NotNil(t, client, "expected a non-nil Mongo client")
}

func TestMongoConnection_Failure(t *testing.T) {
	os.Setenv("MONGO_CONNECTION", "mongodb://invalidhost:27017")

	client, err := MongoConnection()
	assert.Error(t, err, "expected error with invalid URI")
	assert.Nil(t, client, "expected nil client on failure")
}

