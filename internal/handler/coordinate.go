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

type CoordinateHandler struct {
	coordinateUsecase usecase.CoordinateUsecase
	commentRepo       repository.CommentRepository
	likeCoordRepo     repository.LikeCoordinateRepository
}

// NewCoordinateHandler creates a new coordinate handler
func NewCoordinateHandler(
	coordinateUsecase usecase.CoordinateUsecase,
	commentRepo repository.CommentRepository,
	likeCoordRepo repository.LikeCoordinateRepository,
) *CoordinateHandler {
	return &CoordinateHandler{
		coordinateUsecase: coordinateUsecase,
		commentRepo:       commentRepo,
		likeCoordRepo:     likeCoordRepo,
	}
}

// CreateCoordinate POST /api/v1/coordinates
func (h *CoordinateHandler) CreateCoordinate(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var req dto.CreateCoordinateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload
	file, _ := c.FormFile("picture")

	coordinate := &domain.Coordinate{
		Season:         req.Season,
		TPO:            req.TPO,
		SiTopLength:    req.SiTopLength,
		SiTopSleeve:    req.SiTopSleeve,
		SiBottomLength: req.SiBottomLength,
		SiBottomType:   req.SiBottomType,
		SiDressLength:  req.SiDressLength,
		SiDressSleeve:  req.SiDressSleeve,
		SiOuterLength:  req.SiOuterLength,
		SiOuterSleeve:  req.SiOuterSleeve,
		SiShoeSize:     req.SiShoeSize,
		Memo:           req.Memo,
		Rating:         req.Rating,
	}

	err := h.coordinateUsecase.CreateCoordinate(c.Request.Context(), userID, coordinate, req.ItemIDs, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the created coordinate with details
	coordinateWithDetails, _ := h.coordinateUsecase.GetCoordinateWithDetails(c.Request.Context(), coordinate.ID)
	if coordinateWithDetails != nil {
		coordinate = coordinateWithDetails
	}

	// Convert to response
	resp := h.coordinateToResponse(c, coordinate)
	c.JSON(http.StatusCreated, resp)
}

// GetCoordinate GET /api/v1/coordinates/:id
func (h *CoordinateHandler) GetCoordinate(c *gin.Context) {
	coordinateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinate ID"})
		return
	}

	coordinate, err := h.coordinateUsecase.GetCoordinateWithDetails(c.Request.Context(), uint(coordinateID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	resp := h.coordinateToResponse(c, coordinate)
	c.JSON(http.StatusOK, resp)
}

// GetMyCoordinates GET /api/v1/coordinates
func (h *CoordinateHandler) GetMyCoordinates(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	coordinates, total, err := h.coordinateUsecase.GetUserCoordinates(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	coordinateResponses := make([]dto.CoordinateResponse, len(coordinates))
	for i, coordinate := range coordinates {
		coordinateResponses[i] = *h.coordinateToResponse(c, coordinate)
	}

	c.JSON(http.StatusOK, dto.CoordinateListResponse{
		Coordinates: coordinateResponses,
		TotalCount:  total,
		Page:        pagination.Page,
		PerPage:     pagination.PerPage,
	})
}

// GetUserCoordinates GET /api/v1/users/:user_id/coordinates
func (h *CoordinateHandler) GetUserCoordinates(c *gin.Context) {
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

	coordinates, total, err := h.coordinateUsecase.GetUserCoordinates(c.Request.Context(), uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	coordinateResponses := make([]dto.CoordinateResponse, len(coordinates))
	for i, coordinate := range coordinates {
		coordinateResponses[i] = *h.coordinateToResponse(c, coordinate)
	}

	c.JSON(http.StatusOK, dto.CoordinateListResponse{
		Coordinates: coordinateResponses,
		TotalCount:  total,
		Page:        pagination.Page,
		PerPage:     pagination.PerPage,
	})
}

// GetTimeline GET /api/v1/coordinates/timeline
func (h *CoordinateHandler) GetTimeline(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	coordinates, err := h.coordinateUsecase.GetTimelineCoordinates(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	coordinateResponses := make([]dto.CoordinateResponse, len(coordinates))
	for i, coordinate := range coordinates {
		coordinateResponses[i] = *h.coordinateToResponse(c, coordinate)
	}

	c.JSON(http.StatusOK, dto.CoordinateListResponse{
		Coordinates: coordinateResponses,
		TotalCount:  int64(len(coordinates)), // TODO: Add proper count
		Page:        pagination.Page,
		PerPage:     pagination.PerPage,
	})
}

// SearchCoordinates GET /api/v1/coordinates/search
func (h *CoordinateHandler) SearchCoordinates(c *gin.Context) {
	var filter dto.CoordinateFilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to repository filter
	repoFilter := repository.CoordinateFilter{
		Season:    filter.Season,
		TPO:       filter.TPO,
		MinRating: filter.MinRating,
		MaxRating: filter.MaxRating,
		Limit:     filter.PerPage,
		Offset:    (filter.Page - 1) * filter.PerPage,
	}

	coordinates, err := h.coordinateUsecase.SearchCoordinates(c.Request.Context(), repoFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	coordinateResponses := make([]dto.CoordinateResponse, len(coordinates))
	for i, coordinate := range coordinates {
		coordinateResponses[i] = *h.coordinateToResponse(c, coordinate)
	}

	c.JSON(http.StatusOK, dto.CoordinateListResponse{
		Coordinates: coordinateResponses,
		TotalCount:  int64(len(coordinates)), // TODO: Add count to repository method
		Page:        filter.Page,
		PerPage:     filter.PerPage,
	})
}

// UpdateCoordinate PUT /api/v1/coordinates/:id
func (h *CoordinateHandler) UpdateCoordinate(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	coordinateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinate ID"})
		return
	}

	var req dto.UpdateCoordinateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload
	file, _ := c.FormFile("picture")

	// Convert request to map
	updates := make(map[string]interface{})
	if req.Season != nil {
		updates["season"] = *req.Season
	}
	if req.TPO != nil {
		updates["tpo"] = *req.TPO
	}
	if req.SiTopLength != nil {
		updates["si_top_length"] = *req.SiTopLength
	}
	if req.SiTopSleeve != nil {
		updates["si_top_sleeve"] = *req.SiTopSleeve
	}
	if req.SiBottomLength != nil {
		updates["si_bottom_length"] = *req.SiBottomLength
	}
	if req.SiBottomType != nil {
		updates["si_bottom_type"] = *req.SiBottomType
	}
	if req.SiDressLength != nil {
		updates["si_dress_length"] = *req.SiDressLength
	}
	if req.SiDressSleeve != nil {
		updates["si_dress_sleeve"] = *req.SiDressSleeve
	}
	if req.SiOuterLength != nil {
		updates["si_outer_length"] = *req.SiOuterLength
	}
	if req.SiOuterSleeve != nil {
		updates["si_outer_sleeve"] = *req.SiOuterSleeve
	}
	if req.SiShoeSize != nil {
		updates["si_shoe_size"] = *req.SiShoeSize
	}
	if req.Memo != nil {
		updates["memo"] = *req.Memo
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}

	err = h.coordinateUsecase.UpdateCoordinate(c.Request.Context(), userID, uint(coordinateID), updates, req.ItemIDs, file)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coordinate updated successfully"})
}

// DeleteCoordinate DELETE /api/v1/coordinates/:id
func (h *CoordinateHandler) DeleteCoordinate(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	coordinateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinate ID"})
		return
	}

	err = h.coordinateUsecase.DeleteCoordinate(c.Request.Context(), userID, uint(coordinateID))
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coordinate deleted successfully"})
}

// LikeCoordinate POST /api/v1/coordinates/:id/like
func (h *CoordinateHandler) LikeCoordinate(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	coordinateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinate ID"})
		return
	}

	err = h.coordinateUsecase.LikeCoordinate(c.Request.Context(), userID, uint(coordinateID))
	if err != nil {
		if err.Error() == "already liked" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coordinate liked successfully"})
}

// UnlikeCoordinate DELETE /api/v1/coordinates/:id/like
func (h *CoordinateHandler) UnlikeCoordinate(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	coordinateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinate ID"})
		return
	}

	err = h.coordinateUsecase.UnlikeCoordinate(c.Request.Context(), userID, uint(coordinateID))
	if err != nil {
		if err.Error() == "not liked" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coordinate unliked successfully"})
}

// GetCoordinateComments GET /api/v1/coordinates/:id/comments
func (h *CoordinateHandler) GetCoordinateComments(c *gin.Context) {
	coordinateID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coordinate ID"})
		return
	}

	var pagination dto.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := pagination.PerPage
	offset := (pagination.Page - 1) * pagination.PerPage

	comments, err := h.commentRepo.FindByCoordinateID(c.Request.Context(), uint(coordinateID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := h.commentRepo.CountByCoordinateID(c.Request.Context(), uint(coordinateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commentResponses := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = dto.CommentResponse{
			ID:           comment.ID,
			UserID:       comment.UserID,
			CoordinateID: comment.CoordinateID,
			Comment:      comment.Comment,
			User: dto.UserResponse{
				ID:      comment.User.ID,
				Name:    comment.User.Name,
				Email:   comment.User.Email,
				Picture: comment.User.Picture,
			},
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.CommentListResponse{
		Comments:   commentResponses,
		TotalCount: count,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
	})
}

// GetCoordinateStatistics GET /api/v1/coordinates/statistics
func (h *CoordinateHandler) GetCoordinateStatistics(c *gin.Context) {
	userID := c.GetUint("userID") // From auth middleware

	stats, err := h.coordinateUsecase.GetUserCoordinateStatistics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// coordinateToResponse converts domain coordinate to response DTO
func (h *CoordinateHandler) coordinateToResponse(c *gin.Context, coordinate *domain.Coordinate) *dto.CoordinateResponse {
	// Get current user ID if authenticated
	var isLiked bool
	if userID, exists := c.Get("userID"); exists {
		isLiked, _ = h.coordinateUsecase.IsLikedByUser(c.Request.Context(), userID.(uint), coordinate.ID)
	}

	// Get counts
	likeCount, _ := h.likeCoordRepo.CountByCoordinateID(c.Request.Context(), coordinate.ID)
	commentCount, _ := h.commentRepo.CountByCoordinateID(c.Request.Context(), coordinate.ID)

	// Convert items
	itemResponses := make([]dto.ItemResponse, len(coordinate.Items))
	for i, item := range coordinate.Items {
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

	return &dto.CoordinateResponse{
		ID:             coordinate.ID,
		UserID:         coordinate.UserID,
		Season:         coordinate.Season,
		TPO:            coordinate.TPO,
		Picture:        coordinate.Picture,
		SiTopLength:    coordinate.SiTopLength,
		SiTopSleeve:    coordinate.SiTopSleeve,
		SiBottomLength: coordinate.SiBottomLength,
		SiBottomType:   coordinate.SiBottomType,
		SiDressLength:  coordinate.SiDressLength,
		SiDressSleeve:  coordinate.SiDressSleeve,
		SiOuterLength:  coordinate.SiOuterLength,
		SiOuterSleeve:  coordinate.SiOuterSleeve,
		SiShoeSize:     coordinate.SiShoeSize,
		Memo:           coordinate.Memo,
		Rating:         coordinate.Rating,
		Items:          itemResponses,
		LikeCount:      likeCount,
		CommentCount:   commentCount,
		IsLiked:        isLiked,
		User: dto.UserResponse{
			ID:        coordinate.User.ID,
			Name:      coordinate.User.Name,
			Email:     coordinate.User.Email,
			Picture:   coordinate.User.Picture,
			Admin:     coordinate.User.Admin,
			Activated: coordinate.User.Activated,
			CreatedAt: coordinate.User.CreatedAt,
		},
		CreatedAt: coordinate.CreatedAt,
		UpdatedAt: coordinate.UpdatedAt,
	}
}