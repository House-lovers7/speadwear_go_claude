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
)

// Mock social usecase
type mockSocialUsecase struct {
	mock.Mock
}

func (m *mockSocialUsecase) CreateComment(ctx context.Context, userID, coordinateID uint, comment string) (*domain.Comment, error) {
	args := m.Called(ctx, userID, coordinateID, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *mockSocialUsecase) GetComments(ctx context.Context, coordinateID uint, limit, offset int) ([]*domain.Comment, int64, error) {
	args := m.Called(ctx, coordinateID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *mockSocialUsecase) UpdateComment(ctx context.Context, userID, commentID uint, comment string) error {
	args := m.Called(ctx, userID, commentID, comment)
	return args.Error(0)
}

func (m *mockSocialUsecase) DeleteComment(ctx context.Context, userID, commentID uint) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}

func (m *mockSocialUsecase) FollowUser(ctx context.Context, followerID, followedID uint) error {
	args := m.Called(ctx, followerID, followedID)
	return args.Error(0)
}

func (m *mockSocialUsecase) UnfollowUser(ctx context.Context, followerID, followedID uint) error {
	args := m.Called(ctx, followerID, followedID)
	return args.Error(0)
}

func (m *mockSocialUsecase) GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockSocialUsecase) GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockSocialUsecase) IsFollowing(ctx context.Context, followerID, followedID uint) (bool, error) {
	args := m.Called(ctx, followerID, followedID)
	return args.Bool(0), args.Error(1)
}

func (m *mockSocialUsecase) BlockUser(ctx context.Context, blockerID, blockedID uint) error {
	args := m.Called(ctx, blockerID, blockedID)
	return args.Error(0)
}

func (m *mockSocialUsecase) UnblockUser(ctx context.Context, blockerID, blockedID uint) error {
	args := m.Called(ctx, blockerID, blockedID)
	return args.Error(0)
}

func (m *mockSocialUsecase) GetBlockedUsers(ctx context.Context, userID uint) ([]*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *mockSocialUsecase) IsBlocked(ctx context.Context, blockerID, blockedID uint) (bool, error) {
	args := m.Called(ctx, blockerID, blockedID)
	return args.Bool(0), args.Error(1)
}

func (m *mockSocialUsecase) GetNotifications(ctx context.Context, userID uint, limit, offset int) ([]*domain.Notification, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *mockSocialUsecase) GetUnreadNotifications(ctx context.Context, userID uint) ([]*domain.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Notification), args.Error(1)
}

func (m *mockSocialUsecase) MarkNotificationAsRead(ctx context.Context, userID, notificationID uint) error {
	args := m.Called(ctx, userID, notificationID)
	return args.Error(0)
}

func (m *mockSocialUsecase) MarkAllNotificationsAsRead(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockSocialUsecase) GetUnreadNotificationCount(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockSocialUsecase) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func TestSocialHandler_CreateComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		requestBody  map[string]interface{}
		mockSetup    func(*mockSocialUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful comment creation",
			userID: 1,
			requestBody: map[string]interface{}{
				"coordinate_id": uint(1),
				"comment":       "Nice outfit!",
			},
			mockSetup: func(m *mockSocialUsecase) {
				m.On("CreateComment", mock.Anything, uint(1), uint(1), "Nice outfit!").Return(&domain.Comment{
					BaseModel: domain.BaseModel{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					UserID:       1,
					CoordinateID: 1,
					Comment:      "Nice outfit!",
				}, nil)
			},
			expectedCode: http.StatusCreated,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(1), body["id"])
				assert.Equal(t, float64(1), body["user_id"])
				assert.Equal(t, float64(1), body["coordinate_id"])
				assert.Equal(t, "Nice outfit!", body["comment"])
			},
		},
		{
			name:   "blocked by user",
			userID: 1,
			requestBody: map[string]interface{}{
				"coordinate_id": uint(2),
				"comment":       "Great style!",
			},
			mockSetup: func(m *mockSocialUsecase) {
				m.On("CreateComment", mock.Anything, uint(1), uint(2), "Great style!").Return(nil, errors.New("you are blocked by this user"))
			},
			expectedCode: http.StatusForbidden,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "you are blocked by this user", body["error"])
			},
		},
		{
			name:   "empty comment",
			userID: 1,
			requestBody: map[string]interface{}{
				"coordinate_id": uint(1),
				"comment":       "",
			},
			mockSetup:    func(m *mockSocialUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.NotNil(t, body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.CreateComment(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSocialHandler_FollowUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		followUserID string
		mockSetup    func(*mockSocialUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:         "successful follow",
			userID:       1,
			followUserID: "2",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("FollowUser", mock.Anything, uint(1), uint(2)).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Followed successfully", body["message"])
			},
		},
		{
			name:         "follow yourself",
			userID:       1,
			followUserID: "1",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("FollowUser", mock.Anything, uint(1), uint(1)).Return(errors.New("cannot follow yourself"))
			},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "cannot follow yourself", body["error"])
			},
		},
		{
			name:         "already following",
			userID:       1,
			followUserID: "2",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("FollowUser", mock.Anything, uint(1), uint(2)).Return(errors.New("already following"))
			},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "already following", body["error"])
			},
		},
		{
			name:         "blocked by user",
			userID:       1,
			followUserID: "3",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("FollowUser", mock.Anything, uint(1), uint(3)).Return(errors.New("you are blocked by this user"))
			},
			expectedCode: http.StatusForbidden,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "you are blocked by this user", body["error"])
			},
		},
		{
			name:         "invalid user ID",
			userID:       1,
			followUserID: "invalid",
			mockSetup:    func(m *mockSocialUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid user ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/follow/"+tt.followUserID, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			c.Params = gin.Params{
				{Key: "user_id", Value: tt.followUserID},
			}
			
			// Execute
			handler.FollowUser(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSocialHandler_GetFollowers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		query        string
		mockSetup    func(*mockSocialUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get followers",
			userID: 1,
			query:  "page=1&per_page=10",
			mockSetup: func(m *mockSocialUsecase) {
				users := []*domain.User{
					{
						BaseModel: domain.BaseModel{
							ID:        2,
							CreatedAt: time.Now(),
						},
						Name:      "Follower 1",
						Email:     "follower1@example.com",
						Picture:   "/uploads/follower1.jpg",
						Activated: true,
					},
					{
						BaseModel: domain.BaseModel{
							ID:        3,
							CreatedAt: time.Now(),
						},
						Name:      "Follower 2",
						Email:     "follower2@example.com",
						Picture:   "/uploads/follower2.jpg",
						Activated: true,
					},
				}
				m.On("GetFollowers", mock.Anything, uint(1), 10, 0).Return(users, int64(2), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				users := body["users"].([]interface{})
				assert.Len(t, users, 2)
				assert.Equal(t, float64(2), body["total_count"])
				assert.Equal(t, float64(1), body["page"])
				assert.Equal(t, float64(10), body["per_page"])
			},
		},
		{
			name:   "get followers for specific user",
			userID: 1,
			query:  "user_id=5&page=1&per_page=10",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("GetFollowers", mock.Anything, uint(5), 10, 0).Return([]*domain.User{}, int64(0), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				users := body["users"].([]interface{})
				assert.Len(t, users, 0)
				assert.Equal(t, float64(0), body["total_count"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/follow/followers?"+tt.query, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.GetFollowers(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSocialHandler_BlockUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		blockUserID  string
		mockSetup    func(*mockSocialUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:        "successful block",
			userID:      1,
			blockUserID: "2",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("BlockUser", mock.Anything, uint(1), uint(2)).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "User blocked successfully", body["message"])
			},
		},
		{
			name:        "block yourself",
			userID:      1,
			blockUserID: "1",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("BlockUser", mock.Anything, uint(1), uint(1)).Return(errors.New("cannot block yourself"))
			},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "cannot block yourself", body["error"])
			},
		},
		{
			name:        "already blocked",
			userID:      1,
			blockUserID: "2",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("BlockUser", mock.Anything, uint(1), uint(2)).Return(errors.New("already blocked"))
			},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "already blocked", body["error"])
			},
		},
		{
			name:         "invalid user ID",
			userID:       1,
			blockUserID:  "invalid",
			mockSetup:    func(m *mockSocialUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid user ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/blocks/"+tt.blockUserID, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			c.Params = gin.Params{
				{Key: "user_id", Value: tt.blockUserID},
			}
			
			// Execute
			handler.BlockUser(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSocialHandler_GetNotifications(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		query        string
		mockSetup    func(*mockSocialUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get notifications",
			userID: 1,
			query:  "page=1&per_page=10",
			mockSetup: func(m *mockSocialUsecase) {
				notifications := []*domain.Notification{
					{
						BaseModel: domain.BaseModel{
							ID:        1,
							CreatedAt: time.Now(),
						},
						SenderID:   2,
						ReceiverID: 1,
						Action:     domain.NotificationActionFollow,
						Checked:    false,
						Sender: domain.User{
							BaseModel: domain.BaseModel{ID: 2},
							Name:      "Follower",
							Email:     "follower@example.com",
							Picture:   "/uploads/follower.jpg",
						},
					},
					{
						BaseModel: domain.BaseModel{
							ID:        2,
							CreatedAt: time.Now(),
						},
						SenderID:   3,
						ReceiverID: 1,
						Action:     domain.NotificationActionComment,
						Checked:    true,
						Sender: domain.User{
							BaseModel: domain.BaseModel{ID: 3},
							Name:      "Commenter",
							Email:     "commenter@example.com",
							Picture:   "/uploads/commenter.jpg",
						},
						Coordinate: &domain.Coordinate{
							BaseModel: domain.BaseModel{ID: 1},
							UserID:    1,
							Picture:   "/uploads/coord.jpg",
						},
						Comment: &domain.Comment{
							BaseModel: domain.BaseModel{ID: 1},
							Comment:   "Nice outfit!",
						},
					},
				}
				m.On("GetNotifications", mock.Anything, uint(1), 10, 0).Return(notifications, int64(2), nil)
				m.On("GetUnreadNotificationCount", mock.Anything, uint(1)).Return(int64(1), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				notifications := body["notifications"].([]interface{})
				assert.Len(t, notifications, 2)
				assert.Equal(t, float64(2), body["total_count"])
				assert.Equal(t, float64(1), body["unread_count"])
				assert.Equal(t, float64(1), body["page"])
				assert.Equal(t, float64(10), body["per_page"])
				
				// Check first notification details
				firstNotif := notifications[0].(map[string]interface{})
				assert.Equal(t, float64(1), firstNotif["id"])
				assert.Equal(t, float64(2), firstNotif["sender_id"])
				assert.Equal(t, float64(1), firstNotif["receiver_id"])
				assert.Equal(t, domain.NotificationActionFollow, firstNotif["action"])
				assert.Equal(t, false, firstNotif["checked"])
				
				// Check sender info
				sender := firstNotif["sender"].(map[string]interface{})
				assert.Equal(t, float64(2), sender["id"])
				assert.Equal(t, "Follower", sender["name"])
			},
		},
		{
			name:   "empty notifications",
			userID: 1,
			query:  "page=1&per_page=10",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("GetNotifications", mock.Anything, uint(1), 10, 0).Return([]*domain.Notification{}, int64(0), nil)
				m.On("GetUnreadNotificationCount", mock.Anything, uint(1)).Return(int64(0), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				notifications := body["notifications"].([]interface{})
				assert.Len(t, notifications, 0)
				assert.Equal(t, float64(0), body["total_count"])
				assert.Equal(t, float64(0), body["unread_count"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications?"+tt.query, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.GetNotifications(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSocialHandler_MarkAsRead(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		userID         uint
		notificationID string
		mockSetup      func(*mockSocialUsecase)
		expectedCode   int
		checkBody      func(*testing.T, map[string]interface{})
	}{
		{
			name:           "successful mark as read",
			userID:         1,
			notificationID: "1",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("MarkNotificationAsRead", mock.Anything, uint(1), uint(1)).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Notification marked as read", body["message"])
			},
		},
		{
			name:           "unauthorized",
			userID:         1,
			notificationID: "2",
			mockSetup: func(m *mockSocialUsecase) {
				m.On("MarkNotificationAsRead", mock.Anything, uint(1), uint(2)).Return(errors.New("unauthorized"))
			},
			expectedCode: http.StatusForbidden,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "unauthorized", body["error"])
			},
		},
		{
			name:           "invalid notification ID",
			userID:         1,
			notificationID: "invalid",
			mockSetup:      func(m *mockSocialUsecase) {},
			expectedCode:   http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid notification ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/"+tt.notificationID+"/read", nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			c.Params = gin.Params{
				{Key: "id", Value: tt.notificationID},
			}
			
			// Execute
			handler.MarkAsRead(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSocialHandler_MarkAllAsRead(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		mockSetup    func(*mockSocialUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful mark all as read",
			userID: 1,
			mockSetup: func(m *mockSocialUsecase) {
				m.On("MarkAllNotificationsAsRead", mock.Anything, uint(1)).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "All notifications marked as read", body["message"])
			},
		},
		{
			name:   "error marking all as read",
			userID: 1,
			mockSetup: func(m *mockSocialUsecase) {
				m.On("MarkAllNotificationsAsRead", mock.Anything, uint(1)).Return(errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "database error", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/read_all", nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.MarkAllAsRead(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSocialHandler_GetUnreadCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		mockSetup    func(*mockSocialUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get unread count",
			userID: 1,
			mockSetup: func(m *mockSocialUsecase) {
				m.On("GetUnreadNotificationCount", mock.Anything, uint(1)).Return(int64(5), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(5), body["unread_count"])
			},
		},
		{
			name:   "no unread notifications",
			userID: 1,
			mockSetup: func(m *mockSocialUsecase) {
				m.On("GetUnreadNotificationCount", mock.Anything, uint(1)).Return(int64(0), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(0), body["unread_count"])
			},
		},
		{
			name:   "error getting count",
			userID: 1,
			mockSetup: func(m *mockSocialUsecase) {
				m.On("GetUnreadNotificationCount", mock.Anything, uint(1)).Return(int64(0), errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "database error", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockSocialUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewSocialHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/unread/count", nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.GetUnreadCount(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}