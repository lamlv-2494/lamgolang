package services

import (
	"errors"
	"testing"

	"food_delivery/internal/configs"
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newProdRepo() *mockProductRepo {
	return &mockProductRepo{
		createFn:   func(p *entities.Product) error { return nil },
		findByIDFn: func(id uint) (*entities.Product, error) { return nil, errors.New("not found") },
		updateFn:   func(p *entities.Product) error { return nil },
		deleteFn:   func(p *entities.Product) error { return nil },
		listFn: func(s, c string, cID uint, minP, maxP, minR float64, sort string, pg, lim int) ([]*entities.Product, int64, error) {
			return nil, 0, nil
		},
	}
}

func baseProduct(id uint) *entities.Product {
	p := &entities.Product{
		Name:     "Burger",
		Price:    10.0,
		Category: entities.Category{Name: "Food"},
	}
	p.ID = id
	return p
}

// ─── CreateProduct ────────────────────────────────────────────────────────────

func TestProductService_Create_Success(t *testing.T) {
	catRepo := newCatRepo()
	catRepo.findByIDFn = func(id uint) (*entities.Category, error) {
		return &entities.Category{Name: "Food"}, nil
	}
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) {
		return baseProduct(id), nil
	}
	svc := NewProductService(prodRepo, catRepo)

	resp, err := svc.CreateProduct(requests.CreateProductRequest{
		Name: "Burger", Price: 10.0, CategoryID: 1,
	})
	require.NoError(t, err)
	assert.Equal(t, "Burger", resp.Name)
}

func TestProductService_Create_CategoryNotFound(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.CreateProduct(requests.CreateProductRequest{Name: "X", Price: 5, CategoryID: 99})
	assert.Equal(t, configs.CategoryNotFound, err)
}

func TestProductService_Create_CreateFails(t *testing.T) {
	catRepo := newCatRepo()
	catRepo.findByIDFn = func(id uint) (*entities.Category, error) {
		return &entities.Category{}, nil
	}
	prodRepo := newProdRepo()
	prodRepo.createFn = func(p *entities.Product) error { return errors.New("db error") }
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.CreateProduct(requests.CreateProductRequest{Name: "X", Price: 5, CategoryID: 1})
	assert.Equal(t, configs.CreateProductFailed, err)
}

func TestProductService_Create_FindAfterCreateFails(t *testing.T) {
	catRepo := newCatRepo()
	catRepo.findByIDFn = func(id uint) (*entities.Category, error) {
		return &entities.Category{}, nil
	}
	prodRepo := newProdRepo()
	// create succeeds but subsequent FindByID fails
	created := false
	prodRepo.createFn = func(p *entities.Product) error { created = true; return nil }
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) {
		if created {
			return nil, errors.New("not found")
		}
		return nil, errors.New("not found")
	}
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.CreateProduct(requests.CreateProductRequest{Name: "X", Price: 5, CategoryID: 1})
	assert.Equal(t, configs.ProductNotFound, err)
}

// ─── UpdateProduct ────────────────────────────────────────────────────────────

func TestProductService_Update_Success(t *testing.T) {
	catRepo := newCatRepo()
	catRepo.findByIDFn = func(id uint) (*entities.Category, error) { return &entities.Category{Name: "Food"}, nil }
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	call := 0
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) {
		call++
		return prod, nil
	}
	svc := NewProductService(prodRepo, catRepo)

	name := "Updated"
	desc := "Desc"
	price := 20.0
	img := "/img.png"
	catID := uint(2)
	req := requests.UpdateProductRequest{Name: &name, Description: &desc, Price: &price, Image: &img, CategoryID: &catID}
	resp, err := svc.UpdateProduct(1, req)
	require.NoError(t, err)
	assert.Equal(t, "Updated", resp.Name)
}

func TestProductService_Update_NilFields(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return prod, nil }
	svc := NewProductService(prodRepo, catRepo)

	resp, err := svc.UpdateProduct(1, requests.UpdateProductRequest{})
	require.NoError(t, err)
	assert.Equal(t, "Burger", resp.Name)
}

func TestProductService_Update_ProductNotFound(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.UpdateProduct(9999, requests.UpdateProductRequest{})
	assert.Equal(t, configs.ProductNotFound, err)
}

func TestProductService_Update_CategoryNotFound(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return prod, nil }
	svc := NewProductService(prodRepo, catRepo)

	catID := uint(99)
	req := requests.UpdateProductRequest{CategoryID: &catID}
	_, err := svc.UpdateProduct(1, req)
	assert.Equal(t, configs.CategoryNotFound, err)
}

func TestProductService_Update_UpdateFails(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	call := 0
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) {
		call++
		if call == 1 {
			return prod, nil
		}
		return nil, errors.New("not found")
	}
	prodRepo.updateFn = func(p *entities.Product) error { return errors.New("db error") }
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.UpdateProduct(1, requests.UpdateProductRequest{})
	assert.Equal(t, configs.UpdateProductFailed, err)
}

func TestProductService_Update_FindAfterUpdateFails(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	call := 0
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) {
		call++
		if call == 1 {
			return prod, nil
		}
		return nil, errors.New("not found after update")
	}
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.UpdateProduct(1, requests.UpdateProductRequest{})
	assert.Equal(t, configs.ProductNotFound, err)
}

// ─── DeleteProduct ────────────────────────────────────────────────────────────

func TestProductService_Delete_Success(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return prod, nil }
	svc := NewProductService(prodRepo, catRepo)

	err := svc.DeleteProduct(1)
	assert.NoError(t, err)
}

func TestProductService_Delete_NotFound(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	svc := NewProductService(prodRepo, catRepo)

	err := svc.DeleteProduct(9999)
	assert.Equal(t, configs.ProductNotFound, err)
}

func TestProductService_Delete_DeleteFails(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return prod, nil }
	prodRepo.deleteFn = func(p *entities.Product) error { return errors.New("db error") }
	svc := NewProductService(prodRepo, catRepo)

	err := svc.DeleteProduct(1)
	assert.Equal(t, configs.DeleteProductFailed, err)
}

// ─── GetProducts ──────────────────────────────────────────────────────────────

func TestProductService_GetProducts_Empty(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	svc := NewProductService(prodRepo, catRepo)

	resp, err := svc.GetProducts("", "", 0, 0, 0, 0, "", 1, 10)
	require.NoError(t, err)
	assert.Empty(t, resp.Items)
}

func TestProductService_GetProducts_WithRatings(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	stars5 := 5
	stars3 := 3
	prod := &entities.Product{
		Name:     "Test",
		Price:    10.0,
		Category: entities.Category{Name: "Food"},
		Ratings: []*entities.Rating{
			{Stars: stars5},
			{Stars: stars3},
		},
	}
	prod.ID = 1
	prodRepo.listFn = func(s, c string, cID uint, minP, maxP, minR float64, sort string, pg, lim int) ([]*entities.Product, int64, error) {
		return []*entities.Product{prod}, 1, nil
	}
	svc := NewProductService(prodRepo, catRepo)

	resp, err := svc.GetProducts("", "", 0, 0, 0, 0, "", 1, 10)
	require.NoError(t, err)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, 4.0, resp.Items[0].Rating) // (5+3)/2
}

func TestProductService_GetProducts_Fails(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prodRepo.listFn = func(s, c string, cID uint, minP, maxP, minR float64, sort string, pg, lim int) ([]*entities.Product, int64, error) {
		return nil, 0, errors.New("db error")
	}
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.GetProducts("", "", 0, 0, 0, 0, "", 1, 10)
	assert.Equal(t, configs.ProductNotFound, err)
}

// ─── GetProductByID ───────────────────────────────────────────────────────────

func TestProductService_GetProductByID_Success_NoRatings(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return prod, nil }
	svc := NewProductService(prodRepo, catRepo)

	resp, err := svc.GetProductByID(1)
	require.NoError(t, err)
	assert.Equal(t, "Burger", resp.Name)
	assert.Equal(t, 0.0, resp.Rating)
}

func TestProductService_GetProductByID_WithRatings(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	prod.Ratings = []*entities.Rating{{Stars: 4}, {Stars: 2}}
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return prod, nil }
	svc := NewProductService(prodRepo, catRepo)

	resp, err := svc.GetProductByID(1)
	require.NoError(t, err)
	assert.Equal(t, 3.0, resp.Rating) // (4+2)/2
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	catRepo := newCatRepo()
	prodRepo := newProdRepo()
	svc := NewProductService(prodRepo, catRepo)

	_, err := svc.GetProductByID(9999)
	assert.Equal(t, configs.ProductNotFound, err)
}
