package endpoints

import (
	"fmt"
	"os"
	"whisper-api/db"
	"whisper-api/services"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "whisper-api/docs"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	client, err := db.MongoConnection()	
	if err != nil {
		fmt.Println("Cannot connect to mongodb")
		return nil
	}

	collection := client.Database(os.Getenv("WHISPER_DB")).Collection("users")
	userService := services.UserService{collection}

	router.GET("/ping", PingEndpoint{}.Handle)
	router.POST("/register", RegisterEndpoint{&userService}.Handle)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
