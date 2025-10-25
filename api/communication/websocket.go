package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"whisper-api/config"
	"whisper-api/db"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ConnectRequest struct {
	DeviceID string `json:"device_id"`
}

// TODO duration
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[string]*websocket.Conn)

func HandleWebsocket(cfg *config.Config, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer conn.Close()
	for {
		apiToken := c.GetHeader("X-Api-Token")
		_, data, err := conn.ReadMessage()
		if err != nil {
			continue // TODO handle
		}

		var request ConnectRequest
		err = json.Unmarshal(data, &request)
		if err != nil {
			continue // TODO handle
		}

		found, err := db.DoesExists(cfg, apiToken, request.DeviceID)
		if err != nil {
			continue // TODO handle
		}

		if !found {
			continue // TODO handle
		}

		Clients[request.DeviceID] = conn
		fmt.Printf("\nNew device: %s\n", request.DeviceID)
		conn.WriteJSON(gin.H{
			"status":  "ok",
			"message": "device registered successfully",
		})
	}
}
