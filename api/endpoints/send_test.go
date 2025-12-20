package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"whisper-api/middleware"
	"whisper-api/mock"
	"whisper-api/models"
	"whisper-api/repository"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	testDevice := "iphone-xyz"
	testToken := "valid-token"

	clientsMutex.Lock()
	clients = make(map[string]*websocket.Conn)
	clientsMutex.Unlock()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		clientsMutex.Lock()
		clients[testDevice] = conn
		clientsMutex.Unlock()
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	wsClient, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer wsClient.Close()

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	repo.CreateToken(context.Background(), &models.AuthModel{
		ApiToken: utils.HashToken(testToken),
		Devices:  []string{testDevice},
	})

	h := &AuthHandler{Repo: repo}
	r := gin.New()
	r.Use(middleware.AuthMiddleware(repo, resources.Config))
	{
		r.POST("/send", h.Send)
	}

	payload := MessageRequest{
		Message: "Test Alert",
		Device:  testDevice,
		Targets: []string{"user1"},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", testToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var received MessageRequest
	wsClient.SetReadDeadline(time.Now().Add(time.Second * 2))
	err = wsClient.ReadJSON(&received)

	assert.NoError(t, err, "Should not timeout reading from WS")
	assert.Equal(t, "Test Alert", received.Message)

	payload = MessageRequest{Device: "invalid-device"}
	body, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/send", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Token", testToken)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid device id")

	for i := 0; i < resources.Config.RateLimit; i++ {
		payload := MessageRequest{Device: testDevice}
		body, _ = json.Marshal(payload)
		req, _ = http.NewRequest("POST", "/send", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Token", testToken)

		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}

	fmt.Println(w.Body.String())
	assert.Equal(t, 429, w.Code)
}

func TestSend_DeviceNotActive(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	h := &AuthHandler{Repo: repo}

	clientsMutex.Lock()
	clients = make(map[string]*websocket.Conn)
	clientsMutex.Unlock()

	testDevice := "offline-device"
	testToken := "test"
	repo.CreateToken(context.Background(), &models.AuthModel{
		ApiToken: utils.HashToken(testToken),
		Devices:  []string{testDevice},
	})

	r := gin.New()
	r.Use(middleware.AuthMiddleware(repo, resources.Config))
	{
		r.POST("/send", h.Send)
	}

	payload := MessageRequest{Device: testDevice, Message: "hello", Targets: make([]string, 0)}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/send", bytes.NewBuffer(body))
	req.Header.Set("X-API-Token", testToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Device not active", w.Body.String())
}
