package container

import (
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type RoleModule struct {
	Handler *handlers.RoleHandler
}

func NewRoleModule(db *gorm.DB) *RoleModule {
	roleRepo := gormrepo.NewRoleRepository(db)
	permissionRepo := gormrepo.NewPermissionRepository(db)
	uc := usecase.NewRoleUseCase(roleRepo, permissionRepo)
	return &RoleModule{Handler: handlers.NewRoleHandler(uc)}
}
