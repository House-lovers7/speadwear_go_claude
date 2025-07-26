package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/House-lovers7/speadwear-go/internal/domain"
	"github.com/House-lovers7/speadwear-go/internal/dto"
	"github.com/House-lovers7/speadwear-go/internal/repository"
	"github.com/House-lovers7/speadwear-go/internal/usecase"
)

type ItemHandler struct {
	itemUsecase usecase.ItemUsecase
}

// NewItemHandler creates a new item handler
func NewItemHandler(itemUsecase usecase.ItemUsecase) *ItemHandler {
	return &ItemHandler{
		itemUsecase: itemUsecase,
	}
}

// CreateItem POST /api/v1/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var req dto.CreateItemRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload
	file, _ := c.FormFile("picture")

	item := &domain.Item{
		SuperItem: req.SuperItem,
		Season:    req.Season,
		TPO:       req.TPO,
		Color:     req.Color,
		Content:   req.Content,
		Memo:      req.Memo,
		Rating:    req.Rating,
	}

	err := h.itemUsecase.CreateItem(c.Request.Context(), userID, item, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ItemResponse{
		ID:        item.ID,
		UserID:    item.UserID,
		SuperItem: item.SuperItem,
		Season:    item.Season,
		TPO:       item.TPO,
		Color:     item.Color,
		Content:   item.Content,
		Memo:      item.Memo,
		Picture:   item.Picture,
		Rating:    item.Rating,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	})
}

// GetItem GET /api/v1/items/:id
func (h *ItemHandler) GetItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemUsecase.GetItem(c.Request.Context(), uint(itemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ItemResponse{
		ID:           item.ID,
		UserID:       item.UserID,
		CoordinateID: item.CoordinateID,
		SuperItem:    item.SuperItem,
		Season:       item.Season,
		TPO:          item.TPO,
		Color:        item.Color,
		Content:      item.Content,
		Memo:         item.Memo,
		Picture:      item.Picture,
		Rating:       item.Rating,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	})
}

// GetMyItems GET /api/v1/items
func (h *ItemHandler) GetMyItems(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	items, total, err := h.itemUsecase.GetUserItems(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	itemResponses := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = dto.ItemResponse{
			ID:           item.ID,
			UserID:       item.UserID,
			CoordinateID: item.CoordinateID,
			SuperItem:    item.SuperItem,
			Season:       item.Season,
			TPO:          item.TPO,
			Color:        item.Color,
			Content:      item.Content,
			Memo:         item.Memo,
			Picture:      item.Picture,
			Rating:       item.Rating,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.ItemListResponse{
		Items:      itemResponses,
		TotalCount: total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
	})
}

// GetUserItems GET /api/v1/users/:user_id/items
func (h *ItemHandler) GetUserItems(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	items, total, err := h.itemUsecase.GetUserItems(c.Request.Context(), uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	itemResponses := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = dto.ItemResponse{
			ID:           item.ID,
			UserID:       item.UserID,
			CoordinateID: item.CoordinateID,
			SuperItem:    item.SuperItem,
			Season:       item.Season,
			TPO:          item.TPO,
			Color:        item.Color,
			Content:      item.Content,
			Memo:         item.Memo,
			Picture:      item.Picture,
			Rating:       item.Rating,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.ItemListResponse{
		Items:      itemResponses,
		TotalCount: total,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
	})
}

// SearchItems GET /api/v1/items/search
func (h *ItemHandler) SearchItems(c *gin.Context) {
	var filter dto.ItemFilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to repository filter
	repoFilter := repository.ItemFilter{
		Season:    filter.Season,
		TPO:       filter.TPO,
		Color:     filter.Color,
		SuperItem: filter.SuperItem,
		MinRating: filter.MinRating,
		MaxRating: filter.MaxRating,
		Limit:     filter.PerPage,
		Offset:    (filter.Page - 1) * filter.PerPage,
	}

	items, err := h.itemUsecase.SearchItems(c.Request.Context(), repoFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	itemResponses := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = dto.ItemResponse{
			ID:           item.ID,
			UserID:       item.UserID,
			CoordinateID: item.CoordinateID,
			SuperItem:    item.SuperItem,
			Season:       item.Season,
			TPO:          item.TPO,
			Color:        item.Color,
			Content:      item.Content,
			Memo:         item.Memo,
			Picture:      item.Picture,
			Rating:       item.Rating,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.ItemListResponse{
		Items:      itemResponses,
		TotalCount: int64(len(items)), // TODO: Add count to repository method
		Page:       filter.Page,
		PerPage:    filter.PerPage,
	})
}

// UpdateItem PUT /api/v1/items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload
	file, _ := c.FormFile("picture")

	// Convert request to map
	updates := make(map[string]interface{})
	if req.SuperItem != nil {
		updates["super_item"] = *req.SuperItem
	}
	if req.Season != nil {
		updates["season"] = *req.Season
	}
	if req.TPO != nil {
		updates["tpo"] = *req.TPO
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Memo != nil {
		updates["memo"] = *req.Memo
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}

	err = h.itemUsecase.UpdateItem(c.Request.Context(), userID, uint(itemID), updates, file)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// DeleteItem DELETE /api/v1/items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = h.itemUsecase.DeleteItem(c.Request.Context(), userID, uint(itemID))
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// DeleteItems DELETE /api/v1/items
func (h *ItemHandler) DeleteItems(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var req struct {
		ItemIDs []uint `json:"item_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.itemUsecase.DeleteUserItems(c.Request.Context(), userID, req.ItemIDs)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Items deleted successfully"})
}

// GetItemStatistics GET /api/v1/items/statistics
func (h *ItemHandler) GetItemStatistics(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	stats, err := h.itemUsecase.GetUserItemStatistics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}