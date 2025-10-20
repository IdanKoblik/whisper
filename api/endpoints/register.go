package endpoints

import (
	"fmt"
	"whisper-api/config"
	"whisper-api/services"

	"github.com/gin-gonic/gin"
)

type RegisterEndpoint struct {
	service *services.UserService
	cfg *config.Config
}

// RegisterEndpoint godoc
// @Summary Register a new user
// @Description Requires X-Admin-Token header and RawUser JSON body
// @Tags Users
// @Accept json
// @Produce json
// @Param X-Admin-Token header string true "Admin Token"
// @Param rawUser body services.RawUser true "User data"
// @Success 200 {string} string "JWT Token and Signature key"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Router /register [post]
func (endpoint RegisterEndpoint) Handle(c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	if adminToken != endpoint.cfg.AdminToken {
		c.String(401, "Unauthorized: invalid admin token")
		return
	}

	data, err := endpoint.service.RegisterUser()	
	if err != nil {
		fmt.Println(err.Error())
		c.String(400, err.Error())
		return
	}

	c.String(200, data.ApiToken)
}
