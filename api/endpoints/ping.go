package endpoints

import "github.com/gin-gonic/gin"

type PingEndpoint struct{}

// PingEndpoint godoc
// @Summary Ping the server
// @Description Returns "pong"
// @Tags Health
// @Produce plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (PingEndpoint) Handle(c *gin.Context) {
	c.String(200, "pong")
}
