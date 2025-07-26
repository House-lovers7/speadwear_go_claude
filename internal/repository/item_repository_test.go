package repository

import (
	"context"
	"testing"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/testutil"
)

func TestItemRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewItemRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user
	user := fixtures.CreateUser()

	tests := []struct {
		name    string
		item    *domain.Item
		wantErr bool
	}{
		{
			name: "valid item",
			item: &domain.Item{
				UserID:    user.ID,
				SuperItem: "トップス",
				Season:    domain.SeasonSpring,
				TPO:       domain.TPOWork,
				Color:     domain.ColorBlack,
				Content:   "Test content",
				Picture:   "/uploads/item.jpg",
				Rating:    5,
			},
			wantErr: false,
		},
		{
			name: "item without user",
			item: &domain.Item{
				UserID:    99999, // Non-existent user
				SuperItem: "トップス",
				Season:    domain.SeasonSpring,
				TPO:       domain.TPOWork,
				Color:     domain.ColorBlack,
				Rating:    5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.item.ID == 0 {
				t.Error("Create() did not set item ID")
			}
		})
	}
}

func TestItemRepository_FindByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewItemRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test data
	user1, items1 := fixtures.CreateUserWithItems(3)
	user2, items2 := fixtures.CreateUserWithItems(2)

	tests := []struct {
		name    string
		userID  uint
		limit   int
		offset  int
		want    int
		wantErr bool
	}{
		{
			name:   "get user1 items",
			userID: user1.ID,
			limit:  10,
			offset: 0,
			want:   len(items1),
		},
		{
			name:   "get user2 items",
			userID: user2.ID,
			limit:  10,
			offset: 0,
			want:   len(items2),
		},
		{
			name:   "get with pagination",
			userID: user1.ID,
			limit:  2,
			offset: 0,
			want:   2,
		},
		{
			name:   "get with offset",
			userID: user1.ID,
			limit:  10,
			offset: 2,
			want:   1,
		},
		{
			name:   "non-existent user",
			userID: 99999,
			limit:  10,
			offset: 0,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := repo.FindByUserID(ctx, tt.userID, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(items) != tt.want {
				t.Errorf("FindByUserID() returned %d items, want %d", len(items), tt.want)
			}

			// Verify all items belong to the requested user
			for _, item := range items {
				if item.UserID != tt.userID {
					t.Errorf("FindByUserID() returned item with wrong UserID = %d, want %d", item.UserID, tt.userID)
				}
			}
		})
	}
}

func TestItemRepository_FindByFilters(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewItemRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user and items with various attributes
	user := fixtures.CreateUser()
	
	// Create items with different attributes
	springTopsWork := fixtures.CreateItem(user.ID, func(i *domain.Item) {
		i.SuperItem = "トップス"
		i.Season = domain.SeasonSpring
		i.TPO = domain.TPOWork
		i.Color = domain.ColorBlack
		i.Rating = 5
	})

	fixtures.CreateItem(user.ID, func(i *domain.Item) {
		i.SuperItem = "ボトムス"
		i.Season = domain.SeasonSummer
		i.TPO = domain.TPOCasual
		i.Color = domain.ColorWhite
		i.Rating = 4
	})

	springTopsHighRating := fixtures.CreateItem(user.ID, func(i *domain.Item) {
		i.SuperItem = "トップス"
		i.Season = domain.SeasonSpring
		i.TPO = domain.TPOCasual
		i.Color = domain.ColorBlue
		i.Rating = 4.5
	})

	tests := []struct {
		name    string
		filter  ItemFilter
		wantIDs []uint
	}{
		{
			name: "filter by season",
			filter: ItemFilter{
				UserID: &user.ID,
				Season: intPtr(domain.SeasonSpring),
				Limit:  10,
			},
			wantIDs: []uint{springTopsWork.ID, springTopsHighRating.ID},
		},
		{
			name: "filter by TPO",
			filter: ItemFilter{
				UserID: &user.ID,
				TPO:    intPtr(domain.TPOWork),
				Limit:  10,
			},
			wantIDs: []uint{springTopsWork.ID},
		},
		{
			name: "filter by color",
			filter: ItemFilter{
				UserID: &user.ID,
				Color:  intPtr(domain.ColorBlack),
				Limit:  10,
			},
			wantIDs: []uint{springTopsWork.ID},
		},
		{
			name: "filter by super item",
			filter: ItemFilter{
				UserID:    &user.ID,
				SuperItem: stringPtr("トップス"),
				Limit:     10,
			},
			wantIDs: []uint{springTopsWork.ID, springTopsHighRating.ID},
		},
		{
			name: "filter by rating range",
			filter: ItemFilter{
				UserID:    &user.ID,
				MinRating: float32Ptr(4.5),
				MaxRating: float32Ptr(5),
				Limit:     10,
			},
			wantIDs: []uint{springTopsWork.ID, springTopsHighRating.ID},
		},
		{
			name: "multiple filters",
			filter: ItemFilter{
				UserID:    &user.ID,
				Season:    intPtr(domain.SeasonSpring),
				SuperItem: stringPtr("トップス"),
				MinRating: float32Ptr(4),
				Limit:     10,
			},
			wantIDs: []uint{springTopsWork.ID, springTopsHighRating.ID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := repo.FindByFilters(ctx, tt.filter)
			if err != nil {
				t.Fatalf("FindByFilters() error = %v", err)
			}

			if len(items) != len(tt.wantIDs) {
				t.Errorf("FindByFilters() returned %d items, want %d", len(items), len(tt.wantIDs))
			}

			// Check that returned items match expected IDs
			gotIDs := make(map[uint]bool)
			for _, item := range items {
				gotIDs[item.ID] = true
			}

			for _, wantID := range tt.wantIDs {
				if !gotIDs[wantID] {
					t.Errorf("FindByFilters() missing expected item ID %d", wantID)
				}
			}
		})
	}
}

func TestItemRepository_Update(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewItemRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test item
	user := fixtures.CreateUser()
	item := fixtures.CreateItem(user.ID)
	
	originalContent := item.Content
	originalRating := item.Rating

	// Update item
	item.Content = "Updated content"
	item.Rating = 3
	item.Memo = "Updated memo"

	err := repo.Update(ctx, item)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	updated, err := repo.FindByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("FindByID() after update error = %v", err)
	}

	if updated.Content == originalContent {
		t.Error("Update() did not update content")
	}
	if updated.Rating == originalRating {
		t.Error("Update() did not update rating")
	}
	if updated.Memo != item.Memo {
		t.Error("Update() did not update memo")
	}
}

func TestItemRepository_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewItemRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test item
	user := fixtures.CreateUser()
	item := fixtures.CreateItem(user.ID)

	// Delete item
	err := repo.Delete(ctx, item.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deletion
	deleted, err := repo.FindByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("FindByID() after delete error = %v", err)
	}
	if deleted != nil {
		t.Error("Delete() did not delete item")
	}
}

func TestItemRepository_CountByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewItemRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test data
	user1, items1 := fixtures.CreateUserWithItems(5)
	user2, _ := fixtures.CreateUserWithItems(3)

	tests := []struct {
		name    string
		userID  uint
		want    int64
		wantErr bool
	}{
		{
			name:   "count user1 items",
			userID: user1.ID,
			want:   int64(len(items1)),
		},
		{
			name:   "count user2 items",
			userID: user2.ID,
			want:   3,
		},
		{
			name:   "count non-existent user items",
			userID: 99999,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := repo.CountByUserID(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if count != tt.want {
				t.Errorf("CountByUserID() = %d, want %d", count, tt.want)
			}
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func float32Ptr(f float32) *float32 {
	return &f
}

func uintPtr(u uint) *uint {
	return &u
}