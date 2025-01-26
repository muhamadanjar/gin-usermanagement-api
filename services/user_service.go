package services

import "usermanagement-api/repository"

type (
	UserService interface {
	}

	userService struct {
		userRepo   repository.UserRepository
		jwtService JWTService
	}
)
