package communication

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Communication struct {
 	Redis *redis.Client
} 

type HeartbeatRequest struct {
	DeviceID string `json:"device_id"`
}

var Clients = make(map[string]*websocket.Conn)

// TODO duration
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var ctx = context.Background()

func (communication Communication) HandleWebsocket(c *gin.Context) {
	apiToken := c.GetHeader("X-Api-Token")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer conn.Close()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			continue // TODO handle
		}

		var request HeartbeatRequest
		err = json.Unmarshal(data, &request)
		if err != nil {
			continue // TODO handle
		}
		_, err = communication.Redis.Get(ctx, apiToken).Result()
		if err == redis.Nil {
			// Token does not exist → register it
			if err := communication.Redis.Set(ctx, apiToken, request.DeviceID, 0).Err(); err != nil {
				conn.WriteJSON(gin.H{"error": "failed to register device"})
				continue
			}
		} else if err != nil {
			// Redis error
			conn.WriteJSON(gin.H{"error": "redis error"})
			continue
		} else {
			// Token exists → optionally update the device_id
			if err := communication.Redis.Set(ctx, apiToken, request.DeviceID, 0).Err(); err != nil {
				conn.WriteJSON(gin.H{"error": "failed to update device"})
				continue
			}
		}

		Clients[request.DeviceID] = conn
		fmt.Printf("\nReg: %s", request.DeviceID)
		conn.WriteJSON(gin.H{
			"status":  "ok",
			"message": "device registered successfully",
		})
	}
}
