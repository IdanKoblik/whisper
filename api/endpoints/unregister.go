package endpoints

import (
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
)

// UnRegister godoc
// @Summary     Unregister API token
// @Description Remove an API token from the database. Requires admin privileges.
// @Tags        admin
// @Accept      json
// @Produce     plain
// @Security    AdminAuth
// @Param       token path string true "API token to remove"
// @Success     200 {string} string "API token removed successfully"
// @Failure     400 {string} string "Failed to remove token"
// @Failure     401 {string} string "Unauthorized - Admin access required"
// @Router      /admin/unregister/{token} [delete]
func (h *AuthHandler) UnRegister(c *gin.Context) {
	apiToken := c.Param("token")
	err := h.Repo.RemoveToken(c.Request.Context(), utils.HashToken(apiToken))
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(200, "Removed this api token from databse")
}
