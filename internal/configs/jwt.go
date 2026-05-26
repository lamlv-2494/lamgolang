package configs

import (
	"food_delivery/internal/utils/constants"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, role string) (string, error) {
	secretKey := os.Getenv(constants.JWTSecret)

	claims := jwt.MapClaims{
		constants.UserID:   userID,
		constants.UserRole: role,
		constants.Exp:      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (uint, string, error) {
	secretKey := os.Getenv(constants.JWTSecret)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, AlgorithmMismatch
		}
		return []byte(secretKey), nil
	})

	if err != nil || token == nil || !token.Valid {
		return 0, "", InvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", CantGetClaims
	}

	userId := uint(claims[constants.UserID].(float64))
	role, _ := claims[constants.UserRole].(string)
	return userId, role, nil
}
