package endpoints

import (
	"fmt"
	"net/http"
	"slices"
	"time"
	"whisper-api/db"

	"github.com/gin-gonic/gin"
)

type MessageRequest struct {
	Message string   `json:"message"`
	Device  string   `json:"device"`
	Targets []string `json:"targets"`
}

type MessageResponse struct {
	Device  string   `json:"device"`
	Targets []string `json:"targets"`
}

const RATE_LIMIT_DURATION = time.Minute

// Send godoc
// @Summary     Send message to device
// @Description Send a message to a specific device via WebSocket. The device must be active and connected. Rate limiting may apply.
// @Tags        messaging
// @Accept      json
// @Produce     plain
// @Security    ApiKeyAuth
// @Param       request body MessageRequest true "Message request with device and targets"
// @Success     200 {string} string "Message sent successfully"
// @Failure     400 {string} string "Invalid request, device not active, or invalid device id"
// @Failure     401 {string} string "Unauthorized - Invalid or missing token"
// @Failure     429 {string} string "Rate limit exceeded"
// @Header      429 {string} Retry-After "Number of seconds to wait before retrying"
// @Router      /api/send [post]
func (h *AuthHandler) Send(c *gin.Context) {
	var request MessageRequest
	err := c.ShouldBindBodyWithJSON(&request)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	devices := c.GetStringSlice("devices")
	contains := slices.Contains(devices, request.Device)
	if !contains {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid device id"})
		return
	}

	rateLimit := c.GetInt("rate-limit")
	if rateLimit > 0 {
		ttl, err := db.RateLimit(h.Repo.Rdb, request.Device, rateLimit, RATE_LIMIT_DURATION)
		if err != nil {
			c.String(400, err.Error())
			return
		}

		if ttl > 0 {
			c.Header("Retry-After", fmt.Sprintf("%d", ttl))
			c.String(429, "Rate limit exceeded. Try again in %d seconds", ttl)
			return
		}
	}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	ws, ok := clients[request.Device]
	if !ok {
		c.String(400, "Device not active")
		return
	}

	payload := MessageRequest{
		Message: request.Message,
		Targets: request.Targets,
	}

	if err := ws.WriteJSON(payload); err != nil {
		fmt.Printf("WS send failed to %s: %v\n", request.Device, err)
		ws.Close()
		delete(clients, request.Device)
	}

	c.String(200, "Message sent")
}
