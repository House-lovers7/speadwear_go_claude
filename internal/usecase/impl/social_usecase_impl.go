package impl

import (
	"context"
	"errors"
	"fmt"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
	"github.com/House-lovers7/speadwear-go/pkg/config"
)

type socialUsecase struct {
	commentRepo      repository.CommentRepository
	relationshipRepo repository.RelationshipRepository
	blockRepo        repository.BlockRepository
	notificationRepo repository.NotificationRepository
	coordinateRepo   repository.CoordinateRepository
	userRepo         repository.UserRepository
	config           *config.Config
}

// NewSocialUsecase creates a new social usecase
func NewSocialUsecase(
	commentRepo repository.CommentRepository,
	relationshipRepo repository.RelationshipRepository,
	blockRepo repository.BlockRepository,
	notificationRepo repository.NotificationRepository,
	coordinateRepo repository.CoordinateRepository,
	userRepo repository.UserRepository,
	config *config.Config,
) usecase.SocialUsecase {
	return &socialUsecase{
		commentRepo:      commentRepo,
		relationshipRepo: relationshipRepo,
		blockRepo:        blockRepo,
		notificationRepo: notificationRepo,
		coordinateRepo:   coordinateRepo,
		userRepo:         userRepo,
		config:           config,
	}
}

// CreateComment creates a new comment
func (u *socialUsecase) CreateComment(ctx context.Context, userID uint, coordinateID uint, comment string) (*domain.Comment, error) {
	// Check if coordinate exists
	coordinate, err := u.coordinateRepo.FindByID(ctx, coordinateID)
	if err != nil {
		return nil, err
	}
	if coordinate == nil {
		return nil, errors.New("coordinate not found")
	}
	
	// Check if user is blocked by coordinate owner
	isBlocked, err := u.blockRepo.ExistsByBlockerAndBlocked(ctx, coordinate.UserID, userID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, errors.New("you are blocked by this user")
	}
	
	// Create comment
	newComment := &domain.Comment{
		UserID:       userID,
		CoordinateID: coordinateID,
		Comment:      comment,
	}
	
	if err := u.commentRepo.Create(ctx, newComment); err != nil {
		return nil, err
	}
	
	// Create notification if not commenting on own coordinate
	if coordinate.UserID != userID {
		notification := &domain.Notification{
			SenderID:     userID,
			ReceiverID:   coordinate.UserID,
			CoordinateID: &coordinateID,
			CommentID:    &newComment.ID,
			Action:       domain.NotificationActionComment,
		}
		if err := u.notificationRepo.Create(ctx, notification); err != nil {
			// Log error but don't fail the comment operation
			fmt.Printf("Failed to create notification: %v\n", err)
		}
	}
	
	return newComment, nil
}

// GetComments gets comments for a coordinate
func (u *socialUsecase) GetComments(ctx context.Context, coordinateID uint, limit, offset int) ([]*domain.Comment, int64, error) {
	comments, err := u.commentRepo.FindByCoordinateID(ctx, coordinateID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := u.commentRepo.CountByCoordinateID(ctx, coordinateID)
	if err != nil {
		return nil, 0, err
	}
	
	return comments, count, nil
}

// UpdateComment updates a comment
func (u *socialUsecase) UpdateComment(ctx context.Context, userID uint, commentID uint, comment string) error {
	existingComment, err := u.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}
	if existingComment == nil {
		return errors.New("comment not found")
	}
	
	// Check ownership
	if existingComment.UserID != userID {
		return errors.New("unauthorized")
	}
	
	existingComment.Comment = comment
	return u.commentRepo.Update(ctx, existingComment)
}

// DeleteComment deletes a comment
func (u *socialUsecase) DeleteComment(ctx context.Context, userID uint, commentID uint) error {
	comment, err := u.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("comment not found")
	}
	
	// Check ownership
	if comment.UserID != userID {
		return errors.New("unauthorized")
	}
	
	return u.commentRepo.Delete(ctx, commentID)
}

// FollowUser follows a user
func (u *socialUsecase) FollowUser(ctx context.Context, followerID uint, followedID uint) error {
	// Check if user exists
	user, err := u.userRepo.FindByID(ctx, followedID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	// Can't follow yourself
	if followerID == followedID {
		return errors.New("cannot follow yourself")
	}
	
	// Check if already following
	exists, err := u.relationshipRepo.ExistsByFollowerAndFollowed(ctx, followerID, followedID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already following")
	}
	
	// Check if blocked
	isBlocked, err := u.blockRepo.ExistsByBlockerAndBlocked(ctx, followedID, followerID)
	if err != nil {
		return err
	}
	if isBlocked {
		return errors.New("you are blocked by this user")
	}
	
	// Create relationship
	relationship := &domain.Relationship{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	
	if err := u.relationshipRepo.Create(ctx, relationship); err != nil {
		return err
	}
	
	// Create notification
	notification := &domain.Notification{
		SenderID:   followerID,
		ReceiverID: followedID,
		Action:     domain.NotificationActionFollow,
	}
	if err := u.notificationRepo.Create(ctx, notification); err != nil {
		// Log error but don't fail the follow operation
		fmt.Printf("Failed to create notification: %v\n", err)
	}
	
	return nil
}

// UnfollowUser unfollows a user
func (u *socialUsecase) UnfollowUser(ctx context.Context, followerID uint, followedID uint) error {
	// Find relationship
	relationship, err := u.relationshipRepo.FindByFollowerAndFollowed(ctx, followerID, followedID)
	if err != nil {
		return err
	}
	if relationship == nil {
		return errors.New("not following")
	}
	
	// Delete relationship
	return u.relationshipRepo.Delete(ctx, relationship.ID)
}

// GetFollowers gets followers of a user
func (u *socialUsecase) GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, int64, error) {
	followers, err := u.relationshipRepo.FindFollowers(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := u.relationshipRepo.CountFollowers(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	
	return followers, count, nil
}

// GetFollowing gets users that a user is following
func (u *socialUsecase) GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, int64, error) {
	following, err := u.relationshipRepo.FindFollowing(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := u.relationshipRepo.CountFollowing(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	
	return following, count, nil
}

// IsFollowing checks if user A is following user B
func (u *socialUsecase) IsFollowing(ctx context.Context, followerID uint, followedID uint) (bool, error) {
	return u.relationshipRepo.ExistsByFollowerAndFollowed(ctx, followerID, followedID)
}

// BlockUser blocks a user
func (u *socialUsecase) BlockUser(ctx context.Context, blockerID uint, blockedID uint) error {
	// Check if user exists
	user, err := u.userRepo.FindByID(ctx, blockedID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	// Can't block yourself
	if blockerID == blockedID {
		return errors.New("cannot block yourself")
	}
	
	// Check if already blocked
	exists, err := u.blockRepo.ExistsByBlockerAndBlocked(ctx, blockerID, blockedID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already blocked")
	}
	
	// Create block
	block := &domain.Block{
		BlockerID: blockerID,
		BlockedID: blockedID,
	}
	
	if err := u.blockRepo.Create(ctx, block); err != nil {
		return err
	}
	
	// Remove follow relationship if exists
	// TODO: Implement this when unfollow is properly implemented
	
	return nil
}

// UnblockUser unblocks a user
func (u *socialUsecase) UnblockUser(ctx context.Context, blockerID uint, blockedID uint) error {
	block, err := u.blockRepo.FindByBlockerAndBlocked(ctx, blockerID, blockedID)
	if err != nil {
		return err
	}
	if block == nil {
		return errors.New("not blocked")
	}
	
	// Delete block
	return u.blockRepo.Delete(ctx, block.ID)
}

// GetBlockedUsers gets blocked users
func (u *socialUsecase) GetBlockedUsers(ctx context.Context, userID uint) ([]*domain.User, error) {
	return u.blockRepo.FindBlockedUsers(ctx, userID)
}

// IsBlocked checks if user A has blocked user B
func (u *socialUsecase) IsBlocked(ctx context.Context, blockerID uint, blockedID uint) (bool, error) {
	return u.blockRepo.ExistsByBlockerAndBlocked(ctx, blockerID, blockedID)
}

// GetNotifications gets notifications for a user
func (u *socialUsecase) GetNotifications(ctx context.Context, userID uint, limit, offset int) ([]*domain.Notification, int64, error) {
	notifications, err := u.notificationRepo.FindByReceiverID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	// Get total count (not just unread)
	// TODO: Repository needs a method to count all notifications by receiver
	// For now, using unread count
	count, err := u.notificationRepo.CountUnreadByReceiverID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	
	return notifications, count, nil
}

// GetUnreadNotifications gets unread notifications
func (u *socialUsecase) GetUnreadNotifications(ctx context.Context, userID uint) ([]*domain.Notification, error) {
	return u.notificationRepo.FindUnreadByReceiverID(ctx, userID)
}

// MarkNotificationAsRead marks a notification as read
func (u *socialUsecase) MarkNotificationAsRead(ctx context.Context, userID uint, notificationID uint) error {
	notification, err := u.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}
	if notification == nil {
		return errors.New("notification not found")
	}
	
	// Check ownership
	if notification.ReceiverID != userID {
		return errors.New("unauthorized")
	}
	
	return u.notificationRepo.MarkAsRead(ctx, notificationID)
}

// MarkAllNotificationsAsRead marks all notifications as read
func (u *socialUsecase) MarkAllNotificationsAsRead(ctx context.Context, userID uint) error {
	return u.notificationRepo.MarkAllAsReadByReceiver(ctx, userID)
}

// GetUnreadNotificationCount gets unread notification count
func (u *socialUsecase) GetUnreadNotificationCount(ctx context.Context, userID uint) (int64, error) {
	return u.notificationRepo.CountUnreadByReceiverID(ctx, userID)
}

// CreateNotification creates a notification (internal use)
func (u *socialUsecase) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	return u.notificationRepo.Create(ctx, notification)
}