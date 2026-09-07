package container

import (
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type PermissionModule struct {
	Handler *handlers.PermissionHandler
}

func NewPermissionModule(db *gorm.DB) *PermissionModule {
	permissionRepo := gormrepo.NewPermissionRepository(db)
	uc := usecase.NewPermissionUseCase(permissionRepo)
	return &PermissionModule{Handler: handlers.NewPermissionHandler(uc)}
}
