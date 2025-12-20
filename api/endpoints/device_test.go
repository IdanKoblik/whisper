package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/config"
	"whisper-api/middleware"
	"whisper-api/mock"
	"whisper-api/models"
	"whisper-api/repository"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupDeviceRouter(handler *AuthHandler, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.AuthMiddleware(handler.Repo, cfg))
	{
		r.POST("/devices", handler.DeviceHandler)
		r.DELETE("/devices", handler.DeviceHandler)
		r.GET("/devices/:id", handler.DeviceHandler)
	}

	return r
}

func TestDevice_Add(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(
		resources.MongoClient,
		resources.RedisClient,
		resources.DB,
		resources.Collection,
	)

	cfg := mock.ConfigMock(t)
	handler := NewAuthHandler(repo)

	token := "test"
	handler.Repo.CreateToken(context.Background(), &models.AuthModel{
		ApiToken: utils.HashToken(token),
		Devices:  make([]string, 0),
	})

	request := DeviceRequest{
		Device: "testing",
	}
	body, _ := json.Marshal(request)

	r := setupDeviceRouter(handler, cfg)
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	found, err := handler.Repo.ValidateDeviceID(context.Background(), utils.HashToken(token), request.Device)
	assert.NoError(t, err)
	assert.True(t, found)
}

func TestDevice_Remove(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(
		resources.MongoClient,
		resources.RedisClient,
		resources.DB,
		resources.Collection,
	)

	cfg := mock.ConfigMock(t)
	handler := NewAuthHandler(repo)

	token := "test"
	handler.Repo.CreateToken(context.Background(), &models.AuthModel{
		ApiToken: utils.HashToken(token),
		Devices:  make([]string, 0),
	})

	request := DeviceRequest{
		Device: "testing",
	}
	body, _ := json.Marshal(request)

	err := handler.Repo.AddDeviceID(context.Background(), utils.HashToken(token), request.Device)
	assert.NoError(t, err)

	r := setupDeviceRouter(handler, cfg)
	req := httptest.NewRequest(http.MethodDelete, "/devices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	found, err := handler.Repo.ValidateDeviceID(context.Background(), utils.HashToken(token), request.Device)
	assert.NoError(t, err)
	assert.False(t, found)
}

func TestDevice_Get(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(
		resources.MongoClient,
		resources.RedisClient,
		resources.DB,
		resources.Collection,
	)

	cfg := mock.ConfigMock(t)
	handler := NewAuthHandler(repo)

	token := "test"
	handler.Repo.CreateToken(context.Background(), &models.AuthModel{
		ApiToken: utils.HashToken(token),
		Devices:  make([]string, 0),
	})

	request := DeviceRequest{
		Device: "testing",
	}

	err := handler.Repo.AddDeviceID(context.Background(), utils.HashToken(token), request.Device)
	assert.NoError(t, err)

	r := setupDeviceRouter(handler, cfg)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/devices/%s", request.Device), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	found, err := handler.Repo.ValidateDeviceID(context.Background(), utils.HashToken(token), request.Device)
	assert.NoError(t, err)
	assert.True(t, found)
}
