package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"whisper-api/mock"
	"whisper-api/models"
	"whisper-api/repository"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	cfg := resources.Config
	repo := &repository.AuthRepository{
		Col:      resources.MongoClient.Database(resources.DB).Collection(resources.Collection),
		Rdb:      resources.RedisClient,
		CacheTTL: 24 * time.Hour,
	}

	r := gin.Default()
	r.Use(AuthMiddleware(repo, cfg))

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "Success",
			"devices":    c.MustGet("devices"),
			"token":      c.MustGet("token"),
			"rate-limit": c.MustGet("rate-limit"),
		})
	})

	t.Run("Test missing Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "Authorization header missing"}`, w.Body.String())
	})

	t.Run("Test invalid or expired token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(API_HEADER, "invalid-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "Invalid or expired token"}`, w.Body.String())
	})

	t.Run("Test valid token", func(t *testing.T) {
		validToken := "valid-token"
		repo.CreateToken(context.Background(), &models.AuthModel{
			ApiToken: utils.HashToken(validToken),
			Devices:  []string{},
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(API_HEADER, validToken)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
