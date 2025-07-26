package user

import (
	"time"

	"github.com/House-lovers7/speadwear-go/internal/domain"
)

type User struct {
	domain.BaseModel
	Name               string    `gorm:"type:varchar(255);not null" json:"name"`
	Email              string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Picture            string    `gorm:"type:varchar(255)" json:"picture"`
	Admin              bool      `gorm:"default:false" json:"admin"`
	PasswordDigest     string    `gorm:"type:varchar(255)" json:"-"`
	RememberDigest     string    `gorm:"type:varchar(255)" json:"-"`
	ActivationDigest   string    `gorm:"type:varchar(255)" json:"-"`
	Activated          bool      `gorm:"default:false" json:"activated"`
	ActivatedAt        *time.Time `json:"activated_at,omitempty"`
	ResetDigest        string    `gorm:"type:varchar(255)" json:"-"`
	ResetSentAt        *time.Time `json:"reset_sent_at,omitempty"`
	
	// Relations
	Items         []Item         `gorm:"foreignKey:UserID" json:"items,omitempty"`
	Coordinates   []Coordinate   `gorm:"foreignKey:UserID" json:"coordinates,omitempty"`
	Comments      []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Followers     []Relationship `gorm:"foreignKey:FollowedID" json:"followers,omitempty"`
	Following     []Relationship `gorm:"foreignKey:FollowerID" json:"following,omitempty"`
	BlocksAsBlocker []Block      `gorm:"foreignKey:BlockerID" json:"blocks_as_blocker,omitempty"`
	BlocksAsBlocked []Block      `gorm:"foreignKey:BlockedID" json:"blocks_as_blocked,omitempty"`
	SentNotifications []Notification `gorm:"foreignKey:SenderID" json:"sent_notifications,omitempty"`
	ReceivedNotifications []Notification `gorm:"foreignKey:ReceiverID" json:"received_notifications,omitempty"`
}

type Relationship struct {
	domain.BaseModel
	FollowerID uint `gorm:"not null;index" json:"follower_id"`
	FollowedID uint `gorm:"not null;index" json:"followed_id"`
	Follower   User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Followed   User `gorm:"foreignKey:FollowedID" json:"followed,omitempty"`
}

type Block struct {
	domain.BaseModel
	BlockerID uint `gorm:"not null;index" json:"blocker_id"`
	BlockedID uint `gorm:"not null;index" json:"blocked_id"`
	Blocker   User `gorm:"foreignKey:BlockerID" json:"blocker,omitempty"`
	Blocked   User `gorm:"foreignKey:BlockedID" json:"blocked,omitempty"`
}

type Notification struct {
	domain.BaseModel
	SenderID        uint         `gorm:"not null;index" json:"sender_id"`
	ReceiverID      uint         `gorm:"not null;index" json:"receiver_id"`
	CoordinateID    *uint        `json:"coordinate_id,omitempty"`
	CommentID       *uint        `json:"comment_id,omitempty"`
	LikeCoordinateID *uint       `json:"like_coordinate_id,omitempty"`
	Action          string       `gorm:"type:varchar(50);not null" json:"action"`
	Checked         bool         `gorm:"default:false" json:"checked"`
	Sender          User         `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver        User         `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Coordinate      *Coordinate  `gorm:"foreignKey:CoordinateID" json:"coordinate,omitempty"`
	Comment         *Comment     `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
	LikeCoordinate  *LikeCoordinate `gorm:"foreignKey:LikeCoordinateID" json:"like_coordinate,omitempty"`
}

// 以下の型は他のパッケージで定義されるが、ここでは参照のために宣言
type Item struct{}
type Coordinate struct{}
type Comment struct{}
type LikeCoordinate struct{}