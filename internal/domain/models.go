package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	BaseModel
	Name               string         `gorm:"type:varchar(255);not null" json:"name"`
	Email              string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Picture            string         `gorm:"type:varchar(255)" json:"picture"`
	Admin              bool           `gorm:"default:false" json:"admin"`
	PasswordDigest     string         `gorm:"type:varchar(255)" json:"-"`
	RememberDigest     string         `gorm:"type:varchar(255)" json:"-"`
	ActivationDigest   string         `gorm:"type:varchar(255)" json:"-"`
	Activated          bool           `gorm:"default:false" json:"activated"`
	ActivatedAt        *time.Time     `json:"activated_at,omitempty"`
	ResetDigest        string         `gorm:"type:varchar(255)" json:"-"`
	ResetSentAt        *time.Time     `json:"reset_sent_at,omitempty"`
	
	// Relations
	Items              []Item         `gorm:"foreignKey:UserID" json:"items,omitempty"`
	Coordinates        []Coordinate   `gorm:"foreignKey:UserID" json:"coordinates,omitempty"`
	Comments           []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Followers          []Relationship `gorm:"foreignKey:FollowedID" json:"followers,omitempty"`
	Following          []Relationship `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
	BlocksAsBlocker    []Block        `gorm:"foreignKey:BlockerID" json:"blocks_as_blocker,omitempty"`
	BlocksAsBlocked    []Block        `gorm:"foreignKey:BlockedID" json:"blocks_as_blocked,omitempty"`
	SentNotifications     []Notification `gorm:"foreignKey:SenderID" json:"sent_notifications,omitempty"`
	ReceivedNotifications []Notification `gorm:"foreignKey:ReceiverID" json:"received_notifications,omitempty"`
}

// Item represents a clothing item
type Item struct {
	BaseModel
	UserID       uint        `gorm:"not null;index" json:"user_id"`
	CoordinateID *uint       `json:"coordinate_id,omitempty"`
	SuperItem    string      `gorm:"type:varchar(100)" json:"super_item"`
	Season       int         `json:"season"`
	TPO          int         `json:"tpo"`
	Color        int         `json:"color"`
	Content      string      `gorm:"type:text" json:"content"`
	Memo         string      `gorm:"type:text" json:"memo"`
	Picture      string      `gorm:"type:varchar(255)" json:"picture"`
	Rating       float32     `json:"rating"`
	
	// Relations
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coordinate  *Coordinate `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
}

// Coordinate represents an outfit coordination
type Coordinate struct {
	BaseModel
	UserID           uint             `gorm:"not null;index" json:"user_id"`
	Season           int              `json:"season"`
	TPO              int              `json:"tpo"`
	Picture          string           `gorm:"type:varchar(255)" json:"picture"`
	SiTopLength      int              `json:"si_top_length"`
	SiTopSleeve      int              `json:"si_top_sleeve"`
	SiBottomLength   int              `json:"si_bottom_length"`
	SiBottomType     int              `json:"si_bottom_type"`
	SiDressLength    int              `json:"si_dress_length"`
	SiDressSleeve    int              `json:"si_dress_sleeve"`
	SiOuterLength    int              `json:"si_outer_length"`
	SiOuterSleeve    int              `json:"si_outer_sleeve"`
	SiShoeSize       int              `json:"si_shoe_size"`
	Memo             string           `gorm:"type:text" json:"memo"`
	Rating           float32          `json:"rating"`
	
	// Relations
	User            User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items           []Item           `gorm:"foreignKey:CoordinateID" json:"items,omitempty"`
	Comments        []Comment        `gorm:"foreignKey:CoordinateID" json:"comments,omitempty"`
	LikeCoordinates []LikeCoordinate `gorm:"foreignKey:CoordinateID" json:"like_coordinates,omitempty"`
}

// Comment represents a comment on a coordinate
type Comment struct {
	BaseModel
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	CoordinateID uint       `gorm:"not null;index" json:"coordinate_id"`
	Comment      string     `gorm:"type:text;not null" json:"comment"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coordinate   Coordinate `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
}

// LikeCoordinate represents a like on a coordinate
type LikeCoordinate struct {
	BaseModel
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	CoordinateID uint       `gorm:"not null;index" json:"coordinate_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Coordinate   Coordinate `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
}

// Relationship represents a follow relationship between users
type Relationship struct {
	BaseModel
	FollowerID uint `gorm:"not null;index" json:"follower_id"`
	FollowedID uint `gorm:"not null;index" json:"followed_id"`
	Follower   User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Followed   User `gorm:"foreignKey:FollowedID" json:"followed,omitempty"`
}

// Block represents a block relationship between users
type Block struct {
	BaseModel
	BlockerID uint `gorm:"not null;index" json:"blocker_id"`
	BlockedID uint `gorm:"not null;index" json:"blocked_id"`
	Blocker   User `gorm:"foreignKey:BlockerID" json:"blocker,omitempty"`
	Blocked   User `gorm:"foreignKey:BlockedID" json:"blocked,omitempty"`
}

// Notification represents a notification
type Notification struct {
	BaseModel
	SenderID         uint            `gorm:"not null;index" json:"sender_id"`
	ReceiverID       uint            `gorm:"not null;index" json:"receiver_id"`
	CoordinateID     *uint           `json:"coordinate_id,omitempty"`
	CommentID        *uint           `json:"comment_id,omitempty"`
	LikeCoordinateID *uint           `json:"like_coordinate_id,omitempty"`
	Action           string          `gorm:"type:varchar(50);not null" json:"action"`
	Checked          bool            `gorm:"default:false" json:"checked"`
	Sender           User            `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver         User            `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Coordinate       *Coordinate     `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
	Comment          *Comment        `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
	LikeCoordinate   *LikeCoordinate `gorm:"foreignKey:LikeCoordinateID" json:"like_coordinate,omitempty"`
}

// GetAllModels returns all model structs for migration
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Item{},
		&Coordinate{},
		&Comment{},
		&LikeCoordinate{},
		&Relationship{},
		&Block{},
		&Notification{},
	}
}