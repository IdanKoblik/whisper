package endpoints

import (
	"fmt"
	"whisper-api/config"
	"whisper-api/db"
	"whisper-api/services"

	"github.com/gin-gonic/gin"

	_ "whisper-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	client, err := db.MongoConnection(cfg)	
	if err != nil {
		fmt.Println("Cannot connect to mongodb")
		return nil
	}

	collection := client.Database(cfg.Mongo.Database).Collection("users")
	userService := services.UserService{collection}

	router.GET("/ping", PingEndpoint{}.Handle)
	router.POST("/register", RegisterEndpoint{&userService, cfg}.Handle)
	router.DELETE("/unregister/:token", UnregisterEndpoint{&userService, cfg}.Handle)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
