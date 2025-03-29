package main

import (
	"fmt"
	"ivanberries-max/internal/handlers"
	"ivanberries-max/internal/repositories"
	"ivanberries-max/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatalf("Ошибка загрузки .env файла: %s", err)
	//}

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

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.Default()
	product := router.Group("/products")
	{
		product.GET("/:id", productHandler.GetProductByID)
		product.GET("/", productHandler.GetAllProducts)
		product.POST("/", productHandler.CreateProduct)
		product.PUT("/:id", productHandler.UpdateProduct)
		product.DELETE("/:id", productHandler.DeleteProduct)
	}

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	category := router.Group("/categories")
	{
		category.GET("/:id", categoryHandler.GetCategoryByID)
		category.GET("/", categoryHandler.GetAllCategories)
		category.POST("/", categoryHandler.CreateCategory)
		category.PUT("/:id", categoryHandler.UpdateCategory)
		category.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// 6️⃣ Запускаем сервер
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %s", err)
	}
}
