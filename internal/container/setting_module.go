package container

import (
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"
	"usermanagement-api/pkg/cache"

	"gorm.io/gorm"
)

type SettingModule struct {
	Handler *handlers.SettingHandler
}

func NewSettingModule(db *gorm.DB, cache cache.Cache) *SettingModule {
	settingRepo := gormrepo.NewSettingRepository(db)
	uc := usecase.NewSettingUseCase(settingRepo, cache)
	return &SettingModule{Handler: handlers.NewSettingHandler(uc)}
}
