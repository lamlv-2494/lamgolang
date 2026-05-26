package handlers

import (
	"food_delivery/internal/configs"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, data)
}

// ResponseError responds with an error using the error presenter
func ResponseError(ctx *gin.Context, err error) {
	statusCode, message := configs.GetStatusCodeAndMessage(err)

	ctx.JSON(statusCode, gin.H{"code": statusCode, "message": message})
}
