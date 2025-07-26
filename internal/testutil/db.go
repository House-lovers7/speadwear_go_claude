package testutil

import (
	"fmt"
	"testing"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB creates a test database connection
func TestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Test database configuration
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root_password",
		"localhost",
		"3306",
		"speadwear_test",
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Migrate test database
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Item{},
		&domain.Coordinate{},
		&domain.Comment{},
		&domain.LikeCoordinate{},
		&domain.Relationship{},
		&domain.Block{},
		&domain.Notification{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	// Clean up function
	t.Cleanup(func() {
		CleanupDB(t, db)
	})

	return db
}

// CleanupDB cleans up all test data
func CleanupDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Delete all data in reverse order of dependencies
	tables := []interface{}{
		&domain.Notification{},
		&domain.Block{},
		&domain.Relationship{},
		&domain.LikeCoordinate{},
		&domain.Comment{},
		&domain.Coordinate{},
		&domain.Item{},
		&domain.User{},
	}

	for _, table := range tables {
		if err := db.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
			t.Logf("failed to clean up %T: %v", table, err)
		}
	}
}

// TruncateTables truncates all tables for fresh test data
func TruncateTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables := []string{
		"notifications",
		"blocks",
		"relationships",
		"like_coordinates",
		"comments",
		"coordinates",
		"items",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error; err != nil {
			t.Logf("failed to truncate %s: %v", table, err)
		}
	}
}

// TestTransaction creates a transaction for testing
func TestTransaction(t *testing.T, db *gorm.DB, fn func(*gorm.DB)) {
	t.Helper()

	tx := db.Begin()
	defer tx.Rollback()

	fn(tx)
}