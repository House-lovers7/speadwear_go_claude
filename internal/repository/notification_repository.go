package repository

import (
	"context"
	"errors"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

// Create creates a new notification
func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// FindByID finds a notification by ID
func (r *notificationRepository) FindByID(ctx context.Context, id uint) (*domain.Notification, error) {
	var notification domain.Notification
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("Receiver").
		Preload("Coordinate").
		Preload("Comment").
		Preload("LikeCoordinate").
		First(&notification, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

// Update updates a notification
func (r *notificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

// Delete deletes a notification
func (r *notificationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Notification{}, id).Error
}

// FindByReceiverID finds notifications by receiver ID with pagination
func (r *notificationRepository) FindByReceiverID(ctx context.Context, receiverID uint, limit, offset int) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	query := r.db.WithContext(ctx).
		Where("receiver_id = ?", receiverID).
		Preload("Sender").
		Preload("Coordinate").
		Preload("Comment").
		Preload("LikeCoordinate")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// FindUnreadByReceiverID finds unread notifications by receiver ID
func (r *notificationRepository) FindUnreadByReceiverID(ctx context.Context, receiverID uint) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	err := r.db.WithContext(ctx).
		Where("receiver_id = ? AND checked = ?", receiverID, false).
		Preload("Sender").
		Preload("Coordinate").
		Preload("Comment").
		Preload("LikeCoordinate").
		Order("created_at DESC").
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// MarkAsRead marks a notification as read
func (r *notificationRepository) MarkAsRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).Where("id = ?", id).Update("checked", true).Error
}

// MarkAllAsReadByReceiver marks all notifications as read for a receiver
func (r *notificationRepository) MarkAllAsReadByReceiver(ctx context.Context, receiverID uint) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("receiver_id = ? AND checked = ?", receiverID, false).
		Update("checked", true).Error
}

// CountUnreadByReceiverID counts unread notifications by receiver ID
func (r *notificationRepository) CountUnreadByReceiverID(ctx context.Context, receiverID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("receiver_id = ? AND checked = ?", receiverID, false).
		Count(&count).Error
	return count, err
}