package endpoints

import (
	"fmt"
	"whisper-api/config"
	"whisper-api/services"

	"github.com/gin-gonic/gin"
)

// UnRegisterUser godoc
// @Summary      Unregister a user
// @Description  Allows an admin to delete a user by API token
// @Tags         Admin
// @Produce      plain
// @Param        X-Admin-Token header string true "Admin token"
// @Param        ApiToken path string true "User API token to remove"
// @Success      200 {string} string "Deleted {ApiToken}"
// @Failure      400 {string} string "Bad Request"
// @Failure      401 {string} string "Unauthorized: Invalid admin token"
// @Router       /api/admin/unregister/{ApiToken} [delete]
func UnRegisterUser(cfg *config.Config, c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	cfg, err := config.GetConfig()
	if err != nil {
		c.String(400, err.Error())
		return
	}

	if adminToken != cfg.AdminToken {
		c.String(401, "Unauthorized: Invalid admin token")
		return
	}

	token := c.Param("ApiToken")
	err = services.RemoveUser(cfg, token)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(200, fmt.Sprintf("Deleted %s", token))
}
