package handlers

import (
	"net/http"

	"github.com/Dancoi/gogen_backend/internal/service"
	"github.com/Dancoi/gogen_backend/pkg/utils/errors"
	"github.com/gin-gonic/gin"
)

// AuthHandler struct для обработки auth endpoints
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создаёт новый AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest для регистрации
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterResponse ответ при регистрации
type RegisterResponse struct {
	ID        int32  `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

// Register регистрирует нового пользователя
// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		if err == errors.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		if err == errors.ErrUsernameAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		IsActive:  user.IsActive.Bool,
		CreatedAt: user.CreatedAt.Time.String(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user":    response,
	})
}

// LoginRequest для авторизации
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse ответ при авторизации
type LoginResponse struct {
	ID       int32  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// Login авторизует пользователя
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == errors.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := LoginResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Token:    token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    response,
	})
}

// Logout выходит из аккаунта
// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Реализовать отзыв сессии
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
