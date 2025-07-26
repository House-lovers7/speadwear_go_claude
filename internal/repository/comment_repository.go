package repository

import (
	"context"
	"errors"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create creates a new comment
func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// FindByID finds a comment by ID
func (r *commentRepository) FindByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.WithContext(ctx).Preload("User").Preload("Coordinate").First(&comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// Update updates a comment
func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// Delete deletes a comment
func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Comment{}, id).Error
}

// FindByCoordinateID finds comments by coordinate ID with pagination
func (r *commentRepository) FindByCoordinateID(ctx context.Context, coordinateID uint, limit, offset int) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	query := r.db.WithContext(ctx).Where("coordinate_id = ?", coordinateID).Preload("User")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CountByCoordinateID counts comments by coordinate ID
func (r *commentRepository) CountByCoordinateID(ctx context.Context, coordinateID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Comment{}).Where("coordinate_id = ?", coordinateID).Count(&count).Error
	return count, err
}