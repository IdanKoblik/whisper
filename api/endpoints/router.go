package endpoints

import (
	"whisper-api/communication"
	"whisper-api/config"
	_ "whisper-api/docs"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	if cfg == nil {
		return nil
	}

	api := router.Group("/api")
	{
		api.GET("/ping", PingEndpoint)

		api.POST("/send", func(c *gin.Context) {
			SendMessage(cfg, c)
		})
	}

	adminAPI := api.Group("/admin")
	{
		adminAPI.POST("/register", func(c *gin.Context) {
			RegisterUser(cfg, c)
		})

		adminAPI.DELETE("/unregister/:ApiToken", func(c *gin.Context) {
			UnRegisterUser(cfg, c)
		})
	}

	router.GET("/ws", func(c *gin.Context) {
		communication.HandleWebsocket(cfg, c)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
