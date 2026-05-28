package handlers

import (
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
	// Lấy ID user từ middleware xác thực token
	userID := ctx.MustGet(constants.UserID).(uint)

	orderResponse, err := h.service.Checkout(userID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusCreated, orderResponse)
}
