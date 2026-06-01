package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─── ResponseSuccess ──────────────────────────────────────────────────────────

func TestResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	ResponseSuccess(ctx, http.StatusOK, gin.H{"msg": "ok"})
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── ResponseError ────────────────────────────────────────────────────────────

func TestResponseError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	ResponseError(ctx, nil) // nil uses fallback error mapping
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetIDParam ───────────────────────────────────────────────────────────────

func TestGetIDParam_Valid(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "42"}}

	id, err := GetIDParam(ctx)
	assert.NoError(t, err)
	assert.Equal(t, uint(42), id)
}

func TestGetIDParam_Invalid(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "not-a-number"}}

	_, err := GetIDParam(ctx)
	assert.Error(t, err)
}

// ─── GetPaginationParams ──────────────────────────────────────────────────────

func TestGetPaginationParams_Defaults(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	page, limit := GetPaginationParams(ctx)
	assert.Equal(t, 1, page)
	assert.Equal(t, 10, limit)
}

func TestGetPaginationParams_ValidValues(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=3&limit=20", nil)

	page, limit := GetPaginationParams(ctx)
	assert.Equal(t, 3, page)
	assert.Equal(t, 20, limit)
}

func TestGetPaginationParams_InvalidValues(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=-1&limit=bad", nil)

	page, limit := GetPaginationParams(ctx)
	// invalid values should fall back to defaults
	assert.Equal(t, 1, page)
	assert.Equal(t, 10, limit)
}

func TestGetPaginationParams_ZeroValues(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/?page=0&limit=0", nil)

	page, limit := GetPaginationParams(ctx)
	assert.Equal(t, 1, page)
	assert.Equal(t, 10, limit)
}
