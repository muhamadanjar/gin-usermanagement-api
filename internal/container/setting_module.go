package container

import (
	"usermanagement-api/domain/ports"
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type SettingModule struct {
	Handler *handlers.SettingHandler
}

func NewSettingModule(db *gorm.DB, cache ports.Cache) *SettingModule {
	settingRepo := gormrepo.NewSettingRepository(db)
	uc := usecase.NewSettingUseCase(settingRepo, cache)
	return &SettingModule{Handler: handlers.NewSettingHandler(uc)}
}
