package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service services.CartService
}

func NewCartHandler(service services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// Viết lại hàm khởi tạo chuẩn cho CartHandler
func NewActualCartHandler(service services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) AddToCart(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	var req requests.AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := h.service.AddToCart(userID, &req); err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

func (h *CartHandler) GetCart(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	cart, err := h.service.GetCart(userID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, cart)
}

func (h *CartHandler) UpdateCartItem(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	idStr := ctx.Param(constants.IDParam)
	itemID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	var req requests.UpdateCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := h.service.UpdateCartItem(userID, uint(itemID), &req); err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, gin.H{"message": "Cart item updated successfully"})
}

func (h *CartHandler) RemoveFromCart(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	idStr := ctx.Param(constants.IDParam)
	itemID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := h.service.RemoveFromCart(userID, uint(itemID)); err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}
