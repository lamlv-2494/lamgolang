package routes

import (
	"food_delivery/internal/handlers"
	"food_delivery/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	OrderHandler    *handlers.OrderHandler
	CategoryHandler *handlers.CategoryHandler
	ProductHandler  *handlers.ProductHandler
	UserHandler     *handlers.UserHandler
}

func RegisterAdminRoutes(api *gin.RouterGroup, adminHandlers *AdminHandlers) {
	adminStrict := api.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		// Manage Users
		userGroup := adminStrict.Group("/users")
		{
			userGroup.GET("", adminHandlers.UserHandler.GetAllUsers)
			userGroup.PUT("/:id", adminHandlers.UserHandler.AdminUpdateUser)
			userGroup.DELETE("/:id", adminHandlers.UserHandler.AdminDeleteUser)
		}

		// Manage Categories
		categoryGroup := adminStrict.Group("/categories")
		{
			categoryGroup.POST("", adminHandlers.CategoryHandler.CreateCategory)
			categoryGroup.PUT("/:id", adminHandlers.CategoryHandler.UpdateCategory)
			categoryGroup.DELETE("/:id", adminHandlers.CategoryHandler.DeleteCategory)
		}

		// Manage Products
		productGroup := adminStrict.Group("/products")
		{
			productGroup.POST("", adminHandlers.ProductHandler.CreateProduct)
			productGroup.PUT("/:id", adminHandlers.ProductHandler.UpdateProduct)
			productGroup.DELETE("/:id", adminHandlers.ProductHandler.DeleteProduct)
		}

		orderGroup := adminStrict.Group("/orders")
		{
			orderGroup.GET("", adminHandlers.OrderHandler.AdminGetAllOrders)
			orderGroup.PUT("/:id/status", adminHandlers.OrderHandler.AdminUpdateOrderStatus)
		}
	}
}
