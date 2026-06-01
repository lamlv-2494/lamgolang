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

// ─── CreateRating ─────────────────────────────────────────────────────────────

func TestRatingHandler_Create_Success(t *testing.T) {
	svc := &mockRatingService{
		createRatingFn: func(userID, productID uint, req requests.CreateRatingRequest) error { return nil },
	}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"stars": 5, "comment": "Great!"})
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.CreateRating(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRatingHandler_Create_InvalidProductID(t *testing.T) {
	svc := &mockRatingService{}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	h.CreateRating(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestRatingHandler_Create_BindError(t *testing.T) {
	svc := &mockRatingService{}
	h := NewRatingHandler(svc)
	ctx, w := newBadJSONCtx()
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.CreateRating(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestRatingHandler_Create_ServiceError(t *testing.T) {
	svc := &mockRatingService{
		createRatingFn: func(userID, productID uint, req requests.CreateRatingRequest) error {
			return configs.RatingAlreadyExists
		},
	}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"stars": 4})
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.CreateRating(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

// ─── GetRatingByID ────────────────────────────────────────────────────────────

func TestRatingHandler_GetRatingByID_Success(t *testing.T) {
	svc := &mockRatingService{
		getRatingByIDFn: func(userID, productID uint) (*responses.RatingResponse, error) {
			return &responses.RatingResponse{}, nil
		},
	}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.GetRatingByID(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRatingHandler_GetRatingByID_InvalidID(t *testing.T) {
	svc := &mockRatingService{}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
	h.GetRatingByID(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestRatingHandler_GetRatingByID_ServiceError(t *testing.T) {
	svc := &mockRatingService{
		getRatingByIDFn: func(userID, productID uint) (*responses.RatingResponse, error) {
			return nil, configs.RatingNotFound
		},
	}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.GetRatingByID(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetRatingsByUserID ───────────────────────────────────────────────────────

func TestRatingHandler_GetRatingsByUserID_Success(t *testing.T) {
	svc := &mockRatingService{
		getRatingsByUserIDFn: func(userID uint, p, l int) (*responses.ListResponse[*responses.RatingResponseData], error) {
			return &responses.ListResponse[*responses.RatingResponseData]{}, nil
		},
	}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=1&limit=10", nil)
	ctx.Set(constants.UserID, uint(1))
	h.GetRatingsByUserID(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRatingHandler_GetRatingsByUserID_ServiceError(t *testing.T) {
	svc := &mockRatingService{
		getRatingsByUserIDFn: func(userID uint, p, l int) (*responses.ListResponse[*responses.RatingResponseData], error) {
			return nil, configs.RatingNotFound
		},
	}
	h := NewRatingHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Set(constants.UserID, uint(1))
	h.GetRatingsByUserID(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
