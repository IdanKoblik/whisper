package main

import (
	"context"
	"log"
	"time"
	"whisper-api/config"
	"whisper-api/endpoints"
	"whisper-api/middleware"
	"whisper-api/repository"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "whisper-api/docs"
)

// @title           Whisper API
// @version         2025-b2
// @description     API for real-time messaging via WebSocket

// @license.name  GNU 3.0
// @license.url   https://github.com/IdanKoblik/whisper/blob/main/LICENSE.md

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Token
// @description API token for authentication

// @securityDefinitions.apikey AdminAuth
// @in header
// @name X-Admin-Token
// @description Admin token for administrative operations
func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})

	mongoOpts := options.Client().ApplyURI(cfg.Mongo.ConnectionString)
	mongoClient, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		log.Fatal(err)
	}

	authRepo := repository.NewAuthRepository(mongoClient, rdb, cfg.Mongo.Database, cfg.Mongo.Collection)
	authHandler := endpoints.NewAuthHandler(authRepo)

	router := gin.Default()
	router.GET("/health", authHandler.Health)

	ws := router.Group("/ws")
	ws.Use(middleware.AuthMiddleware(authRepo, cfg))
	{
		ws.GET("/", authHandler.Websocket)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(authRepo, cfg))
	{
		api.POST("/send", authHandler.Send)
		api.POST("/devices", authHandler.AddDevice)
		api.DELETE("/devices", authHandler.RemoveDevice)
		api.GET("/devices/:id", authHandler.GetDevice)
	}

	admin := router.Group("/admin")
	admin.Use(middleware.AdminMiddleware(cfg))
	{
		admin.POST("/register", authHandler.Register)
		admin.DELETE("/unregister/:token", authHandler.UnRegister)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	go endpoints.HandleHeartbeat(authRepo)

	err = router.Run(cfg.Addr)
	if err != nil {
		panic(err)
	}
}
