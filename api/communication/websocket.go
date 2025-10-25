package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"whisper-api/config"
	"whisper-api/db"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ConnectRequest struct {
	DeviceID string `json:"device_id"`
}

type Conn interface {
	WriteJSON(v interface{}) error
}

var (
	Clients      = make(map[string]Conn)
	clientsMutex = &sync.Mutex{}
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

		var request ConnectRequest
		if err := json.Unmarshal(data, &request); err != nil {
			fmt.Printf("Unmarshal error: %v\n", err)
			continue
		}

		found, err := db.DoesExists(cfg, apiToken, request.DeviceID)
		if err != nil || !found {
			fmt.Printf("Device validation failed for %s\n", request.DeviceID)
			continue
		}

		clientsMutex.Lock()
		Clients[request.DeviceID] = conn
		clientsMutex.Unlock()

		fmt.Printf("New device connected: %s\n", request.DeviceID)
	}
}
