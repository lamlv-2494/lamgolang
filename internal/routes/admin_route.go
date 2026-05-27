package routes

import (
	"food_delivery/internal/handlers"
	"food_delivery/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	CategoryHandler *handlers.CategoryHandler
	ProductHandler  *handlers.ProductHandler
}

func RegisterAdminRoutes(api *gin.RouterGroup, adminHandlers *AdminHandlers) {

	api.GET("/categories", adminHandlers.CategoryHandler.GetCategories)
	api.GET("/products", adminHandlers.ProductHandler.GetProducts)

	adminStrict := api.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		// CRUD Categories
		categorygroup := adminStrict.Group("/categories")
		{
			categorygroup.POST("", adminHandlers.CategoryHandler.CreateCategory)
			categorygroup.PUT("/:id", adminHandlers.CategoryHandler.UpdateCategory)
			categorygroup.DELETE("/:id", adminHandlers.CategoryHandler.DeleteCategory)
		}

		// CRUD Products
		productgroup := adminStrict.Group("/products")
		{
			productgroup.POST("", adminHandlers.ProductHandler.CreateProduct)
			productgroup.PUT("/:id", adminHandlers.ProductHandler.UpdateProduct)
			productgroup.DELETE("/:id", adminHandlers.ProductHandler.DeleteProduct)
		}
	}
}
