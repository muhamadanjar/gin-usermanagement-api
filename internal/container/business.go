package container

import (
	"usermanagement-api/config"
	"usermanagement-api/domain/ports"
	"usermanagement-api/infrastructure/logger"
	"usermanagement-api/infrastructure/security"
	"usermanagement-api/internal/presentation/http/handlers"
	"usermanagement-api/internal/presentation/http/middleware"
	"usermanagement-api/pkg/auth"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// BusinessContainer is the composition root: it assembles every module and
// hands the finished handlers to the routing layer. Per-module factories in
// this package own their own dependency wiring.
type BusinessContainer struct {
	// Middleware
	AuthMiddleware middleware.AuthMiddleware
	CORSMiddleware middleware.CORSMiddleware

	// Handlers
	UserHandler         *handlers.UserHandler
	RoleHandler         *handlers.RoleHandler
	PermissionHandler   *handlers.PermissionHandler
	MenuHandler         *handlers.MenuHandler
	AuthHandler         *handlers.AuthHandler
	UserMetaHandler     *handlers.UserMetaHandler
	SettingHandler      *handlers.SettingHandler
	NotificationHandler *handlers.NotificationHandler
}

func NewBusinessContainer(
	db *gorm.DB,
	cache ports.Cache,
	fcmClient ports.Notifier,
	corsConfig config.CORSConfig,
	jwtService *auth.JWTService,
	zapLogger *zap.Logger,
) *BusinessContainer {
	tokenManager := security.NewTokenManager(jwtService)
	hasher := security.NewPasswordHasher()
	log := logger.NewPortLogger(zapLogger)

	return &BusinessContainer{
		CORSMiddleware:      middleware.NewCORSMiddleware(corsConfig),
		AuthMiddleware:      NewAuthzModule(db, tokenManager).Middleware,
		UserHandler:         NewUserModule(db, hasher).Handler,
		RoleHandler:         NewRoleModule(db).Handler,
		PermissionHandler:   NewPermissionModule(db).Handler,
		MenuHandler:         NewMenuModule(db).Handler,
		AuthHandler:         NewAuthModule(db, fcmClient, tokenManager, hasher, log).Handler,
		UserMetaHandler:     NewUserMetaModule(db, cache).Handler,
		SettingHandler:      NewSettingModule(db, cache).Handler,
		NotificationHandler: NewNotificationModule(db, fcmClient).Handler,
	}
}
