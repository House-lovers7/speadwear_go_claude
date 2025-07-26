package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/House-lovers7/speadwear-go/internal/dto"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// GetUser GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userUsecase.GetUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

// GetMe GET /api/v1/users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware
	
	// Check if userID is set (not zero)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userUsecase.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

// ListUsers GET /api/v1/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	users, total, err := h.userUsecase.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Picture:   user.Picture,
			Admin:     user.Admin,
			Activated: user.Activated,
			CreatedAt: user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		Users:      userResponses,
		TotalCount: total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
	})
}

// UpdateProfile PUT /api/v1/users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var req struct {
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.UpdateProfile(c.Request.Context(), userID, req.Name, req.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// UpdateUser PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user can update (must be same user or admin)
	currentUserID := c.GetUint("userID")
	if currentUserID != uint(userID) {
		// TODO: Check admin permission
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userUsecase.UpdateUser(c.Request.Context(), uint(userID), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user can delete (must be same user or admin)
	currentUserID := c.GetUint("userID")
	if currentUserID != uint(userID) {
		// TODO: Check admin permission
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err = h.userUsecase.DeleteUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ChangePassword PUT /api/v1/users/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "invalid old password" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ResetPasswordRequest POST /api/v1/users/password/reset
func (h *UserHandler) ResetPasswordRequest(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.ResetPasswordRequest(c.Request.Context(), req.Email)
	if err != nil {
		// Don't reveal if email exists or not
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
}

// ResetPassword PUT /api/v1/users/password/reset
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dto.PasswordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ActivateAccount POST /api/v1/users/activate
func (h *UserHandler) ActivateAccount(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.ActivateAccount(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account activated successfully"})
}

// ResendActivationEmail POST /api/v1/users/activate/resend
func (h *UserHandler) ResendActivationEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.ResendActivationEmail(c.Request.Context(), req.Email)
	if err != nil {
		// Don't reveal if email exists or not
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists and account is not activated, an activation link has been sent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists and account is not activated, an activation link has been sent"})
}