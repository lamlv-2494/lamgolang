package handlers

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/utils/constants"
	"strconv"

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

func GetIDParam(ctx *gin.Context) (uint, error) {
	idStr := ctx.Param(constants.IDParam)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func GetPaginationParams(ctx *gin.Context) (page int, limit int) {
	pageStr := ctx.Query(constants.PageParam)
	limitStr := ctx.Query(constants.LimitParam)

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	return page, limit
}
