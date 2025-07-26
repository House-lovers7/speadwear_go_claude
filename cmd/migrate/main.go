package main

import (
	"log"
	"os"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/pkg/config"
	"github.com/House-lovers7/speadwear-go/pkg/database"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// データベース接続
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// コマンドライン引数の確認
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go [up|down|reset]")
	}

	command := os.Args[1]

	switch command {
	case "up":
		log.Println("Running migrations...")
		if err := runMigrations(); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("Migrations completed successfully")

	case "down":
		log.Println("Dropping all tables...")
		if err := dropTables(); err != nil {
			log.Fatal("Failed to drop tables:", err)
		}
		log.Println("All tables dropped successfully")

	case "reset":
		log.Println("Resetting database...")
		if err := dropTables(); err != nil {
			log.Fatal("Failed to drop tables:", err)
		}
		if err := runMigrations(); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("Database reset completed successfully")

	default:
		log.Fatal("Unknown command. Use: up, down, or reset")
	}
}

func runMigrations() error {
	models := domain.GetAllModels()
	return database.Migrate(models...)
}

func dropTables() error {
	// 逆順でテーブルをドロップ（外部キー制約を考慮）
	tables := []string{
		"notifications",
		"blocks",
		"relationships",
		"like_coordinates",
		"comments",
		"items",
		"coordinates",
		"users",
	}

	for _, table := range tables {
		if err := database.DB.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
			return err
		}
		log.Printf("Dropped table: %s", table)
	}

	return nil
}