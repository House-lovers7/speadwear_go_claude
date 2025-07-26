package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/House-lovers7/speadwear-go/internal/dto"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
)

type SocialHandler struct {
	socialUsecase usecase.SocialUsecase
}

// NewSocialHandler creates a new social handler
func NewSocialHandler(socialUsecase usecase.SocialUsecase) *SocialHandler {
	return &SocialHandler{
		socialUsecase: socialUsecase,
	}
}

// CreateComment POST /api/v1/comments
func (h *SocialHandler) CreateComment(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var req struct {
		CoordinateID uint   `json:"coordinate_id" binding:"required"`
		Comment      string `json:"comment" binding:"required,min=1,max=1000"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.socialUsecase.CreateComment(c.Request.Context(), userID, req.CoordinateID, req.Comment)
	if err != nil {
		if err.Error() == "you are blocked by this user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CommentResponse{
		ID:           comment.ID,
		UserID:       comment.UserID,
		CoordinateID: comment.CoordinateID,
		Comment:      comment.Comment,
		User: dto.UserResponse{
			ID:      userID,
			// Other fields will be populated if user is preloaded
		},
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}

// UpdateComment PUT /api/v1/comments/:id
func (h *SocialHandler) UpdateComment(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.socialUsecase.UpdateComment(c.Request.Context(), userID, uint(commentID), req.Comment)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

// DeleteComment DELETE /api/v1/comments/:id
func (h *SocialHandler) DeleteComment(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	err = h.socialUsecase.DeleteComment(c.Request.Context(), userID, uint(commentID))
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// FollowUser POST /api/v1/follow/:user_id
func (h *SocialHandler) FollowUser(c *gin.Context) {
	followerID := c.GetUint("userID") // From auth middleware

	followedID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.socialUsecase.FollowUser(c.Request.Context(), followerID, uint(followedID))
	if err != nil {
		if err.Error() == "cannot follow yourself" || err.Error() == "already following" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "you are blocked by this user" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

// UnfollowUser DELETE /api/v1/follow/:user_id
func (h *SocialHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetUint("userID") // From auth middleware

	followedID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.socialUsecase.UnfollowUser(c.Request.Context(), followerID, uint(followedID))
	if err != nil {
		if err.Error() == "not following" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

// GetFollowers GET /api/v1/follow/followers
func (h *SocialHandler) GetFollowers(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	// Check if getting followers for a specific user
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			userID = uint(id)
		}
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	followers, total, err := h.socialUsecase.GetFollowers(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponses := make([]dto.UserResponse, len(followers))
	for i, user := range followers {
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

// GetFollowing GET /api/v1/follow/following
func (h *SocialHandler) GetFollowing(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	// Check if getting following for a specific user
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			userID = uint(id)
		}
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	following, total, err := h.socialUsecase.GetFollowing(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponses := make([]dto.UserResponse, len(following))
	for i, user := range following {
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

// CheckFollowStatus GET /api/v1/follow/status/:user_id
func (h *SocialHandler) CheckFollowStatus(c *gin.Context) {
	followerID := c.GetUint("userID") // From auth middleware

	followedID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	isFollowing, err := h.socialUsecase.IsFollowing(c.Request.Context(), followerID, uint(followedID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.FollowResponse{
		IsFollowing: isFollowing,
	})
}

// BlockUser POST /api/v1/blocks/:user_id
func (h *SocialHandler) BlockUser(c *gin.Context) {
	blockerID := c.GetUint("userID") // From auth middleware

	blockedID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.socialUsecase.BlockUser(c.Request.Context(), blockerID, uint(blockedID))
	if err != nil {
		if err.Error() == "cannot block yourself" || err.Error() == "already blocked" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// UnblockUser DELETE /api/v1/blocks/:user_id
func (h *SocialHandler) UnblockUser(c *gin.Context) {
	blockerID := c.GetUint("userID") // From auth middleware

	blockedID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.socialUsecase.UnblockUser(c.Request.Context(), blockerID, uint(blockedID))
	if err != nil {
		if err.Error() == "not blocked" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unblocked successfully"})
}

// GetBlockedUsers GET /api/v1/blocks
func (h *SocialHandler) GetBlockedUsers(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	users, err := h.socialUsecase.GetBlockedUsers(c.Request.Context(), userID)
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

	c.JSON(http.StatusOK, userResponses)
}

// CheckBlockStatus GET /api/v1/blocks/status/:user_id
func (h *SocialHandler) CheckBlockStatus(c *gin.Context) {
	blockerID := c.GetUint("userID") // From auth middleware

	blockedID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	isBlocked, err := h.socialUsecase.IsBlocked(c.Request.Context(), blockerID, uint(blockedID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.BlockResponse{
		IsBlocked: isBlocked,
	})
}

// GetNotifications GET /api/v1/notifications
func (h *SocialHandler) GetNotifications(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	notifications, total, err := h.socialUsecase.GetNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	unreadCount, _ := h.socialUsecase.GetUnreadNotificationCount(c.Request.Context(), userID)

	notificationResponses := make([]dto.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		resp := dto.NotificationResponse{
			ID:         notification.ID,
			SenderID:   notification.SenderID,
			ReceiverID: notification.ReceiverID,
			Action:     notification.Action,
			Checked:    notification.Checked,
			CreatedAt:  notification.CreatedAt,
		}

		// Add sender info if available
		if notification.Sender.ID != 0 {
			resp.Sender = dto.UserResponse{
				ID:      notification.Sender.ID,
				Name:    notification.Sender.Name,
				Email:   notification.Sender.Email,
				Picture: notification.Sender.Picture,
			}
		}

		// Add coordinate info if available
		if notification.Coordinate != nil && notification.Coordinate.ID != 0 {
			resp.Coordinate = &dto.CoordinateResponse{
				ID:      notification.Coordinate.ID,
				UserID:  notification.Coordinate.UserID,
				Picture: notification.Coordinate.Picture,
			}
		}

		// Add comment info if available
		if notification.Comment != nil && notification.Comment.ID != 0 {
			resp.Comment = &dto.CommentResponse{
				ID:      notification.Comment.ID,
				Comment: notification.Comment.Comment,
			}
		}

		notificationResponses[i] = resp
	}

	c.JSON(http.StatusOK, dto.NotificationListResponse{
		Notifications: notificationResponses,
		TotalCount:    total,
		UnreadCount:   unreadCount,
		Page:          pagination.Page,
		PerPage:       pagination.PerPage,
	})
}

// GetUnreadNotifications GET /api/v1/notifications/unread
func (h *SocialHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	notifications, err := h.socialUsecase.GetUnreadNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	notificationResponses := make([]dto.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		resp := dto.NotificationResponse{
			ID:         notification.ID,
			SenderID:   notification.SenderID,
			ReceiverID: notification.ReceiverID,
			Action:     notification.Action,
			Checked:    notification.Checked,
			CreatedAt:  notification.CreatedAt,
		}

		// Add sender info if available
		if notification.Sender.ID != 0 {
			resp.Sender = dto.UserResponse{
				ID:      notification.Sender.ID,
				Name:    notification.Sender.Name,
				Email:   notification.Sender.Email,
				Picture: notification.Sender.Picture,
			}
		}

		notificationResponses[i] = resp
	}

	c.JSON(http.StatusOK, notificationResponses)
}

// MarkAsRead PUT /api/v1/notifications/:id/read
func (h *SocialHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	notificationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	err = h.socialUsecase.MarkNotificationAsRead(c.Request.Context(), userID, uint(notificationID))
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead PUT /api/v1/notifications/read_all
func (h *SocialHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	err := h.socialUsecase.MarkAllNotificationsAsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

// GetUnreadCount GET /api/v1/notifications/unread/count
func (h *SocialHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	count, err := h.socialUsecase.GetUnreadNotificationCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}