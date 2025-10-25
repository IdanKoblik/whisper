package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"whisper-api/communication"
	"whisper-api/mock"
	"whisper-api/services"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	deviceID := "device-420"
	registerW := httptest.NewRecorder()
	registerReq, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	registerReq.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(registerW, registerReq)

	token := registerW.Body.String()

	mockConn := &mock.MockConn{}
	communication.Clients = map[string]communication.Conn{
		deviceID: mockConn,
	}

	validRequestBody := services.MessageRequest{
		DeviceID:    deviceID,
		Message:     "test",
		Subscribers: []string{"sub1", "sub2", "sub3"},
	}

	validBodyBytes, err := json.Marshal(validRequestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	validW := httptest.NewRecorder()
	validReq, _ := http.NewRequest("POST", "/api/send", bytes.NewReader(validBodyBytes))
	validReq.Header.Set("X-Api-Token", token)
	router.ServeHTTP(validW, validReq)

	assert.Equal(t, http.StatusOK, validW.Code)
	assert.Equal(t, "Message sent", validW.Body.String())

	invalidRequestBody := services.MessageRequest{
		DeviceID: deviceID,
	}

	invalidBodyBytes, err := json.Marshal(invalidRequestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	invalidW := httptest.NewRecorder()
	invalidReq, _ := http.NewRequest("POST", "/api/send", bytes.NewReader(invalidBodyBytes))
	validReq.Header.Set("X-Api-Token", token)
	router.ServeHTTP(invalidW, invalidReq)

	assert.Equal(t, http.StatusBadRequest, invalidW.Code)
}

func TestSendMessage_Unauthorized(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	requestBody := services.MessageRequest{
		DeviceID:    "000",
		Message:     "test",
		Subscribers: []string{"sub1", "sub2", "sub3"},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/send", bytes.NewReader(bodyBytes))
	req.Header.Set("X-Api-Token", "invalid")
	router.ServeHTTP(w, req)
}

func TestSendMessage_RateLimit(t *testing.T) {
	cfg := mock.ConfigMock(t)
	router := SetupRouter(cfg)

	deviceID := "device-420"
	registerW := httptest.NewRecorder()
	registerReq, _ := http.NewRequest("POST", "/api/admin/register", bytes.NewBuffer([]byte{}))
	registerReq.Header.Set("X-Admin-Token", cfg.AdminToken)
	router.ServeHTTP(registerW, registerReq)

	token := registerW.Body.String()

	mockConn := &mock.MockConn{}
	communication.Clients = map[string]communication.Conn{
		deviceID: mockConn,
	}

	validRequestBody := services.MessageRequest{
		DeviceID:    deviceID,
		Message:     "test",
		Subscribers: []string{"sub1", "sub2", "sub3"},
	}

	validBodyBytes, err := json.Marshal(validRequestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	for i := 0; i < *cfg.RateLimit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/send", bytes.NewReader(validBodyBytes))
		req.Header.Set("X-Api-Token", token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Message sent", w.Body.String())
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/send", bytes.NewReader(validBodyBytes))
	req.Header.Set("X-Api-Token", token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 429, w.Code)
	assert.Equal(t, "Rate limit exceeded. Try again in 60 seconds", w.Body.String())
}
