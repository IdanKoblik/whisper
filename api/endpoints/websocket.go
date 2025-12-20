package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"whisper-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type WebsocketRequest struct {
	Device string `json:"device"`
}

const (
	DEVICE_CACHE        = 15 * time.Second
	DEVICE_CACHE_PREFIX = "heartbeat:"
	MSG_DEADLINE        = 5 * time.Second
)

var (
	clients      = make(map[string]*websocket.Conn)
	clientsMutex sync.Mutex

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// Websocket godoc
// @Summary     WebSocket connection endpoint
// @Description Establishes a WebSocket connection for real-time messaging. Requires authentication token. Client must send device information in JSON format.
// @Tags        websocket
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       device body WebsocketRequest true "Device information"
// @Success     101 {string} string "WebSocket connection upgraded successfully"
// @Failure     400 {object} map[string]string "Invalid payload or device validation failed"
// @Failure     401 {string} string "Unauthorized - Invalid or missing token"
// @Router      /ws/ [get]
func (h *AuthHandler) Websocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Upgrade error: %v\n", err)
		return
	}
	defer ws.Close()

	ctx := c.Request.Context()
	token := c.GetString("token")

	ws.SetPongHandler(func(string) error {
		return nil
	})

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			h.cleanupConnection(ws)
			return
		}

		var req WebsocketRequest
		if err := json.Unmarshal(data, &req); err != nil {
			_ = ws.WriteJSON(gin.H{"error": "invalid payload"})
			continue
		}

		if req.Device == "" {
			_ = ws.WriteJSON(gin.H{"error": "device required"})
			continue
		}

		found, err := h.Repo.ValidateDeviceID(ctx, token, req.Device)
		if err != nil || !found {
			_ = ws.WriteJSON(gin.H{"error": "device validation failed"})
			continue
		}

		h.registerClient(req.Device, ws)
		h.setHeartbeat(ctx, req.Device)
		_ = ws.WriteJSON(gin.H{"status": "ok"})
	}
}

func (h *AuthHandler) registerClient(device string, ws *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if _, exists := clients[device]; exists {
		return
	}

	clients[device] = ws
	fmt.Printf("WS active: %s\n", device)
}

func (h *AuthHandler) cleanupConnection(ws *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for device, client := range clients {
		if client == ws {
			delete(clients, device)
			_ = h.Repo.Rdb.Del(
				context.Background(),
				DEVICE_CACHE_PREFIX+device,
			).Err()

			fmt.Printf("WS disconnected: %s\n", device)
			break
		}
	}
}

func (h *AuthHandler) setHeartbeat(ctx context.Context, device string) {
	err := h.Repo.Rdb.Set(
		ctx,
		DEVICE_CACHE_PREFIX+device,
		"alive",
		DEVICE_CACHE,
	).Err()

	if err != nil {
		fmt.Printf("Heartbeat error: %v\n", err)
	}
}

func HandleHeartbeat(repo *repository.AuthRepository) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		clientsMutex.Lock()

		for device, ws := range clients {
			_, err := repo.Rdb.Get(
				context.Background(),
				DEVICE_CACHE_PREFIX+device,
			).Result()

			if errors.Is(err, redis.Nil) {
				fmt.Printf("Disconnecting %s (heartbeat timeout)\n", device)
				_ = ws.WriteJSON(gin.H{"error": "heartbeat timeout"})
				ws.Close()
				delete(clients, device)
				continue
			}

			if err != nil {
				fmt.Printf("Redis error: %v\n", err)
				continue
			}

			_ = ws.WriteControl(
				websocket.PingMessage,
				[]byte("ping"),
				time.Now().Add(MSG_DEADLINE),
			)
		}

		clientsMutex.Unlock()
	}
}
