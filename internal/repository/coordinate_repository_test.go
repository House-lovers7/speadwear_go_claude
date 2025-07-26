package repository

import (
	"context"
	"testing"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/testutil"
)

func TestCoordinateRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewCoordinateRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user
	user := fixtures.CreateUser()

	tests := []struct {
		name    string
		coord   *domain.Coordinate
		wantErr bool
	}{
		{
			name: "valid coordinate",
			coord: &domain.Coordinate{
				UserID:  user.ID,
				Season:  domain.SeasonSpring,
				TPO:     domain.TPOWork,
				Picture: "/uploads/coordinate.jpg",
				Rating:  5,
				Memo:    "Test coordinate",
			},
			wantErr: false,
		},
		{
			name: "coordinate without user",
			coord: &domain.Coordinate{
				UserID: 99999, // Non-existent user
				Season: domain.SeasonSpring,
				TPO:    domain.TPOWork,
				Rating: 5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.coord)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.coord.ID == 0 {
				t.Error("Create() did not set coordinate ID")
			}
		})
	}
}

func TestCoordinateRepository_FindByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewCoordinateRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test data
	user1, coords1 := fixtures.CreateUserWithCoordinates(3)
	user2, coords2 := fixtures.CreateUserWithCoordinates(2)

	tests := []struct {
		name    string
		userID  uint
		limit   int
		offset  int
		want    int
		wantErr bool
	}{
		{
			name:   "get user1 coordinates",
			userID: user1.ID,
			limit:  10,
			offset: 0,
			want:   len(coords1),
		},
		{
			name:   "get user2 coordinates",
			userID: user2.ID,
			limit:  10,
			offset: 0,
			want:   len(coords2),
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
			coords, err := repo.FindByUserID(ctx, tt.userID, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(coords) != tt.want {
				t.Errorf("FindByUserID() returned %d coordinates, want %d", len(coords), tt.want)
			}

			// Verify all coordinates belong to the requested user
			for _, coord := range coords {
				if coord.UserID != tt.userID {
					t.Errorf("FindByUserID() returned coordinate with wrong UserID = %d, want %d", coord.UserID, tt.userID)
				}
			}
		})
	}
}

func TestCoordinateRepository_FindByFilters(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewCoordinateRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test user and coordinates with various attributes
	user := fixtures.CreateUser()
	
	// Create coordinates with different attributes
	springWork5 := fixtures.CreateCoordinate(user.ID, func(c *domain.Coordinate) {
		c.Season = domain.SeasonSpring
		c.TPO = domain.TPOWork
		c.Rating = 5
		c.Memo = "Spring work outfit"
	})

	fixtures.CreateCoordinate(user.ID, func(c *domain.Coordinate) {
		c.Season = domain.SeasonSummer
		c.TPO = domain.TPOCasual
		c.Rating = 4
		c.Memo = "Summer casual outfit"
	})

	springCasual45 := fixtures.CreateCoordinate(user.ID, func(c *domain.Coordinate) {
		c.Season = domain.SeasonSpring
		c.TPO = domain.TPOCasual
		c.Rating = 4.5
		c.Memo = "Spring casual outfit"
	})

	tests := []struct {
		name    string
		filter  CoordinateFilter
		wantIDs []uint
	}{
		{
			name: "filter by season",
			filter: CoordinateFilter{
				UserID: &user.ID,
				Season: intPtr(domain.SeasonSpring),
				Limit:  10,
			},
			wantIDs: []uint{springWork5.ID, springCasual45.ID},
		},
		{
			name: "filter by TPO",
			filter: CoordinateFilter{
				UserID: &user.ID,
				TPO:    intPtr(domain.TPOWork),
				Limit:  10,
			},
			wantIDs: []uint{springWork5.ID},
		},
		{
			name: "filter by rating range",
			filter: CoordinateFilter{
				UserID:    &user.ID,
				MinRating: float32Ptr(4.5),
				MaxRating: float32Ptr(5),
				Limit:     10,
			},
			wantIDs: []uint{springWork5.ID, springCasual45.ID},
		},
		{
			name: "multiple filters",
			filter: CoordinateFilter{
				UserID:    &user.ID,
				Season:    intPtr(domain.SeasonSpring),
				MinRating: float32Ptr(4),
				Limit:     10,
			},
			wantIDs: []uint{springWork5.ID, springCasual45.ID},
		},
		{
			name: "no matches",
			filter: CoordinateFilter{
				UserID: &user.ID,
				Season: intPtr(domain.SeasonWinter),
				Limit:  10,
			},
			wantIDs: []uint{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coords, err := repo.FindByFilters(ctx, tt.filter)
			if err != nil {
				t.Fatalf("FindByFilters() error = %v", err)
			}

			if len(coords) != len(tt.wantIDs) {
				t.Errorf("FindByFilters() returned %d coordinates, want %d", len(coords), len(tt.wantIDs))
			}

			// Check that returned coordinates match expected IDs
			gotIDs := make(map[uint]bool)
			for _, coord := range coords {
				gotIDs[coord.ID] = true
			}

			for _, wantID := range tt.wantIDs {
				if !gotIDs[wantID] {
					t.Errorf("FindByFilters() missing expected coordinate ID %d", wantID)
				}
			}
		})
	}
}

func TestCoordinateRepository_FindWithItems(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewCoordinateRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test coordinate with items
	user := fixtures.CreateUser()
	coord := fixtures.CreateCoordinate(user.ID)

	// Find coordinate with items
	found, err := repo.FindWithItems(ctx, coord.ID)
	if err != nil {
		t.Fatalf("FindWithItems() error = %v", err)
	}

	if found == nil {
		t.Fatal("FindWithItems() returned nil")
	}

	if found.ID != coord.ID {
		t.Errorf("FindWithItems() returned wrong coordinate ID = %d, want %d", found.ID, coord.ID)
	}

	// Check that items are loaded
	if len(found.Items) == 0 {
		t.Error("FindWithItems() did not load items")
	}

	// Verify items belong to the coordinate
	for _, item := range found.Items {
		if item.CoordinateID == nil || *item.CoordinateID != coord.ID {
			t.Errorf("FindWithItems() loaded item with wrong CoordinateID")
		}
	}

	// Test non-existent coordinate
	notFound, err := repo.FindWithItems(ctx, 99999)
	if err != nil {
		t.Fatalf("FindWithItems() with non-existent ID error = %v", err)
	}
	if notFound != nil {
		t.Error("FindWithItems() should return nil for non-existent coordinate")
	}
}

func TestCoordinateRepository_Update(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewCoordinateRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test coordinate
	user := fixtures.CreateUser()
	coord := fixtures.CreateCoordinate(user.ID)
	
	originalMemo := coord.Memo
	originalRating := coord.Rating

	// Update coordinate
	coord.Memo = "Updated memo"
	coord.Rating = 3
	coord.SiTopLength = domain.TopLengthLong
	coord.SiTopSleeve = domain.TopSleeveLong

	err := repo.Update(ctx, coord)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	updated, err := repo.FindByID(ctx, coord.ID)
	if err != nil {
		t.Fatalf("FindByID() after update error = %v", err)
	}

	if updated.Memo == originalMemo {
		t.Error("Update() did not update memo")
	}
	if updated.Rating == originalRating {
		t.Error("Update() did not update rating")
	}
	if updated.SiTopLength != coord.SiTopLength {
		t.Error("Update() did not update SiTopLength")
	}
}

func TestCoordinateRepository_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewCoordinateRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test coordinate
	user := fixtures.CreateUser()
	coord := fixtures.CreateCoordinate(user.ID)

	// Delete coordinate
	err := repo.Delete(ctx, coord.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deletion
	deleted, err := repo.FindByID(ctx, coord.ID)
	if err != nil {
		t.Fatalf("FindByID() after delete error = %v", err)
	}
	if deleted != nil {
		t.Error("Delete() did not delete coordinate")
	}
}

func TestCoordinateRepository_CountByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewCoordinateRepository(db)
	ctx := context.Background()
	fixtures := testutil.NewFixtures(t, db)

	// Create test data
	user1, coords1 := fixtures.CreateUserWithCoordinates(4)
	user2, _ := fixtures.CreateUserWithCoordinates(2)

	tests := []struct {
		name    string
		userID  uint
		want    int64
		wantErr bool
	}{
		{
			name:   "count user1 coordinates",
			userID: user1.ID,
			want:   int64(len(coords1)),
		},
		{
			name:   "count user2 coordinates",
			userID: user2.ID,
			want:   2,
		},
		{
			name:   "count non-existent user coordinates",
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