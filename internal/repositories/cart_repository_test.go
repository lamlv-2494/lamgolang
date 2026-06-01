package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedCartFixtures(t *testing.T, db interface{ AutoMigrate(...interface{}) error }) {}

func createUserAndProduct(t *testing.T, userRepo UserRepository, catRepo CategoryRepository, prodRepo ProductRepository, username, email, productName string) (*entities.User, *entities.Product) {
	t.Helper()
	user := &entities.User{Username: username, Email: email, Password: "pass"}
	require.NoError(t, userRepo.CreateUser(user))

	cat := &entities.Category{Name: "CartCat_" + productName}
	require.NoError(t, catRepo.Create(cat))

	product := &entities.Product{Name: productName, Price: 5.0, CategoryID: cat.ID}
	require.NoError(t, prodRepo.Create(product))

	return user, product
}

func TestCartRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)

	user, product := createUserAndProduct(t, userRepo, catRepo, prodRepo, "cartuser1", "cu1@example.com", "CartItem1")
	item := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 2}

	err := repo.Create(item)
	require.NoError(t, err)
	assert.NotZero(t, item.ID)
}

func TestCartRepository_FindByID_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)

	user, product := createUserAndProduct(t, userRepo, catRepo, prodRepo, "cartuser2", "cu2@example.com", "CartItem2")
	item := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 3}
	require.NoError(t, repo.Create(item))

	found, err := repo.FindByID(item.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, found.Quantity)
}

func TestCartRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)

	_, err := repo.FindByID(9999)
	assert.Error(t, err)
}

func TestCartRepository_FindByUserAndProduct_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)

	user, product := createUserAndProduct(t, userRepo, catRepo, prodRepo, "cartuser3", "cu3@example.com", "CartItem3")
	item := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 1}
	require.NoError(t, repo.Create(item))

	found, err := repo.FindByUserAndProduct(user.ID, product.ID)
	require.NoError(t, err)
	assert.Equal(t, item.ID, found.ID)
}

func TestCartRepository_FindByUserAndProduct_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)

	_, err := repo.FindByUserAndProduct(9999, 9999)
	assert.Error(t, err)
}

func TestCartRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)

	user, product := createUserAndProduct(t, userRepo, catRepo, prodRepo, "cartuser4", "cu4@example.com", "CartItem4")
	item := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 1}
	require.NoError(t, repo.Create(item))

	item.Quantity = 5
	err := repo.Update(item)
	require.NoError(t, err)

	found, err := repo.FindByID(item.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, found.Quantity)
}

func TestCartRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)

	user, product := createUserAndProduct(t, userRepo, catRepo, prodRepo, "cartuser5", "cu5@example.com", "CartItem5")
	item := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 1}
	require.NoError(t, repo.Create(item))

	err := repo.Delete(item)
	require.NoError(t, err)

	_, err = repo.FindByID(item.ID)
	assert.Error(t, err)
}

func TestCartRepository_ListByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)

	user, _ := createUserAndProduct(t, userRepo, catRepo, prodRepo, "cartuser6", "cu6@example.com", "CartBase")

	cat := &entities.Category{Name: "ExtraCartCat"}
	require.NoError(t, catRepo.Create(cat))

	for i := 0; i < 3; i++ {
		p := &entities.Product{Name: "CartProd" + string(rune('A'+i)), Price: float64(i + 1), CategoryID: cat.ID}
		require.NoError(t, prodRepo.Create(p))
		item := &entities.CartItem{UserID: user.ID, ProductID: p.ID, Quantity: i + 1}
		require.NoError(t, repo.Create(item))
	}

	items, err := repo.ListByUserID(user.ID)
	require.NoError(t, err)
	assert.Len(t, items, 3)
}

func TestCartRepository_ListByUserID_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCartRepository(db)

	items, err := repo.ListByUserID(9999)
	require.NoError(t, err)
	assert.Empty(t, items)
}
