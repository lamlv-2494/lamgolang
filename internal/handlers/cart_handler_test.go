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

// ─── AddToCart ────────────────────────────────────────────────────────────────

func TestCartHandler_AddToCart_Success(t *testing.T) {
	svc := &mockCartService{
		addToCartFn: func(userID uint, req *requests.AddToCartRequest) error { return nil },
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"product_id": 1, "quantity": 2})
	ctx.Set(constants.UserID, uint(1))
	h.AddToCart(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCartHandler_AddToCart_BindError(t *testing.T) {
	svc := &mockCartService{}
	h := NewCartHandler(svc)
	ctx, w := newBadJSONCtx()
	ctx.Set(constants.UserID, uint(1))
	h.AddToCart(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCartHandler_AddToCart_ServiceError(t *testing.T) {
	svc := &mockCartService{
		addToCartFn: func(userID uint, req *requests.AddToCartRequest) error {
			return configs.ProductNotFound
		},
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"product_id": 99, "quantity": 1})
	ctx.Set(constants.UserID, uint(1))
	h.AddToCart(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetCart ──────────────────────────────────────────────────────────────────

func TestCartHandler_GetCart_Success(t *testing.T) {
	svc := &mockCartService{
		getCartFn: func(userID uint) (*responses.CartResponse, error) {
			return &responses.CartResponse{}, nil
		},
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	h.GetCart(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCartHandler_GetCart_ServiceError(t *testing.T) {
	svc := &mockCartService{
		getCartFn: func(userID uint) (*responses.CartResponse, error) {
			return nil, configs.CartItemNotFound
		},
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	h.GetCart(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── UpdateCartItem ───────────────────────────────────────────────────────────

func TestCartHandler_UpdateCartItem_Success(t *testing.T) {
	svc := &mockCartService{
		updateCartItemFn: func(userID, itemID uint, req *requests.UpdateCartItemRequest) error { return nil },
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"quantity": 5})
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateCartItem(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCartHandler_UpdateCartItem_InvalidID(t *testing.T) {
	svc := &mockCartService{}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
	h.UpdateCartItem(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCartHandler_UpdateCartItem_BindError(t *testing.T) {
	svc := &mockCartService{}
	h := NewCartHandler(svc)
	ctx, w := newBadJSONCtx()
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateCartItem(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCartHandler_UpdateCartItem_ServiceError(t *testing.T) {
	svc := &mockCartService{
		updateCartItemFn: func(userID, itemID uint, req *requests.UpdateCartItemRequest) error {
			return configs.CartItemNotFound
		},
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"quantity": 3})
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateCartItem(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── RemoveFromCart ───────────────────────────────────────────────────────────

func TestCartHandler_RemoveFromCart_Success(t *testing.T) {
	svc := &mockCartService{
		removeFromCartFn: func(userID, itemID uint) error { return nil },
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.RemoveFromCart(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCartHandler_RemoveFromCart_InvalidID(t *testing.T) {
	svc := &mockCartService{}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
	h.RemoveFromCart(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCartHandler_RemoveFromCart_ServiceError(t *testing.T) {
	svc := &mockCartService{
		removeFromCartFn: func(userID, itemID uint) error { return configs.CartItemNotFound },
	}
	h := NewCartHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.RemoveFromCart(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
