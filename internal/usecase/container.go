package usecase

// Container holds all usecases
type Container struct {
	User       UserUsecase
	Item       ItemUsecase
	Coordinate CoordinateUsecase
	Social     SocialUsecase
}