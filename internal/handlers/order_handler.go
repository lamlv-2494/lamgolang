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

func (h *OrderHandler) Checkout(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	orderResponse, err := h.service.Checkout(userID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusCreated, orderResponse)
}

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

func (h *OrderHandler) AdminGetAllOrders(ctx *gin.Context) {
	page, limit := GetPaginationParams(ctx)

	orders, err := h.service.AdminGetAllOrders(page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, orders)
}

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
