package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestUserForRating(t *testing.T, repo UserRepository, username, email string) *entities.User {
	t.Helper()
	user := &entities.User{Username: username, Email: email, Password: "pass"}
	require.NoError(t, repo.CreateUser(user))
	return user
}

func createTestProductForRating(t *testing.T, pRepo ProductRepository, cRepo CategoryRepository, name string) *entities.Product {
	t.Helper()
	cat := &entities.Category{Name: "RatingCat_" + name}
	require.NoError(t, cRepo.Create(cat))
	p := &entities.Product{Name: name, Price: 10.0, CategoryID: cat.ID}
	require.NoError(t, pRepo.Create(p))
	return p
}

func TestRatingRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRatingRepository(db)
	userRepo := NewUserRepository(db)
	productRepo := NewProductRepository(db)
	catRepo := NewCategoryRepository(db)

	user := createTestUserForRating(t, userRepo, "rater1", "rater1@example.com")
	product := createTestProductForRating(t, productRepo, catRepo, "RatedProduct1")

	rating := &entities.Rating{
		UserID:    user.ID,
		ProductID: product.ID,
		Stars:     5,
		Comment:   "Excellent!",
	}

	err := repo.Create(rating)
	require.NoError(t, err)
	assert.NotZero(t, rating.ID)
}

func TestRatingRepository_FindByUserAndProduct_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRatingRepository(db)
	userRepo := NewUserRepository(db)
	productRepo := NewProductRepository(db)
	catRepo := NewCategoryRepository(db)

	user := createTestUserForRating(t, userRepo, "rater2", "rater2@example.com")
	product := createTestProductForRating(t, productRepo, catRepo, "RatedProduct2")

	rating := &entities.Rating{UserID: user.ID, ProductID: product.ID, Stars: 4, Comment: "Good"}
	require.NoError(t, repo.Create(rating))

	found, err := repo.FindByUserAndProduct(user.ID, product.ID)
	require.NoError(t, err)
	assert.Equal(t, 4, found.Stars)
	assert.Equal(t, "Good", found.Comment)
}

func TestRatingRepository_FindByUserAndProduct_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRatingRepository(db)

	_, err := repo.FindByUserAndProduct(9999, 9999)
	assert.Error(t, err)
}

func TestRatingRepository_GetRatingsByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRatingRepository(db)
	userRepo := NewUserRepository(db)
	productRepo := NewProductRepository(db)
	catRepo := NewCategoryRepository(db)

	user := createTestUserForRating(t, userRepo, "rater3", "rater3@example.com")
	product1 := createTestProductForRating(t, productRepo, catRepo, "ProductA")
	product2 := createTestProductForRating(t, productRepo, catRepo, "ProductB")
	product3 := createTestProductForRating(t, productRepo, catRepo, "ProductC")

	require.NoError(t, repo.Create(&entities.Rating{UserID: user.ID, ProductID: product1.ID, Stars: 5}))
	require.NoError(t, repo.Create(&entities.Rating{UserID: user.ID, ProductID: product2.ID, Stars: 3}))
	require.NoError(t, repo.Create(&entities.Rating{UserID: user.ID, ProductID: product3.ID, Stars: 4}))

	ratings, total, err := repo.GetRatingsByUserID(user.ID, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, ratings, 2)

	// page 2
	ratings2, _, err := repo.GetRatingsByUserID(user.ID, 2, 2)
	require.NoError(t, err)
	assert.Len(t, ratings2, 1)
}

func TestRatingRepository_GetRatingsByUserID_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRatingRepository(db)

	ratings, total, err := repo.GetRatingsByUserID(9999, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, ratings)
}
