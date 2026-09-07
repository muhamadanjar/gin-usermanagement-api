package container

import (
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/middleware"
	"usermanagement-api/pkg/auth"

	"gorm.io/gorm"
)

// AuthzModule assembles the Authorizer use case and the HTTP middleware that
// adapts it to gin. This is the single place authentication/authorization is
// wired.
type AuthzModule struct {
	Middleware middleware.AuthMiddleware
}

func NewAuthzModule(db *gorm.DB, jwtService *auth.JWTService) *AuthzModule {
	authorizer := usecase.NewAuthorizer(
		gormrepo.NewUserRepository(db),
		gormrepo.NewRoleRepository(db),
		gormrepo.NewPermissionRepository(db),
		gormrepo.NewModelPermissionRepository(db),
		jwtService,
	)
	return &AuthzModule{Middleware: middleware.NewAuthMiddleware(authorizer)}
}
