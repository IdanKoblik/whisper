package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/middleware"
	"whisper-api/mock"
	"whisper-api/models"
	"whisper-api/repository"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUnRegister(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	cfg := resources.Config
	repo := repository.NewAuthRepository(
		resources.MongoClient,
		resources.RedisClient,
		resources.DB,
		resources.Collection,
	)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewAuthHandler(repo)
	r.Use(middleware.AdminMiddleware(cfg))
	{
		r.DELETE("/unregister/:token", handler.UnRegister)
	}

	token := "test-token"
	user := &models.AuthModel{
		ApiToken: utils.HashToken(token),
		Devices:  []string{"device1", "device2"},
	}

	handler.Repo.CreateToken(context.Background(), user)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/unregister/%s", token), nil)
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req2 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/unregister/%s", token), nil)
	req2.Header.Set("X-Admin-Token", cfg.AdminToken)
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}
