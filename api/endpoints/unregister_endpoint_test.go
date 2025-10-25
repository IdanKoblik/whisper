package endpoints

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

func TestUnRegisterUser_Success(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	registerW := httptest.NewRecorder()
	registerReq, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	registerReq.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(registerW, registerReq)

	token := registerW.Body.String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/unregister/%s", token), nil)
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("Deleted %s", token), w.Body.String())
}

func TestUnRegisterUser_Unauthorized(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/unregister/test-token", nil)
	req.Header.Set("X-Admin-Token", "wrong-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "Unauthorized: Invalid admin token", w.Body.String())
}

func TestUnRegisterUser_BadRequest(t *testing.T) {
	cfg := mock.ConfigMock(t)
	cfg.Mongo.Database = ""
	router := SetupRouter(cfg)

	registerW := httptest.NewRecorder()
	registerReq, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	registerReq.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(registerW, registerReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/unregister/%s", registerW.Body.String()), nil)
	req.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
