package repository

import (
	"context"
	"github.com/House-lovers7/speadwear-go/internal/domain"
)

// BaseRepository defines common methods for all repositories
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
}

// UserRepository defines methods for user data access
type UserRepository interface {
	BaseRepository[domain.User]
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error)
	Count(ctx context.Context) (int64, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// ItemRepository defines methods for item data access
type ItemRepository interface {
	BaseRepository[domain.Item]
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*domain.Item, error)
	FindByFilters(ctx context.Context, filters ItemFilter) ([]*domain.Item, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

// CoordinateRepository defines methods for coordinate data access
type CoordinateRepository interface {
	BaseRepository[domain.Coordinate]
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, error)
	FindByFilters(ctx context.Context, filters CoordinateFilter) ([]*domain.Coordinate, error)
	FindWithItems(ctx context.Context, id uint) (*domain.Coordinate, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

// CommentRepository defines methods for comment data access
type CommentRepository interface {
	BaseRepository[domain.Comment]
	FindByCoordinateID(ctx context.Context, coordinateID uint, limit, offset int) ([]*domain.Comment, error)
	CountByCoordinateID(ctx context.Context, coordinateID uint) (int64, error)
}

// LikeCoordinateRepository defines methods for like data access
type LikeCoordinateRepository interface {
	BaseRepository[domain.LikeCoordinate]
	FindByCoordinateID(ctx context.Context, coordinateID uint) ([]*domain.LikeCoordinate, error)
	FindByUserAndCoordinate(ctx context.Context, userID, coordinateID uint) (*domain.LikeCoordinate, error)
	ExistsByUserAndCoordinate(ctx context.Context, userID, coordinateID uint) (bool, error)
	CountByCoordinateID(ctx context.Context, coordinateID uint) (int64, error)
}

// RelationshipRepository defines methods for follow relationship data access
type RelationshipRepository interface {
	BaseRepository[domain.Relationship]
	FindFollowers(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, error)
	FindFollowing(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, error)
	FindByFollowerAndFollowed(ctx context.Context, followerID, followedID uint) (*domain.Relationship, error)
	ExistsByFollowerAndFollowed(ctx context.Context, followerID, followedID uint) (bool, error)
	CountFollowers(ctx context.Context, userID uint) (int64, error)
	CountFollowing(ctx context.Context, userID uint) (int64, error)
}

// BlockRepository defines methods for block data access
type BlockRepository interface {
	BaseRepository[domain.Block]
	FindBlockedUsers(ctx context.Context, blockerID uint) ([]*domain.User, error)
	FindByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint) (*domain.Block, error)
	ExistsByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint) (bool, error)
}

// NotificationRepository defines methods for notification data access
type NotificationRepository interface {
	BaseRepository[domain.Notification]
	FindByReceiverID(ctx context.Context, receiverID uint, limit, offset int) ([]*domain.Notification, error)
	FindUnreadByReceiverID(ctx context.Context, receiverID uint) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsReadByReceiver(ctx context.Context, receiverID uint) error
	CountUnreadByReceiverID(ctx context.Context, receiverID uint) (int64, error)
}

// Filter structures
type ItemFilter struct {
	UserID   *uint
	Season   *int
	TPO      *int
	Color    *int
	SuperItem *string
	MinRating *float32
	MaxRating *float32
	Limit    int
	Offset   int
}

type CoordinateFilter struct {
	UserID    *uint
	Season    *int
	TPO       *int
	MinRating *float32
	MaxRating *float32
	Limit     int
	Offset    int
}