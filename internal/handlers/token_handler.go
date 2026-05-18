package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Dancoi/gogen_backend/internal/service"
	"github.com/Dancoi/gogen_backend/pkg/utils/errors"
	"github.com/gin-gonic/gin"
)

// TokenHandler struct для обработки token endpoints
type TokenHandler struct {
	tokenService service.TokenService
}

// NewTokenHandler создаёт новый TokenHandler
func NewTokenHandler(tokenService service.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
	}
}

// GenerateTokenRequest для генерации нового токена
type GenerateTokenRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// GenerateTokenResponse ответ при генерации токена
type GenerateTokenResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Note  string `json:"note"`
}

// GenerateToken генерирует новый API токен
// POST /api/tokens
func (h *TokenHandler) GenerateToken(c *gin.Context) {
	var req GenerateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из контекста (должен быть установлен middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDInterface.(int32)

	token, err := h.tokenService.GenerateApiToken(c.Request.Context(), userID, req.Name)
	if err != nil {
		if err == errors.ErrSubscriptionNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "subscription not found"})
			return
		}
		if err == errors.ErrSubscriptionExpired {
			c.JSON(http.StatusForbidden, gin.H{"error": "subscription has expired"})
			return
		}
		if err == errors.ErrTokenLimitExceeded {
			c.JSON(http.StatusForbidden, gin.H{"error": "token generation limit exceeded"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := GenerateTokenResponse{
		Token: token,
		Name:  req.Name,
		Note:  "Save this token in a secure place. You won't be able to see it again!",
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "token generated successfully",
		"data":    response,
	})
}

// ListTokensResponse структура для ответа списка токенов
type ListTokensResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	LastUsed  string `json:"last_used"`
}

// ListTokens получает все активные токены пользователя (отозванные скрыты)
// GET /api/tokens
func (h *TokenHandler) ListTokens(c *gin.Context) {
	// Получаем userID из контекста
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDInterface.(int32)

	tokens, err := h.tokenService.ListUserTokens(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []ListTokensResponse
	for _, token := range tokens {
		// Пропускаем отозванные токены
		if !token.IsActive.Bool || token.RevokedAt.Valid {
			continue
		}

		lastUsed := ""
		if token.LastUsedAt.Valid {
			lastUsed = token.LastUsedAt.Time.String()
		}

		response = append(response, ListTokensResponse{
			ID:        token.ID,
			Name:      token.Name,
			CreatedAt: token.CreatedAt.Time.String(),
			LastUsed:  lastUsed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": response,
	})
}

// RevokeToken отзывает токен
// DELETE /api/tokens/:id
func (h *TokenHandler) RevokeToken(c *gin.Context) {
	tokenIDStr := c.Param("id")
	tokenID, err := strconv.ParseInt(tokenIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token id"})
		return
	}

	// Получаем userID из контекста (установлен AuthMiddleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(int32)

	// Проверяем что токен принадлежит текущему пользователю
	token, err := h.tokenService.GetTokenByID(c.Request.Context(), int32(tokenID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}

	if token.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to revoke this token"})
		return
	}

	err = h.tokenService.RevokeApiToken(c.Request.Context(), int32(tokenID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token revoked successfully",
	})
}
