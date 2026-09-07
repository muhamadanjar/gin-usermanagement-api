package container

import (
	"usermanagement-api/domain/ports"
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type AuthModule struct {
	Handler *handlers.AuthHandler
}

func NewAuthModule(
	db *gorm.DB,
	notifier ports.Notifier,
	tokenManager ports.TokenManager,
	hasher ports.PasswordHasher,
	log ports.Logger,
) *AuthModule {
	userRepo := gormrepo.NewUserRepository(db)
	roleRepo := gormrepo.NewRoleRepository(db)
	menuRepo := gormrepo.NewMenuRepository(db)
	userMetaRepo := gormrepo.NewUserMetaRepository(db)
	uc := usecase.NewAuthUseCase(userRepo, roleRepo, menuRepo, userMetaRepo, notifier, tokenManager, hasher, log)
	return &AuthModule{Handler: handlers.NewAuthHandler(uc)}
}
