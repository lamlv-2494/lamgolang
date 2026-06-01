package handlers

import (
	"food_delivery/internal/models/dto/requests"
	"food_delivery/internal/models/dto/responses"
)

// ─── Mock UserService ─────────────────────────────────────────────────────────

type mockUserService struct {
	registerFn        func(requests.RegisterRequest) (*responses.UserResponse, error)
	loginFn           func(requests.LoginRequest) (*responses.UserResponse, error)
	getCurrentUserFn  func(uint) (*responses.UserResponse, error)
	updateUserFn      func(uint, *requests.UpdateUserRequest) (*responses.UserResponse, error)
	getAllUsersFn     func(int, int) (*responses.ListResponse[*responses.UserData], error)
	adminUpdateUserFn func(uint, *requests.AdminUpdateUserRequest) (*responses.UserResponse, error)
	adminDeleteUserFn func(uint) error
}

func (m *mockUserService) Register(req requests.RegisterRequest) (*responses.UserResponse, error) {
	return m.registerFn(req)
}
func (m *mockUserService) Login(req requests.LoginRequest) (*responses.UserResponse, error) {
	return m.loginFn(req)
}
func (m *mockUserService) GetCurrentUser(id uint) (*responses.UserResponse, error) {
	return m.getCurrentUserFn(id)
}
func (m *mockUserService) UpdateUser(id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error) {
	return m.updateUserFn(id, req)
}
func (m *mockUserService) GetAllUsers(p, l int) (*responses.ListResponse[*responses.UserData], error) {
	return m.getAllUsersFn(p, l)
}
func (m *mockUserService) AdminUpdateUser(id uint, req *requests.AdminUpdateUserRequest) (*responses.UserResponse, error) {
	return m.adminUpdateUserFn(id, req)
}
func (m *mockUserService) AdminDeleteUser(id uint) error { return m.adminDeleteUserFn(id) }

// ─── Mock CategoryService ─────────────────────────────────────────────────────

type mockCategoryService struct {
	createCategoryFn func(requests.CreateCategoryRequest) (*responses.CategoryResponse, error)
	updateCategoryFn func(uint, requests.UpdateCategoryRequest) (*responses.CategoryResponse, error)
	deleteCategoryFn func(uint) error
	getCategoriesFn  func(int, int) (responses.ListResponse[*responses.CategoryResponse], error)
}

func (m *mockCategoryService) CreateCategory(req requests.CreateCategoryRequest) (*responses.CategoryResponse, error) {
	return m.createCategoryFn(req)
}
func (m *mockCategoryService) UpdateCategory(id uint, req requests.UpdateCategoryRequest) (*responses.CategoryResponse, error) {
	return m.updateCategoryFn(id, req)
}
func (m *mockCategoryService) DeleteCategory(id uint) error { return m.deleteCategoryFn(id) }
func (m *mockCategoryService) GetCategories(p, l int) (responses.ListResponse[*responses.CategoryResponse], error) {
	return m.getCategoriesFn(p, l)
}

// ─── Mock ProductService ──────────────────────────────────────────────────────

type mockProductService struct {
	createProductFn  func(requests.CreateProductRequest) (*responses.ProductData, error)
	updateProductFn  func(uint, requests.UpdateProductRequest) (*responses.ProductData, error)
	deleteProductFn  func(uint) error
	getProductsFn    func(string, string, uint, float64, float64, float64, string, int, int) (*responses.ListResponse[*responses.ProductData], error)
	getProductByIDFn func(uint) (*responses.ProductData, error)
}

func (m *mockProductService) CreateProduct(req requests.CreateProductRequest) (*responses.ProductData, error) {
	return m.createProductFn(req)
}
func (m *mockProductService) UpdateProduct(id uint, req requests.UpdateProductRequest) (*responses.ProductData, error) {
	return m.updateProductFn(id, req)
}
func (m *mockProductService) DeleteProduct(id uint) error { return m.deleteProductFn(id) }
func (m *mockProductService) GetProducts(search, classify string, catID uint, minP, maxP, minR float64, sort string, p, l int) (*responses.ListResponse[*responses.ProductData], error) {
	return m.getProductsFn(search, classify, catID, minP, maxP, minR, sort, p, l)
}
func (m *mockProductService) GetProductByID(id uint) (*responses.ProductData, error) {
	return m.getProductByIDFn(id)
}

// ─── Mock CartService ─────────────────────────────────────────────────────────

type mockCartService struct {
	addToCartFn      func(uint, *requests.AddToCartRequest) error
	getCartFn        func(uint) (*responses.CartResponse, error)
	updateCartItemFn func(uint, uint, *requests.UpdateCartItemRequest) error
	removeFromCartFn func(uint, uint) error
}

func (m *mockCartService) AddToCart(userID uint, req *requests.AddToCartRequest) error {
	return m.addToCartFn(userID, req)
}
func (m *mockCartService) GetCart(userID uint) (*responses.CartResponse, error) {
	return m.getCartFn(userID)
}
func (m *mockCartService) UpdateCartItem(userID, itemID uint, req *requests.UpdateCartItemRequest) error {
	return m.updateCartItemFn(userID, itemID, req)
}
func (m *mockCartService) RemoveFromCart(userID, itemID uint) error {
	return m.removeFromCartFn(userID, itemID)
}

// ─── Mock OrderService ────────────────────────────────────────────────────────

type mockOrderService struct {
	checkoutFn               func(uint) (*responses.OrderResponse, error)
	getOrderHistoryFn        func(uint, int, int) (*responses.OrderListResponse, error)
	adminGetAllOrdersFn      func(int, int) (*responses.ListResponse[*responses.OrderResponse], error)
	adminUpdateOrderStatusFn func(uint, requests.UpdateOrderStatusRequest) (*responses.OrderResponse, error)
}

func (m *mockOrderService) Checkout(userID uint) (*responses.OrderResponse, error) {
	return m.checkoutFn(userID)
}
func (m *mockOrderService) GetOrderHistory(userID uint, p, l int) (*responses.OrderListResponse, error) {
	return m.getOrderHistoryFn(userID, p, l)
}
func (m *mockOrderService) AdminGetAllOrders(p, l int) (*responses.ListResponse[*responses.OrderResponse], error) {
	return m.adminGetAllOrdersFn(p, l)
}
func (m *mockOrderService) AdminUpdateOrderStatus(id uint, req requests.UpdateOrderStatusRequest) (*responses.OrderResponse, error) {
	return m.adminUpdateOrderStatusFn(id, req)
}

// ─── Mock RatingService ───────────────────────────────────────────────────────

type mockRatingService struct {
	createRatingFn       func(uint, uint, requests.CreateRatingRequest) error
	getRatingByIDFn      func(uint, uint) (*responses.RatingResponse, error)
	getRatingsByUserIDFn func(uint, int, int) (*responses.ListResponse[*responses.RatingResponseData], error)
}

func (m *mockRatingService) CreateRating(userID, productID uint, req requests.CreateRatingRequest) error {
	return m.createRatingFn(userID, productID, req)
}
func (m *mockRatingService) GetRatingByID(userID, productID uint) (*responses.RatingResponse, error) {
	return m.getRatingByIDFn(userID, productID)
}
func (m *mockRatingService) GetRatingsByUserID(userID uint, p, l int) (*responses.ListResponse[*responses.RatingResponseData], error) {
	return m.getRatingsByUserIDFn(userID, p, l)
}

// ─── Mock SuggestionService ───────────────────────────────────────────────────

type mockSuggestionService struct {
	createSuggestionFn  func(uint, requests.CreateSuggestionRequest) error
	getAllSuggestionsFn func(int, int) (*responses.ListResponse[*responses.SuggestionData], error)
}

func (m *mockSuggestionService) CreateSuggestion(userID uint, req requests.CreateSuggestionRequest) error {
	return m.createSuggestionFn(userID, req)
}
func (m *mockSuggestionService) GetAllSuggestions(p, l int) (*responses.ListResponse[*responses.SuggestionData], error) {
	return m.getAllSuggestionsFn(p, l)
}
