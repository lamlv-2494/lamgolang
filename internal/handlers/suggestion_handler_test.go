package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
	"food_delivery/internal/utils/constants"

	"github.com/stretchr/testify/assert"
)

// ─── CreateSuggestion ─────────────────────────────────────────────────────────

func TestSuggestionHandler_Create_Success(t *testing.T) {
	svc := &mockSuggestionService{
		createSuggestionFn: func(userID uint, req requests.CreateSuggestionRequest) error { return nil },
	}
	h := NewSuggestionHandler(svc)
	ctx, w := newTestCtx(map[string]string{"title": "More Choices", "description": "Add more items"})
	ctx.Set(constants.UserID, uint(1))
	h.CreateSuggestion(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSuggestionHandler_Create_BindError(t *testing.T) {
	svc := &mockSuggestionService{}
	h := NewSuggestionHandler(svc)
	ctx, w := newBadJSONCtx()
	ctx.Set(constants.UserID, uint(1))
	h.CreateSuggestion(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestSuggestionHandler_Create_ServiceError(t *testing.T) {
	svc := &mockSuggestionService{
		createSuggestionFn: func(userID uint, req requests.CreateSuggestionRequest) error {
			return configs.CreateSuggestionFailed
		},
	}
	h := NewSuggestionHandler(svc)
	ctx, w := newTestCtx(map[string]string{"title": "Better Menu", "description": "Add more options"})
	ctx.Set(constants.UserID, uint(1))
	h.CreateSuggestion(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

// ─── GetAllSuggestions ────────────────────────────────────────────────────────

func TestSuggestionHandler_GetAll_Success(t *testing.T) {
	svc := &mockSuggestionService{
		getAllSuggestionsFn: func(p, l int) (*responses.ListResponse[*responses.SuggestionData], error) {
			return &responses.ListResponse[*responses.SuggestionData]{}, nil
		},
	}
	h := NewSuggestionHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=1&limit=10", nil)
	h.GetAllSuggestions(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuggestionHandler_GetAll_ServiceError(t *testing.T) {
	svc := &mockSuggestionService{
		getAllSuggestionsFn: func(p, l int) (*responses.ListResponse[*responses.SuggestionData], error) {
			return nil, configs.FetchSuggestionsFailed
		},
	}
	h := NewSuggestionHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetAllSuggestions(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
