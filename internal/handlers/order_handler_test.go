package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/utils/constants"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─── Checkout ─────────────────────────────────────────────────────────────────

func TestOrderHandler_Checkout_Success(t *testing.T) {
	svc := &mockOrderService{
		checkoutFn: func(userID uint) (*responses.OrderResponse, error) {
			return &responses.OrderResponse{Status: "pending", TotalPrice: 30.0}, nil
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	h.Checkout(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestOrderHandler_Checkout_ServiceError(t *testing.T) {
	svc := &mockOrderService{
		checkoutFn: func(userID uint) (*responses.OrderResponse, error) {
			return nil, configs.CartIsEmpty
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	h.Checkout(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

// ─── GetOrderHistory ──────────────────────────────────────────────────────────

func TestOrderHandler_GetOrderHistory_Success(t *testing.T) {
	svc := &mockOrderService{
		getOrderHistoryFn: func(userID uint, p, l int) (*responses.OrderListResponse, error) {
			return &responses.OrderListResponse{TotalCount: 1}, nil
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=1&limit=10", nil)
	ctx.Set(constants.UserID, uint(1))
	h.GetOrderHistory(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_GetOrderHistory_ServiceError(t *testing.T) {
	svc := &mockOrderService{
		getOrderHistoryFn: func(userID uint, p, l int) (*responses.OrderListResponse, error) {
			return nil, configs.FetchOrdersFailed
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	h.GetOrderHistory(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── AdminGetAllOrders ────────────────────────────────────────────────────────

func TestOrderHandler_AdminGetAllOrders_Success(t *testing.T) {
	svc := &mockOrderService{
		adminGetAllOrdersFn: func(p, l int) (*responses.ListResponse[*responses.OrderResponse], error) {
			return &responses.ListResponse[*responses.OrderResponse]{}, nil
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=1&limit=10", nil)
	h.AdminGetAllOrders(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_AdminGetAllOrders_ServiceError(t *testing.T) {
	svc := &mockOrderService{
		adminGetAllOrdersFn: func(p, l int) (*responses.ListResponse[*responses.OrderResponse], error) {
			return nil, configs.FetchOrdersFailed
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.AdminGetAllOrders(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── AdminUpdateOrderStatus ───────────────────────────────────────────────────

func TestOrderHandler_AdminUpdateOrderStatus_Success(t *testing.T) {
	svc := &mockOrderService{
		adminUpdateOrderStatusFn: func(id uint, req requests.UpdateOrderStatusRequest) (*responses.OrderResponse, error) {
			return &responses.OrderResponse{Status: "completed"}, nil
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(map[string]string{"status": "completed"})
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminUpdateOrderStatus(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrderHandler_AdminUpdateOrderStatus_InvalidID(t *testing.T) {
	svc := &mockOrderService{}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	h.AdminUpdateOrderStatus(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestOrderHandler_AdminUpdateOrderStatus_BindError(t *testing.T) {
	svc := &mockOrderService{}
	h := NewOrderHandler(svc)
	ctx, w := newBadJSONCtx()
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminUpdateOrderStatus(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestOrderHandler_AdminUpdateOrderStatus_ServiceError(t *testing.T) {
	svc := &mockOrderService{
		adminUpdateOrderStatusFn: func(id uint, req requests.UpdateOrderStatusRequest) (*responses.OrderResponse, error) {
			return nil, configs.OrderNotFound
		},
	}
	h := NewOrderHandler(svc)
	ctx, w := newTestCtx(map[string]string{"status": "completed"})
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.AdminUpdateOrderStatus(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
