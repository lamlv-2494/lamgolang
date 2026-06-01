package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─── CreateCategory ───────────────────────────────────────────────────────────

func TestCategoryHandler_Create_Success(t *testing.T) {
	svc := &mockCategoryService{
		createCategoryFn: func(req requests.CreateCategoryRequest) (*responses.CategoryResponse, error) {
			return &responses.CategoryResponse{Name: "Food"}, nil
		},
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(map[string]string{"name": "Food"})
	h.CreateCategory(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCategoryHandler_Create_BindError(t *testing.T) {
	svc := &mockCategoryService{}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx.Request.Header.Set("Content-Type", "application/json")
	h.CreateCategory(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestCategoryHandler_Create_ServiceError(t *testing.T) {
	svc := &mockCategoryService{
		createCategoryFn: func(req requests.CreateCategoryRequest) (*responses.CategoryResponse, error) {
			return nil, configs.CategoryAlreadyExists
		},
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(map[string]interface{}{"name": "Food"})
	h.CreateCategory(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

// ─── UpdateCategory ───────────────────────────────────────────────────────────

func TestCategoryHandler_Update_Success(t *testing.T) {
	svc := &mockCategoryService{
		updateCategoryFn: func(id uint, req requests.UpdateCategoryRequest) (*responses.CategoryResponse, error) {
			return &responses.CategoryResponse{Name: "Updated"}, nil
		},
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(map[string]string{"name": "Updated"})
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateCategory(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_Update_InvalidID(t *testing.T) {
	svc := &mockCategoryService{}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	h.UpdateCategory(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_Update_BindError(t *testing.T) {
	svc := &mockCategoryService{}
	h := NewCategoryHandler(svc)
	ctx, w := newBadJSONCtx()
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateCategory(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_Update_ServiceError(t *testing.T) {
	svc := &mockCategoryService{
		updateCategoryFn: func(id uint, req requests.UpdateCategoryRequest) (*responses.CategoryResponse, error) {
			return nil, configs.CategoryNotFound
		},
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(map[string]string{"name": "X"})
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateCategory(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── DeleteCategory ───────────────────────────────────────────────────────────

func TestCategoryHandler_Delete_Success(t *testing.T) {
	svc := &mockCategoryService{
		deleteCategoryFn: func(id uint) error { return nil },
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.DeleteCategory(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_Delete_InvalidID(t *testing.T) {
	svc := &mockCategoryService{}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
	h.DeleteCategory(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_Delete_ServiceError(t *testing.T) {
	svc := &mockCategoryService{
		deleteCategoryFn: func(id uint) error { return configs.CategoryNotFound },
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.DeleteCategory(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetCategories ────────────────────────────────────────────────────────────

func TestCategoryHandler_GetAll_Success(t *testing.T) {
	svc := &mockCategoryService{
		getCategoriesFn: func(p, l int) (responses.ListResponse[*responses.CategoryResponse], error) {
			return responses.ListResponse[*responses.CategoryResponse]{}, nil
		},
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=1&limit=10", nil)
	h.GetCategories(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_GetAll_ServiceError(t *testing.T) {
	svc := &mockCategoryService{
		getCategoriesFn: func(p, l int) (responses.ListResponse[*responses.CategoryResponse], error) {
			return responses.ListResponse[*responses.CategoryResponse]{}, errors.New("db error")
		},
	}
	h := NewCategoryHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetCategories(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// helper for bad JSON requests
func newBadJSONCtx() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	ctx.Request.Header.Set("Content-Type", "application/json")
	// empty body with content-type json triggers bind error
	return ctx, w
}
