package handler

import (
	"bytes"
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
)

func TestUserHandler_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       string
		mockSetup    func(*mockUserUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get user",
			userID: "1",
			mockSetup: func(m *mockUserUsecase) {
				m.On("GetUser", mock.Anything, uint(1)).Return(&domain.User{
					BaseModel: domain.BaseModel{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Name:      "Test User",
					Email:     "test@example.com",
					Picture:   "/uploads/avatar.jpg",
					Admin:     false,
					Activated: true,
				}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(1), body["id"])
				assert.Equal(t, "Test User", body["name"])
				assert.Equal(t, "test@example.com", body["email"])
				assert.Equal(t, "/uploads/avatar.jpg", body["picture"])
				assert.Equal(t, false, body["admin"])
				assert.Equal(t, true, body["activated"])
			},
		},
		{
			name:   "user not found",
			userID: "999",
			mockSetup: func(m *mockUserUsecase) {
				m.On("GetUser", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "user not found", body["error"])
			},
		},
		{
			name:         "invalid user ID",
			userID:       "invalid",
			mockSetup:    func(m *mockUserUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid user ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockUserUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewUserHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{
				{Key: "id", Value: tt.userID},
			}
			
			// Execute
			handler.GetUser(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		setUserID    bool
		mockSetup    func(*mockUserUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:      "successful get current user",
			userID:    1,
			setUserID: true,
			mockSetup: func(m *mockUserUsecase) {
				m.On("GetUser", mock.Anything, uint(1)).Return(&domain.User{
					BaseModel: domain.BaseModel{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Name:      "Current User",
					Email:     "current@example.com",
					Picture:   "/uploads/me.jpg",
					Admin:     true,
					Activated: true,
				}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(1), body["id"])
				assert.Equal(t, "Current User", body["name"])
				assert.Equal(t, "current@example.com", body["email"])
				assert.Equal(t, "/uploads/me.jpg", body["picture"])
				assert.Equal(t, true, body["admin"])
			},
		},
		{
			name:         "user ID not in context",
			userID:       0,
			setUserID:    false,
			mockSetup:    func(m *mockUserUsecase) {},
			expectedCode: http.StatusUnauthorized,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "User not authenticated", body["error"])
			},
		},
		{
			name:      "user not found",
			userID:    999,
			setUserID: true,
			mockSetup: func(m *mockUserUsecase) {
				m.On("GetUser", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "user not found", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockUserUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewUserHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			if tt.setUserID {
				c.Set("userID", tt.userID)
			}
			
			// Execute
			handler.GetMe(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		requestBody  map[string]string
		mockSetup    func(*mockUserUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful profile update",
			userID: 1,
			requestBody: map[string]string{
				"name":    "Updated Name",
				"picture": "/uploads/new-avatar.jpg",
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("UpdateProfile", mock.Anything, uint(1), "Updated Name", "/uploads/new-avatar.jpg").Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Profile updated successfully", body["message"])
			},
		},
		{
			name:   "update name only",
			userID: 1,
			requestBody: map[string]string{
				"name": "Only Name Update",
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("UpdateProfile", mock.Anything, uint(1), "Only Name Update", "").Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Profile updated successfully", body["message"])
			},
		},
		{
			name:   "user not found",
			userID: 1,
			requestBody: map[string]string{
				"name": "Ghost User",
			},
			mockSetup: func(m *mockUserUsecase) {
				m.On("UpdateProfile", mock.Anything, uint(1), "Ghost User", "").Return(errors.New("user not found"))
			},
			expectedCode: http.StatusInternalServerError,
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
			
			handler := NewUserHandler(mockUsecase)
			
			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.UpdateProfile(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}