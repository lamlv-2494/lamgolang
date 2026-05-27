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
	// Bọc toàn bộ cụm quản lý giỏ hàng qua bộ lọc Token
	cartGroup := api.Group("/cart", middlewares.AuthMiddleware())
	{
		cartGroup.GET("", cartHandler.GetCart)               // Xem giỏ hàng
		cartGroup.POST("", cartHandler.AddToCart)            // Thêm món vào giỏ
		cartGroup.PUT("/:id", cartHandler.UpdateCartItem)    // Tăng giảm số lượng
		cartGroup.DELETE("/:id", cartHandler.RemoveFromCart) // Xoá món khỏi giỏ
	}
}
