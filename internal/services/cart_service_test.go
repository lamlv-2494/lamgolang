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

func newCartRepo() *mockCartRepo {
	return &mockCartRepo{
		createFn:               func(i *entities.CartItem) error { return nil },
		findByIDFn:             func(id uint) (*entities.CartItem, error) { return nil, errors.New("not found") },
		findByUserAndProductFn: func(uID, pID uint) (*entities.CartItem, error) { return nil, errors.New("not found") },
		updateFn:               func(i *entities.CartItem) error { return nil },
		deleteFn:               func(i *entities.CartItem) error { return nil },
		listByUserIDFn:         func(uID uint) ([]*entities.CartItem, error) { return nil, nil },
	}
}

// ─── AddToCart ────────────────────────────────────────────────────────────────

func TestCartService_AddToCart_NewItem(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return baseProduct(id), nil }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.AddToCart(1, &requests.AddToCartRequest{ProductID: 1, Quantity: 2})
	assert.NoError(t, err)
}

func TestCartService_AddToCart_ExistingItem(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return baseProduct(id), nil }
	existing := &entities.CartItem{UserID: 1, ProductID: 1, Quantity: 2}
	existing.ID = 1
	cartRepo.findByUserAndProductFn = func(uID, pID uint) (*entities.CartItem, error) {
		return existing, nil
	}
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.AddToCart(1, &requests.AddToCartRequest{ProductID: 1, Quantity: 3})
	assert.NoError(t, err)
	assert.Equal(t, 5, existing.Quantity)
}

func TestCartService_AddToCart_ExistingItem_UpdateFails(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return baseProduct(id), nil }
	existing := &entities.CartItem{UserID: 1, ProductID: 1, Quantity: 1}
	cartRepo.findByUserAndProductFn = func(uID, pID uint) (*entities.CartItem, error) { return existing, nil }
	cartRepo.updateFn = func(i *entities.CartItem) error { return errors.New("db error") }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.AddToCart(1, &requests.AddToCartRequest{ProductID: 1, Quantity: 1})
	assert.Equal(t, configs.UpdateCartItemFailed, err)
}

func TestCartService_AddToCart_ProductNotFound(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.AddToCart(1, &requests.AddToCartRequest{ProductID: 99, Quantity: 1})
	assert.Equal(t, configs.ProductNotFound, err)
}

func TestCartService_AddToCart_CreateFails(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return baseProduct(id), nil }
	cartRepo.createFn = func(i *entities.CartItem) error { return errors.New("db error") }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.AddToCart(1, &requests.AddToCartRequest{ProductID: 1, Quantity: 1})
	assert.Equal(t, configs.AddToCartFailed, err)
}

// ─── GetCart ──────────────────────────────────────────────────────────────────

func TestCartService_GetCart_Empty(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	svc := NewCartService(cartRepo, prodRepo)

	resp, err := svc.GetCart(1)
	require.NoError(t, err)
	assert.Empty(t, resp.Items)
	assert.Equal(t, 0.0, resp.TotalAmount)
}

func TestCartService_GetCart_WithItems(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	prod := baseProduct(1)
	prod.Price = 15.0
	cartRepo.listByUserIDFn = func(uID uint) ([]*entities.CartItem, error) {
		item := &entities.CartItem{UserID: uID, ProductID: 1, Quantity: 2}
		item.ID = 1
		return []*entities.CartItem{item}, nil
	}
	prodRepo.findByIDFn = func(id uint) (*entities.Product, error) { return prod, nil }
	svc := NewCartService(cartRepo, prodRepo)

	resp, err := svc.GetCart(1)
	require.NoError(t, err)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, 30.0, resp.TotalAmount)
}

func TestCartService_GetCart_ProductMissing_Skipped(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	cartRepo.listByUserIDFn = func(uID uint) ([]*entities.CartItem, error) {
		item := &entities.CartItem{UserID: uID, ProductID: 99, Quantity: 1}
		return []*entities.CartItem{item}, nil
	}
	// product not found → item skipped
	svc := NewCartService(cartRepo, prodRepo)

	resp, err := svc.GetCart(1)
	require.NoError(t, err)
	assert.Empty(t, resp.Items)
	assert.Equal(t, 0.0, resp.TotalAmount)
}

func TestCartService_GetCart_ListFails(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	cartRepo.listByUserIDFn = func(uID uint) ([]*entities.CartItem, error) {
		return nil, errors.New("db error")
	}
	svc := NewCartService(cartRepo, prodRepo)

	_, err := svc.GetCart(1)
	assert.Error(t, err)
}

// ─── UpdateCartItem ───────────────────────────────────────────────────────────

func TestCartService_UpdateCartItem_Success(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	item := &entities.CartItem{UserID: 1, ProductID: 1, Quantity: 1}
	item.ID = 1
	cartRepo.findByIDFn = func(id uint) (*entities.CartItem, error) { return item, nil }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.UpdateCartItem(1, 1, &requests.UpdateCartItemRequest{Quantity: 5})
	assert.NoError(t, err)
	assert.Equal(t, 5, item.Quantity)
}

func TestCartService_UpdateCartItem_NotFound(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.UpdateCartItem(1, 9999, &requests.UpdateCartItemRequest{Quantity: 1})
	assert.Equal(t, configs.CartItemNotFound, err)
}

func TestCartService_UpdateCartItem_AccessDenied(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	item := &entities.CartItem{UserID: 2, ProductID: 1, Quantity: 1}
	item.ID = 1
	cartRepo.findByIDFn = func(id uint) (*entities.CartItem, error) { return item, nil }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.UpdateCartItem(1, 1, &requests.UpdateCartItemRequest{Quantity: 2})
	assert.Equal(t, configs.AccessDenied, err)
}

func TestCartService_UpdateCartItem_UpdateFails(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	item := &entities.CartItem{UserID: 1, ProductID: 1, Quantity: 1}
	item.ID = 1
	cartRepo.findByIDFn = func(id uint) (*entities.CartItem, error) { return item, nil }
	cartRepo.updateFn = func(i *entities.CartItem) error { return errors.New("db error") }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.UpdateCartItem(1, 1, &requests.UpdateCartItemRequest{Quantity: 3})
	assert.Equal(t, configs.UpdateCartItemFailed, err)
}

// ─── RemoveFromCart ───────────────────────────────────────────────────────────

func TestCartService_RemoveFromCart_Success(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	item := &entities.CartItem{UserID: 1, ProductID: 1}
	item.ID = 1
	cartRepo.findByIDFn = func(id uint) (*entities.CartItem, error) { return item, nil }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.RemoveFromCart(1, 1)
	assert.NoError(t, err)
}

func TestCartService_RemoveFromCart_NotFound(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.RemoveFromCart(1, 9999)
	assert.Equal(t, configs.CartItemNotFound, err)
}

func TestCartService_RemoveFromCart_AccessDenied(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	item := &entities.CartItem{UserID: 2}
	item.ID = 1
	cartRepo.findByIDFn = func(id uint) (*entities.CartItem, error) { return item, nil }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.RemoveFromCart(1, 1)
	assert.Equal(t, configs.AccessDenied, err)
}

func TestCartService_RemoveFromCart_DeleteFails(t *testing.T) {
	cartRepo := newCartRepo()
	prodRepo := newProdRepo()
	item := &entities.CartItem{UserID: 1}
	item.ID = 1
	cartRepo.findByIDFn = func(id uint) (*entities.CartItem, error) { return item, nil }
	cartRepo.deleteFn = func(i *entities.CartItem) error { return errors.New("db error") }
	svc := NewCartService(cartRepo, prodRepo)

	err := svc.RemoveFromCart(1, 1)
	assert.Equal(t, configs.RemoveFromCartFailed, err)
}
