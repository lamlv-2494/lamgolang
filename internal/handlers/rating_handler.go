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

// CreateRating godoc
// @Summary Create rating for product
// @Description Add a rating/review for a product
// @Tags Ratings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param body body requests.CreateRatingRequest true "Rating data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Router /products/{id}/rating [post]
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

// GetRatingByID godoc
// @Summary Get user's rating for product
// @Description Get the authenticated user's rating for a specific product
// @Tags Ratings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Failure 404 "Rating not found"
// @Router /products/{id}/rating [get]
func (h *RatingHandler) GetRatingByID(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)
	productID, err := GetIDParam(ctx)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	rating, err := h.service.GetRatingByID(userID, productID)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, rating)
}

// GetRatingsByUserID godoc
// @Summary Get user's ratings
// @Description Get all ratings/reviews submitted by the authenticated user
// @Tags Ratings
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Router /products/ratings [get]
func (h *RatingHandler) GetRatingsByUserID(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	page, limit := GetPaginationParams(ctx)

	ratings, err := h.service.GetRatingsByUserID(userID, page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}

	ResponseSuccess(ctx, http.StatusOK, ratings)
}
