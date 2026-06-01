package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// Checkout godoc
// @Summary Create order (checkout)
// @Description Create an order from user's cart
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Router /orders [post]
func (h *OrderHandler) Checkout(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	orderResponse, err := h.service.Checkout(userID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusCreated, orderResponse)
}

// GetOrderHistory godoc
// @Summary Get user order history
// @Description Get orders of the authenticated user
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Router /orders [get]
func (h *OrderHandler) GetOrderHistory(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	page, limit := GetPaginationParams(ctx)

	orderHistory, err := h.service.GetOrderHistory(userID, page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, orderHistory)
}

// AdminGetAllOrders godoc
// @Summary Get all orders (Admin only)
// @Description Get all orders with pagination (Admin only)
// @Tags Admin - Orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/orders [get]
func (h *OrderHandler) AdminGetAllOrders(ctx *gin.Context) {
	page, limit := GetPaginationParams(ctx)

	orders, err := h.service.AdminGetAllOrders(page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, orders)
}

// AdminUpdateOrderStatus godoc
// @Summary Update order status (Admin only)
// @Description Update the status of an order (Admin only)
// @Tags Admin - Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param body body requests.UpdateOrderStatusRequest true "Order status data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/orders/{id}/status [put]
func (h *OrderHandler) AdminUpdateOrderStatus(ctx *gin.Context) {
	orderID, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	var req requests.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	orderResponse, err := h.service.AdminUpdateOrderStatus(orderID, req)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, orderResponse)
}
