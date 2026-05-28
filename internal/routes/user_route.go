package routes

import (
	"food_delivery/internal/handlers"
	"food_delivery/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserGroup(api *gin.RouterGroup, userController *handlers.UserHandler) {
	userOptional := api.Group("/user")
	{
		userOptional.POST("/register", userController.Register)
		userOptional.POST("/login", userController.Login)
	}

	userStrict := api.Group("/user", middlewares.AuthMiddleware())
	{
		userStrict.GET("", userController.GetCurrentUser)
		userStrict.PUT("", userController.UpdateUser)
	}
}

func RegisterCartRoutes(api *gin.RouterGroup, cartHandler *handlers.CartHandler) {
	cartGroup := api.Group("/cart", middlewares.AuthMiddleware())
	{
		cartGroup.GET("", cartHandler.GetCart)
		cartGroup.POST("", cartHandler.AddToCart)
		cartGroup.PUT("/:id", cartHandler.UpdateCartItem)
		cartGroup.DELETE("/:id", cartHandler.RemoveFromCart)
	}
}

func RegisterOrderRoutes(api *gin.RouterGroup, orderHandler *handlers.OrderHandler) {
	orderGroup := api.Group("/orders", middlewares.AuthMiddleware())
	{
		orderGroup.POST("", orderHandler.Checkout)
		orderGroup.GET("", orderHandler.GetOrderHistory)
	}
}
