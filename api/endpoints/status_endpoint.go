package endpoints

import (
	"whisper-api/config"
	"whisper-api/db"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
)

// DeviceStatus godoc
// @Summary      Check device status
// @Description  Checks if a given device is registered and active based on its heartbeat record.
// @Tags         Devices
// @Accept       json
// @Produce      plain
// @Param        X-Api-Token  header    string  true  "API token for authentication"
// @Param        DeviceID     path      string  true  "Device ID to check status for"
// @Success      200          {string}  string  "Device is active"
// @Failure      401          {string}  string  "Unauthorized"
// @Failure      404          {string}  string  "Device not found"
// @Failure      500          {string}  string  "Internal server error"
// @Router       /api/status/{DeviceID} [get]
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
