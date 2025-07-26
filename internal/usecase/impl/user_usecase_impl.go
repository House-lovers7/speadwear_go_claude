package impl

import (
	"context"
	"errors"
	"time"
	
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/dto"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
	"github.com/House-lovers7/speadwear-go/pkg/config"
	"github.com/House-lovers7/speadwear-go/pkg/utils"
)

type userUsecase struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository, config *config.Config) usecase.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		config:   config,
	}
}

// Login handles user login
func (u *userUsecase) Login(ctx context.Context, email, password string) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}
	
	// Check password
	if !utils.CheckPassword(password, user.PasswordDigest) {
		return nil, errors.New("invalid email or password")
	}
	
	// Check if account is activated
	if !user.Activated {
		return nil, errors.New("account not activated")
	}
	
	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, u.config.JWT.Secret, u.config.JWT.Expiration)
	if err != nil {
		return nil, err
	}
	
	duration, _ := time.ParseDuration(u.config.JWT.Expiration)
	expiresAt := time.Now().Add(duration)
	
	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Picture:   user.Picture,
			Admin:     user.Admin,
			Activated: user.Activated,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// Signup handles user registration
func (u *userUsecase) Signup(ctx context.Context, req *dto.SignupRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}
	
	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	
	// Create user
	now := time.Now()
	user := &domain.User{
		Name:           req.Name,
		Email:          req.Email,
		PasswordDigest: hashedPassword,
		Activated:      true, // TODO: Implement email activation
		ActivatedAt:    &now,
	}
	
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	
	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, u.config.JWT.Secret, u.config.JWT.Expiration)
	if err != nil {
		return nil, err
	}
	
	duration, _ := time.ParseDuration(u.config.JWT.Expiration)
	expiresAt := time.Now().Add(duration)
	
	return &dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Picture:   user.Picture,
			Admin:     user.Admin,
			Activated: user.Activated,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// RefreshToken refreshes JWT token
func (u *userUsecase) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	newToken, err := utils.RefreshToken(tokenString, u.config.JWT.Secret, u.config.JWT.Expiration)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

// GetUser gets user by ID
func (u *userUsecase) GetUser(ctx context.Context, userID uint) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByEmail gets user by email
func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// UpdateUser updates user information
func (u *userUsecase) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	// Apply updates
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if email, ok := updates["email"].(string); ok {
		// Check if new email already exists
		if email != user.Email {
			exists, err := u.userRepo.ExistsByEmail(ctx, email)
			if err != nil {
				return err
			}
			if exists {
				return errors.New("email already exists")
			}
		}
		user.Email = email
	}
	if picture, ok := updates["picture"].(string); ok {
		user.Picture = picture
	}
	
	return u.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user
func (u *userUsecase) DeleteUser(ctx context.Context, userID uint) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	return u.userRepo.Delete(ctx, userID)
}

// ListUsers lists all users with pagination
func (u *userUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	users, err := u.userRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	
	return users, count, nil
}

// ChangePassword changes user password
func (u *userUsecase) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	// Check old password
	if !utils.CheckPassword(oldPassword, user.PasswordDigest) {
		return errors.New("invalid old password")
	}
	
	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	user.PasswordDigest = hashedPassword
	return u.userRepo.Update(ctx, user)
}

// ResetPasswordRequest initiates password reset
func (u *userUsecase) ResetPasswordRequest(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	// TODO: Generate reset token and send email
	// For now, just return nil
	return nil
}

// ResetPassword resets password with token
func (u *userUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// TODO: Validate reset token and update password
	// For now, just return nil
	return nil
}

// ActivateAccount activates user account
func (u *userUsecase) ActivateAccount(ctx context.Context, token string) error {
	// TODO: Validate activation token and activate account
	// For now, just return nil
	return nil
}

// ResendActivationEmail resends activation email
func (u *userUsecase) ResendActivationEmail(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	if user.Activated {
		return errors.New("account already activated")
	}
	
	// TODO: Generate activation token and send email
	// For now, just return nil
	return nil
}

// UpdateProfile updates user profile
func (u *userUsecase) UpdateProfile(ctx context.Context, userID uint, name, picture string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	
	if name != "" {
		user.Name = name
	}
	if picture != "" {
		user.Picture = picture
	}
	
	return u.userRepo.Update(ctx, user)
}