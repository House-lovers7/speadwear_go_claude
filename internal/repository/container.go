package repository

import "gorm.io/gorm"

// Container holds all repositories
type Container struct {
	User             UserRepository
	Item             ItemRepository
	Coordinate       CoordinateRepository
	Comment          CommentRepository
	LikeCoordinate   LikeCoordinateRepository
	Relationship     RelationshipRepository
	Block            BlockRepository
	Notification     NotificationRepository
}

// NewContainer creates a new repository container
func NewContainer(db *gorm.DB) *Container {
	return &Container{
		User:           NewUserRepository(db),
		Item:           NewItemRepository(db),
		Coordinate:     NewCoordinateRepository(db),
		Comment:        NewCommentRepository(db),
		LikeCoordinate: NewLikeCoordinateRepository(db),
		Relationship:   NewRelationshipRepository(db),
		Block:          NewBlockRepository(db),
		Notification:   NewNotificationRepository(db),
	}
}