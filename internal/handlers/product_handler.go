package handlers

import (
	"fmt"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct godoc
// @Summary Create product (Admin only)
// @Description Create a new product (Admin only)
// @Tags Admin - Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Product name"
// @Param description formData string false "Product description"
// @Param price formData number true "Product price"
// @Param category_id formData int true "Category ID"
// @Param image formData file true "Product image"
// @Success 201 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/products [post]
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req requests.CreateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		ResponseError(ctx, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	uploadDir := "uploads/products"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		_ = os.MkdirAll(uploadDir, os.ModePerm)
	}

	filename := filepath.Base(file.Filename)
	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)

	dst := filepath.Join(uploadDir, uniqueFilename)

	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ResponseError(ctx, err)
		return
	}

	req.Image = "/" + dst

	product, err := h.service.CreateProduct(req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update product (Admin only)
// @Description Update a product (Admin only)
// @Tags Admin - Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param name formData string false "Product name"
// @Param description formData string false "Product description"
// @Param price formData number false "Product price"
// @Param category_id formData int false "Category ID"
// @Param image formData file false "Product image"
// @Success 200 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/products/{id} [put]
func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	id, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	var req requests.UpdateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		ResponseError(ctx, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err == nil && file != nil {
		uploadDir := "uploads/products"

		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			_ = os.MkdirAll(uploadDir, os.ModePerm)
		}

		filename := filepath.Base(file.Filename)
		uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)

		dst := filepath.Join(uploadDir, uniqueFilename)

		if err := ctx.SaveUploadedFile(file, dst); err != nil {
			ResponseError(ctx, err)
			return
		}

		imagePath := "/" + dst
		req.Image = &imagePath
	}

	product, err := h.service.UpdateProduct(uint(id), req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete product (Admin only)
// @Description Delete a product (Admin only)
// @Tags Admin - Products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 "Product deleted"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	id, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetProducts godoc
// @Summary List all products
// @Description Get products with pagination and filtering
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category_id query int false "Filter by category ID"
// @Param search query string false "Search by name"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param rating query number false "Minimum rating"
// @Param sort query string false "Sort order"
// @Success 200 {array} map[string]interface{}
// @Router /products [get]
func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	classify := ctx.Query(constants.ClassifyParam)
	search := ctx.Query(constants.SearchParam)
	sort := ctx.Query(constants.SortParam)

	var categoryID uint
	if idStr := ctx.Query(constants.CategoryIDParam); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			categoryID = uint(id)
		}
	}

	var minPrice float64
	if minPriceStr := ctx.Query(constants.MinPriceParam); minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = val
		}
	}

	var maxPrice float64
	if maxPriceStr := ctx.Query(constants.MaxPriceParam); maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = val
		}
	}

	var minRating float64
	if ratingStr := ctx.Query(constants.RatingParam); ratingStr != "" {
		if val, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			minRating = val
		}
	}

	page, limit := GetPaginationParams(ctx)

	products, err := h.service.GetProducts(search, classify, categoryID, minPrice, maxPrice, minRating, sort, page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 "Product not found"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(ctx *gin.Context) {
	id, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, product)
}
