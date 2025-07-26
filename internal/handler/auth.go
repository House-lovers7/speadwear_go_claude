package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/House-lovers7/speadwear-go/internal/dto"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
	"github.com/House-lovers7/speadwear-go/pkg/config"
)

type AuthHandler struct {
	config      *config.Config
	userUsecase usecase.UserUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(cfg *config.Config, userUsecase usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{
		config:      cfg,
		userUsecase: userUsecase,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid email or password" || err.Error() == "account not activated" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// Signup handles user registration
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := h.userUsecase.Signup(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, authResponse)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// クライアント側でトークンを削除する必要があるため、
	// サーバー側では特に処理は不要
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// "Bearer " プレフィックスを削除
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	newToken, err := h.userUsecase.RefreshToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Calculate expiration time
	duration, _ := time.ParseDuration(h.config.JWT.Expiration)
	expiresAt := time.Now().Add(duration)

	c.JSON(http.StatusOK, gin.H{
		"token":      newToken,
		"expires_at": expiresAt,
	})
}

// Me returns current user information
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, err := h.userUsecase.GetUser(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Picture:   user.Picture,
		Admin:     user.Admin,
		Activated: user.Activated,
		CreatedAt: user.CreatedAt,
	})
}