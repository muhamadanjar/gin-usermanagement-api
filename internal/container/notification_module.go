package container

import (
	"usermanagement-api/domain/ports"
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"

	"gorm.io/gorm"
)

type NotificationModule struct {
	Handler *handlers.NotificationHandler
}

func NewNotificationModule(db *gorm.DB, notifier ports.Notifier) *NotificationModule {
	userMetaRepo := gormrepo.NewUserMetaRepository(db)
	uc := usecase.NewNotificationUseCase(userMetaRepo, notifier)
	return &NotificationModule{Handler: handlers.NewNotificationHandler(uc)}
}
