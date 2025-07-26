package usecase

import (
	"context"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/dto"
)

// UserUsecase defines user-related business logic
type UserUsecase interface {
	// Authentication
	Login(ctx context.Context, email, password string) (*dto.AuthResponse, error)
	Signup(ctx context.Context, req *dto.SignupRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
	
	// User management
	GetUser(ctx context.Context, userID uint) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID uint) error
	
	// User listing
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
	
	// Password management
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	ResetPasswordRequest(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	
	// Account activation
	ActivateAccount(ctx context.Context, token string) error
	ResendActivationEmail(ctx context.Context, email string) error
	
	// Profile
	UpdateProfile(ctx context.Context, userID uint, name, picture string) error
}