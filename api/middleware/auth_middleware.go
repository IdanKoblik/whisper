package middleware

import (
	"net/http"
	"whisper-api/config"
	"whisper-api/repository"
	"whisper-api/utils"

	"github.com/gin-gonic/gin"
)

const API_HEADER = "X-API-Token"

func AuthMiddleware(repo *repository.AuthRepository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(API_HEADER)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		authData, err := repo.ValidateToken(c.Request.Context(), utils.HashToken(authHeader))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("devices", authData.Devices)
		c.Set("token", authData.ApiToken)
		c.Set("rate-limit", cfg.RateLimit)
		c.Next()
	}
}
