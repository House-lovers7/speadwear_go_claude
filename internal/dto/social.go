package dto

import "time"

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Comment string `json:"comment" binding:"required,min=1,max=1000"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Comment string `json:"comment" binding:"required,min=1,max=1000"`
}

// CommentResponse represents comment data in responses
type CommentResponse struct {
	ID           uint         `json:"id"`
	UserID       uint         `json:"user_id"`
	CoordinateID uint         `json:"coordinate_id"`
	Comment      string       `json:"comment"`
	User         UserResponse `json:"user"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// CommentListResponse represents paginated comment list response
type CommentListResponse struct {
	Comments   []CommentResponse `json:"comments"`
	TotalCount int64             `json:"total_count"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
}

// NotificationResponse represents notification data in responses
type NotificationResponse struct {
	ID               uint                 `json:"id"`
	SenderID         uint                 `json:"sender_id"`
	ReceiverID       uint                 `json:"receiver_id"`
	Action           string               `json:"action"`
	Checked          bool                 `json:"checked"`
	Sender           UserResponse         `json:"sender"`
	Coordinate       *CoordinateResponse  `json:"coordinate,omitempty"`
	Comment          *CommentResponse     `json:"comment,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
}

// NotificationListResponse represents paginated notification list response
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	TotalCount    int64                  `json:"total_count"`
	UnreadCount   int64                  `json:"unread_count"`
	Page          int                    `json:"page"`
	PerPage       int                    `json:"per_page"`
}

// FollowResponse represents follow relationship response
type FollowResponse struct {
	IsFollowing bool `json:"is_following"`
}

// BlockResponse represents block relationship response
type BlockResponse struct {
	IsBlocked bool `json:"is_blocked"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Page    int `form:"page,default=1" binding:"min=1"`
	PerPage int `form:"per_page,default=20" binding:"min=1,max=100"`
}