package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"net/http"

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

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a product to user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body requests.AddToCartRequest true "Cart item data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Router /cart [post]
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

// GetCart godoc
// @Summary Get user cart
// @Description Get the current user's shopping cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Router /cart [get]
func (h *CartHandler) GetCart(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	cart, err := h.service.GetCart(userID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, cart)
}

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Description Update the quantity of an item in user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart item ID"
// @Param body body requests.UpdateCartItemRequest true "Updated cart item data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Router /cart/{id} [put]
func (h *CartHandler) UpdateCartItem(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	itemID, err := GetIDParam(ctx)
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

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Description Remove a product from user's cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart item ID"
// @Success 200 "Item removed from cart"
// @Failure 401 "Unauthorized"
// @Router /cart/{id} [delete]
func (h *CartHandler) RemoveFromCart(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	itemID, err := GetIDParam(ctx)
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
