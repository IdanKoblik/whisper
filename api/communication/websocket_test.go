package communication

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"whisper-api/db"
	"whisper-api/mock"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type MockUser struct {
	ApiToken string `bson:"_id" json:"id"`
}

func newTestServer(t *testing.T) (*httptest.Server, func()) {
	cfg := mock.ConfigMock(t)
	router := gin.Default()

	db.InsertData(cfg, MockUser{ApiToken: "good-token"}, "good-token")
	router.GET("/ws", func(c *gin.Context) {
		HandleWebsocket(cfg, c)
	})

	server := httptest.NewServer(router)
	return server, func() { server.Close() }
}

func TestHandleWebsocket_MissingToken(t *testing.T) {
	server, cleanup := newTestServer(t)
	defer cleanup()

	resp, err := http.Get(server.URL + "/ws")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestHandleWebsocket_InvalidDevice(t *testing.T) {
	server, cleanup := newTestServer(t)
	defer cleanup()

	wsURL := "ws" + server.URL[len("http"):] + "/ws"
	header := http.Header{}
	header.Add("X-Api-Token", "good-token")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	assert.NoError(t, err)
	defer conn.Close()

	req := ConnectRequest{DeviceID: "bad-device"}
	data, _ := json.Marshal(req)
	err = conn.WriteMessage(websocket.TextMessage, data)
	assert.NoError(t, err)

	clientsMutex.Lock()
	_, exists := Clients["bad-device"]
	clientsMutex.Unlock()
	assert.False(t, exists, "invalid device should not be added to Clients map")
}

func TestHandleWebsocket_Heartbeat(t *testing.T) {
	server, cleanup := newTestServer(t)
	defer cleanup()

	wsURL := "ws" + server.URL[len("http"):] + "/ws"
	header := http.Header{}
	header.Add("X-Api-Token", "good-token")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	assert.NoError(t, err)
	defer conn.Close()

	req := HeartbeatRequest{
		DeviceID: "test-device",
		Ping:     "ping",
	}
	data, _ := json.Marshal(req)
	err = conn.WriteMessage(websocket.TextMessage, data)
	assert.NoError(t, err)

	cfg := mock.ConfigMock(t)
	val, err := db.RedisConnection(cfg).Get(ctx, "heartbeat"+"test-device").Result()
	assert.NoError(t, err)
	assert.Equal(t, "alive", val)

	time.Sleep(16 * time.Second)

	val, err = db.RedisConnection(cfg).Get(ctx, "heartbeat"+"test-device").Result()
	assert.Error(t, err, "expected redis key to expire")
	assert.Empty(t, val)
}
