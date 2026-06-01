package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategory godoc
// @Summary Create category (Admin only)
// @Description Create a new category (Admin only)
// @Tags Admin - Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body requests.CreateCategoryRequest true "Category data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/categories [post]
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

// UpdateCategory godoc
// @Summary Update category (Admin only)
// @Description Update a category (Admin only)
// @Tags Admin - Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param body body requests.UpdateCategoryRequest true "Updated category data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	id, err := GetIDParam(ctx)
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

// DeleteCategory godoc
// @Summary Delete category (Admin only)
// @Description Delete a category (Admin only)
// @Tags Admin - Categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 "Category deleted"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	id, err := GetIDParam(ctx)
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

// GetCategories godoc
// @Summary List all categories
// @Description Get all categories with pagination
// @Tags Categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} map[string]interface{}
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
	page, limit := GetPaginationParams(ctx)
	categories, err := h.service.GetCategories(page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	// TODO: create response model later
	ResponseSuccess(ctx, http.StatusOK, categories)
}
