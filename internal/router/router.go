package router

import (
	"github.com/gin-gonic/gin"
	"github.com/House-lovers7/speadwear-go/internal/handler"
	"github.com/House-lovers7/speadwear-go/internal/middleware"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
	"github.com/House-lovers7/speadwear-go/pkg/config"
)

// SetupRouter configures and returns the main router
func SetupRouter(usecases *usecase.Container, repos *repository.Container, cfg *config.Config) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS middleware (add this in production)
	// r.Use(middleware.CORS())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(cfg, usecases.User)
	userHandler := handler.NewUserHandler(usecases.User)
	itemHandler := handler.NewItemHandler(usecases.Item)
	coordinateHandler := handler.NewCoordinateHandler(
		usecases.Coordinate,
		repos.Comment,
		repos.LikeCoordinate,
	)
	socialHandler := handler.NewSocialHandler(usecases.Social)

	// Static files for uploaded images
	r.Static("/uploads", "./uploads")

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Speadwear API is running",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := v1.Group("")
		{
			// Authentication
			public.POST("/auth/login", authHandler.Login)
			public.POST("/auth/signup", authHandler.Signup)
			
			// Password reset
			public.POST("/users/password/reset", userHandler.ResetPasswordRequest)
			public.PUT("/users/password/reset", userHandler.ResetPassword)
			
			// Account activation
			public.POST("/users/activate", userHandler.ActivateAccount)
			public.POST("/users/activate/resend", userHandler.ResendActivationEmail)
			
			// Public user profiles
			public.GET("/users/:id", userHandler.GetUser)
			public.GET("/users/:user_id/items", itemHandler.GetUserItems)
			public.GET("/users/:user_id/coordinates", coordinateHandler.GetUserCoordinates)
			
			// Public item and coordinate viewing
			public.GET("/items/:id", itemHandler.GetItem)
			public.GET("/coordinates/:id", coordinateHandler.GetCoordinate)
			public.GET("/coordinates/:id/comments", coordinateHandler.GetCoordinateComments)
			
			// Search (public)
			public.GET("/items/search", itemHandler.SearchItems)
			public.GET("/coordinates/search", coordinateHandler.SearchCoordinates)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired(cfg))
		{
			// Authentication
			protected.POST("/auth/logout", authHandler.Logout)
			protected.POST("/auth/refresh", authHandler.RefreshToken)
			protected.GET("/auth/me", authHandler.Me)

			// User management
			protected.GET("/users", userHandler.ListUsers)
			protected.GET("/users/me", userHandler.GetMe)
			protected.PUT("/users/profile", userHandler.UpdateProfile)
			protected.PUT("/users/password", userHandler.ChangePassword)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.DELETE("/users/:id", userHandler.DeleteUser)

			// Item management
			protected.POST("/items", itemHandler.CreateItem)
			protected.GET("/items", itemHandler.GetMyItems)
			protected.PUT("/items/:id", itemHandler.UpdateItem)
			protected.DELETE("/items/:id", itemHandler.DeleteItem)
			protected.DELETE("/items", itemHandler.DeleteItems) // Batch delete
			protected.GET("/items/statistics", itemHandler.GetItemStatistics)

			// Coordinate management
			protected.POST("/coordinates", coordinateHandler.CreateCoordinate)
			protected.GET("/coordinates", coordinateHandler.GetMyCoordinates)
			protected.GET("/coordinates/timeline", coordinateHandler.GetTimeline)
			protected.PUT("/coordinates/:id", coordinateHandler.UpdateCoordinate)
			protected.DELETE("/coordinates/:id", coordinateHandler.DeleteCoordinate)
			protected.GET("/coordinates/statistics", coordinateHandler.GetCoordinateStatistics)

			// Like functionality
			protected.POST("/coordinates/:id/like", coordinateHandler.LikeCoordinate)
			protected.DELETE("/coordinates/:id/like", coordinateHandler.UnlikeCoordinate)

			// Comment functionality
			protected.POST("/comments", socialHandler.CreateComment)
			protected.PUT("/comments/:id", socialHandler.UpdateComment)
			protected.DELETE("/comments/:id", socialHandler.DeleteComment)

			// Follow functionality
			protected.POST("/follow/:user_id", socialHandler.FollowUser)
			protected.DELETE("/follow/:user_id", socialHandler.UnfollowUser)
			protected.GET("/follow/followers", socialHandler.GetFollowers)
			protected.GET("/follow/following", socialHandler.GetFollowing)
			protected.GET("/follow/status/:user_id", socialHandler.CheckFollowStatus)

			// Block functionality
			protected.POST("/blocks/:user_id", socialHandler.BlockUser)
			protected.DELETE("/blocks/:user_id", socialHandler.UnblockUser)
			protected.GET("/blocks", socialHandler.GetBlockedUsers)
			protected.GET("/blocks/status/:user_id", socialHandler.CheckBlockStatus)

			// Notification functionality
			protected.GET("/notifications", socialHandler.GetNotifications)
			protected.GET("/notifications/unread", socialHandler.GetUnreadNotifications)
			protected.GET("/notifications/unread/count", socialHandler.GetUnreadCount)
			protected.PUT("/notifications/:id/read", socialHandler.MarkAsRead)
			protected.PUT("/notifications/read_all", socialHandler.MarkAllAsRead)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired(cfg))
		admin.Use(middleware.AdminRequired())
		{
			// Add admin-specific routes here
			// admin.GET("/users", adminHandler.ListAllUsers)
			// admin.DELETE("/users/:id", adminHandler.DeleteUser)
			// admin.DELETE("/items/:id", adminHandler.DeleteItem)
			// admin.DELETE("/coordinates/:id", adminHandler.DeleteCoordinate)
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Route not found"})
	})

	return r
}

// SetupDevRouter sets up router with additional development routes
func SetupDevRouter(usecases *usecase.Container, repos *repository.Container, cfg *config.Config) *gin.Engine {
	r := SetupRouter(usecases, repos, cfg)

	// Development-only routes
	dev := r.Group("/dev")
	{
		// Database seed endpoint
		dev.POST("/seed", func(c *gin.Context) {
			// TODO: Implement database seeding
			c.JSON(200, gin.H{"message": "Database seeded"})
		})

		// Clear cache endpoint
		dev.POST("/clear-cache", func(c *gin.Context) {
			// TODO: Implement cache clearing
			c.JSON(200, gin.H{"message": "Cache cleared"})
		})
	}

	return r
}