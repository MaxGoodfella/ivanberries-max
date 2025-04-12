package main

import (
	cache1 "catalog-service/internal/cache"
	"catalog-service/internal/handler"
	"catalog-service/internal/kafka"
	repo1 "catalog-service/internal/repository"
	logic1 "catalog-service/internal/service/logic"

	cache2 "github.com/MaxGoodfella/ivanberries-max/users-service/pkg/cache"
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/middleware"
	repo2 "github.com/MaxGoodfella/ivanberries-max/users-service/pkg/repository"
	logic2 "github.com/MaxGoodfella/ivanberries-max/users-service/pkg/service/logic"

	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

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

	redisClient := cache1.NewRedisClient(os.Getenv("REDIS_ADDR"))

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	producer := kafka.NewKafkaProducer(kafkaBroker, kafkaTopic)
	defer producer.Close()

	productRepo := repo1.NewProductRepository(db)
	productService := logic1.NewProductService(productRepo, producer)
	productHandler := handler.NewProductHandler(productService)

	//
	//router := gin.Default()
	//product := router.Group("/products")
	//{
	//	product.GET("/:id", productHandler.GetByID)
	//	product.GET("/", productHandler.GetAll)
	//	product.POST("/", productHandler.Create)
	//	product.PUT("/:id", productHandler.Update)
	//	product.DELETE("/:id", productHandler.Delete)
	//}
	//
	categoryRepo := repo1.NewCategoryRepository(db)
	categoryService := logic1.NewCategoryService(categoryRepo, redisClient)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	//category := router.Group("/categories")
	//{
	//	category.GET("/:id", categoryHandler.GetByID)
	//	category.GET("/", categoryHandler.GetAll)
	//	category.POST("/", categoryHandler.Create)
	//	category.PUT("/:id", categoryHandler.Update)
	//	category.DELETE("/:id", categoryHandler.Delete)
	//}

	authRepo := repo2.NewAuthRepository(db)
	usersRedisClient := cache2.NewRedisClient(os.Getenv("REDIS_ADDR"), 1)
	authService := logic2.NewAuthService(authRepo, os.Getenv("JWT_SECRET"), usersRedisClient)

	router := gin.Default()

	product := router.Group("/products")
	{
		product.GET("/:id", middleware.PermissionMiddleware(authService, "product.getbyid"), productHandler.GetByID)
		product.GET("/", middleware.PermissionMiddleware(authService, "product.getall"), productHandler.GetAll)
		product.POST("/", middleware.PermissionMiddleware(authService, "product.create"), productHandler.Create)
		product.PUT("/:id", middleware.PermissionMiddleware(authService, "product.update"), productHandler.Update)
		product.DELETE("/:id", middleware.PermissionMiddleware(authService, "product.delete"), productHandler.Delete)
	}

	category := router.Group("/categories")
	{
		category.GET("/:id", middleware.PermissionMiddleware(authService, "category.getbyid"), categoryHandler.GetByID)
		category.GET("/", middleware.PermissionMiddleware(authService, "category.getall"), categoryHandler.GetAll)
		category.POST("/", middleware.PermissionMiddleware(authService, "category.create"), categoryHandler.Create)
		category.PUT("/:id", middleware.PermissionMiddleware(authService, "category.update"), categoryHandler.Update)
		category.DELETE("/:id", middleware.PermissionMiddleware(authService, "category.delete"), categoryHandler.Delete)
	}

	go kafka.StartConsumer(kafkaBroker, kafkaTopic, "catalog-service")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server start error: %s", err)
	}
}
