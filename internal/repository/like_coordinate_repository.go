package repository

import (
	"context"
	"errors"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/gorm"
)

type likeCoordinateRepository struct {
	db *gorm.DB
}

// NewLikeCoordinateRepository creates a new like coordinate repository
func NewLikeCoordinateRepository(db *gorm.DB) LikeCoordinateRepository {
	return &likeCoordinateRepository{db: db}
}

// Create creates a new like
func (r *likeCoordinateRepository) Create(ctx context.Context, like *domain.LikeCoordinate) error {
	return r.db.WithContext(ctx).Create(like).Error
}

// FindByID finds a like by ID
func (r *likeCoordinateRepository) FindByID(ctx context.Context, id uint) (*domain.LikeCoordinate, error) {
	var like domain.LikeCoordinate
	err := r.db.WithContext(ctx).First(&like, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

// Update updates a like
func (r *likeCoordinateRepository) Update(ctx context.Context, like *domain.LikeCoordinate) error {
	return r.db.WithContext(ctx).Save(like).Error
}

// Delete deletes a like
func (r *likeCoordinateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.LikeCoordinate{}, id).Error
}

// FindByCoordinateID finds likes by coordinate ID
func (r *likeCoordinateRepository) FindByCoordinateID(ctx context.Context, coordinateID uint) ([]*domain.LikeCoordinate, error) {
	var likes []*domain.LikeCoordinate
	err := r.db.WithContext(ctx).Where("coordinate_id = ?", coordinateID).Preload("User").Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

// FindByUserAndCoordinate finds a like by user ID and coordinate ID
func (r *likeCoordinateRepository) FindByUserAndCoordinate(ctx context.Context, userID, coordinateID uint) (*domain.LikeCoordinate, error) {
	var like domain.LikeCoordinate
	err := r.db.WithContext(ctx).Where("user_id = ? AND coordinate_id = ?", userID, coordinateID).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

// ExistsByUserAndCoordinate checks if a like exists by user ID and coordinate ID
func (r *likeCoordinateRepository) ExistsByUserAndCoordinate(ctx context.Context, userID, coordinateID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.LikeCoordinate{}).
		Where("user_id = ? AND coordinate_id = ?", userID, coordinateID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByCoordinateID counts likes by coordinate ID
func (r *likeCoordinateRepository) CountByCoordinateID(ctx context.Context, coordinateID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.LikeCoordinate{}).Where("coordinate_id = ?", coordinateID).Count(&count).Error
	return count, err
}