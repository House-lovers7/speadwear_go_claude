package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/dto"
	"github.com/House-lovers7/speadwear-go/pkg/config"
)

// Mock usecase
type mockUserUsecase struct {
	mock.Mock
}

func (m *mockUserUsecase) Login(ctx context.Context, email, password string) (*dto.AuthResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *mockUserUsecase) Signup(ctx context.Context, req *dto.SignupRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *mockUserUsecase) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	args := m.Called(ctx, tokenString)
	return args.String(0), args.Error(1)
}

func (m *mockUserUsecase) GetUser(ctx context.Context, userID uint) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserUsecase) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *mockUserUsecase) DeleteUser(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserUsecase) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *mockUserUsecase) ResetPasswordRequest(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *mockUserUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

func (m *mockUserUsecase) ActivateAccount(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockUserUsecase) ResendActivationEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *mockUserUsecase) UpdateProfile(ctx context.Context, userID uint, name, picture string) error {
	args := m.Called(ctx, userID, name, picture)
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		requestBody  dto.LoginRequest
		mockSetup    func(*mockUserUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful login",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("Login", mock.Anything, "test@example.com", "password123").Return(&dto.AuthResponse{
					Token:     "test-token",
					ExpiresAt: time.Now().Add(24 * time.Hour),
					User: dto.UserResponse{
						ID:    1,
						Name:  "Test User",
						Email: "test@example.com",
					},
				}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "test-token", body["token"])
				user := body["user"].(map[string]interface{})
				assert.Equal(t, float64(1), user["id"])
				assert.Equal(t, "Test User", user["name"])
				assert.Equal(t, "test@example.com", user["email"])
			},
		},
		{
			name: "invalid credentials",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("Login", mock.Anything, "test@example.com", "wrongpassword").Return(nil, errors.New("invalid email or password"))
			},
			expectedCode: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "invalid email or password", body["error"])
			},
		},
		{
			name: "invalid request body",
			requestBody: dto.LoginRequest{
				Email: "", // Empty email
			},
			mockSetup:    func(m *mockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.NotNil(t, body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockUserUsecase)
			tt.mockSetup(mockUsecase)
			
			cfg := &config.Config{
				JWT: config.JWTConfig{
					Secret:     "test-secret",
					Expiration: "24h",
				},
			}
			
			handler := NewAuthHandler(cfg, mockUsecase)
			
			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			
			// Execute
			handler.Login(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Signup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		requestBody  dto.SignupRequest
		mockSetup    func(*mockUserUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful signup",
			requestBody: dto.SignupRequest{
				Name:     "New User",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("Signup", mock.Anything, &dto.SignupRequest{
					Name:     "New User",
					Email:    "newuser@example.com",
					Password: "password123",
				}).Return(&dto.AuthResponse{
					Token:     "new-token",
					ExpiresAt: time.Now().Add(24 * time.Hour),
					User: dto.UserResponse{
						ID:    2,
						Name:  "New User",
						Email: "newuser@example.com",
					},
				}, nil)
			},
			expectedCode: http.StatusCreated,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "new-token", body["token"])
				user := body["user"].(map[string]interface{})
				assert.Equal(t, float64(2), user["id"])
				assert.Equal(t, "New User", user["name"])
			},
		},
		{
			name: "email already exists",
			requestBody: dto.SignupRequest{
				Name:     "Duplicate User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("Signup", mock.Anything, &dto.SignupRequest{
					Name:     "Duplicate User",
					Email:    "existing@example.com",
					Password: "password123",
				}).Return(nil, errors.New("email already exists"))
			},
			expectedCode: http.StatusConflict,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "email already exists", body["error"])
			},
		},
		{
			name: "invalid request body",
			requestBody: dto.SignupRequest{
				Name:  "User",
				Email: "invalid-email", // Invalid email format
			},
			mockSetup:    func(m *mockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.NotNil(t, body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockUserUsecase)
			tt.mockSetup(mockUsecase)
			
			cfg := &config.Config{
				JWT: config.JWTConfig{
					Secret:     "test-secret",
					Expiration: "24h",
				},
			}
			
			handler := NewAuthHandler(cfg, mockUsecase)
			
			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			
			// Execute
			handler.Signup(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		token        string
		mockSetup    func(*mockUserUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:  "successful token refresh",
			token: "old-token",
			mockSetup: func(m *mockUserUsecase) {
				m.On("RefreshToken", mock.Anything, "old-token").Return("new-refreshed-token", nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "new-refreshed-token", body["token"])
				assert.NotNil(t, body["expires_at"])
			},
		},
		{
			name:  "invalid token",
			token: "invalid-token",
			mockSetup: func(m *mockUserUsecase) {
				m.On("RefreshToken", mock.Anything, "invalid-token").Return("", errors.New("invalid token"))
			},
			expectedCode: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid token", body["error"])
			},
		},
		{
			name:         "missing authorization header",
			token:        "",
			mockSetup:    func(m *mockUserUsecase) {},
			expectedCode: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Token required", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockUserUsecase)
			tt.mockSetup(mockUsecase)
			
			cfg := &config.Config{
				JWT: config.JWTConfig{
					Secret:     "test-secret",
					Expiration: "24h",
				},
			}
			
			handler := NewAuthHandler(cfg, mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			
			// Execute
			handler.RefreshToken(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}