package usecase

import (
	"context"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
)

// SocialUsecase defines social features business logic
type SocialUsecase interface {
	// Comment features
	CreateComment(ctx context.Context, userID uint, coordinateID uint, comment string) (*domain.Comment, error)
	GetComments(ctx context.Context, coordinateID uint, limit, offset int) ([]*domain.Comment, int64, error)
	UpdateComment(ctx context.Context, userID uint, commentID uint, comment string) error
	DeleteComment(ctx context.Context, userID uint, commentID uint) error
	
	// Follow features
	FollowUser(ctx context.Context, followerID uint, followedID uint) error
	UnfollowUser(ctx context.Context, followerID uint, followedID uint) error
	GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, int64, error)
	GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, int64, error)
	IsFollowing(ctx context.Context, followerID uint, followedID uint) (bool, error)
	
	// Block features
	BlockUser(ctx context.Context, blockerID uint, blockedID uint) error
	UnblockUser(ctx context.Context, blockerID uint, blockedID uint) error
	GetBlockedUsers(ctx context.Context, userID uint) ([]*domain.User, error)
	IsBlocked(ctx context.Context, blockerID uint, blockedID uint) (bool, error)
	
	// Notification features
	GetNotifications(ctx context.Context, userID uint, limit, offset int) ([]*domain.Notification, int64, error)
	GetUnreadNotifications(ctx context.Context, userID uint) ([]*domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, userID uint, notificationID uint) error
	MarkAllNotificationsAsRead(ctx context.Context, userID uint) error
	GetUnreadNotificationCount(ctx context.Context, userID uint) (int64, error)
	
	// Internal notification creation (called by other usecases)
	CreateNotification(ctx context.Context, notification *domain.Notification) error
}