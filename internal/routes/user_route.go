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
	}
}
