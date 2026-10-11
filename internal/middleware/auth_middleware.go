package middleware

import (
	"net/http"
	"strings"

	"objexrg-backend/internal/services/auth"
	common_result "objexrg-backend/internal/shared/common"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenService *auth.TokenService
}

func NewAuthMiddleware(tokenService *auth.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}

func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			res := common_result.Fail[any]("Authorization header required", http.StatusUnauthorized)
			c.JSON(res.Code, res)
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			res := common_result.Fail[any]("Invalid authorization header format", http.StatusUnauthorized)
			c.JSON(res.Code, res)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := m.tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			res := common_result.Fail[any]("Invalid or expired token", http.StatusUnauthorized)
			c.JSON(res.Code, res)
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.Subject)
		c.Next()
	}
}

func (m *AuthMiddleware) OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}
