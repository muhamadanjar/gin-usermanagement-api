package usecase

import (
	"usermanagement-api/internal/user/domain"
	"usermanagement-api/internal/user/repository"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (uc *UserUsecase) GetUserByID(id uint) (*domain.User, error) {
	return uc.userRepo.GetUserById()
}
