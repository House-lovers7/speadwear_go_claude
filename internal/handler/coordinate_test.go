package handler

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/repository"
)

// Mock usecase
type mockCoordinateUsecase struct {
	mock.Mock
}

func (m *mockCoordinateUsecase) CreateCoordinate(ctx context.Context, userID uint, coordinate *domain.Coordinate, itemIDs []uint, image *multipart.FileHeader) error {
	args := m.Called(ctx, userID, coordinate, itemIDs, image)
	return args.Error(0)
}

func (m *mockCoordinateUsecase) GetCoordinate(ctx context.Context, coordinateID uint) (*domain.Coordinate, error) {
	args := m.Called(ctx, coordinateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Coordinate), args.Error(1)
}

func (m *mockCoordinateUsecase) GetCoordinateWithDetails(ctx context.Context, coordinateID uint) (*domain.Coordinate, error) {
	args := m.Called(ctx, coordinateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Coordinate), args.Error(1)
}

func (m *mockCoordinateUsecase) UpdateCoordinate(ctx context.Context, userID uint, coordinateID uint, updates map[string]interface{}, itemIDs []uint, image *multipart.FileHeader) error {
	args := m.Called(ctx, userID, coordinateID, updates, itemIDs, image)
	return args.Error(0)
}

func (m *mockCoordinateUsecase) DeleteCoordinate(ctx context.Context, userID uint, coordinateID uint) error {
	args := m.Called(ctx, userID, coordinateID)
	return args.Error(0)
}

func (m *mockCoordinateUsecase) GetUserCoordinates(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Coordinate), args.Get(1).(int64), args.Error(2)
}

func (m *mockCoordinateUsecase) GetTimelineCoordinates(ctx context.Context, userID uint, limit, offset int) ([]*domain.Coordinate, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Coordinate), args.Error(1)
}

func (m *mockCoordinateUsecase) SearchCoordinates(ctx context.Context, filters repository.CoordinateFilter) ([]*domain.Coordinate, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Coordinate), args.Error(1)
}

func (m *mockCoordinateUsecase) LikeCoordinate(ctx context.Context, userID uint, coordinateID uint) error {
	args := m.Called(ctx, userID, coordinateID)
	return args.Error(0)
}

func (m *mockCoordinateUsecase) UnlikeCoordinate(ctx context.Context, userID uint, coordinateID uint) error {
	args := m.Called(ctx, userID, coordinateID)
	return args.Error(0)
}

func (m *mockCoordinateUsecase) IsLikedByUser(ctx context.Context, userID uint, coordinateID uint) (bool, error) {
	args := m.Called(ctx, userID, coordinateID)
	return args.Bool(0), args.Error(1)
}

func (m *mockCoordinateUsecase) GetUserCoordinateStatistics(ctx context.Context, userID uint) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockCoordinateUsecase) GetCoordinateLikes(ctx context.Context, coordinateID uint) ([]*domain.LikeCoordinate, error) {
	args := m.Called(ctx, coordinateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LikeCoordinate), args.Error(1)
}

// Mock repositories
type mockCommentRepository struct {
	mock.Mock
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *mockCommentRepository) FindByID(ctx context.Context, id uint) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *mockCommentRepository) FindByCoordinateID(ctx context.Context, coordinateID uint, limit, offset int) ([]*domain.Comment, error) {
	args := m.Called(ctx, coordinateID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Comment), args.Error(1)
}

func (m *mockCommentRepository) CountByCoordinateID(ctx context.Context, coordinateID uint) (int64, error) {
	args := m.Called(ctx, coordinateID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *mockCommentRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockLikeCoordinateRepository struct {
	mock.Mock
}

func (m *mockLikeCoordinateRepository) Create(ctx context.Context, like *domain.LikeCoordinate) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *mockLikeCoordinateRepository) FindByID(ctx context.Context, id uint) (*domain.LikeCoordinate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LikeCoordinate), args.Error(1)
}

func (m *mockLikeCoordinateRepository) Update(ctx context.Context, like *domain.LikeCoordinate) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *mockLikeCoordinateRepository) FindByCoordinateID(ctx context.Context, coordinateID uint) ([]*domain.LikeCoordinate, error) {
	args := m.Called(ctx, coordinateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LikeCoordinate), args.Error(1)
}

func (m *mockLikeCoordinateRepository) FindByUserAndCoordinate(ctx context.Context, userID, coordinateID uint) (*domain.LikeCoordinate, error) {
	args := m.Called(ctx, userID, coordinateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LikeCoordinate), args.Error(1)
}

func (m *mockLikeCoordinateRepository) CountByCoordinateID(ctx context.Context, coordinateID uint) (int64, error) {
	args := m.Called(ctx, coordinateID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockLikeCoordinateRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockLikeCoordinateRepository) ExistsByUserAndCoordinate(ctx context.Context, userID, coordinateID uint) (bool, error) {
	args := m.Called(ctx, userID, coordinateID)
	return args.Bool(0), args.Error(1)
}

func TestCoordinateHandler_GetCoordinate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name          string
		coordinateID  string
		userID        uint
		setUserID     bool
		mockSetup     func(*mockCoordinateUsecase, *mockLikeCoordinateRepository, *mockCommentRepository)
		expectedCode  int
		checkBody     func(*testing.T, map[string]interface{})
	}{
		{
			name:         "successful get coordinate",
			coordinateID: "1",
			userID:       1,
			setUserID:    true,
			mockSetup: func(m *mockCoordinateUsecase, likeRepo *mockLikeCoordinateRepository, commentRepo *mockCommentRepository) {
				coordinate := &domain.Coordinate{
					BaseModel: domain.BaseModel{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					UserID:  1,
					Season:  domain.SeasonSummer,
					TPO:     domain.TPOCasual,
					Picture: "/uploads/coordinate.jpg",
					Memo:    "Summer casual outfit",
					Rating:  5,
					Items: []domain.Item{
						{
							BaseModel: domain.BaseModel{
								ID: 1,
							},
							SuperItem: "Tシャツ",
							Color:     domain.ColorBlue,
						},
					},
					User: domain.User{
						BaseModel: domain.BaseModel{
							ID: 1,
						},
						Name:    "Test User",
						Email:   "test@example.com",
						Picture: "/uploads/avatar.jpg",
					},
				}
				m.On("GetCoordinateWithDetails", mock.Anything, uint(1)).Return(coordinate, nil)
				m.On("IsLikedByUser", mock.Anything, uint(1), uint(1)).Return(true, nil)
				likeRepo.On("CountByCoordinateID", mock.Anything, uint(1)).Return(int64(5), nil)
				commentRepo.On("CountByCoordinateID", mock.Anything, uint(1)).Return(int64(3), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(1), body["id"])
				assert.Equal(t, float64(1), body["user_id"])
				assert.Equal(t, float64(domain.SeasonSummer), body["season"])
				assert.Equal(t, float64(domain.TPOCasual), body["tpo"])
				assert.Equal(t, "/uploads/coordinate.jpg", body["picture"])
				assert.Equal(t, "Summer casual outfit", body["memo"])
				assert.Equal(t, float64(5), body["rating"])
				assert.Equal(t, float64(5), body["like_count"])
				assert.Equal(t, float64(3), body["comment_count"])
				assert.Equal(t, true, body["is_liked"])
				
				items := body["items"].([]interface{})
				assert.Len(t, items, 1)
				
				user := body["user"].(map[string]interface{})
				assert.Equal(t, float64(1), user["id"])
				assert.Equal(t, "Test User", user["name"])
			},
		},
		{
			name:         "coordinate not found",
			coordinateID: "999",
			userID:       1,
			setUserID:    true,
			mockSetup: func(m *mockCoordinateUsecase, likeRepo *mockLikeCoordinateRepository, commentRepo *mockCommentRepository) {
				m.On("GetCoordinateWithDetails", mock.Anything, uint(999)).Return(nil, errors.New("coordinate not found"))
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "coordinate not found", body["error"])
			},
		},
		{
			name:         "invalid coordinate ID",
			coordinateID: "invalid",
			userID:       1,
			setUserID:    true,
			mockSetup:    func(m *mockCoordinateUsecase, likeRepo *mockLikeCoordinateRepository, commentRepo *mockCommentRepository) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid coordinate ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockCoordinateUsecase)
			mockCommentRepo := new(mockCommentRepository)
			mockLikeRepo := new(mockLikeCoordinateRepository)
			
			tt.mockSetup(mockUsecase, mockLikeRepo, mockCommentRepo)
			
			handler := NewCoordinateHandler(mockUsecase, mockCommentRepo, mockLikeRepo)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/coordinates/"+tt.coordinateID, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{
				{Key: "id", Value: tt.coordinateID},
			}
			if tt.setUserID {
				c.Set("userID", tt.userID)
			}
			
			// Execute
			handler.GetCoordinate(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
			mockCommentRepo.AssertExpectations(t)
			mockLikeRepo.AssertExpectations(t)
		})
	}
}

func TestCoordinateHandler_GetMyCoordinates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		query        string
		mockSetup    func(*mockCoordinateUsecase, *mockLikeCoordinateRepository, *mockCommentRepository)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get my coordinates",
			userID: 1,
			query:  "page=1&per_page=10",
			mockSetup: func(m *mockCoordinateUsecase, likeRepo *mockLikeCoordinateRepository, commentRepo *mockCommentRepository) {
				coordinates := []*domain.Coordinate{
					{
						BaseModel: domain.BaseModel{
							ID:        1,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						UserID:  1,
						Season:  domain.SeasonSummer,
						TPO:     domain.TPOCasual,
						Picture: "/uploads/coord1.jpg",
						Rating:  5,
						User: domain.User{
							BaseModel: domain.BaseModel{ID: 1},
							Name:      "Test User",
						},
					},
					{
						BaseModel: domain.BaseModel{
							ID:        2,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						UserID:  1,
						Season:  domain.SeasonWinter,
						TPO:     domain.TPOWork,
						Picture: "/uploads/coord2.jpg",
						Rating:  4,
						User: domain.User{
							BaseModel: domain.BaseModel{ID: 1},
							Name:      "Test User",
						},
					},
				}
				m.On("GetUserCoordinates", mock.Anything, uint(1), 10, 0).Return(coordinates, int64(2), nil)
				m.On("IsLikedByUser", mock.Anything, uint(1), uint(1)).Return(false, nil)
				m.On("IsLikedByUser", mock.Anything, uint(1), uint(2)).Return(true, nil)
				likeRepo.On("CountByCoordinateID", mock.Anything, uint(1)).Return(int64(3), nil)
				likeRepo.On("CountByCoordinateID", mock.Anything, uint(2)).Return(int64(5), nil)
				commentRepo.On("CountByCoordinateID", mock.Anything, uint(1)).Return(int64(2), nil)
				commentRepo.On("CountByCoordinateID", mock.Anything, uint(2)).Return(int64(4), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				coordinates := body["coordinates"].([]interface{})
				assert.Len(t, coordinates, 2)
				assert.Equal(t, float64(2), body["total_count"])
				assert.Equal(t, float64(1), body["page"])
				assert.Equal(t, float64(10), body["per_page"])
			},
		},
		{
			name:   "empty result",
			userID: 1,
			query:  "page=1&per_page=10",
			mockSetup: func(m *mockCoordinateUsecase, likeRepo *mockLikeCoordinateRepository, commentRepo *mockCommentRepository) {
				m.On("GetUserCoordinates", mock.Anything, uint(1), 10, 0).Return([]*domain.Coordinate{}, int64(0), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				coordinates := body["coordinates"].([]interface{})
				assert.Len(t, coordinates, 0)
				assert.Equal(t, float64(0), body["total_count"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockCoordinateUsecase)
			mockCommentRepo := new(mockCommentRepository)
			mockLikeRepo := new(mockLikeCoordinateRepository)
			
			tt.mockSetup(mockUsecase, mockLikeRepo, mockCommentRepo)
			
			handler := NewCoordinateHandler(mockUsecase, mockCommentRepo, mockLikeRepo)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/coordinates?"+tt.query, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.GetMyCoordinates(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
			mockCommentRepo.AssertExpectations(t)
			mockLikeRepo.AssertExpectations(t)
		})
	}
}

func TestCoordinateHandler_LikeCoordinate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name          string
		userID        uint
		coordinateID  string
		mockSetup     func(*mockCoordinateUsecase)
		expectedCode  int
		checkBody     func(*testing.T, map[string]interface{})
	}{
		{
			name:         "successful like",
			userID:       1,
			coordinateID: "1",
			mockSetup: func(m *mockCoordinateUsecase) {
				m.On("LikeCoordinate", mock.Anything, uint(1), uint(1)).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Coordinate liked successfully", body["message"])
			},
		},
		{
			name:         "already liked",
			userID:       1,
			coordinateID: "1",
			mockSetup: func(m *mockCoordinateUsecase) {
				m.On("LikeCoordinate", mock.Anything, uint(1), uint(1)).Return(errors.New("already liked"))
			},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "already liked", body["error"])
			},
		},
		{
			name:         "coordinate not found",
			userID:       1,
			coordinateID: "999",
			mockSetup: func(m *mockCoordinateUsecase) {
				m.On("LikeCoordinate", mock.Anything, uint(1), uint(999)).Return(errors.New("coordinate not found"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "coordinate not found", body["error"])
			},
		},
		{
			name:         "invalid coordinate ID",
			userID:       1,
			coordinateID: "invalid",
			mockSetup:    func(m *mockCoordinateUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid coordinate ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockCoordinateUsecase)
			mockCommentRepo := new(mockCommentRepository)
			mockLikeRepo := new(mockLikeCoordinateRepository)
			
			tt.mockSetup(mockUsecase)
			
			handler := NewCoordinateHandler(mockUsecase, mockCommentRepo, mockLikeRepo)
			
			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/coordinates/"+tt.coordinateID+"/like", nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			c.Params = gin.Params{
				{Key: "id", Value: tt.coordinateID},
			}
			
			// Execute
			handler.LikeCoordinate(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestCoordinateHandler_DeleteCoordinate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name          string
		userID        uint
		coordinateID  string
		mockSetup     func(*mockCoordinateUsecase)
		expectedCode  int
		checkBody     func(*testing.T, map[string]interface{})
	}{
		{
			name:         "successful delete",
			userID:       1,
			coordinateID: "1",
			mockSetup: func(m *mockCoordinateUsecase) {
				m.On("DeleteCoordinate", mock.Anything, uint(1), uint(1)).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Coordinate deleted successfully", body["message"])
			},
		},
		{
			name:         "unauthorized delete",
			userID:       2,
			coordinateID: "1",
			mockSetup: func(m *mockCoordinateUsecase) {
				m.On("DeleteCoordinate", mock.Anything, uint(2), uint(1)).Return(errors.New("unauthorized"))
			},
			expectedCode: http.StatusForbidden,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "unauthorized", body["error"])
			},
		},
		{
			name:         "coordinate not found",
			userID:       1,
			coordinateID: "999",
			mockSetup: func(m *mockCoordinateUsecase) {
				m.On("DeleteCoordinate", mock.Anything, uint(1), uint(999)).Return(errors.New("coordinate not found"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "coordinate not found", body["error"])
			},
		},
		{
			name:         "invalid coordinate ID",
			userID:       1,
			coordinateID: "invalid",
			mockSetup:    func(m *mockCoordinateUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid coordinate ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockCoordinateUsecase)
			mockCommentRepo := new(mockCommentRepository)
			mockLikeRepo := new(mockLikeCoordinateRepository)
			
			tt.mockSetup(mockUsecase)
			
			handler := NewCoordinateHandler(mockUsecase, mockCommentRepo, mockLikeRepo)
			
			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/coordinates/"+tt.coordinateID, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			c.Params = gin.Params{
				{Key: "id", Value: tt.coordinateID},
			}
			
			// Execute
			handler.DeleteCoordinate(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}