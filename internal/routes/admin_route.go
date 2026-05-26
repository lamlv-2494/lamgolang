package routes

import (
	"food_delivery/internal/handlers"
	"food_delivery/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(api *gin.RouterGroup, categoryHandler *handlers.CategoryHandler) {

	api.GET("/categories", categoryHandler.GetCategories)

	adminStrict := api.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		adminStrict.POST("/categories", categoryHandler.CreateCategory)
		adminStrict.PUT("/categories/:id", categoryHandler.UpdateCategory)
		adminStrict.DELETE("/categories/:id", categoryHandler.DeleteCategory)
	}
}
