package container

import (
	"usermanagement-api/domain/ports"
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type UserModule struct {
	Handler *handlers.UserHandler
}

func NewUserModule(db *gorm.DB, hasher ports.PasswordHasher) *UserModule {
	userRepo := gormrepo.NewUserRepository(db)
	roleRepo := gormrepo.NewRoleRepository(db)
	userMetaRepo := gormrepo.NewUserMetaRepository(db)
	uc := usecase.NewUserUseCase(userRepo, roleRepo, userMetaRepo, hasher)
	return &UserModule{Handler: handlers.NewUserHandler(uc)}
}
