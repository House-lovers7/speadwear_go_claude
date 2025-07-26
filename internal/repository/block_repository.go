package repository

import (
	"context"
	"errors"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/gorm"
)

type blockRepository struct {
	db *gorm.DB
}

// NewBlockRepository creates a new block repository
func NewBlockRepository(db *gorm.DB) BlockRepository {
	return &blockRepository{db: db}
}

// Create creates a new block
func (r *blockRepository) Create(ctx context.Context, block *domain.Block) error {
	return r.db.WithContext(ctx).Create(block).Error
}

// FindByID finds a block by ID
func (r *blockRepository) FindByID(ctx context.Context, id uint) (*domain.Block, error) {
	var block domain.Block
	err := r.db.WithContext(ctx).First(&block, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

// Update updates a block
func (r *blockRepository) Update(ctx context.Context, block *domain.Block) error {
	return r.db.WithContext(ctx).Save(block).Error
}

// Delete deletes a block
func (r *blockRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Block{}, id).Error
}

// FindBlockedUsers finds users blocked by a user
func (r *blockRepository) FindBlockedUsers(ctx context.Context, blockerID uint) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("INNER JOIN blocks ON users.id = blocks.blocked_id").
		Where("blocks.blocker_id = ?", blockerID).
		Order("blocks.created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindByBlockerAndBlocked finds a block by blocker and blocked IDs
func (r *blockRepository) FindByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint) (*domain.Block, error) {
	var block domain.Block
	err := r.db.WithContext(ctx).
		Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		First(&block).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

// ExistsByBlockerAndBlocked checks if a block relationship exists
func (r *blockRepository) ExistsByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Block{}).
		Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}