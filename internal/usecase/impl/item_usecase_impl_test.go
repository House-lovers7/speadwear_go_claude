package impl

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/testutil"
	"github.com/House-lovers7/speadwear-go/pkg/config"
)

func setupItemUsecase(t *testing.T) (*itemUsecase, *testutil.Fixtures) {
	db := testutil.TestDB(t)
	
	repos := repository.NewContainer(db)
	cfg := &config.Config{
		Upload: config.UploadConfig{
			Path:        "./test-uploads",
			MaxFileSize: 5 * 1024 * 1024,
		},
	}
	
	usecase := NewItemUsecase(
		repos.Item,
		cfg,
	).(*itemUsecase)
	
	fixtures := testutil.NewFixtures(t, db)
	
	return usecase, fixtures
}

func TestItemUsecase_CreateItem(t *testing.T) {
	usecase, fixtures := setupItemUsecase(t)
	ctx := context.Background()
	
	// Create test user
	user := fixtures.CreateUser()
	
	tests := []struct {
		name    string
		userID  uint
		item    *domain.Item
		image   *multipart.FileHeader
		wantErr bool
		errMsg  string
	}{
		{
			name:   "valid item creation",
			userID: user.ID,
			item: &domain.Item{
				SuperItem: "シャツ",
				Color:     domain.ColorBlue,
				Season:    domain.SeasonSpring,
				TPO:       domain.TPOWork,
				Content:   "Test item content",
				Memo:      "Test item",
				Rating:    5,
			},
			image:   nil,
			wantErr: false,
		},
		{
			name:   "item with empty super item",
			userID: user.ID,
			item: &domain.Item{
				SuperItem: "",
				Color:     domain.ColorBlue,
			},
			image:   nil,
			wantErr: false, // No validation in usecase
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.CreateItem(ctx, tt.userID, tt.item, tt.image)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("CreateItem() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr && tt.item != nil {
				if tt.item.ID == 0 {
					t.Error("CreateItem() did not set item ID")
				}
				if tt.item.UserID != tt.userID {
					t.Errorf("CreateItem() UserID = %v, want %v", tt.item.UserID, tt.userID)
				}
			}
		})
	}
}

func TestItemUsecase_GetItem(t *testing.T) {
	usecase, fixtures := setupItemUsecase(t)
	ctx := context.Background()
	
	// Create test user and item
	user := fixtures.CreateUser()
	item := fixtures.CreateItem(user.ID)
	
	tests := []struct {
		name    string
		itemID  uint
		wantErr bool
	}{
		{
			name:    "existing item",
			itemID:  item.ID,
			wantErr: false,
		},
		{
			name:    "non-existent item",
			itemID:  99999,
			wantErr: true,
		},
		{
			name:    "zero item ID",
			itemID:  0,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundItem, err := usecase.GetItem(ctx, tt.itemID)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && foundItem != nil {
				if foundItem.ID != tt.itemID {
					t.Errorf("GetItem() ID = %v, want %v", foundItem.ID, tt.itemID)
				}
			}
		})
	}
}

func TestItemUsecase_UpdateItem(t *testing.T) {
	usecase, fixtures := setupItemUsecase(t)
	ctx := context.Background()
	
	// Create test user and item
	user := fixtures.CreateUser()
	item := fixtures.CreateItem(user.ID, func(i *domain.Item) {
		i.SuperItem = "Original Item"
		i.Content = "Original content"
		i.Rating = 3
	})
	
	tests := []struct {
		name    string
		userID  uint
		itemID  uint
		updates map[string]interface{}
		image   *multipart.FileHeader
		wantErr bool
		errMsg  string
	}{
		{
			name:   "valid update",
			userID: user.ID,
			itemID: item.ID,
			updates: map[string]interface{}{
				"super_item": "Updated Item",
				"content":    "Updated content",
				"rating":     float32(5),
			},
			image:   nil,
			wantErr: false,
		},
		{
			name:   "partial update",
			userID: user.ID,
			itemID: item.ID,
			updates: map[string]interface{}{
				"memo": "New memo",
			},
			image:   nil,
			wantErr: false,
		},
		{
			name:   "non-existent item",
			userID: user.ID,
			itemID: 99999,
			updates: map[string]interface{}{
				"super_item": "Ghost Item",
			},
			image:   nil,
			wantErr: true,
			errMsg:  "item not found",
		},
		{
			name:   "unauthorized update",
			userID: 99999,
			itemID: item.ID,
			updates: map[string]interface{}{
				"super_item": "Hacked Item",
			},
			image:   nil,
			wantErr: true,
			errMsg:  "unauthorized",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.UpdateItem(ctx, tt.userID, tt.itemID, tt.updates, tt.image)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("UpdateItem() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr {
				// Verify update
				updated, _ := usecase.GetItem(ctx, tt.itemID)
				if superItem, ok := tt.updates["super_item"].(string); ok && updated.SuperItem != superItem {
					t.Errorf("UpdateItem() did not update super_item: got %v, want %v", updated.SuperItem, superItem)
				}
				if content, ok := tt.updates["content"].(string); ok && updated.Content != content {
					t.Errorf("UpdateItem() did not update content: got %v, want %v", updated.Content, content)
				}
				if rating, ok := tt.updates["rating"].(float32); ok && updated.Rating != rating {
					t.Errorf("UpdateItem() did not update rating: got %v, want %v", updated.Rating, rating)
				}
			}
		})
	}
}

func TestItemUsecase_DeleteItem(t *testing.T) {
	usecase, fixtures := setupItemUsecase(t)
	ctx := context.Background()
	
	// Create test user and items
	user := fixtures.CreateUser()
	item := fixtures.CreateItem(user.ID)
	other := fixtures.CreateUser()
	
	tests := []struct {
		name    string
		userID  uint
		itemID  uint
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid deletion",
			userID:  user.ID,
			itemID:  item.ID,
			wantErr: false,
		},
		{
			name:    "non-existent item",
			userID:  user.ID,
			itemID:  99999,
			wantErr: true,
			errMsg:  "item not found",
		},
		{
			name:    "unauthorized deletion",
			userID:  other.ID,
			itemID:  item.ID,
			wantErr: true,
			errMsg:  "unauthorized",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.DeleteItem(ctx, tt.userID, tt.itemID)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("DeleteItem() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr {
				// Verify deletion
				_, err := usecase.GetItem(ctx, tt.itemID)
				if err == nil {
					t.Error("DeleteItem() did not delete the item")
				}
			}
		})
	}
}

func TestItemUsecase_GetUserItems(t *testing.T) {
	usecase, fixtures := setupItemUsecase(t)
	ctx := context.Background()
	
	// Create test users with items
	user1 := fixtures.CreateUser()
	fixtures.CreateItem(user1.ID)
	fixtures.CreateItem(user1.ID)
	fixtures.CreateItem(user1.ID)
	
	user2 := fixtures.CreateUser()
	fixtures.CreateItem(user2.ID)
	fixtures.CreateItem(user2.ID)
	
	tests := []struct {
		name      string
		userID    uint
		limit     int
		offset    int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "get all user1 items",
			userID:    user1.ID,
			limit:     10,
			offset:    0,
			wantCount: 3,
		},
		{
			name:      "get user2 items with pagination",
			userID:    user2.ID,
			limit:     1,
			offset:    0,
			wantCount: 1,
		},
		{
			name:      "get items with offset",
			userID:    user1.ID,
			limit:     10,
			offset:    2,
			wantCount: 1,
		},
		{
			name:      "non-existent user",
			userID:    99999,
			limit:     10,
			offset:    0,
			wantCount: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, total, err := usecase.GetUserItems(ctx, tt.userID, tt.limit, tt.offset)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if len(items) != tt.wantCount {
				t.Errorf("GetUserItems() returned %d items, want %d", len(items), tt.wantCount)
			}
			
			// Verify all items belong to the requested user
			for _, item := range items {
				if item.UserID != tt.userID {
					t.Errorf("GetUserItems() returned item with wrong UserID = %d, want %d", item.UserID, tt.userID)
				}
			}
			
			// Check total count
			if tt.userID == user1.ID && total != 3 {
				t.Errorf("GetUserItems() total = %d, want 3", total)
			}
			if tt.userID == user2.ID && total != 2 {
				t.Errorf("GetUserItems() total = %d, want 2", total)
			}
		})
	}
}

func TestItemUsecase_SearchItems(t *testing.T) {
	usecase, fixtures := setupItemUsecase(t)
	ctx := context.Background()
	
	// Create test user with various items
	user := fixtures.CreateUser()
	
	fixtures.CreateItem(user.ID, func(i *domain.Item) {
		i.SuperItem = "Tシャツ"
		i.Content = "Nike T-shirt"
		i.Color = domain.ColorBlue
		i.Season = domain.SeasonSummer
		i.TPO = domain.TPOCasual
	})
	
	fixtures.CreateItem(user.ID, func(i *domain.Item) {
		i.SuperItem = "パンツ"
		i.Content = "Adidas pants"
		i.Color = domain.ColorBlack
		i.Season = domain.SeasonWinter
		i.TPO = domain.TPOWork
	})
	
	fixtures.CreateItem(user.ID, func(i *domain.Item) {
		i.SuperItem = "スニーカー"
		i.Content = "Nike sneakers"
		i.Color = domain.ColorWhite
		i.Season = domain.SeasonAllSeason
		i.TPO = domain.TPOCasual
	})
	
	tests := []struct {
		name      string
		filter    repository.ItemFilter
		wantCount int
	}{
		{
			name: "filter by super item",
			filter: repository.ItemFilter{
				UserID:    &user.ID,
				SuperItem: strPtr("Tシャツ"),
				Limit:     10,
			},
			wantCount: 1,
		},
		{
			name: "filter by color",
			filter: repository.ItemFilter{
				UserID: &user.ID,
				Color:  intPtr(domain.ColorBlue),
				Limit:  10,
			},
			wantCount: 1,
		},
		{
			name: "filter by season",
			filter: repository.ItemFilter{
				UserID: &user.ID,
				Season: intPtr(domain.SeasonSummer),
				Limit:  10,
			},
			wantCount: 1,
		},
		{
			name: "filter by TPO",
			filter: repository.ItemFilter{
				UserID: &user.ID,
				TPO:    intPtr(domain.TPOCasual),
				Limit:  10,
			},
			wantCount: 2,
		},
		{
			name: "multiple filters",
			filter: repository.ItemFilter{
				UserID: &user.ID,
				TPO:    intPtr(domain.TPOCasual),
				Color:  intPtr(domain.ColorBlue),
				Limit:  10,
			},
			wantCount: 1,
		},
		{
			name: "no matches",
			filter: repository.ItemFilter{
				UserID:    &user.ID,
				SuperItem: strPtr("ジャケット"),
				Limit:     10,
			},
			wantCount: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := usecase.SearchItems(ctx, tt.filter)
			if err != nil {
				t.Fatalf("SearchItems() error = %v", err)
			}
			
			if len(items) != tt.wantCount {
				t.Errorf("SearchItems() returned %d items, want %d", len(items), tt.wantCount)
			}
		})
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}