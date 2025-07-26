package impl

import (
	"context"
	"errors"
	"testing"

	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/dto"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/testutil"
	"github.com/House-lovers7/speadwear-go/pkg/config"
	"github.com/House-lovers7/speadwear-go/pkg/utils"
)

func setupUserUsecase(t *testing.T) (*userUsecase, *testutil.Fixtures) {
	db := testutil.TestDB(t)
	
	repos := repository.NewContainer(db)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: "24h",
		},
		Upload: config.UploadConfig{
			Path:        "./test-uploads",
			MaxFileSize: 5 * 1024 * 1024,
		},
	}
	
	usecase := NewUserUsecase(
		repos.User,
		cfg,
	).(*userUsecase)
	
	fixtures := testutil.NewFixtures(t, db)
	
	return usecase, fixtures
}

func TestUserUsecase_Signup(t *testing.T) {
	usecase, _ := setupUserUsecase(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *dto.SignupRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid registration",
			req: &dto.SignupRequest{
				Name:     "New User",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &dto.SignupRequest{
				Name:     "Another User",
				Email:    "newuser@example.com", // Same as above
				Password: "password456",
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
		{
			name: "empty email",
			req: &dto.SignupRequest{
				Name:     "User",
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "empty password",
			req: &dto.SignupRequest{
				Name:     "User",
				Email:    "another@example.com",
				Password: "",
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "short password",
			req: &dto.SignupRequest{
				Name:     "User",
				Email:    "short@example.com",
				Password: "123",
			},
			wantErr: true,
			errMsg:  "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := usecase.Signup(ctx, tt.req)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Signup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("Signup() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr {
				if resp == nil {
					t.Error("Signup() returned nil response")
				}
				if resp != nil && resp.Token == "" {
					t.Error("Signup() returned empty token")
				}
				if resp != nil && resp.User.Email != tt.req.Email {
					t.Errorf("Signup() user email = %v, want %v", resp.User.Email, tt.req.Email)
				}
			}
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	usecase, fixtures := setupUserUsecase(t)
	ctx := context.Background()

	// Create test user
	plainPassword := "testpassword123"
	hashedPassword, _ := utils.HashPassword(plainPassword)
	testUser := fixtures.CreateUser(func(u *domain.User) {
		u.Email = "login@example.com"
		u.PasswordDigest = hashedPassword
		u.Activated = true
	})

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid login",
			email:    testUser.Email,
			password: plainPassword,
			wantErr:  false,
		},
		{
			name:     "wrong password",
			email:    testUser.Email,
			password: "wrongpassword",
			wantErr:  true,
			errMsg:   "invalid email or password",
		},
		{
			name:     "non-existent email",
			email:    "nonexistent@example.com",
			password: plainPassword,
			wantErr:  true,
			errMsg:   "invalid email or password",
		},
		{
			name:     "empty email",
			email:    "",
			password: plainPassword,
			wantErr:  true,
			errMsg:   "email is required",
		},
		{
			name:     "empty password",
			email:    testUser.Email,
			password: "",
			wantErr:  true,
			errMsg:   "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := usecase.Login(ctx, tt.email, tt.password)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("Login() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr {
				if resp == nil {
					t.Error("Login() returned nil response")
				}
				if resp != nil && resp.Token == "" {
					t.Error("Login() returned empty token")
				}
				if resp != nil && resp.User.Email != tt.email {
					t.Errorf("Login() user email = %v, want %v", resp.User.Email, tt.email)
				}
			}
		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	usecase, fixtures := setupUserUsecase(t)
	ctx := context.Background()

	// Create test user
	testUser := fixtures.CreateUser()

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "existing user",
			userID:  testUser.ID,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			userID:  99999,
			wantErr: true,
		},
		{
			name:    "zero user ID",
			userID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := usecase.GetUser(ctx, tt.userID)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if user == nil {
					t.Error("GetUser() returned nil user")
				}
				if user != nil && user.ID != tt.userID {
					t.Errorf("GetUser() user ID = %v, want %v", user.ID, tt.userID)
				}
			}
		})
	}
}

func TestUserUsecase_UpdateProfile(t *testing.T) {
	usecase, fixtures := setupUserUsecase(t)
	ctx := context.Background()

	// Create test user
	testUser := fixtures.CreateUser(func(u *domain.User) {
		u.Name = "Original Name"
		u.Email = "original@example.com"
	})

	tests := []struct {
		name        string
		userID      uint
		updatedName string
		picture     string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid update",
			userID:      testUser.ID,
			updatedName: "Updated Name",
			picture:     "/uploads/avatar.jpg",
			wantErr:     false,
		},
		{
			name:        "update name only",
			userID:      testUser.ID,
			updatedName: "Another Update",
			picture:     "",
			wantErr:     false,
		},
		{
			name:        "non-existent user",
			userID:      99999,
			updatedName: "Ghost User",
			picture:     "",
			wantErr:     true,
			errMsg:      "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.UpdateProfile(ctx, tt.userID, tt.updatedName, tt.picture)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("UpdateProfile() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr {
				// Verify update
				updated, _ := usecase.GetUser(ctx, tt.userID)
				if updated.Name != tt.updatedName {
					t.Errorf("UpdateProfile() did not update name: got %v, want %v", updated.Name, tt.updatedName)
				}
				if tt.picture != "" && updated.Picture != tt.picture {
					t.Errorf("UpdateProfile() did not update picture: got %v, want %v", updated.Picture, tt.picture)
				}
			}
		})
	}
}

func TestUserUsecase_ChangePassword(t *testing.T) {
	usecase, fixtures := setupUserUsecase(t)
	ctx := context.Background()

	// Create test user with known password
	currentPassword := "currentpassword123"
	hashedPassword, _ := utils.HashPassword(currentPassword)
	testUser := fixtures.CreateUser(func(u *domain.User) {
		u.PasswordDigest = hashedPassword
	})

	tests := []struct {
		name            string
		userID          uint
		currentPassword string
		newPassword     string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "valid password change",
			userID:          testUser.ID,
			currentPassword: currentPassword,
			newPassword:     "newpassword123",
			wantErr:         false,
		},
		{
			name:            "wrong current password",
			userID:          testUser.ID,
			currentPassword: "wrongpassword",
			newPassword:     "newpassword123",
			wantErr:         true,
			errMsg:          "current password is incorrect",
		},
		{
			name:            "short new password",
			userID:          testUser.ID,
			currentPassword: currentPassword,
			newPassword:     "short",
			wantErr:         true,
			errMsg:          "password must be at least 8 characters",
		},
		{
			name:            "non-existent user",
			userID:          99999,
			currentPassword: "anypassword",
			newPassword:     "newpassword123",
			wantErr:         true,
			errMsg:          "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.ChangePassword(ctx, tt.userID, tt.currentPassword, tt.newPassword)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("ChangePassword() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr {
				// Try logging in with new password
				_, err := usecase.Login(ctx, testUser.Email, tt.newPassword)
				if err != nil {
					t.Error("ChangePassword() failed: cannot login with new password")
				}
				
				// Try logging in with old password (should fail)
				_, err = usecase.Login(ctx, testUser.Email, currentPassword)
				if err == nil {
					t.Error("ChangePassword() failed: can still login with old password")
				}
			}
		})
	}
}


// Mock repository for testing error cases
type mockUserRepository struct {
	repository.UserRepository
	findByEmailErr error
	createErr      error
	existsByEmailResult bool
	existsByEmailErr    error
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	return nil, nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.createErr
}

func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.existsByEmailErr != nil {
		return false, m.existsByEmailErr
	}
	return m.existsByEmailResult, nil
}

func TestUserUsecase_SignupWithRepositoryError(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: "24h",
		},
	}
	
	mockRepo := &mockUserRepository{
		existsByEmailErr: errors.New("database error"),
	}
	
	usecase := &userUsecase{
		userRepo: mockRepo,
		config:   cfg,
	}
	
	ctx := context.Background()
	_, err := usecase.Signup(ctx, &dto.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})
	
	if err == nil {
		t.Error("Signup() should return error when repository fails")
	}
}