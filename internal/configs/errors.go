package configs

import (
	"errors"
	"net/http"
)

var (
	AlgorithmMismatch = errors.New("Algorithm mismatch")
	InvalidToken      = errors.New("Invalid token")
	CantGetClaims     = errors.New("Can't get claims")

	UserNotFound       = errors.New("user not found")
	InvalidPassword    = errors.New("Invalid email or password")
	CreateUserFailed   = errors.New("Failed to create user")
	GenerateTokenFail  = errors.New("Failed to generate token")
	EmailOrUserTaken   = errors.New("Email or Username has already been taken")
	UsernameTaken      = errors.New("Username has already been taken")
	EmailTaken         = errors.New("Email has already been taken")
	MissingAuthHeader  = errors.New("Authorization header is required")
	InvalidTokenFormat = errors.New("Invalid token format")
	AccessDenied       = errors.New("Access denied. Admin role required.")
	UpdateUserFailed   = errors.New("Failed to update user")
)

var errorStatusMap = map[error]int{
	UserNotFound:       http.StatusNotFound,
	InvalidPassword:    http.StatusUnauthorized,
	EmailOrUserTaken:   http.StatusUnprocessableEntity,
	UsernameTaken:      http.StatusUnprocessableEntity,
	EmailTaken:         http.StatusUnprocessableEntity,
	CreateUserFailed:   http.StatusInternalServerError,
	GenerateTokenFail:  http.StatusInternalServerError,
	UpdateUserFailed:   http.StatusInternalServerError,
	InvalidToken:       http.StatusUnauthorized,
	MissingAuthHeader:  http.StatusUnauthorized,
	InvalidTokenFormat: http.StatusUnauthorized,
	AccessDenied:       http.StatusForbidden,
	AlgorithmMismatch:  http.StatusUnauthorized,
	CantGetClaims:      http.StatusUnauthorized,
}

func GetStatusCodeAndMessage(err error) (int, string) {
	statusCode, exists := errorStatusMap[err]
	if !exists {
		return http.StatusInternalServerError, "Internal Server Error"
	}
	return statusCode, err.Error()
}
