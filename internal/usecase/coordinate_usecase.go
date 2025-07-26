package usecase

import (
	"context"
	"mime/multipart"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/repository"
)

// CoordinateUsecase defines coordinate-related business logic
type CoordinateUsecase interface {
	// CRUD operations
	CreateCoordinate(ctx context.Context, userID uint, coordinate *domain.Coordinate, itemIDs []uint, image *multipart.FileHeader) error
	GetCoordinate(ctx context.Context, coordinateID uint) (*domain.Coordinate, error)
	GetCoordinateWithDetails(ctx context.Context, coordinateID uint) (*domain.Coordinate, error)
	UpdateCoordinate(ctx context.Context, userID uint, coordinateID uint, updates map[string]interface{}, itemIDs []uint, image *multipart.FileHeader) error
	DeleteCoordinate(ctx context.Context, userID uint, coordinateID uint) error
	
	// Listing and searching
	GetUserCoordinates(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, int64, error)
	SearchCoordinates(ctx context.Context, filters repository.CoordinateFilter) ([]*domain.Coordinate, error)
	GetTimelineCoordinates(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, error)
	
	// Social features
	LikeCoordinate(ctx context.Context, userID uint, coordinateID uint) error
	UnlikeCoordinate(ctx context.Context, userID uint, coordinateID uint) error
	GetCoordinateLikes(ctx context.Context, coordinateID uint) ([]*domain.LikeCoordinate, error)
	IsLikedByUser(ctx context.Context, userID uint, coordinateID uint) (bool, error)
	
	// Statistics
	GetUserCoordinateStatistics(ctx context.Context, userID uint) (map[string]interface{}, error)
}