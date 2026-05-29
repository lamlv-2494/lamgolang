package main

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/handlers"
	"food_delivery/internal/repositories"
	"food_delivery/internal/routes"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"log"
	"net/http"
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
	cartRepo := repositories.NewCartRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	ratingRepo := repositories.NewRatingRepository(db)
	suggestionRepo := repositories.NewSuggestionRepository(db)

	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo)
	ratingService := services.NewRatingService(ratingRepo, productRepo)
	suggestionService := services.NewSuggestionService(suggestionRepo)

	userHandlers := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)
	ratingHandler := handlers.NewRatingHandler(ratingService)
	suggestionHandler := handlers.NewSuggestionHandler(suggestionService)

	adminHandlers := &routes.AdminHandlers{
		OrderHandler:      orderHandler,
		CategoryHandler:   categoryHandler,
		ProductHandler:    productHandler,
		UserHandler:       userHandlers,
		SuggestionHandler: suggestionHandler,
	}

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		// Common
		api.GET("/categories", adminHandlers.CategoryHandler.GetCategories)
		api.GET("/products", adminHandlers.ProductHandler.GetProducts)
		api.GET("/products/:id", adminHandlers.ProductHandler.GetProductByID)

		// User
		routes.RegisterUserGroup(api, userHandlers)
		routes.RegisterCartRoutes(api, cartHandler)
		routes.RegisterOrderRoutes(api, orderHandler)
		routes.RegisterRatingRoutes(api, ratingHandler)
		routes.RegisterSuggestionRoutes(api, suggestionHandler)

		// Admin
		routes.RegisterAdminRoutes(api, adminHandlers)
	}

	staticDir := "./static"
	// Fallback routing to serve static frontend files when no API route matches
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir(staticDir))))

	port := os.Getenv(constants.Port)
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
