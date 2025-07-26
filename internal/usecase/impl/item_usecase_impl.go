package impl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
	"github.com/House-lovers7/speadwear-go/pkg/config"
)

type itemUsecase struct {
	itemRepo repository.ItemRepository
	config   *config.Config
}

// NewItemUsecase creates a new item usecase
func NewItemUsecase(itemRepo repository.ItemRepository, config *config.Config) usecase.ItemUsecase {
	return &itemUsecase{
		itemRepo: itemRepo,
		config:   config,
	}
}

// CreateItem creates a new item
func (u *itemUsecase) CreateItem(ctx context.Context, userID uint, item *domain.Item, image *multipart.FileHeader) error {
	item.UserID = userID
	
	// Upload image if provided
	if image != nil {
		filename, err := u.uploadImage(image, "items")
		if err != nil {
			return err
		}
		item.Picture = filename
	}
	
	return u.itemRepo.Create(ctx, item)
}

// GetItem gets an item by ID
func (u *itemUsecase) GetItem(ctx context.Context, itemID uint) (*domain.Item, error) {
	item, err := u.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}
	return item, nil
}

// UpdateItem updates an item
func (u *itemUsecase) UpdateItem(ctx context.Context, userID uint, itemID uint, updates map[string]interface{}, image *multipart.FileHeader) error {
	item, err := u.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}
	
	// Check ownership
	if item.UserID != userID {
		return errors.New("unauthorized")
	}
	
	// Apply updates
	if superItem, ok := updates["super_item"].(string); ok {
		item.SuperItem = superItem
	}
	if season, ok := updates["season"].(int); ok {
		item.Season = season
	}
	if tpo, ok := updates["tpo"].(int); ok {
		item.TPO = tpo
	}
	if color, ok := updates["color"].(int); ok {
		item.Color = color
	}
	if content, ok := updates["content"].(string); ok {
		item.Content = content
	}
	if memo, ok := updates["memo"].(string); ok {
		item.Memo = memo
	}
	if rating, ok := updates["rating"].(float32); ok {
		item.Rating = rating
	}
	
	// Upload new image if provided
	if image != nil {
		// Delete old image if exists
		if item.Picture != "" {
			u.deleteImage(item.Picture)
		}
		
		filename, err := u.uploadImage(image, "items")
		if err != nil {
			return err
		}
		item.Picture = filename
	}
	
	return u.itemRepo.Update(ctx, item)
}

// DeleteItem deletes an item
func (u *itemUsecase) DeleteItem(ctx context.Context, userID uint, itemID uint) error {
	item, err := u.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}
	
	// Check ownership
	if item.UserID != userID {
		return errors.New("unauthorized")
	}
	
	// Delete image if exists
	if item.Picture != "" {
		u.deleteImage(item.Picture)
	}
	
	return u.itemRepo.Delete(ctx, itemID)
}

// GetUserItems gets items for a user
func (u *itemUsecase) GetUserItems(ctx context.Context, userID uint, limit, offset int) ([]*domain.Item, int64, error) {
	items, err := u.itemRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := u.itemRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	
	return items, count, nil
}

// SearchItems searches items with filters
func (u *itemUsecase) SearchItems(ctx context.Context, filters repository.ItemFilter) ([]*domain.Item, error) {
	return u.itemRepo.FindByFilters(ctx, filters)
}

// DeleteUserItems deletes multiple items for a user
func (u *itemUsecase) DeleteUserItems(ctx context.Context, userID uint, itemIDs []uint) error {
	for _, itemID := range itemIDs {
		item, err := u.itemRepo.FindByID(ctx, itemID)
		if err != nil {
			return err
		}
		if item == nil {
			continue
		}
		
		// Check ownership
		if item.UserID != userID {
			return errors.New("unauthorized")
		}
		
		// Delete image if exists
		if item.Picture != "" {
			u.deleteImage(item.Picture)
		}
		
		if err := u.itemRepo.Delete(ctx, itemID); err != nil {
			return err
		}
	}
	
	return nil
}

// GetUserItemStatistics gets item statistics for a user
func (u *itemUsecase) GetUserItemStatistics(ctx context.Context, userID uint) (map[string]interface{}, error) {
	items, err := u.itemRepo.FindByUserID(ctx, userID, 0, 0)
	if err != nil {
		return nil, err
	}
	
	// Calculate statistics
	stats := make(map[string]interface{})
	stats["total_count"] = len(items)
	
	// Count by category
	categoryCount := make(map[string]int)
	seasonCount := make(map[int]int)
	tpoCount := make(map[int]int)
	colorCount := make(map[int]int)
	
	var totalRating float32
	ratedCount := 0
	
	for _, item := range items {
		categoryCount[item.SuperItem]++
		seasonCount[item.Season]++
		tpoCount[item.TPO]++
		colorCount[item.Color]++
		
		if item.Rating > 0 {
			totalRating += item.Rating
			ratedCount++
		}
	}
	
	stats["category_count"] = categoryCount
	stats["season_count"] = seasonCount
	stats["tpo_count"] = tpoCount
	stats["color_count"] = colorCount
	
	if ratedCount > 0 {
		stats["average_rating"] = totalRating / float32(ratedCount)
	} else {
		stats["average_rating"] = 0
	}
	
	return stats, nil
}

// uploadImage uploads an image file
func (u *itemUsecase) uploadImage(file *multipart.FileHeader, folder string) (string, error) {
	// Check file size
	if file.Size > u.config.Upload.MaxFileSize {
		return "", errors.New("file size exceeds limit")
	}
	
	// Open file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s/%d_%d%s", folder, time.Now().Unix(), time.Now().Nanosecond(), ext)
	
	// Create upload directory if not exists
	uploadPath := filepath.Join(u.config.Upload.Path, folder)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}
	
	// Create destination file
	dstPath := filepath.Join(u.config.Upload.Path, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	// Copy file
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	
	return filename, nil
}

// deleteImage deletes an image file
func (u *itemUsecase) deleteImage(filename string) error {
	if filename == "" {
		return nil
	}
	
	path := filepath.Join(u.config.Upload.Path, filename)
	return os.Remove(path)
}