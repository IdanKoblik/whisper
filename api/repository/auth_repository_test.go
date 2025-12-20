package repository_test

import (
	"context"
	"testing"
	"time"
	"whisper-api/mock"
	"whisper-api/models"
	"whisper-api/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateToken(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)

	data := &models.AuthModel{
		ApiToken: "token123",
		Devices:  []string{"device1", "device2"},
	}

	err := repo.CreateToken(context.Background(), data)
	assert.Nil(t, err)

	var result models.AuthModel
	err = repo.Col.FindOne(context.Background(), map[string]interface{}{"_id": "token123"}).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, "token123", result.ApiToken)
}

func TestValidateTokenFromCache(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	data := &models.AuthModel{
		ApiToken: "token123",
		Devices:  []string{"device1", "device2"},
	}

	err := repo.CreateToken(context.Background(), data)
	assert.Nil(t, err)

	key := "auth:token123"
	resources.RedisClient.Set(context.Background(), key, `{"_id":"token123","devices":["device1","device2"]}`, 24*time.Hour)

	validatedData, err := repo.ValidateToken(context.Background(), "token123")
	assert.Nil(t, err)
	assert.Equal(t, "token123", validatedData.ApiToken)
}

func TestAddDeviceID(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	data := &models.AuthModel{
		ApiToken: "token123",
		Devices:  []string{"device1", "device2"},
	}

	err := repo.CreateToken(context.Background(), data)
	assert.Nil(t, err)

	err = repo.AddDeviceID(context.Background(), "token123", "device3")
	assert.Nil(t, err)

	var result models.AuthModel
	err = repo.Col.FindOne(context.Background(), map[string]interface{}{"_id": "token123"}).Decode(&result)
	assert.Nil(t, err)
	assert.Contains(t, result.Devices, "device3")
}

func TestRemoveDeviceID(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	data := &models.AuthModel{
		ApiToken: "token123",
		Devices:  []string{"device1", "device2"},
	}

	err := repo.CreateToken(context.Background(), data)
	assert.Nil(t, err)

	err = repo.RemoveDeviceID(context.Background(), "token123", "device2")
	assert.Nil(t, err)

	var result models.AuthModel
	err = repo.Col.FindOne(context.Background(), map[string]interface{}{"_id": "token123"}).Decode(&result)
	assert.Nil(t, err)
	assert.NotContains(t, result.Devices, "device2")
}

func TestRemoveToken(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	data := &models.AuthModel{
		ApiToken: "token123",
		Devices:  []string{"device1", "device2"},
	}

	err := repo.CreateToken(context.Background(), data)
	assert.Nil(t, err)

	err = repo.RemoveToken(context.Background(), "token123")
	assert.Nil(t, err)

	var result models.AuthModel
	err = repo.Col.FindOne(context.Background(), map[string]interface{}{"_id": "token123"}).Decode(&result)
	assert.NotNil(t, err)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

func TestValidateDeviceID(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	data := &models.AuthModel{
		ApiToken: "token123",
		Devices:  []string{"device1", "device2"},
	}

	err := repo.CreateToken(context.Background(), data)
	assert.Nil(t, err)

	valid, err := repo.ValidateDeviceID(context.Background(), "token123", "device1")
	assert.Nil(t, err)
	assert.True(t, valid)

	valid, err = repo.ValidateDeviceID(context.Background(), "token123", "device3")
	assert.Nil(t, err)
	assert.False(t, valid)
}
