package container

import (
	"usermanagement-api/config"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/delivery/http/handlers"
	"usermanagement-api/internal/delivery/http/middleware"
	"usermanagement-api/internal/usecase"
	"usermanagement-api/pkg/cache"
	"usermanagement-api/pkg/firebase"

	"gorm.io/gorm"
)

// BusinessContainer holds all business logic dependencies
type BusinessContainer struct {
	// Repositories
	UserRepository            repositories.UserRepository
	RoleRepository            repositories.RoleRepository
	PermissionRepository      repositories.PermissionRepository
	MenuRepository            repositories.MenuRepository
	ModelPermissionRepository repositories.ModelPermissionRepository
	UserMetaRepository        repositories.UserMetaRepository
	SettingRepository         repositories.SettingRepository

	// Use Cases
	UserUseCase         usecase.UserUseCase
	RoleUseCase         usecase.RoleUseCase
	PermissionUseCase   usecase.PermissionUseCase
	MenuUseCase         usecase.MenuUseCase
	AuthUseCase         usecase.AuthUseCase
	UserMetaUseCase     usecase.UserMetaUseCase
	SettingUseCase      usecase.SettingUseCase
	NotificationUseCase usecase.NotificationUseCase

	// Handlers
	UserHandler         *handlers.UserHandler
	RoleHandler         *handlers.RoleHandler
	PermissionHandler   *handlers.PermissionHandler
	MenuHandler         *handlers.MenuHandler
	AuthHandler         *handlers.AuthHandler
	UserMetaHandler     *handlers.UserMetaHandler
	SettingHandler      *handlers.SettingHandler
	NotificationHandler *handlers.NotificationHandler

	// Middleware
	AuthMiddleware middleware.AuthMiddleware
	CORSMiddleware middleware.CORSMiddleware
}

// NewBusinessContainer creates and initializes a new BusinessContainer
func NewBusinessContainer(db *gorm.DB, cache cache.Cache, fcmClient firebase.FCMClient, corsConfig config.CORSConfig) *BusinessContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	menuRepo := repositories.NewMenuRepository(db)
	modelPermissionRepo := repositories.NewModelPermissionRepository(db)
	userMetaRepo := repositories.NewUserMetaRepository(db)
	settingRepo := repositories.NewSettingRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo, userMetaRepo)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)
	menuUseCase := usecase.NewMenuUseCase(menuRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, roleRepo, menuRepo, modelPermissionRepo, userMetaRepo, fcmClient)
	userMetaUseCase := usecase.NewUserMetaUseCase(userMetaRepo, cache)
	settingUseCase := usecase.NewSettingUseCase(settingRepo, cache)
	notificationUseCase := usecase.NewNotificationUseCase(userMetaRepo, fcmClient)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo, roleRepo, permissionRepo, modelPermissionRepo)
	corsMiddleware := middleware.NewCORSMiddleware(corsConfig)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	roleHandler := handlers.NewRoleHandler(roleUseCase)
	permissionHandler := handlers.NewPermissionHandler(permissionUseCase)
	menuHandler := handlers.NewMenuHandler(menuUseCase)
	authHandler := handlers.NewAuthHandler(authUseCase)
	userMetaHandler := handlers.NewUserMetaHandler(userMetaUseCase)
	settingHandler := handlers.NewSettingHandler(settingUseCase)
	notificationHandler := handlers.NewNotificationHandler(notificationUseCase)

	return &BusinessContainer{
		// Repositories
		UserRepository:            userRepo,
		RoleRepository:            roleRepo,
		PermissionRepository:      permissionRepo,
		MenuRepository:            menuRepo,
		ModelPermissionRepository: modelPermissionRepo,
		UserMetaRepository:        userMetaRepo,
		SettingRepository:         settingRepo,

		// Use Cases
		UserUseCase:         userUseCase,
		RoleUseCase:         roleUseCase,
		PermissionUseCase:   permissionUseCase,
		MenuUseCase:         menuUseCase,
		AuthUseCase:         authUseCase,
		UserMetaUseCase:     userMetaUseCase,
		SettingUseCase:      settingUseCase,
		NotificationUseCase: notificationUseCase,

		// Handlers
		UserHandler:         userHandler,
		RoleHandler:         roleHandler,
		PermissionHandler:   permissionHandler,
		MenuHandler:         menuHandler,
		AuthHandler:         authHandler,
		UserMetaHandler:     userMetaHandler,
		SettingHandler:      settingHandler,
		NotificationHandler: notificationHandler,

		// Middleware
		AuthMiddleware: authMiddleware,
		CORSMiddleware: corsMiddleware,
	}
}
