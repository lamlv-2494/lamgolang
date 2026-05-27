package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var req requests.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	idStr := ctx.Param(constants.IDParam)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	var req requests.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	category, err := h.service.UpdateCategory(uint(id), req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param(constants.IDParam)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, categories)
}
