package configs

import (
	"errors"
	"net/http"
)

var (
	UserNotFound     = errors.New("user not found")
	CategoryNotFound = errors.New("category not found")
	ProductNotFound  = errors.New("product not found")
	CartItemNotFound = errors.New("cart item not found")

	InvalidPassword = errors.New("Invalid email or password")

	EmailOrUserTaken = errors.New("Email or Username has already been taken")
	UsernameTaken    = errors.New("Username has already been taken")
	EmailTaken       = errors.New("Email has already been taken")

	CreateUserFailed     = errors.New("Failed to create user")
	GenerateTokenFail    = errors.New("Failed to generate token")
	UpdateUserFailed     = errors.New("Failed to update user")
	CreateCategoryFailed = errors.New("failed to create category")
	UpdateCategoryFailed = errors.New("failed to update category")
	DeleteCategoryFailed = errors.New("failed to delete category")
	CreateProductFailed  = errors.New("failed to create product")
	UpdateProductFailed  = errors.New("failed to update product")
	DeleteProductFailed  = errors.New("failed to delete product")
	FetchUsersFailed     = errors.New("failed to fetch users list")
	DeleteUserFailed     = errors.New("failed to delete user")
	AddToCartFailed      = errors.New("failed to add item to cart")
	UpdateCartItemFailed = errors.New("failed to update cart item")
	RemoveFromCartFailed = errors.New("failed to remove item from cart")

	InvalidToken       = errors.New("Invalid token")
	MissingAuthHeader  = errors.New("Authorization header is required")
	InvalidTokenFormat = errors.New("Invalid token format")
	AlgorithmMismatch  = errors.New("Algorithm mismatch")
	CantGetClaims      = errors.New("Can't get claims")

	AccessDenied = errors.New("Access denied. Admin role required.")

	CategoryAlreadyExists = errors.New("category already exists")
)

var errorStatusMap = map[error]int{
	UserNotFound:     http.StatusNotFound,
	CategoryNotFound: http.StatusNotFound,
	ProductNotFound:  http.StatusNotFound,
	CartItemNotFound: http.StatusNotFound,

	InvalidPassword: http.StatusUnauthorized,

	EmailOrUserTaken: http.StatusUnprocessableEntity,
	UsernameTaken:    http.StatusUnprocessableEntity,
	EmailTaken:       http.StatusUnprocessableEntity,

	CreateUserFailed:     http.StatusInternalServerError,
	GenerateTokenFail:    http.StatusInternalServerError,
	UpdateUserFailed:     http.StatusInternalServerError,
	CreateCategoryFailed: http.StatusInternalServerError,
	UpdateCategoryFailed: http.StatusInternalServerError,
	DeleteCategoryFailed: http.StatusInternalServerError,
	CreateProductFailed:  http.StatusInternalServerError,
	UpdateProductFailed:  http.StatusInternalServerError,
	DeleteProductFailed:  http.StatusInternalServerError,
	FetchUsersFailed:     http.StatusInternalServerError,
	DeleteUserFailed:     http.StatusInternalServerError,
	AddToCartFailed:      http.StatusInternalServerError,
	UpdateCartItemFailed: http.StatusInternalServerError,
	RemoveFromCartFailed: http.StatusInternalServerError,

	InvalidToken:       http.StatusUnauthorized,
	MissingAuthHeader:  http.StatusUnauthorized,
	InvalidTokenFormat: http.StatusUnauthorized,
	AlgorithmMismatch:  http.StatusUnauthorized,
	CantGetClaims:      http.StatusUnauthorized,

	AccessDenied: http.StatusForbidden,

	CategoryAlreadyExists: http.StatusConflict,
}

func GetStatusCodeAndMessage(err error) (int, string) {
	statusCode, exists := errorStatusMap[err]
	if !exists {
		return http.StatusInternalServerError, "Internal Server Error"
	}
	return statusCode, err.Error()
}
