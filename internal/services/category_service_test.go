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

func newCatRepo() *mockCategoryRepo {
	return &mockCategoryRepo{
		createFn:     func(c *entities.Category) error { return nil },
		findByIDFn:   func(id uint) (*entities.Category, error) { return nil, errors.New("not found") },
		findByNameFn: func(s string) (*entities.Category, error) { return nil, errors.New("not found") },
		updateFn:     func(c *entities.Category) error { return nil },
		deleteFn:     func(c *entities.Category) error { return nil },
		listFn:       func(p, l int) ([]*entities.Category, int64, error) { return nil, 0, nil },
	}
}

// ─── CreateCategory ───────────────────────────────────────────────────────────

func TestCategoryService_Create_Success(t *testing.T) {
	repo := newCatRepo()
	svc := NewCategoryService(repo)

	resp, err := svc.CreateCategory(requests.CreateCategoryRequest{Name: "Burgers", Description: "Tasty"})
	require.NoError(t, err)
	assert.Equal(t, "Burgers", resp.Name)
}

func TestCategoryService_Create_AlreadyExists(t *testing.T) {
	repo := newCatRepo()
	repo.findByNameFn = func(s string) (*entities.Category, error) {
		return &entities.Category{Name: s}, nil
	}
	svc := NewCategoryService(repo)

	_, err := svc.CreateCategory(requests.CreateCategoryRequest{Name: "Burgers"})
	assert.Equal(t, configs.CategoryAlreadyExists, err)
}

func TestCategoryService_Create_CreateFails(t *testing.T) {
	repo := newCatRepo()
	repo.createFn = func(c *entities.Category) error { return errors.New("db error") }
	svc := NewCategoryService(repo)

	_, err := svc.CreateCategory(requests.CreateCategoryRequest{Name: "New"})
	assert.Equal(t, configs.CreateCategoryFailed, err)
}

// ─── UpdateCategory ───────────────────────────────────────────────────────────

func TestCategoryService_Update_Success(t *testing.T) {
	repo := newCatRepo()
	cat := &entities.Category{Name: "Old", Description: "Old desc"}
	cat.ID = 1
	repo.findByIDFn = func(id uint) (*entities.Category, error) { return cat, nil }
	svc := NewCategoryService(repo)

	name := "New"
	desc := "New desc"
	req := requests.UpdateCategoryRequest{Name: &name, Description: &desc}
	resp, err := svc.UpdateCategory(1, req)
	require.NoError(t, err)
	assert.Equal(t, "New", resp.Name)
	assert.Equal(t, "New desc", resp.Description)
}

func TestCategoryService_Update_NilFields(t *testing.T) {
	repo := newCatRepo()
	cat := &entities.Category{Name: "Same"}
	cat.ID = 1
	repo.findByIDFn = func(id uint) (*entities.Category, error) { return cat, nil }
	svc := NewCategoryService(repo)

	resp, err := svc.UpdateCategory(1, requests.UpdateCategoryRequest{})
	require.NoError(t, err)
	assert.Equal(t, "Same", resp.Name)
}

func TestCategoryService_Update_NotFound(t *testing.T) {
	repo := newCatRepo()
	svc := NewCategoryService(repo)

	_, err := svc.UpdateCategory(9999, requests.UpdateCategoryRequest{})
	assert.Equal(t, configs.CategoryNotFound, err)
}

func TestCategoryService_Update_UpdateFails(t *testing.T) {
	repo := newCatRepo()
	cat := &entities.Category{Name: "Cat"}
	cat.ID = 1
	repo.findByIDFn = func(id uint) (*entities.Category, error) { return cat, nil }
	repo.updateFn = func(c *entities.Category) error { return errors.New("db error") }
	svc := NewCategoryService(repo)

	_, err := svc.UpdateCategory(1, requests.UpdateCategoryRequest{})
	assert.Equal(t, configs.UpdateCategoryFailed, err)
}

// ─── DeleteCategory ───────────────────────────────────────────────────────────

func TestCategoryService_Delete_Success(t *testing.T) {
	repo := newCatRepo()
	cat := &entities.Category{}
	cat.ID = 1
	repo.findByIDFn = func(id uint) (*entities.Category, error) { return cat, nil }
	svc := NewCategoryService(repo)

	err := svc.DeleteCategory(1)
	assert.NoError(t, err)
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	repo := newCatRepo()
	svc := NewCategoryService(repo)

	err := svc.DeleteCategory(9999)
	assert.Equal(t, configs.CategoryNotFound, err)
}

// ─── GetCategories ────────────────────────────────────────────────────────────

func TestCategoryService_GetCategories_Success(t *testing.T) {
	repo := newCatRepo()
	repo.listFn = func(p, l int) ([]*entities.Category, int64, error) {
		return []*entities.Category{
			{Name: "A"},
			{Name: "B"},
		}, 2, nil
	}
	svc := NewCategoryService(repo)

	resp, err := svc.GetCategories(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), resp.TotalCount)
	assert.Len(t, resp.Items, 2)
}

func TestCategoryService_GetCategories_Fails(t *testing.T) {
	repo := newCatRepo()
	repo.listFn = func(p, l int) ([]*entities.Category, int64, error) {
		return nil, 0, errors.New("db error")
	}
	svc := NewCategoryService(repo)

	_, err := svc.GetCategories(1, 10)
	assert.Error(t, err)
}
