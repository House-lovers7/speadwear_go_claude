package dto

import "time"

// CreateCoordinateRequest represents coordinate creation request
type CreateCoordinateRequest struct {
	Season         int     `json:"season" binding:"required,min=1,max=5"`
	TPO            int     `json:"tpo" binding:"required,min=1,max=5"`
	SiTopLength    int     `json:"si_top_length" binding:"min=0,max=3"`
	SiTopSleeve    int     `json:"si_top_sleeve" binding:"min=0,max=5"`
	SiBottomLength int     `json:"si_bottom_length" binding:"min=0,max=6"`
	SiBottomType   int     `json:"si_bottom_type" binding:"min=0,max=2"`
	SiDressLength  int     `json:"si_dress_length" binding:"min=0,max=6"`
	SiDressSleeve  int     `json:"si_dress_sleeve" binding:"min=0,max=5"`
	SiOuterLength  int     `json:"si_outer_length" binding:"min=0,max=3"`
	SiOuterSleeve  int     `json:"si_outer_sleeve" binding:"min=0,max=3"`
	SiShoeSize     int     `json:"si_shoe_size"`
	Memo           string  `json:"memo"`
	Rating         float32 `json:"rating" binding:"min=0,max=5"`
	ItemIDs        []uint  `json:"item_ids" binding:"required,min=1"`
}

// UpdateCoordinateRequest represents coordinate update request
type UpdateCoordinateRequest struct {
	Season         *int     `json:"season" binding:"omitempty,min=1,max=5"`
	TPO            *int     `json:"tpo" binding:"omitempty,min=1,max=5"`
	SiTopLength    *int     `json:"si_top_length" binding:"omitempty,min=0,max=3"`
	SiTopSleeve    *int     `json:"si_top_sleeve" binding:"omitempty,min=0,max=5"`
	SiBottomLength *int     `json:"si_bottom_length" binding:"omitempty,min=0,max=6"`
	SiBottomType   *int     `json:"si_bottom_type" binding:"omitempty,min=0,max=2"`
	SiDressLength  *int     `json:"si_dress_length" binding:"omitempty,min=0,max=6"`
	SiDressSleeve  *int     `json:"si_dress_sleeve" binding:"omitempty,min=0,max=5"`
	SiOuterLength  *int     `json:"si_outer_length" binding:"omitempty,min=0,max=3"`
	SiOuterSleeve  *int     `json:"si_outer_sleeve" binding:"omitempty,min=0,max=3"`
	SiShoeSize     *int     `json:"si_shoe_size"`
	Memo           *string  `json:"memo"`
	Rating         *float32 `json:"rating" binding:"omitempty,min=0,max=5"`
	ItemIDs        []uint   `json:"item_ids"`
}

// CoordinateResponse represents coordinate data in responses
type CoordinateResponse struct {
	ID             uint           `json:"id"`
	UserID         uint           `json:"user_id"`
	Season         int            `json:"season"`
	TPO            int            `json:"tpo"`
	Picture        string         `json:"picture"`
	SiTopLength    int            `json:"si_top_length"`
	SiTopSleeve    int            `json:"si_top_sleeve"`
	SiBottomLength int            `json:"si_bottom_length"`
	SiBottomType   int            `json:"si_bottom_type"`
	SiDressLength  int            `json:"si_dress_length"`
	SiDressSleeve  int            `json:"si_dress_sleeve"`
	SiOuterLength  int            `json:"si_outer_length"`
	SiOuterSleeve  int            `json:"si_outer_sleeve"`
	SiShoeSize     int            `json:"si_shoe_size"`
	Memo           string         `json:"memo"`
	Rating         float32        `json:"rating"`
	Items          []ItemResponse `json:"items"`
	LikeCount      int64          `json:"like_count"`
	CommentCount   int64          `json:"comment_count"`
	IsLiked        bool           `json:"is_liked"`
	User           UserResponse   `json:"user"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// CoordinateListResponse represents paginated coordinate list response
type CoordinateListResponse struct {
	Coordinates []CoordinateResponse `json:"coordinates"`
	TotalCount  int64                `json:"total_count"`
	Page        int                  `json:"page"`
	PerPage     int                  `json:"per_page"`
}

// CoordinateFilterRequest represents coordinate search filters
type CoordinateFilterRequest struct {
	Season    *int     `form:"season" binding:"omitempty,min=1,max=5"`
	TPO       *int     `form:"tpo" binding:"omitempty,min=1,max=5"`
	MinRating *float32 `form:"min_rating" binding:"omitempty,min=0,max=5"`
	MaxRating *float32 `form:"max_rating" binding:"omitempty,min=0,max=5"`
	Page      int      `form:"page,default=1" binding:"min=1"`
	PerPage   int      `form:"per_page,default=20" binding:"min=1,max=100"`
}