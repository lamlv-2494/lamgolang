package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestCategory(t *testing.T, repo CategoryRepository, name string) *entities.Category {
	t.Helper()
	cat := &entities.Category{Name: name}
	require.NoError(t, repo.Create(cat))
	return cat
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "FastFood")
	product := &entities.Product{
		Name:       "Burger",
		Price:      9.99,
		CategoryID: cat.ID,
	}

	err := repo.Create(product)
	require.NoError(t, err)
	assert.NotZero(t, product.ID)
}

func TestProductRepository_FindByID_Found(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "Pizza")
	product := &entities.Product{Name: "Margherita", Price: 12.50, CategoryID: cat.ID}
	require.NoError(t, repo.Create(product))

	found, err := repo.FindByID(product.ID)
	require.NoError(t, err)
	assert.Equal(t, "Margherita", found.Name)
	assert.Equal(t, cat.ID, found.CategoryID)
}

func TestProductRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	_, err := repo.FindByID(9999)
	assert.Error(t, err)
}

func TestProductRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "Drinks")
	product := &entities.Product{Name: "Cola", Price: 2.50, CategoryID: cat.ID}
	require.NoError(t, repo.Create(product))

	product.Price = 3.00
	err := repo.Update(product)
	require.NoError(t, err)

	found, err := repo.FindByID(product.ID)
	require.NoError(t, err)
	assert.Equal(t, 3.00, found.Price)
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "Deletable")
	product := &entities.Product{Name: "DeleteMe", Price: 1.00, CategoryID: cat.ID}
	require.NoError(t, repo.Create(product))

	err := repo.Delete(product)
	require.NoError(t, err)

	_, err = repo.FindByID(product.ID)
	assert.Error(t, err)
}

func TestProductRepository_List_NoFilter(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "Snacks")
	for i := 0; i < 4; i++ {
		p := &entities.Product{
			Name:       "Snack" + string(rune('A'+i)),
			Price:      float64(i+1) * 2,
			CategoryID: cat.ID,
		}
		require.NoError(t, repo.Create(p))
	}

	products, total, err := repo.List("", "", 0, 0, 0, 0, "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(4), total)
	assert.Len(t, products, 4)
}

func TestProductRepository_List_SearchFilter(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "Search")
	require.NoError(t, repo.Create(&entities.Product{Name: "Apple Juice", Price: 3.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Orange Juice", Price: 4.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Burger", Price: 8.0, CategoryID: cat.ID}))

	products, total, err := repo.List("juice", "", 0, 0, 0, 0, "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, products, 2)
}

func TestProductRepository_List_CategoryFilter(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat1 := createTestCategory(t, catRepo, "Italian")
	cat2 := createTestCategory(t, catRepo, "Mexican")
	require.NoError(t, repo.Create(&entities.Product{Name: "Pasta", Price: 10.0, CategoryID: cat1.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Pizza", Price: 12.0, CategoryID: cat1.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Taco", Price: 7.0, CategoryID: cat2.ID}))

	products, total, err := repo.List("", "", cat1.ID, 0, 0, 0, "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, products, 2)
}

func TestProductRepository_List_PriceFilter(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "Priced")
	require.NoError(t, repo.Create(&entities.Product{Name: "Cheap", Price: 5.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Mid", Price: 15.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Expensive", Price: 30.0, CategoryID: cat.ID}))

	products, total, err := repo.List("", "", 0, 10.0, 20.0, 0, "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, products, 1)
	assert.Equal(t, "Mid", products[0].Name)
}

func TestProductRepository_List_SortNameAsc(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "SortTest")
	require.NoError(t, repo.Create(&entities.Product{Name: "Zebra", Price: 1.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Apple", Price: 2.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Mango", Price: 3.0, CategoryID: cat.ID}))

	products, _, err := repo.List("", "", 0, 0, 0, 0, "name_asc", 1, 10)
	require.NoError(t, err)
	require.Len(t, products, 3)
	assert.Equal(t, "Apple", products[0].Name)
	assert.Equal(t, "Zebra", products[2].Name)
}

func TestProductRepository_List_SortNameDesc(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "SortDesc")
	require.NoError(t, repo.Create(&entities.Product{Name: "Zebra", Price: 1.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Apple", Price: 2.0, CategoryID: cat.ID}))

	products, _, err := repo.List("", "", 0, 0, 0, 0, "name_desc", 1, 10)
	require.NoError(t, err)
	require.Len(t, products, 2)
	assert.Equal(t, "Zebra", products[0].Name)
}

func TestProductRepository_List_SortPriceAsc(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "SortPrice")
	require.NoError(t, repo.Create(&entities.Product{Name: "Cheap", Price: 5.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Expensive", Price: 50.0, CategoryID: cat.ID}))

	products, _, err := repo.List("", "", 0, 0, 0, 0, "price_asc", 1, 10)
	require.NoError(t, err)
	require.Len(t, products, 2)
	assert.Equal(t, "Cheap", products[0].Name)
}

func TestProductRepository_List_SortPriceDesc(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "SortPriceDesc")
	require.NoError(t, repo.Create(&entities.Product{Name: "Cheap", Price: 5.0, CategoryID: cat.ID}))
	require.NoError(t, repo.Create(&entities.Product{Name: "Expensive", Price: 50.0, CategoryID: cat.ID}))

	products, _, err := repo.List("", "", 0, 0, 0, 0, "price_desc", 1, 10)
	require.NoError(t, err)
	require.Len(t, products, 2)
	assert.Equal(t, "Expensive", products[0].Name)
}

func TestProductRepository_List_Pagination(t *testing.T) {
	db := setupTestDB(t)
	catRepo := NewCategoryRepository(db)
	repo := NewProductRepository(db)

	cat := createTestCategory(t, catRepo, "Page")
	for i := 0; i < 6; i++ {
		p := &entities.Product{Name: "Item" + string(rune('A'+i)), Price: float64(i + 1), CategoryID: cat.ID}
		require.NoError(t, repo.Create(p))
	}

	p1, total, err := repo.List("", "", 0, 0, 0, 0, "", 1, 4)
	require.NoError(t, err)
	assert.Equal(t, int64(6), total)
	assert.Len(t, p1, 4)

	p2, _, err := repo.List("", "", 0, 0, 0, 0, "", 2, 4)
	require.NoError(t, err)
	assert.Len(t, p2, 2)
}
