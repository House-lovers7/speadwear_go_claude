package database

import (
	"fmt"
	"log"

	"github.com/House-lovers7/speadwear-go/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.DatabaseConfig) error {
	var err error
	
	// ログレベルの設定
	logLevel := logger.Info
	if cfg.Host == "production" {
		logLevel = logger.Error
	}

	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func Migrate(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}