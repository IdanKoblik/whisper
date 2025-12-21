package endpoints

import (
	"context"
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

func TestWebsocket_Success(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	clientsMutex.Lock()
	clients = make(map[string]*websocket.Conn)
	clientsMutex.Unlock()

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	h := &AuthHandler{Repo: repo}

	testToken := "valid-token"
	testDevice := "device-1"
	repo.CreateToken(context.Background(), &models.AuthModel{
		ApiToken: utils.HashToken(testToken),
		Devices:  []string{testDevice},
	})

	r := gin.New()
	r.Use(middleware.AuthMiddleware(repo, resources.Config))
	{
		r.GET("/ws", h.Websocket)
	}

	server := httptest.NewServer(r)
	defer server.Close()

	headers := http.Header{}
	headers.Set("X-API-Token", testToken)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
	assert.NoError(t, err)
	defer ws.Close()

	err = ws.WriteJSON(WebsocketRequest{Device: testDevice})
	assert.NoError(t, err)

	var resp map[string]string
	err = ws.ReadJSON(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])

	val, _ := resources.RedisClient.Get(context.Background(), DEVICE_CACHE_PREFIX+testDevice).Result()
	assert.Equal(t, "alive", val)

	clientsMutex.Lock()
	_, exists := clients[testDevice]
	clientsMutex.Unlock()
	assert.True(t, exists)
}

func TestHandleHeartbeat_Timeout(t *testing.T) {
	resources := mock.Setup(t)
	defer mock.Teardown(resources)

	clientsMutex.Lock()
	clients = make(map[string]*websocket.Conn)
	clientsMutex.Unlock()

	testDevice := "timeout-device"

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
	assert.NoError(t, err, "Dial should succeed with correct URL")
	defer wsClient.Close()

	repo := repository.NewAuthRepository(resources.MongoClient, resources.RedisClient, resources.DB, resources.Collection)
	go HandleHeartbeat(repo)

	resources.RedisClient.Del(context.Background(), DEVICE_CACHE_PREFIX+testDevice)
	time.Sleep(6 * time.Second)

	clientsMutex.Lock()
	_, exists := clients[testDevice]
	clientsMutex.Unlock()

	assert.False(t, exists, "Device should be removed from map after heartbeat timeout")

	var msg map[string]interface{}
	wsClient.SetReadDeadline(time.Now().Add(time.Second * 2))
	err = wsClient.ReadJSON(&msg)
	assert.NoError(t, err)
	assert.Equal(t, "heartbeat timeout", msg["error"])
}
