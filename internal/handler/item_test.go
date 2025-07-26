package handler

import (
	"bytes"
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
type mockItemUsecase struct {
	mock.Mock
}

func (m *mockItemUsecase) CreateItem(ctx context.Context, userID uint, item *domain.Item, image *multipart.FileHeader) error {
	args := m.Called(ctx, userID, item, image)
	return args.Error(0)
}

func (m *mockItemUsecase) GetItem(ctx context.Context, itemID uint) (*domain.Item, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *mockItemUsecase) UpdateItem(ctx context.Context, userID uint, itemID uint, updates map[string]interface{}, image *multipart.FileHeader) error {
	args := m.Called(ctx, userID, itemID, updates, image)
	return args.Error(0)
}

func (m *mockItemUsecase) DeleteItem(ctx context.Context, userID uint, itemID uint) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}

func (m *mockItemUsecase) GetUserItems(ctx context.Context, userID uint, limit, offset int) ([]*domain.Item, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Item), args.Get(1).(int64), args.Error(2)
}

func (m *mockItemUsecase) SearchItems(ctx context.Context, filters repository.ItemFilter) ([]*domain.Item, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Item), args.Error(1)
}

func (m *mockItemUsecase) DeleteUserItems(ctx context.Context, userID uint, itemIDs []uint) error {
	args := m.Called(ctx, userID, itemIDs)
	return args.Error(0)
}

func (m *mockItemUsecase) GetUserItemStatistics(ctx context.Context, userID uint) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestItemHandler_GetItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		itemID       string
		mockSetup    func(*mockItemUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get item",
			itemID: "1",
			mockSetup: func(m *mockItemUsecase) {
				m.On("GetItem", mock.Anything, uint(1)).Return(&domain.Item{
					BaseModel: domain.BaseModel{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					UserID:    1,
					SuperItem: "Tシャツ",
					Season:    domain.SeasonSummer,
					TPO:       domain.TPOCasual,
					Color:     domain.ColorBlue,
					Content:   "Nice T-shirt",
					Memo:      "Favorite item",
					Picture:   "/uploads/tshirt.jpg",
					Rating:    5,
				}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(1), body["id"])
				assert.Equal(t, float64(1), body["user_id"])
				assert.Equal(t, "Tシャツ", body["super_item"])
				assert.Equal(t, float64(domain.SeasonSummer), body["season"])
				assert.Equal(t, float64(domain.TPOCasual), body["tpo"])
				assert.Equal(t, float64(domain.ColorBlue), body["color"])
				assert.Equal(t, "Nice T-shirt", body["content"])
				assert.Equal(t, "Favorite item", body["memo"])
				assert.Equal(t, "/uploads/tshirt.jpg", body["picture"])
				assert.Equal(t, float64(5), body["rating"])
			},
		},
		{
			name:   "item not found",
			itemID: "999",
			mockSetup: func(m *mockItemUsecase) {
				m.On("GetItem", mock.Anything, uint(999)).Return(nil, errors.New("item not found"))
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "item not found", body["error"])
			},
		},
		{
			name:         "invalid item ID",
			itemID:       "invalid",
			mockSetup:    func(m *mockItemUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid item ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockItemUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewItemHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/items/"+tt.itemID, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{
				gin.Param{Key: "id", Value: tt.itemID},
			}
			
			// Execute
			handler.GetItem(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestItemHandler_GetMyItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		query        string
		mockSetup    func(*mockItemUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get my items",
			userID: 1,
			query:  "page=1&per_page=10",
			mockSetup: func(m *mockItemUsecase) {
				items := []*domain.Item{
					{
						BaseModel: domain.BaseModel{
							ID:        1,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						UserID:    1,
						SuperItem: "Tシャツ",
						Season:    domain.SeasonSummer,
						TPO:       domain.TPOCasual,
						Color:     domain.ColorBlue,
						Content:   "Blue T-shirt",
						Picture:   "/uploads/tshirt1.jpg",
						Rating:    5,
					},
					{
						BaseModel: domain.BaseModel{
							ID:        2,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						UserID:    1,
						SuperItem: "パンツ",
						Season:    domain.SeasonAllSeason,
						TPO:       domain.TPOWork,
						Color:     domain.ColorBlack,
						Content:   "Black pants",
						Picture:   "/uploads/pants1.jpg",
						Rating:    4,
					},
				}
				m.On("GetUserItems", mock.Anything, uint(1), 10, 0).Return(items, int64(2), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				items := body["items"].([]interface{})
				assert.Len(t, items, 2)
				assert.Equal(t, float64(2), body["total_count"])
				assert.Equal(t, float64(1), body["page"])
				assert.Equal(t, float64(10), body["per_page"])
			},
		},
		{
			name:   "empty result",
			userID: 1,
			query:  "page=1&per_page=10",
			mockSetup: func(m *mockItemUsecase) {
				m.On("GetUserItems", mock.Anything, uint(1), 10, 0).Return([]*domain.Item{}, int64(0), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				items := body["items"].([]interface{})
				assert.Len(t, items, 0)
				assert.Equal(t, float64(0), body["total_count"])
			},
		},
		{
			name:         "invalid pagination",
			userID:       1,
			query:        "page=invalid&per_page=10",
			mockSetup:    func(m *mockItemUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.NotNil(t, body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockItemUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewItemHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/items?"+tt.query, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.GetMyItems(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestItemHandler_DeleteItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		itemID       string
		mockSetup    func(*mockItemUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful delete",
			userID: 1,
			itemID: "1",
			mockSetup: func(m *mockItemUsecase) {
				m.On("DeleteItem", mock.Anything, uint(1), uint(1)).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Item deleted successfully", body["message"])
			},
		},
		{
			name:   "unauthorized delete",
			userID: 2,
			itemID: "1",
			mockSetup: func(m *mockItemUsecase) {
				m.On("DeleteItem", mock.Anything, uint(2), uint(1)).Return(errors.New("unauthorized"))
			},
			expectedCode: http.StatusForbidden,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "unauthorized", body["error"])
			},
		},
		{
			name:   "item not found",
			userID: 1,
			itemID: "999",
			mockSetup: func(m *mockItemUsecase) {
				m.On("DeleteItem", mock.Anything, uint(1), uint(999)).Return(errors.New("item not found"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "item not found", body["error"])
			},
		},
		{
			name:         "invalid item ID",
			userID:       1,
			itemID:       "invalid",
			mockSetup:    func(m *mockItemUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Invalid item ID", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockItemUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewItemHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/items/"+tt.itemID, nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			c.Params = gin.Params{
				gin.Param{Key: "id", Value: tt.itemID},
			}
			
			// Execute
			handler.DeleteItem(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestItemHandler_DeleteItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		requestBody  map[string]interface{}
		mockSetup    func(*mockItemUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful batch delete",
			userID: 1,
			requestBody: map[string]interface{}{
				"item_ids": []uint{1, 2, 3},
			},
			mockSetup: func(m *mockItemUsecase) {
				m.On("DeleteUserItems", mock.Anything, uint(1), []uint{1, 2, 3}).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Items deleted successfully", body["message"])
			},
		},
		{
			name:   "empty item IDs",
			userID: 1,
			requestBody: map[string]interface{}{
				"item_ids": []uint{},
			},
			mockSetup:    func(m *mockItemUsecase) {},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.NotNil(t, body["error"])
			},
		},
		{
			name:   "unauthorized batch delete",
			userID: 2,
			requestBody: map[string]interface{}{
				"item_ids": []uint{1, 2},
			},
			mockSetup: func(m *mockItemUsecase) {
				m.On("DeleteUserItems", mock.Anything, uint(2), []uint{1, 2}).Return(errors.New("unauthorized"))
			},
			expectedCode: http.StatusForbidden,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "unauthorized", body["error"])
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUsecase := new(mockItemUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewItemHandler(mockUsecase)
			
			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/items", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.DeleteItems(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestItemHandler_GetItemStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name         string
		userID       uint
		mockSetup    func(*mockItemUsecase)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:   "successful get statistics",
			userID: 1,
			mockSetup: func(m *mockItemUsecase) {
				stats := map[string]interface{}{
					"total_count": 50,
					"category_count": map[string]int{
						"Tシャツ": 10,
						"パンツ":  8,
						"シューズ": 5,
					},
					"season_count": map[int]int{
						domain.SeasonSpring: 15,
						domain.SeasonSummer: 20,
						domain.SeasonAutumn: 10,
						domain.SeasonWinter: 5,
					},
					"average_rating": 4.5,
				}
				m.On("GetUserItemStatistics", mock.Anything, uint(1)).Return(stats, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(50), body["total_count"])
				assert.NotNil(t, body["category_count"])
				assert.NotNil(t, body["season_count"])
				assert.Equal(t, 4.5, body["average_rating"])
			},
		},
		{
			name:   "statistics error",
			userID: 1,
			mockSetup: func(m *mockItemUsecase) {
				m.On("GetUserItemStatistics", mock.Anything, uint(1)).Return(nil, errors.New("database error"))
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
			mockUsecase := new(mockItemUsecase)
			tt.mockSetup(mockUsecase)
			
			handler := NewItemHandler(mockUsecase)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/items/statistics", nil)
			
			// Create response recorder
			w := httptest.NewRecorder()
			
			// Setup gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)
			
			// Execute
			handler.GetItemStatistics(c)
			
			// Assert
			assert.Equal(t, tt.expectedCode, w.Code)
			
			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			tt.checkBody(t, responseBody)
			
			mockUsecase.AssertExpectations(t)
		})
	}
}