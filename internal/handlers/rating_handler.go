package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RatingHandler struct {
	service services.RatingService
}

func NewRatingHandler(service services.RatingService) *RatingHandler {
	return &RatingHandler{service: service}
}

func (h *RatingHandler) CreateRating(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	productID, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	var req requests.CreateRatingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := h.service.CreateRating(userID, productID, req); err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusCreated, gin.H{"message": "Rating submitted successfully"})

}

func (h *RatingHandler) GetRatingsByID(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)
	productID, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	rating, err := h.service.GetRatingsByID(userID, productID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, rating)
}

func (h *RatingHandler) GetRatingsByUserID(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	ratings, err := h.service.GetRatingsByUserID(userID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, ratings)
}
