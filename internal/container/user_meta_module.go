package container

import (
	"usermanagement-api/domain/ports"
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type UserMetaModule struct {
	Handler *handlers.UserMetaHandler
}

func NewUserMetaModule(db *gorm.DB, cache ports.Cache) *UserMetaModule {
	userMetaRepo := gormrepo.NewUserMetaRepository(db)
	uc := usecase.NewUserMetaUseCase(userMetaRepo, cache)
	return &UserMetaModule{Handler: handlers.NewUserMetaHandler(uc)}
}
