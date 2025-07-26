package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/router"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
	"github.com/House-lovers7/speadwear-go/internal/usecase/impl"
	"github.com/House-lovers7/speadwear-go/pkg/config"
	"github.com/House-lovers7/speadwear-go/pkg/database"
	"gorm.io/gorm"
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

	// リポジトリの初期化
	repos := repository.NewContainer(database.DB)

	// ユースケースの初期化
	usecases := createUsecaseContainer(repos, cfg, database.DB)

	// ルーターの設定
	var r *gin.Engine
	if cfg.App.Env == "development" {
		r = router.SetupDevRouter(usecases, repos, cfg)
	} else {
		r = router.SetupRouter(usecases, repos, cfg)
	}

	// サーバーの起動
	log.Printf("Server starting on port %s in %s mode", cfg.App.Port, cfg.App.Env)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// createUsecaseContainer creates a usecase container with actual implementations
func createUsecaseContainer(repos *repository.Container, cfg *config.Config, db *gorm.DB) *usecase.Container {
	return &usecase.Container{
		User:       impl.NewUserUsecase(repos.User, cfg),
		Item:       impl.NewItemUsecase(repos.Item, cfg),
		Coordinate: impl.NewCoordinateUsecase(
			repos.Coordinate,
			repos.Item,
			repos.LikeCoordinate,
			repos.Relationship,
			repos.Block,
			repos.Notification,
			cfg,
			db,
		),
		Social: impl.NewSocialUsecase(
			repos.Comment,
			repos.Relationship,
			repos.Block,
			repos.Notification,
			repos.Coordinate,
			repos.User,
			cfg,
		),
	}
}