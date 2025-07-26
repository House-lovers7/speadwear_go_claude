package usecase

import (
	"context"
	"mime/multipart"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/repository"
)

// ItemUsecase defines item-related business logic
type ItemUsecase interface {
	// CRUD operations
	CreateItem(ctx context.Context, userID uint, item *domain.Item, image *multipart.FileHeader) error
	GetItem(ctx context.Context, itemID uint) (*domain.Item, error)
	UpdateItem(ctx context.Context, userID uint, itemID uint, updates map[string]interface{}, image *multipart.FileHeader) error
	DeleteItem(ctx context.Context, userID uint, itemID uint) error
	
	// Listing and searching
	GetUserItems(ctx context.Context, userID uint, limit, offset int) ([]*domain.Item, int64, error)
	SearchItems(ctx context.Context, filters repository.ItemFilter) ([]*domain.Item, error)
	
	// Batch operations
	DeleteUserItems(ctx context.Context, userID uint, itemIDs []uint) error
	
	// Statistics
	GetUserItemStatistics(ctx context.Context, userID uint) (map[string]interface{}, error)
}