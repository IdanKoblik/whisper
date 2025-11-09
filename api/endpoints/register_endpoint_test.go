package endpoints

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, w.Body.String())
}

func TestRegisterUser_Unauthorized(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	req.Header.Set("X-Admin-Token", "wrong-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Unauthorized: Invalid admin token", w.Body.String())
}

func TestRegisterUser_BadRequest(t *testing.T) {
	cfg := mock.ConfigMock(t)
	cfg.Mongo.Database = ""
	router := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
