package repository

import (
	"context"
	"errors"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/gorm"
)

type relationshipRepository struct {
	db *gorm.DB
}

// NewRelationshipRepository creates a new relationship repository
func NewRelationshipRepository(db *gorm.DB) RelationshipRepository {
	return &relationshipRepository{db: db}
}

// Create creates a new relationship
func (r *relationshipRepository) Create(ctx context.Context, relationship *domain.Relationship) error {
	return r.db.WithContext(ctx).Create(relationship).Error
}

// FindByID finds a relationship by ID
func (r *relationshipRepository) FindByID(ctx context.Context, id uint) (*domain.Relationship, error) {
	var relationship domain.Relationship
	err := r.db.WithContext(ctx).First(&relationship, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &relationship, nil
}

// Update updates a relationship
func (r *relationshipRepository) Update(ctx context.Context, relationship *domain.Relationship) error {
	return r.db.WithContext(ctx).Save(relationship).Error
}

// Delete deletes a relationship
func (r *relationshipRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Relationship{}, id).Error
}

// FindFollowers finds followers of a user with pagination
func (r *relationshipRepository) FindFollowers(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	query := r.db.WithContext(ctx).
		Table("users").
		Joins("INNER JOIN relationships ON users.id = relationships.follower_id").
		Where("relationships.followed_id = ?", userID)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("relationships.created_at DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindFollowing finds users that a user is following with pagination
func (r *relationshipRepository) FindFollowing(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	query := r.db.WithContext(ctx).
		Table("users").
		Joins("INNER JOIN relationships ON users.id = relationships.followed_id").
		Where("relationships.follower_id = ?", userID)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("relationships.created_at DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindByFollowerAndFollowed finds a relationship by follower and followed IDs
func (r *relationshipRepository) FindByFollowerAndFollowed(ctx context.Context, followerID, followedID uint) (*domain.Relationship, error) {
	var relationship domain.Relationship
	err := r.db.WithContext(ctx).
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		First(&relationship).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &relationship, nil
}

// ExistsByFollowerAndFollowed checks if a relationship exists
func (r *relationshipRepository) ExistsByFollowerAndFollowed(ctx context.Context, followerID, followedID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Relationship{}).
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountFollowers counts followers of a user
func (r *relationshipRepository) CountFollowers(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Relationship{}).Where("followed_id = ?", userID).Count(&count).Error
	return count, err
}

// CountFollowing counts users that a user is following
func (r *relationshipRepository) CountFollowing(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Relationship{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}