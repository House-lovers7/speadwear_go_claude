package dto

import "time"

// CreateItemRequest represents item creation request
type CreateItemRequest struct {
	SuperItem    string  `json:"super_item" binding:"required"`
	Season       int     `json:"season" binding:"required,min=1,max=5"`
	TPO          int     `json:"tpo" binding:"required,min=1,max=5"`
	Color        int     `json:"color" binding:"required,min=1,max=15"`
	Content      string  `json:"content"`
	Memo         string  `json:"memo"`
	Rating       float32 `json:"rating" binding:"min=0,max=5"`
}

// UpdateItemRequest represents item update request
type UpdateItemRequest struct {
	SuperItem    *string  `json:"super_item"`
	Season       *int     `json:"season" binding:"omitempty,min=1,max=5"`
	TPO          *int     `json:"tpo" binding:"omitempty,min=1,max=5"`
	Color        *int     `json:"color" binding:"omitempty,min=1,max=15"`
	Content      *string  `json:"content"`
	Memo         *string  `json:"memo"`
	Rating       *float32 `json:"rating" binding:"omitempty,min=0,max=5"`
}

// ItemResponse represents item data in responses
type ItemResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	CoordinateID *uint     `json:"coordinate_id,omitempty"`
	SuperItem    string    `json:"super_item"`
	Season       int       `json:"season"`
	TPO          int       `json:"tpo"`
	Color        int       `json:"color"`
	Content      string    `json:"content"`
	Memo         string    `json:"memo"`
	Picture      string    `json:"picture"`
	Rating       float32   `json:"rating"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ItemListResponse represents paginated item list response
type ItemListResponse struct {
	Items      []ItemResponse `json:"items"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}

// ItemFilterRequest represents item search filters
type ItemFilterRequest struct {
	Season    *int     `form:"season" binding:"omitempty,min=1,max=5"`
	TPO       *int     `form:"tpo" binding:"omitempty,min=1,max=5"`
	Color     *int     `form:"color" binding:"omitempty,min=1,max=15"`
	SuperItem *string  `form:"super_item"`
	MinRating *float32 `form:"min_rating" binding:"omitempty,min=0,max=5"`
	MaxRating *float32 `form:"max_rating" binding:"omitempty,min=0,max=5"`
	Page      int      `form:"page,default=1" binding:"min=1"`
	PerPage   int      `form:"per_page,default=20" binding:"min=1,max=100"`
}