package middleware

import (
	"net/http"
	"strings"

	"github.com/Dancoi/gogen_backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен и устанавливает userID в контекст
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Извлекаем токен (Bearer scheme)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Проверяем токен
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Устанавливаем userID в контекст
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// ApiTokenMiddleware проверяет API токен и устанавливает userID в контекст
func ApiTokenMiddleware(tokenService service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка X-API-Token
		apiToken := c.GetHeader("X-API-Token")
		if apiToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API token"})
			c.Abort()
			return
		}

		// Проверяем токен
		token, err := tokenService.ValidateApiToken(c.Request.Context(), apiToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API token"})
			c.Abort()
			return
		}

		// Устанавливаем userID в контекст
		c.Set("userID", token.UserID)
		c.Set("tokenID", token.ID)

		c.Next()
	}
}
