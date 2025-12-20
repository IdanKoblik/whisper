package middleware

import (
	"net/http"
	"whisper-api/config"

	"github.com/gin-gonic/gin"
)

const ADMIN_HEADER = "X-Admin-Token"

func AdminMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(ADMIN_HEADER)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		if authHeader != cfg.AdminToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Next()
	}
}
