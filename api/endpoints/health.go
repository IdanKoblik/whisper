package endpoints

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ServerState int

const (
	OK ServerState = iota
	DOWN
)

var stateName = map[ServerState]string{
	OK:   "ok",
	DOWN: "down",
}

// Health godoc
// @Summary     Health check endpoint
// @Description Check the health status of the server and its dependencies (MongoDB and Redis)
// @Tags        health
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]interface{} "Server and dependencies are healthy"
// @Success     503 {object} map[string]interface{} "One or more dependencies are down"
// @Router      /health [get]
func (h *AuthHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	mongoStatus := stateName[OK]
	redisStatus := stateName[OK]

	if err := h.Repo.Col.Database().Client().Ping(ctx, readpref.Primary()); err != nil {
		mongoStatus = stateName[DOWN]
	}

	if err := h.Repo.Rdb.Ping(ctx).Err(); err != nil {
		redisStatus = stateName[DOWN]
	}

	statusCode := http.StatusOK
	if mongoStatus != stateName[OK] || redisStatus != stateName[OK] {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"dependencies": gin.H{
			"mongodb": mongoStatus,
			"redis":   redisStatus,
		},
	})
}
