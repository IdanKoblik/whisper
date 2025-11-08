package communication

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"whisper-api/config"
	"whisper-api/db"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ConnectRequest struct {
	DeviceID string `json:"device_id"`
}

type HeartbeatRequest struct {
	DeviceID string `json:"device_id"`
	Ping     string `json:"ping"`
}

type Conn interface {
	WriteJSON(v interface{}) error
}

var (
	Clients      = make(map[string]Conn)
	clientsMutex = &sync.Mutex{}
	ctx          = context.Background()
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func HandleWebsocket(cfg *config.Config, c *gin.Context) {
	apiToken := c.GetHeader("X-Api-Token")
	if apiToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	for {
		_, data, err := conn.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			break
		}

		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			break
		}

		var heartbeatRequest HeartbeatRequest
		if err := json.Unmarshal(data, &heartbeatRequest); err == nil {
			db.RedisConnection(cfg).Set(ctx, "heartbeat"+heartbeatRequest.DeviceID, "alive", 15*time.Second)
			continue
		}

		var connRequest ConnectRequest
		if err := json.Unmarshal(data, &connRequest); err != nil {
			fmt.Printf("Unmarshal error: %v\n", err)
			continue
		}

		found, err := db.DoesExists(cfg, utils.HashToken(apiToken), connRequest.DeviceID)
		if err != nil {
			fmt.Printf("Device validation failed for %s\n", err.Error())
			continue
		}

		if !found {
			fmt.Printf("Device validation failed for %s\n", connRequest.DeviceID)
			continue
		}

		clientsMutex.Lock()
		Clients[connRequest.DeviceID] = conn
		clientsMutex.Unlock()

		fmt.Printf("New device connected: %s\n", connRequest.DeviceID)
	}
}

func HandleHeartbeat(cfg *config.Config) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		clientsMutex.Lock()
		for deviceID, client := range Clients {
			exists, err := db.DoesExists(cfg, "heartbeat:"+deviceID, "")
			if err != nil {
				fmt.Printf("Redis error: %v\n", err)
				continue
			}

			if exists {
				fmt.Printf("Disconnecting device %s due to timeout\n", deviceID)
				client.WriteJSON(map[string]string{"error": "timeout"})
				if wsConn, ok := client.(*websocket.Conn); ok {
					wsConn.Close()
				}
				delete(Clients, deviceID)

				err := db.RedisConnection(cfg).Del(ctx, "heartbeat:"+deviceID).Err()
				if err != nil {
					fmt.Printf("Failed to remove heartbeat key for %s: %v\n", deviceID, err)
				}
			}
		}
	}

	clientsMutex.Unlock()
}
