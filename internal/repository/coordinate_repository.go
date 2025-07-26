package repository

import (
	"context"
	"errors"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/gorm"
)

type coordinateRepository struct {
	db *gorm.DB
}

// NewCoordinateRepository creates a new coordinate repository
func NewCoordinateRepository(db *gorm.DB) CoordinateRepository {
	return &coordinateRepository{db: db}
}

// Create creates a new coordinate
func (r *coordinateRepository) Create(ctx context.Context, coordinate *domain.Coordinate) error {
	return r.db.WithContext(ctx).Create(coordinate).Error
}

// FindByID finds a coordinate by ID
func (r *coordinateRepository) FindByID(ctx context.Context, id uint) (*domain.Coordinate, error) {
	var coordinate domain.Coordinate
	err := r.db.WithContext(ctx).Preload("User").First(&coordinate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &coordinate, nil
}

// Update updates a coordinate
func (r *coordinateRepository) Update(ctx context.Context, coordinate *domain.Coordinate) error {
	return r.db.WithContext(ctx).Save(coordinate).Error
}

// Delete deletes a coordinate
func (r *coordinateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Coordinate{}, id).Error
}

// FindByUserID finds coordinates by user ID with pagination
func (r *coordinateRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, error) {
	var coordinates []*domain.Coordinate
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User").Preload("Items")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("created_at DESC").Find(&coordinates).Error
	if err != nil {
		return nil, err
	}
	return coordinates, nil
}

// FindByFilters finds coordinates by filters
func (r *coordinateRepository) FindByFilters(ctx context.Context, filters CoordinateFilter) ([]*domain.Coordinate, error) {
	var coordinates []*domain.Coordinate
	query := r.db.WithContext(ctx).Preload("User").Preload("Items")
	
	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Season != nil {
		query = query.Where("season = ?", *filters.Season)
	}
	if filters.TPO != nil {
		query = query.Where("tpo = ?", *filters.TPO)
	}
	if filters.MinRating != nil {
		query = query.Where("rating >= ?", *filters.MinRating)
	}
	if filters.MaxRating != nil {
		query = query.Where("rating <= ?", *filters.MaxRating)
	}
	
	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	
	// Order by created_at descending
	err := query.Order("created_at DESC").Find(&coordinates).Error
	if err != nil {
		return nil, err
	}
	return coordinates, nil
}

// FindWithItems finds a coordinate with its items
func (r *coordinateRepository) FindWithItems(ctx context.Context, id uint) (*domain.Coordinate, error) {
	var coordinate domain.Coordinate
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.User").
		First(&coordinate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &coordinate, nil
}

// CountByUserID counts coordinates by user ID
func (r *coordinateRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Coordinate{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}