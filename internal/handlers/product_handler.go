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

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	idStr := ctx.Param(constants.IDParam)
	id, err := strconv.ParseUint(idStr, 10, 64)
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

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	idStr := ctx.Param(constants.IDParam)
	id, err := strconv.ParseUint(idStr, 10, 64)
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

func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	products, err := h.service.GetProducts()
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, products)
}
