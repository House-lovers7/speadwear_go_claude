package testutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/pkg/utils"
	"gorm.io/gorm"
)

// Fixtures provides test data creation helpers
type Fixtures struct {
	db *gorm.DB
	t  *testing.T
}

// NewFixtures creates a new fixtures instance
func NewFixtures(t *testing.T, db *gorm.DB) *Fixtures {
	return &Fixtures{
		db: db,
		t:  t,
	}
}

// CreateUser creates a test user
func (f *Fixtures) CreateUser(opts ...func(*domain.User)) *domain.User {
	f.t.Helper()

	// Default user
	user := &domain.User{
		Name:  fmt.Sprintf("Test User %d", time.Now().UnixNano()),
		Email: fmt.Sprintf("test%d@example.com", time.Now().UnixNano()),
	}

	// Apply options
	for _, opt := range opts {
		opt(user)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		f.t.Fatalf("failed to hash password: %v", err)
	}
	user.PasswordDigest = hashedPassword

	if err := f.db.Create(user).Error; err != nil {
		f.t.Fatalf("failed to create user: %v", err)
	}

	return user
}

// CreateItem creates a test item
func (f *Fixtures) CreateItem(userID uint, opts ...func(*domain.Item)) *domain.Item {
	f.t.Helper()

	// Default item
	item := &domain.Item{
		UserID:    userID,
		SuperItem: "トップス",
		Season:    1,
		TPO:       1,
		Color:     1,
		Content:   "Test content",
		Picture:   "/uploads/test-item.jpg",
		Rating:    5,
	}

	// Apply options
	for _, opt := range opts {
		opt(item)
	}

	if err := f.db.Create(item).Error; err != nil {
		f.t.Fatalf("failed to create item: %v", err)
	}

	return item
}

// CreateCoordinate creates a test coordinate
func (f *Fixtures) CreateCoordinate(userID uint, opts ...func(*domain.Coordinate)) *domain.Coordinate {
	f.t.Helper()

	// Default coordinate
	coordinate := &domain.Coordinate{
		UserID:  userID,
		Picture: "/uploads/test-coordinate.jpg",
		Season:  1,
		TPO:     1,
		Rating:  5,
		Memo:    "Test coordinate",
	}

	// Apply options
	for _, opt := range opts {
		opt(coordinate)
	}

	if err := f.db.Create(coordinate).Error; err != nil {
		f.t.Fatalf("failed to create coordinate: %v", err)
	}

	// Create required items and associate with coordinate
	shoes := f.CreateItem(userID, func(i *domain.Item) {
		i.SuperItem = "シューズ"
		i.CoordinateID = &coordinate.ID
	})
	bottoms := f.CreateItem(userID, func(i *domain.Item) {
		i.SuperItem = "ボトムス"
		i.CoordinateID = &coordinate.ID
	})
	tops := f.CreateItem(userID, func(i *domain.Item) {
		i.SuperItem = "トップス"
		i.CoordinateID = &coordinate.ID
	})

	// Preload items for return value
	coordinate.Items = []domain.Item{*shoes, *bottoms, *tops}

	return coordinate
}

// CreateComment creates a test comment
func (f *Fixtures) CreateComment(userID, coordinateID uint, opts ...func(*domain.Comment)) *domain.Comment {
	f.t.Helper()

	comment := &domain.Comment{
		UserID:       userID,
		CoordinateID: coordinateID,
		Comment:      "Test comment",
	}

	// Apply options
	for _, opt := range opts {
		opt(comment)
	}

	if err := f.db.Create(comment).Error; err != nil {
		f.t.Fatalf("failed to create comment: %v", err)
	}

	return comment
}

// CreateLike creates a test like
func (f *Fixtures) CreateLike(userID, coordinateID uint) *domain.LikeCoordinate {
	f.t.Helper()

	like := &domain.LikeCoordinate{
		UserID:       userID,
		CoordinateID: coordinateID,
	}

	if err := f.db.Create(like).Error; err != nil {
		f.t.Fatalf("failed to create like: %v", err)
	}

	return like
}

// CreateRelationship creates a test follow relationship
func (f *Fixtures) CreateRelationship(followerID, followedID uint) *domain.Relationship {
	f.t.Helper()

	relationship := &domain.Relationship{
		FollowerID: followerID,
		FollowedID: followedID,
	}

	if err := f.db.Create(relationship).Error; err != nil {
		f.t.Fatalf("failed to create relationship: %v", err)
	}

	return relationship
}

// CreateBlock creates a test block
func (f *Fixtures) CreateBlock(blockerID, blockedID uint) *domain.Block {
	f.t.Helper()

	block := &domain.Block{
		BlockerID: blockerID,
		BlockedID: blockedID,
	}

	if err := f.db.Create(block).Error; err != nil {
		f.t.Fatalf("failed to create block: %v", err)
	}

	return block
}

// CreateNotification creates a test notification
func (f *Fixtures) CreateNotification(senderID, receiverID uint, action string, opts ...func(*domain.Notification)) *domain.Notification {
	f.t.Helper()

	notification := &domain.Notification{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Action:     action,
	}

	// Apply options
	for _, opt := range opts {
		opt(notification)
	}

	if err := f.db.Create(notification).Error; err != nil {
		f.t.Fatalf("failed to create notification: %v", err)
	}

	return notification
}

// Helper functions for common test scenarios

// CreateUserWithItems creates a user with multiple items
func (f *Fixtures) CreateUserWithItems(itemCount int) (*domain.User, []*domain.Item) {
	f.t.Helper()

	user := f.CreateUser()
	items := make([]*domain.Item, itemCount)

	for i := 0; i < itemCount; i++ {
		items[i] = f.CreateItem(user.ID, func(item *domain.Item) {
			item.SuperItem = domain.SuperItemCategories[i%len(domain.SuperItemCategories)]
			item.Content = fmt.Sprintf("Item %d", i+1)
		})
	}

	return user, items
}

// CreateUserWithCoordinates creates a user with coordinates
func (f *Fixtures) CreateUserWithCoordinates(coordinateCount int) (*domain.User, []*domain.Coordinate) {
	f.t.Helper()

	user := f.CreateUser()
	coordinates := make([]*domain.Coordinate, coordinateCount)

	for i := 0; i < coordinateCount; i++ {
		coordinates[i] = f.CreateCoordinate(user.ID, func(c *domain.Coordinate) {
			c.Memo = fmt.Sprintf("Coordinate %d", i+1)
		})
	}

	return user, coordinates
}

// CreateFollowingRelationship creates a following relationship between users
func (f *Fixtures) CreateFollowingRelationship() (*domain.User, *domain.User, *domain.Relationship) {
	f.t.Helper()

	follower := f.CreateUser(func(u *domain.User) {
		u.Name = "Follower"
		u.Email = "follower@example.com"
	})

	followed := f.CreateUser(func(u *domain.User) {
		u.Name = "Followed"
		u.Email = "followed@example.com"
	})

	relationship := f.CreateRelationship(follower.ID, followed.ID)

	return follower, followed, relationship
}