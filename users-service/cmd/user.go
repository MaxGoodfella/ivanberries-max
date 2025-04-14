package main

import (
	"fmt"
	"github.com/MaxGoodfella/ivanberries-max/users-service/config"
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/cache"
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/handler"
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/middleware"
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/repository"
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/service/logic"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	log.Println("DSN:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection error: %s", err)
	}

	redisClient := cache.NewRedisClient(os.Getenv("REDIS_ADDR"), 1)

	authRepo := repository.NewAuthRepository(db)
	authService := logic.NewAuthService(authRepo, cfg.JWTSecret, redisClient)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", middleware.RoleMiddleware(authService, "admin"), authHandler.RefreshToken)
		auth.GET("/me", middleware.JWTMiddleware(authService), authHandler.Me)
		auth.POST("/logout", authHandler.Logout)
	}

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("server start error: %s", err)
	}
}
