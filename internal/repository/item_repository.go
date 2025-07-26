package repository

import (
	"context"
	"errors"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/gorm"
)

type itemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

// Create creates a new item
func (r *itemRepository) Create(ctx context.Context, item *domain.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// FindByID finds an item by ID
func (r *itemRepository) FindByID(ctx context.Context, id uint) (*domain.Item, error) {
	var item domain.Item
	err := r.db.WithContext(ctx).Preload("User").First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Update updates an item
func (r *itemRepository) Update(ctx context.Context, item *domain.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete deletes an item
func (r *itemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Item{}, id).Error
}

// FindByUserID finds items by user ID with pagination
func (r *itemRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*domain.Item, error) {
	var items []*domain.Item
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Order("created_at DESC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// FindByFilters finds items by filters
func (r *itemRepository) FindByFilters(ctx context.Context, filters ItemFilter) ([]*domain.Item, error) {
	var items []*domain.Item
	query := r.db.WithContext(ctx).Preload("User")
	
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
	if filters.Color != nil {
		query = query.Where("color = ?", *filters.Color)
	}
	if filters.SuperItem != nil {
		query = query.Where("super_item = ?", *filters.SuperItem)
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
	err := query.Order("created_at DESC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// CountByUserID counts items by user ID
func (r *itemRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Item{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}