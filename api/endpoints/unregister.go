package endpoints

import (
	"os"
	"whisper-api/services"

	"github.com/gin-gonic/gin"
)

type UnregisterEndpoint struct{
	service *services.UserService	
}

// UnregisterEndpoint godoc
// @Summary Unregister an existing user
// @Description Requires X-Admin-Token header and user token as a URL parameter
// @Tags Users
// @Produce plain
// @Param X-Admin-Token header string true "Admin Token"
// @Param token path string true "User JWT Token"
// @Success 200 {string} string "Successfully removed token"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Router /unregister/{token} [delete]
func (endpoint UnregisterEndpoint) Handle(c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	if adminToken != os.Getenv("WHISPER_ADMIN_TOKEN") {
		c.String(401, "Unauthorized: invalid admin token")
		return
	}

	token := c.Param("token")
	err := endpoint.service.UnregisterUser(token)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(200, "Successfully removed token")
}
