package container

import (
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"
	"usermanagement-api/pkg/auth"
	"usermanagement-api/pkg/firebase"

	"gorm.io/gorm"
)

type AuthModule struct {
	Handler *handlers.AuthHandler
}

func NewAuthModule(db *gorm.DB, fcmClient firebase.FCMClient, jwtService *auth.JWTService) *AuthModule {
	userRepo := gormrepo.NewUserRepository(db)
	roleRepo := gormrepo.NewRoleRepository(db)
	menuRepo := gormrepo.NewMenuRepository(db)
	modelPermissionRepo := gormrepo.NewModelPermissionRepository(db)
	userMetaRepo := gormrepo.NewUserMetaRepository(db)
	uc := usecase.NewAuthUseCase(userRepo, roleRepo, menuRepo, modelPermissionRepo, userMetaRepo, fcmClient, jwtService)
	return &AuthModule{Handler: handlers.NewAuthHandler(uc)}
}
