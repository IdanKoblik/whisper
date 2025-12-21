package endpoints

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/config"
	"whisper-api/middleware"
	"whisper-api/mock"
	"whisper-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRegisterRouter(repo *repository.AuthRepository, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewAuthHandler(repo)
	r.Use(middleware.AdminMiddleware(cfg))
	{
		r.GET("/register", handler.Register)
	}
	return r
}

func TestRegister_Valid(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	cfg := resources.Config
	repo := repository.NewAuthRepository(
		resources.MongoClient,
		resources.RedisClient,
		resources.DB,
		resources.Collection,
	)

	r := setupRegisterRouter(repo, cfg)
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegister_Invalid(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	cfg := resources.Config
	repo := repository.NewAuthRepository(
		resources.MongoClient,
		resources.RedisClient,
		"",
		"",
	)

	r := setupRegisterRouter(repo, cfg)
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
