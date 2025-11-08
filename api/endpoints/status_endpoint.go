package endpoints

import (
	"whisper-api/config"
	"whisper-api/db"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
)

func DeviceStatus(cfg *config.Config, c *gin.Context) {
	apiToken := c.GetHeader("X-Api-Token")
	token := utils.HashToken(apiToken)
	found, err := db.DoesExists(cfg, token, "")
	if err != nil || !found {
		c.String(401, "Unauthorized")
		return
	}

	deviceID := c.Param("DeviceID")
	exists, err := db.DoesExists(cfg, "heartbeat:"+deviceID, "")
	if err != nil || !exists {
		c.String(404, "Device not found")
		return
	}

	c.Status(200)
}
