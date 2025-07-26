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
	"gorm.io/gorm"
)

type coordinateUsecase struct {
	coordinateRepo       repository.CoordinateRepository
	itemRepo            repository.ItemRepository
	likeCoordinateRepo  repository.LikeCoordinateRepository
	relationshipRepo    repository.RelationshipRepository
	blockRepo           repository.BlockRepository
	notificationRepo    repository.NotificationRepository
	config              *config.Config
	db                  *gorm.DB
}

// NewCoordinateUsecase creates a new coordinate usecase
func NewCoordinateUsecase(
	coordinateRepo repository.CoordinateRepository,
	itemRepo repository.ItemRepository,
	likeCoordinateRepo repository.LikeCoordinateRepository,
	relationshipRepo repository.RelationshipRepository,
	blockRepo repository.BlockRepository,
	notificationRepo repository.NotificationRepository,
	config *config.Config,
	db *gorm.DB,
) usecase.CoordinateUsecase {
	return &coordinateUsecase{
		coordinateRepo:      coordinateRepo,
		itemRepo:           itemRepo,
		likeCoordinateRepo: likeCoordinateRepo,
		relationshipRepo:   relationshipRepo,
		blockRepo:          blockRepo,
		notificationRepo:   notificationRepo,
		config:             config,
		db:                 db,
	}
}

// CreateCoordinate creates a new coordinate
func (u *coordinateUsecase) CreateCoordinate(ctx context.Context, userID uint, coordinate *domain.Coordinate, itemIDs []uint, image *multipart.FileHeader) error {
	coordinate.UserID = userID
	
	// Verify all items belong to the user
	for _, itemID := range itemIDs {
		item, err := u.itemRepo.FindByID(ctx, itemID)
		if err != nil {
			return err
		}
		if item == nil {
			return errors.New("item not found")
		}
		if item.UserID != userID {
			return errors.New("unauthorized: item does not belong to user")
		}
	}
	
	// Upload image if provided
	if image != nil {
		filename, err := u.uploadImage(image, "coordinates")
		if err != nil {
			return err
		}
		coordinate.Picture = filename
	}
	
	// Use transaction to create coordinate and update items
	return u.db.Transaction(func(tx *gorm.DB) error {
		// Create coordinate
		if err := tx.Create(coordinate).Error; err != nil {
			return err
		}
		
		// Update items to link to coordinate
		for _, itemID := range itemIDs {
			if err := tx.Model(&domain.Item{}).Where("id = ?", itemID).Update("coordinate_id", coordinate.ID).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// GetCoordinate gets a coordinate by ID
func (u *coordinateUsecase) GetCoordinate(ctx context.Context, coordinateID uint) (*domain.Coordinate, error) {
	coordinate, err := u.coordinateRepo.FindByID(ctx, coordinateID)
	if err != nil {
		return nil, err
	}
	if coordinate == nil {
		return nil, errors.New("coordinate not found")
	}
	return coordinate, nil
}

// GetCoordinateWithDetails gets a coordinate with all related data
func (u *coordinateUsecase) GetCoordinateWithDetails(ctx context.Context, coordinateID uint) (*domain.Coordinate, error) {
	coordinate, err := u.coordinateRepo.FindWithItems(ctx, coordinateID)
	if err != nil {
		return nil, err
	}
	if coordinate == nil {
		return nil, errors.New("coordinate not found")
	}
	return coordinate, nil
}

// UpdateCoordinate updates a coordinate
func (u *coordinateUsecase) UpdateCoordinate(ctx context.Context, userID uint, coordinateID uint, updates map[string]interface{}, itemIDs []uint, image *multipart.FileHeader) error {
	coordinate, err := u.coordinateRepo.FindByID(ctx, coordinateID)
	if err != nil {
		return err
	}
	if coordinate == nil {
		return errors.New("coordinate not found")
	}
	
	// Check ownership
	if coordinate.UserID != userID {
		return errors.New("unauthorized")
	}
	
	// Apply updates
	if season, ok := updates["season"].(int); ok {
		coordinate.Season = season
	}
	if tpo, ok := updates["tpo"].(int); ok {
		coordinate.TPO = tpo
	}
	if memo, ok := updates["memo"].(string); ok {
		coordinate.Memo = memo
	}
	if rating, ok := updates["rating"].(float32); ok {
		coordinate.Rating = rating
	}
	
	// Update silhouette information
	if siTopLength, ok := updates["si_top_length"].(int); ok {
		coordinate.SiTopLength = siTopLength
	}
	if siTopSleeve, ok := updates["si_top_sleeve"].(int); ok {
		coordinate.SiTopSleeve = siTopSleeve
	}
	if siBottomLength, ok := updates["si_bottom_length"].(int); ok {
		coordinate.SiBottomLength = siBottomLength
	}
	if siBottomType, ok := updates["si_bottom_type"].(int); ok {
		coordinate.SiBottomType = siBottomType
	}
	if siDressLength, ok := updates["si_dress_length"].(int); ok {
		coordinate.SiDressLength = siDressLength
	}
	if siDressSleeve, ok := updates["si_dress_sleeve"].(int); ok {
		coordinate.SiDressSleeve = siDressSleeve
	}
	if siOuterLength, ok := updates["si_outer_length"].(int); ok {
		coordinate.SiOuterLength = siOuterLength
	}
	if siOuterSleeve, ok := updates["si_outer_sleeve"].(int); ok {
		coordinate.SiOuterSleeve = siOuterSleeve
	}
	if siShoeSize, ok := updates["si_shoe_size"].(int); ok {
		coordinate.SiShoeSize = siShoeSize
	}
	
	// Upload new image if provided
	if image != nil {
		// Delete old image if exists
		if coordinate.Picture != "" {
			u.deleteImage(coordinate.Picture)
		}
		
		filename, err := u.uploadImage(image, "coordinates")
		if err != nil {
			return err
		}
		coordinate.Picture = filename
	}
	
	// Use transaction to update coordinate and items
	return u.db.Transaction(func(tx *gorm.DB) error {
		// Update coordinate
		if err := u.coordinateRepo.Update(ctx, coordinate); err != nil {
			return err
		}
		
		// Update items if provided
		if len(itemIDs) > 0 {
			// Remove all items from this coordinate
			if err := tx.Model(&domain.Item{}).Where("coordinate_id = ?", coordinateID).Update("coordinate_id", nil).Error; err != nil {
				return err
			}
			
			// Add new items
			for _, itemID := range itemIDs {
				item, err := u.itemRepo.FindByID(ctx, itemID)
				if err != nil {
					return err
				}
				if item == nil {
					return errors.New("item not found")
				}
				if item.UserID != userID {
					return errors.New("unauthorized: item does not belong to user")
				}
				
				if err := tx.Model(&domain.Item{}).Where("id = ?", itemID).Update("coordinate_id", coordinateID).Error; err != nil {
					return err
				}
			}
		}
		
		return nil
	})
}

// DeleteCoordinate deletes a coordinate
func (u *coordinateUsecase) DeleteCoordinate(ctx context.Context, userID uint, coordinateID uint) error {
	coordinate, err := u.coordinateRepo.FindByID(ctx, coordinateID)
	if err != nil {
		return err
	}
	if coordinate == nil {
		return errors.New("coordinate not found")
	}
	
	// Check ownership
	if coordinate.UserID != userID {
		return errors.New("unauthorized")
	}
	
	// Delete image if exists
	if coordinate.Picture != "" {
		u.deleteImage(coordinate.Picture)
	}
	
	// Use transaction to delete coordinate and update items
	return u.db.Transaction(func(tx *gorm.DB) error {
		// Remove items from coordinate
		if err := tx.Model(&domain.Item{}).Where("coordinate_id = ?", coordinateID).Update("coordinate_id", nil).Error; err != nil {
			return err
		}
		
		// Delete coordinate
		return u.coordinateRepo.Delete(ctx, coordinateID)
	})
}

// GetUserCoordinates gets coordinates for a user
func (u *coordinateUsecase) GetUserCoordinates(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, int64, error) {
	coordinates, err := u.coordinateRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := u.coordinateRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	
	return coordinates, count, nil
}

// SearchCoordinates searches coordinates with filters
func (u *coordinateUsecase) SearchCoordinates(ctx context.Context, filters repository.CoordinateFilter) ([]*domain.Coordinate, error) {
	return u.coordinateRepo.FindByFilters(ctx, filters)
}

// GetTimelineCoordinates gets timeline coordinates for a user (from followed users)
func (u *coordinateUsecase) GetTimelineCoordinates(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, error) {
	// Get followed users
	followedUsers, err := u.relationshipRepo.FindFollowing(ctx, userID, 0, 0)
	if err != nil {
		return nil, err
	}
	
	// Get coordinates from followed users
	var coordinates []*domain.Coordinate
	for _, user := range followedUsers {
		// Check if user is blocked
		isBlocked, err := u.blockRepo.ExistsByBlockerAndBlocked(ctx, userID, user.ID)
		if err != nil {
			return nil, err
		}
		if isBlocked {
			continue
		}
		
		userCoordinates, err := u.coordinateRepo.FindByUserID(ctx, user.ID, 10, 0) // Get recent 10
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, userCoordinates...)
	}
	
	// Sort by created_at desc and apply pagination
	// TODO: Implement proper sorting and pagination
	
	return coordinates, nil
}

// LikeCoordinate likes a coordinate
func (u *coordinateUsecase) LikeCoordinate(ctx context.Context, userID uint, coordinateID uint) error {
	// Check if coordinate exists
	coordinate, err := u.coordinateRepo.FindByID(ctx, coordinateID)
	if err != nil {
		return err
	}
	if coordinate == nil {
		return errors.New("coordinate not found")
	}
	
	// Check if already liked
	exists, err := u.likeCoordinateRepo.ExistsByUserAndCoordinate(ctx, userID, coordinateID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already liked")
	}
	
	// Create like
	like := &domain.LikeCoordinate{
		UserID:       userID,
		CoordinateID: coordinateID,
	}
	if err := u.likeCoordinateRepo.Create(ctx, like); err != nil {
		return err
	}
	
	// Create notification if not liking own coordinate
	if coordinate.UserID != userID {
		notification := &domain.Notification{
			SenderID:         userID,
			ReceiverID:       coordinate.UserID,
			CoordinateID:     &coordinateID,
			LikeCoordinateID: &like.ID,
			Action:           domain.NotificationActionLike,
		}
		if err := u.notificationRepo.Create(ctx, notification); err != nil {
			// Log error but don't fail the like operation
			fmt.Printf("Failed to create notification: %v\n", err)
		}
	}
	
	return nil
}

// UnlikeCoordinate unlikes a coordinate
func (u *coordinateUsecase) UnlikeCoordinate(ctx context.Context, userID uint, coordinateID uint) error {
	like, err := u.likeCoordinateRepo.FindByUserAndCoordinate(ctx, userID, coordinateID)
	if err != nil {
		return err
	}
	if like == nil {
		return errors.New("not liked")
	}
	
	return u.likeCoordinateRepo.Delete(ctx, like.ID)
}

// GetCoordinateLikes gets likes for a coordinate
func (u *coordinateUsecase) GetCoordinateLikes(ctx context.Context, coordinateID uint) ([]*domain.LikeCoordinate, error) {
	return u.likeCoordinateRepo.FindByCoordinateID(ctx, coordinateID)
}

// IsLikedByUser checks if a coordinate is liked by a user
func (u *coordinateUsecase) IsLikedByUser(ctx context.Context, userID uint, coordinateID uint) (bool, error) {
	return u.likeCoordinateRepo.ExistsByUserAndCoordinate(ctx, userID, coordinateID)
}

// GetUserCoordinateStatistics gets coordinate statistics for a user
func (u *coordinateUsecase) GetUserCoordinateStatistics(ctx context.Context, userID uint) (map[string]interface{}, error) {
	coordinates, err := u.coordinateRepo.FindByUserID(ctx, userID, 0, 0)
	if err != nil {
		return nil, err
	}
	
	// Calculate statistics
	stats := make(map[string]interface{})
	stats["total_count"] = len(coordinates)
	
	// Count by season and TPO
	seasonCount := make(map[int]int)
	tpoCount := make(map[int]int)
	
	var totalRating float32
	var totalLikes int64
	ratedCount := 0
	
	for _, coordinate := range coordinates {
		seasonCount[coordinate.Season]++
		tpoCount[coordinate.TPO]++
		
		if coordinate.Rating > 0 {
			totalRating += coordinate.Rating
			ratedCount++
		}
		
		// Count likes
		likeCount, err := u.likeCoordinateRepo.CountByCoordinateID(ctx, coordinate.ID)
		if err == nil {
			totalLikes += likeCount
		}
	}
	
	stats["season_count"] = seasonCount
	stats["tpo_count"] = tpoCount
	stats["total_likes"] = totalLikes
	
	if ratedCount > 0 {
		stats["average_rating"] = totalRating / float32(ratedCount)
	} else {
		stats["average_rating"] = 0
	}
	
	return stats, nil
}

// uploadImage uploads an image file
func (u *coordinateUsecase) uploadImage(file *multipart.FileHeader, folder string) (string, error) {
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
func (u *coordinateUsecase) deleteImage(filename string) error {
	if filename == "" {
		return nil
	}
	
	path := filepath.Join(u.config.Upload.Path, filename)
	return os.Remove(path)
}