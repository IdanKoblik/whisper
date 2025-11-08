package endpoints

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/db"
	"whisper-api/mock"

	"github.com/stretchr/testify/assert"
)

var (
	ctx      = context.Background()
	deviceID = "test123"
)

func TestKnownDeviceStatus(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	registerW := httptest.NewRecorder()
	registerReq, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	registerReq.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(registerW, registerReq)

	token := registerW.Body.String()

	db.RedisConnection(cfg).Set(ctx, "heartbeat:"+deviceID, "alive", 0)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/status/%s", deviceID), nil)
	req.Header.Set("X-Api-Token", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	db.RedisConnection(cfg).Del(ctx, "heartbeat:"+deviceID)
}

func TestUnknownDeviceStatus(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	registerW := httptest.NewRecorder()
	registerReq, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	registerReq.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(registerW, registerReq)

	token := registerW.Body.String()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/status/%s", "unknown"), nil)
	req.Header.Set("X-Api-Token", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeviceStatusUnauthorized(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/status/test-token", nil)
	req.Header.Set("X-Api-Token", "wrong-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
