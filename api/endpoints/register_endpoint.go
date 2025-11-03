package endpoints

import (
	"whisper-api/config"
	"whisper-api/services"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary      Register a new API user
// @Description  Allows an admin to create a new API user and receive an API token
// @Tags         Admin
// @Accept       json
// @Produce      plain
// @Param        X-Admin-Token header string true "Admin token"
// @Success      200 {string} string "API token for the new user"
// @Failure      400 {string} string "Bad Request"
// @Failure      401 {string} string "Unauthorized: Invalid admin token"
// @Router       /api/admin/register [post]
func RegisterUser(cfg *config.Config, c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	if adminToken != cfg.AdminToken {
		c.String(401, "Unauthorized: Invalid admin token")
		return
	}

	response, err := services.RegisterUser(cfg)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(201, response.ApiToken)
}
