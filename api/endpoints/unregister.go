package endpoints

import (
	"os"
	"whisper-api/services"

	"github.com/gin-gonic/gin"
)

type UnregisterEndpoint struct{
	service *services.UserService	
}

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
