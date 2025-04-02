package main

import (
	"fmt"
	"ivanberries-max/internal/cache"
	"ivanberries-max/internal/handler"
	"ivanberries-max/internal/kafka"
	"ivanberries-max/internal/repository"
	"ivanberries-max/internal/service/logic"
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

	redisClient := cache.NewRedisClient(os.Getenv("REDIS_ADDR"))

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	producer := kafka.NewKafkaProducer(kafkaBroker, kafkaTopic)
	defer producer.Close()

	productRepo := repository.NewProductRepository(db)
	productService := logic.NewProductService(productRepo, producer)
	productHandler := handler.NewProductHandler(productService)

	router := gin.Default()
	product := router.Group("/products")
	{
		product.GET("/:id", productHandler.GetByID)
		product.GET("/", productHandler.GetAll)
		product.POST("/", productHandler.Create)
		product.PUT("/:id", productHandler.Update)
		product.DELETE("/:id", productHandler.Delete)
	}

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := logic.NewCategoryService(categoryRepo, redisClient)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	category := router.Group("/categories")
	{
		category.GET("/:id", categoryHandler.GetByID)
		category.GET("/", categoryHandler.GetAll)
		category.POST("/", categoryHandler.Create)
		category.PUT("/:id", categoryHandler.Update)
		category.DELETE("/:id", categoryHandler.Delete)
	}

	go kafka.StartConsumer(kafkaBroker, kafkaTopic, "catalog-service")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server start error: %s", err)
	}
}
