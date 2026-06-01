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

// CreateSuggestion godoc
// @Summary Create a suggestion
// @Description Submit a suggestion or feedback
// @Tags Suggestions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body requests.CreateSuggestionRequest true "Suggestion data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 "Invalid request"
// @Failure 401 "Unauthorized"
// @Router /suggestions [post]
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

// GetAllSuggestions godoc
// @Summary Get all suggestions (Admin only)
// @Description Get all submitted suggestions (Admin only)
// @Tags Admin - Suggestions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} map[string]interface{}
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Router /admin/suggestions [get]
func (h *SuggestionHandler) GetAllSuggestions(ctx *gin.Context) {

	page, limit := GetPaginationParams(ctx)

	suggestions, err := h.service.GetAllSuggestions(page, limit)
	if err != nil {
		ResponseError(ctx, err)
		return
	}
	ResponseSuccess(ctx, http.StatusOK, suggestions)
}
