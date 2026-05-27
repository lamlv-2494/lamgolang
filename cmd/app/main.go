package main

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/handlers"
	"food_delivery/internal/repositories"
	"food_delivery/internal/routes"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Can't found env file")
	}

	db := configs.ConnectDB()

	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)

	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)

	userHandlers := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)

	adminHandlers := &routes.AdminHandlers{
		CategoryHandler: categoryHandler,
		ProductHandler:  productHandler,
		UserHandler:     userHandlers,
	}

	r := gin.Default()

	api := r.Group("/api")
	{
		routes.RegisterUserGroup(api, userHandlers)
		routes.RegisterAdminRoutes(api, adminHandlers)
	}

	port := os.Getenv(constants.Port)
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
