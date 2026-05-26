package middlewares

import (
	"food_delivery/internal/configs"
	"food_delivery/internal/utils/constants"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(constants.Authorization)
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			ctx.Abort()
			return
		}

		part := strings.Split(authHeader, " ")
		if len(part) != 2 || part[0] != constants.Token {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}

		userID, role, err := configs.ValidateToken(part[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set(constants.UserID, userID)
		ctx.Set(constants.UserRole, role)
		ctx.Next()
	}
}

// AdminMiddleware ensures only admin users can access the endpoint
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get(constants.UserRole)
		if !exists || role.(string) != constants.AdminRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
