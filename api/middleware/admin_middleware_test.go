package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminMiddleware(t *testing.T) {
	cfg := mock.ConfigMock(t)

	r := gin.Default()
	r.Use(AdminMiddleware(cfg))

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	t.Run("Test missing Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "Authorization header missing"}`, w.Body.String())
	})

	t.Run("Test invalid Admin token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(ADMIN_HEADER, "invalid-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error": "Invalid or expired token"}`, w.Body.String())
	})

	t.Run("Test valid Admin token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(ADMIN_HEADER, cfg.AdminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Success"}`, w.Body.String())
	})
}
