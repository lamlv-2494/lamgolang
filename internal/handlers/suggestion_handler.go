package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/services"
	"food_delivery/internal/utils/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuggestionHandler struct {
	service services.SuggestionService
}

func NewSuggestionHandler(service services.SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{service: service}
}

func (h *SuggestionHandler) CreateSuggestion(ctx *gin.Context) {
	userID := ctx.MustGet(constants.UserID).(uint)

	var req requests.CreateSuggestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseError(ctx, err)
		return
	}

	if err := h.service.CreateSuggestion(userID, req); err != nil {
		ResponseError(ctx, err)
		return
	}
	ResponseSuccess(ctx, http.StatusCreated, gin.H{"message": "Suggestion submitted successfully"})
}

func (h *SuggestionHandler) GetAllSuggestions(ctx *gin.Context) {
	suggestions, err := h.service.GetAllSuggestions()
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	ResponseSuccess(ctx, http.StatusOK, suggestions)
}
