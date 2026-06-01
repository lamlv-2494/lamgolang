package repositories

import (
	"testing"

	"food_delivery/internal/models/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedOrderFixtures(t *testing.T, userRepo UserRepository, catRepo CategoryRepository, prodRepo ProductRepository, cartRepo CartRepository, username, email string) (*entities.User, *entities.Product) {
	t.Helper()
	user := &entities.User{Username: username, Email: email, Password: "pass"}
	require.NoError(t, userRepo.CreateUser(user))

	cat := &entities.Category{Name: "OrderCat_" + username}
	require.NoError(t, catRepo.Create(cat))

	product := &entities.Product{Name: "OrderProd_" + username, Price: 20.0, CategoryID: cat.ID}
	require.NoError(t, prodRepo.Create(product))

	cartItem := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 2}
	require.NoError(t, cartRepo.Create(cartItem))

	return user, product
}

func buildOrder(userID, productID uint, price float64, status string) *entities.Order {
	return &entities.Order{
		UserID:     userID,
		TotalPrice: price,
		Status:     status,
		OrderItems: []*entities.OrderItem{
			{ProductID: productID, Quantity: 1, Price: price},
		},
	}
}

func TestOrderRepository_CreateOrderWithTransaction_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)
	cartRepo := NewCartRepository(db)

	user, product := seedOrderFixtures(t, userRepo, catRepo, prodRepo, cartRepo, "orderuser1", "ou1@example.com")

	// Verify cart item exists before order
	items, err := cartRepo.ListByUserID(user.ID)
	require.NoError(t, err)
	assert.Len(t, items, 1)

	order := buildOrder(user.ID, product.ID, 20.0, "pending")
	err = repo.CreateOrderWithTransaction(order)
	require.NoError(t, err)
	assert.NotZero(t, order.ID)

	// Cart should be cleared after order
	itemsAfter, err := cartRepo.ListByUserID(user.ID)
	require.NoError(t, err)
	assert.Empty(t, itemsAfter)
}

func TestOrderRepository_ListByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)
	cartRepo := NewCartRepository(db)

	user, product := seedOrderFixtures(t, userRepo, catRepo, prodRepo, cartRepo, "orderuser2", "ou2@example.com")

	statuses := []string{"pending", "completed", "cancelled"}
	for _, s := range statuses {
		order := buildOrder(user.ID, product.ID, 10.0, s)
		require.NoError(t, repo.CreateOrderWithTransaction(order))
		// re-add cart for next order (cleared by transaction)
		cartItem := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 1}
		require.NoError(t, cartRepo.Create(cartItem))
	}

	orders, total, totalAmount, statusCounts, err := repo.ListByUserID(user.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, orders, 3)
	assert.Equal(t, 30.0, totalAmount)
	assert.Equal(t, int64(1), statusCounts["pending"])
	assert.Equal(t, int64(1), statusCounts["completed"])
	assert.Equal(t, int64(1), statusCounts["cancelled"])
}

func TestOrderRepository_ListByUserID_Pagination(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)
	cartRepo := NewCartRepository(db)

	user, product := seedOrderFixtures(t, userRepo, catRepo, prodRepo, cartRepo, "orderuser3", "ou3@example.com")

	for i := 0; i < 5; i++ {
		order := buildOrder(user.ID, product.ID, float64(i+1)*5.0, "pending")
		require.NoError(t, repo.CreateOrderWithTransaction(order))
		cartItem := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 1}
		require.NoError(t, cartRepo.Create(cartItem))
	}

	orders, total, _, _, err := repo.ListByUserID(user.ID, 1, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, orders, 3)
}

func TestOrderRepository_FindByOrderID_Found(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)
	cartRepo := NewCartRepository(db)

	user, product := seedOrderFixtures(t, userRepo, catRepo, prodRepo, cartRepo, "orderuser4", "ou4@example.com")
	order := buildOrder(user.ID, product.ID, 50.0, "processing")
	require.NoError(t, repo.CreateOrderWithTransaction(order))

	found, err := repo.FindByOrderID(order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, found.ID)
	assert.Equal(t, "processing", found.Status)
}

func TestOrderRepository_FindByOrderID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)

	_, err := repo.FindByOrderID(9999)
	assert.Error(t, err)
}

func TestOrderRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)
	cartRepo := NewCartRepository(db)

	user, product := seedOrderFixtures(t, userRepo, catRepo, prodRepo, cartRepo, "orderuser5", "ou5@example.com")
	order := buildOrder(user.ID, product.ID, 30.0, "pending")
	require.NoError(t, repo.CreateOrderWithTransaction(order))

	order.Status = "completed"
	err := repo.Update(order)
	require.NoError(t, err)

	found, err := repo.FindByOrderID(order.ID)
	require.NoError(t, err)
	assert.Equal(t, "completed", found.Status)
}

func TestOrderRepository_ListAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)
	userRepo := NewUserRepository(db)
	catRepo := NewCategoryRepository(db)
	prodRepo := NewProductRepository(db)
	cartRepo := NewCartRepository(db)

	user, product := seedOrderFixtures(t, userRepo, catRepo, prodRepo, cartRepo, "orderuser6", "ou6@example.com")

	for i := 0; i < 4; i++ {
		order := buildOrder(user.ID, product.ID, float64(i+1)*10.0, "pending")
		require.NoError(t, repo.CreateOrderWithTransaction(order))
		cartItem := &entities.CartItem{UserID: user.ID, ProductID: product.ID, Quantity: 1}
		require.NoError(t, cartRepo.Create(cartItem))
	}

	orders, total, err := repo.ListAll(1, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(4), total)
	assert.Len(t, orders, 3)

	orders2, _, err := repo.ListAll(2, 3)
	require.NoError(t, err)
	assert.Len(t, orders2, 1)
}
