package services

import (
	"food_delivery/internal/models/entities"
)

// ─── User Repo Mock ──────────────────────────────────────────────────────────

type mockUserRepo struct {
	createUserFn     func(*entities.User) error
	findByEmailFn    func(string) (*entities.User, error)
	findByIDFn       func(uint) (*entities.User, error)
	findByUsernameFn func(string) (*entities.User, error)
	updateUserFn     func(*entities.User) error
	findAllUsersFn   func(int, int) ([]*entities.User, int64, error)
	deleteUserFn     func(*entities.User) error
}

func (m *mockUserRepo) CreateUser(u *entities.User) error            { return m.createUserFn(u) }
func (m *mockUserRepo) FindByEmail(e string) (*entities.User, error) { return m.findByEmailFn(e) }
func (m *mockUserRepo) FindByID(id uint) (*entities.User, error)     { return m.findByIDFn(id) }
func (m *mockUserRepo) FindByUsername(n string) (*entities.User, error) {
	return m.findByUsernameFn(n)
}
func (m *mockUserRepo) UpdateUser(u *entities.User) error { return m.updateUserFn(u) }
func (m *mockUserRepo) FindAllUsers(p, l int) ([]*entities.User, int64, error) {
	return m.findAllUsersFn(p, l)
}
func (m *mockUserRepo) DeleteUser(u *entities.User) error { return m.deleteUserFn(u) }

// ─── Category Repo Mock ──────────────────────────────────────────────────────

type mockCategoryRepo struct {
	createFn     func(*entities.Category) error
	findByIDFn   func(uint) (*entities.Category, error)
	findByNameFn func(string) (*entities.Category, error)
	updateFn     func(*entities.Category) error
	deleteFn     func(*entities.Category) error
	listFn       func(int, int) ([]*entities.Category, int64, error)
}

func (m *mockCategoryRepo) Create(c *entities.Category) error            { return m.createFn(c) }
func (m *mockCategoryRepo) FindByID(id uint) (*entities.Category, error) { return m.findByIDFn(id) }
func (m *mockCategoryRepo) FindByName(n string) (*entities.Category, error) {
	return m.findByNameFn(n)
}
func (m *mockCategoryRepo) Update(c *entities.Category) error { return m.updateFn(c) }
func (m *mockCategoryRepo) Delete(c *entities.Category) error { return m.deleteFn(c) }
func (m *mockCategoryRepo) List(p, l int) ([]*entities.Category, int64, error) {
	return m.listFn(p, l)
}

// ─── Product Repo Mock ───────────────────────────────────────────────────────

type mockProductRepo struct {
	createFn   func(*entities.Product) error
	findByIDFn func(uint) (*entities.Product, error)
	updateFn   func(*entities.Product) error
	deleteFn   func(*entities.Product) error
	listFn     func(string, string, uint, float64, float64, float64, string, int, int) ([]*entities.Product, int64, error)
}

func (m *mockProductRepo) Create(p *entities.Product) error            { return m.createFn(p) }
func (m *mockProductRepo) FindByID(id uint) (*entities.Product, error) { return m.findByIDFn(id) }
func (m *mockProductRepo) Update(p *entities.Product) error            { return m.updateFn(p) }
func (m *mockProductRepo) Delete(p *entities.Product) error            { return m.deleteFn(p) }
func (m *mockProductRepo) List(search, classify string, catID uint, minP, maxP, minR float64, sort string, page, limit int) ([]*entities.Product, int64, error) {
	return m.listFn(search, classify, catID, minP, maxP, minR, sort, page, limit)
}

// ─── Rating Repo Mock ─────────────────────────────────────────────────────────

type mockRatingRepo struct {
	createFn               func(*entities.Rating) error
	findByUserAndProductFn func(uint, uint) (*entities.Rating, error)
	getRatingsByUserIDFn   func(uint, int, int) ([]*entities.Rating, int64, error)
}

func (m *mockRatingRepo) Create(r *entities.Rating) error { return m.createFn(r) }
func (m *mockRatingRepo) FindByUserAndProduct(uID, pID uint) (*entities.Rating, error) {
	return m.findByUserAndProductFn(uID, pID)
}
func (m *mockRatingRepo) GetRatingsByUserID(uID uint, p, l int) ([]*entities.Rating, int64, error) {
	return m.getRatingsByUserIDFn(uID, p, l)
}

// ─── Cart Repo Mock ───────────────────────────────────────────────────────────

type mockCartRepo struct {
	createFn               func(*entities.CartItem) error
	findByIDFn             func(uint) (*entities.CartItem, error)
	findByUserAndProductFn func(uint, uint) (*entities.CartItem, error)
	updateFn               func(*entities.CartItem) error
	deleteFn               func(*entities.CartItem) error
	listByUserIDFn         func(uint) ([]*entities.CartItem, error)
}

func (m *mockCartRepo) Create(i *entities.CartItem) error            { return m.createFn(i) }
func (m *mockCartRepo) FindByID(id uint) (*entities.CartItem, error) { return m.findByIDFn(id) }
func (m *mockCartRepo) FindByUserAndProduct(uID, pID uint) (*entities.CartItem, error) {
	return m.findByUserAndProductFn(uID, pID)
}
func (m *mockCartRepo) Update(i *entities.CartItem) error { return m.updateFn(i) }
func (m *mockCartRepo) Delete(i *entities.CartItem) error { return m.deleteFn(i) }
func (m *mockCartRepo) ListByUserID(uID uint) ([]*entities.CartItem, error) {
	return m.listByUserIDFn(uID)
}

// ─── Order Repo Mock ──────────────────────────────────────────────────────────

type mockOrderRepo struct {
	createOrderWithTransactionFn func(*entities.Order) error
	listByUserIDFn               func(uint, int, int) ([]*entities.Order, int64, float64, map[string]int64, error)
	findByOrderIDFn              func(uint) (*entities.Order, error)
	updateFn                     func(*entities.Order) error
	listAllFn                    func(int, int) ([]*entities.Order, int64, error)
}

func (m *mockOrderRepo) CreateOrderWithTransaction(o *entities.Order) error {
	return m.createOrderWithTransactionFn(o)
}
func (m *mockOrderRepo) ListByUserID(uID uint, p, l int) ([]*entities.Order, int64, float64, map[string]int64, error) {
	return m.listByUserIDFn(uID, p, l)
}
func (m *mockOrderRepo) FindByOrderID(id uint) (*entities.Order, error) {
	return m.findByOrderIDFn(id)
}
func (m *mockOrderRepo) Update(o *entities.Order) error { return m.updateFn(o) }
func (m *mockOrderRepo) ListAll(p, l int) ([]*entities.Order, int64, error) {
	return m.listAllFn(p, l)
}

// ─── Suggestion Repo Mock ─────────────────────────────────────────────────────

type mockSuggestionRepo struct {
	createFn  func(*entities.Suggestion) error
	listAllFn func(int, int) ([]*entities.Suggestion, int64, error)
}

func (m *mockSuggestionRepo) Create(s *entities.Suggestion) error { return m.createFn(s) }
func (m *mockSuggestionRepo) ListAll(p, l int) ([]*entities.Suggestion, int64, error) {
	return m.listAllFn(p, l)
}
