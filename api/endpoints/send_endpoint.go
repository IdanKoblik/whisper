package endpoints

import (
	"context"
	"fmt"
	"whisper-api/config"
	"whisper-api/db"
	"whisper-api/services"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// SendMessage godoc
// @Summary      Send a message
// @Description  Sends a message through the API using the user's token
// @Tags         Messages
// @Accept       json
// @Produce      plain
// @Param        X-Api-Token header string true "User API token"
// @Param        message body services.MessageRequest true "Message request payload"
// @Success      200 {string} string "Message sent"
// @Failure      400 {string} string "Bad Request"
// @Failure      401 {string} string "Unauthorized"
// @Failure      429 {string} string "Rate limit exceeded"
// @Router       /api/send [post]
func SendMessage(cfg *config.Config, c *gin.Context) {
	apiToken := c.GetHeader("X-Api-Token")
	token := utils.HashToken(apiToken)
	found, err := db.DoesExists(cfg, token, "")
	if err != nil {
		c.String(400, err.Error())
		return
	}

	if !found {
		c.String(401, "Unauthorized")
		return
	}

	if cfg.RateLimit != nil {
		ttl, err := db.RateLimit(cfg, fmt.Sprintf("rate-%s", apiToken), *cfg.RateLimit)
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

	var request services.MessageRequest
	err = c.ShouldBindBodyWithJSON(&request)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	err = services.SendMessage(request)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(200, "Message sent")
}
