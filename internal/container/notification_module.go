package container

import (
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/handlers"
	"usermanagement-api/pkg/firebase"

	"gorm.io/gorm"
)

type NotificationModule struct {
	Handler *handlers.NotificationHandler
}

func NewNotificationModule(db *gorm.DB, fcmClient firebase.FCMClient) *NotificationModule {
	userMetaRepo := gormrepo.NewUserMetaRepository(db)
	uc := usecase.NewNotificationUseCase(userMetaRepo, fcmClient)
	return &NotificationModule{Handler: handlers.NewNotificationHandler(uc)}
}
