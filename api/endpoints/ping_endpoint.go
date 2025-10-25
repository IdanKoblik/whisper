package endpoints

import "github.com/gin-gonic/gin"

// PingEndpoint godoc
// @Summary      Ping the server
// @Description  Simple health check endpoint
// @Tags         Health
// @Produce      plain
// @Success      200 {string} string "pong"
// @Router       /api/ping [get]
func PingEndpoint(c *gin.Context) {
	c.String(200, "pong")
}
