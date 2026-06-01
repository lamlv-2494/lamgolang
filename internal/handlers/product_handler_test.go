package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// helper: create a multipart form request with an image file field
func newMultipartCtx(fields map[string]string, includeImage bool) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	if includeImage {
		part, _ := writer.CreateFormFile("image", "test.jpg")
		_, _ = io.WriteString(part, "fake image data")
	}

	writer.Close()

	ctx.Request = httptest.NewRequest(http.MethodPost, "/", body)
	ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return ctx, w
}

// ─── CreateProduct ────────────────────────────────────────────────────────────

func TestProductHandler_Create_Success(t *testing.T) {
	svc := &mockProductService{
		createProductFn: func(req requests.CreateProductRequest) (*responses.ProductData, error) {
			return &responses.ProductData{Name: "Burger"}, nil
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newMultipartCtx(map[string]string{
		"name": "Burger", "price": "10.5", "category_id": "1",
	}, true)
	h.CreateProduct(ctx)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductHandler_Create_NoImage(t *testing.T) {
	svc := &mockProductService{}
	h := NewProductHandler(svc)
	// Provide all required fields but no image → FormFile fails
	ctx, w := newMultipartCtx(map[string]string{"name": "Burger", "price": "10.5", "category_id": "1"}, false)
	h.CreateProduct(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestProductHandler_Create_BindError(t *testing.T) {
	svc := &mockProductService{}
	h := NewProductHandler(svc)
	// Provide a non-numeric price so binding fails with a type mismatch
	ctx, w := newMultipartCtx(map[string]string{"name": "Burger", "price": "notanumber", "category_id": "1"}, true)
	h.CreateProduct(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestProductHandler_Create_ServiceError(t *testing.T) {
	svc := &mockProductService{
		createProductFn: func(req requests.CreateProductRequest) (*responses.ProductData, error) {
			return nil, configs.CategoryNotFound
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newMultipartCtx(map[string]string{
		"name": "Burger", "price": "10.5", "category_id": "99",
	}, true)
	h.CreateProduct(ctx)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

// ─── UpdateProduct ────────────────────────────────────────────────────────────

func TestProductHandler_Update_Success(t *testing.T) {
	svc := &mockProductService{
		updateProductFn: func(id uint, req requests.UpdateProductRequest) (*responses.ProductData, error) {
			return &responses.ProductData{Name: "Updated"}, nil
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newMultipartCtx(map[string]string{"name": "Updated"}, false)
	ctx.Request.Method = http.MethodPut
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateProduct(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_Update_WithImage(t *testing.T) {
	svc := &mockProductService{
		updateProductFn: func(id uint, req requests.UpdateProductRequest) (*responses.ProductData, error) {
			return &responses.ProductData{Name: "Updated"}, nil
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newMultipartCtx(map[string]string{"name": "Updated"}, true)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateProduct(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_Update_InvalidID(t *testing.T) {
	svc := &mockProductService{}
	h := NewProductHandler(svc)
	ctx, w := newMultipartCtx(map[string]string{"name": "X"}, false)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
	h.UpdateProduct(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestProductHandler_Update_ServiceError(t *testing.T) {
	svc := &mockProductService{
		updateProductFn: func(id uint, req requests.UpdateProductRequest) (*responses.ProductData, error) {
			return nil, configs.ProductNotFound
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newMultipartCtx(map[string]string{"name": "X"}, false)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.UpdateProduct(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── DeleteProduct ────────────────────────────────────────────────────────────

func TestProductHandler_Delete_Success(t *testing.T) {
	svc := &mockProductService{
		deleteProductFn: func(id uint) error { return nil },
	}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.DeleteProduct(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_Delete_InvalidID(t *testing.T) {
	svc := &mockProductService{}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "bad"}}
	h.DeleteProduct(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestProductHandler_Delete_ServiceError(t *testing.T) {
	svc := &mockProductService{
		deleteProductFn: func(id uint) error { return configs.ProductNotFound },
	}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.DeleteProduct(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetProducts ──────────────────────────────────────────────────────────────

func TestProductHandler_GetProducts_Success(t *testing.T) {
	svc := &mockProductService{
		getProductsFn: func(search, classify string, catID uint, minP, maxP, minR float64, sort string, p, l int) (*responses.ListResponse[*responses.ProductData], error) {
			return &responses.ListResponse[*responses.ProductData]{}, nil
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?search=x&classify=hot&category_id=1&min_price=5&max_price=50&rating=3&sort=price&page=1&limit=10"), nil)
	h.GetProducts(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetProducts_DefaultParams(t *testing.T) {
	svc := &mockProductService{
		getProductsFn: func(search, classify string, catID uint, minP, maxP, minR float64, sort string, p, l int) (*responses.ListResponse[*responses.ProductData], error) {
			return &responses.ListResponse[*responses.ProductData]{}, nil
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetProducts(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetProducts_ServiceError(t *testing.T) {
	svc := &mockProductService{
		getProductsFn: func(search, classify string, catID uint, minP, maxP, minR float64, sort string, p, l int) (*responses.ListResponse[*responses.ProductData], error) {
			return nil, errors.New("db error")
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	h.GetProducts(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ─── GetProductByID ───────────────────────────────────────────────────────────

func TestProductHandler_GetProductByID_Success(t *testing.T) {
	svc := &mockProductService{
		getProductByIDFn: func(id uint) (*responses.ProductData, error) {
			return &responses.ProductData{Name: "Burger"}, nil
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.GetProductByID(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetProductByID_InvalidID(t *testing.T) {
	svc := &mockProductService{}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "xyz"}}
	h.GetProductByID(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestProductHandler_GetProductByID_ServiceError(t *testing.T) {
	svc := &mockProductService{
		getProductByIDFn: func(id uint) (*responses.ProductData, error) {
			return nil, configs.ProductNotFound
		},
	}
	h := NewProductHandler(svc)
	ctx, w := newTestCtx(nil)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	h.GetProductByID(ctx)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
