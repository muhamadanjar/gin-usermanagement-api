package container

import (
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type MenuModule struct {
	Handler *handlers.MenuHandler
}

func NewMenuModule(db *gorm.DB) *MenuModule {
	menuRepo := gormrepo.NewMenuRepository(db)
	uc := usecase.NewMenuUseCase(menuRepo)
	return &MenuModule{Handler: handlers.NewMenuHandler(uc)}
}
