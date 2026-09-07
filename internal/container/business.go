package container

import (
	"usermanagement-api/config"
	"usermanagement-api/internal/presentation/http/handlers"
	"usermanagement-api/internal/presentation/http/middleware"
	"usermanagement-api/pkg/auth"
	"usermanagement-api/pkg/cache"
	"usermanagement-api/pkg/firebase"

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
	cache cache.Cache,
	fcmClient firebase.FCMClient,
	corsConfig config.CORSConfig,
	jwtService *auth.JWTService,
) *BusinessContainer {
	return &BusinessContainer{
		CORSMiddleware:      middleware.NewCORSMiddleware(corsConfig),
		AuthMiddleware:      NewAuthzModule(db, jwtService).Middleware,
		UserHandler:         NewUserModule(db).Handler,
		RoleHandler:         NewRoleModule(db).Handler,
		PermissionHandler:   NewPermissionModule(db).Handler,
		MenuHandler:         NewMenuModule(db).Handler,
		AuthHandler:         NewAuthModule(db, fcmClient, jwtService).Handler,
		UserMetaHandler:     NewUserMetaModule(db, cache).Handler,
		SettingHandler:      NewSettingModule(db, cache).Handler,
		NotificationHandler: NewNotificationModule(db, fcmClient).Handler,
	}
}
