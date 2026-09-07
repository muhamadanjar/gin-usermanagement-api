package container

import (
	"usermanagement-api/domain/ports"
	gormrepo "usermanagement-api/infrastructure/database/repositories"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/presentation/http/middleware"

	"gorm.io/gorm"
)

// AuthzModule assembles the Authorizer use case and the HTTP middleware that
// adapts it to gin. This is the single place authentication/authorization is
// wired.
type AuthzModule struct {
	Middleware middleware.AuthMiddleware
}

func NewAuthzModule(db *gorm.DB, tokenManager ports.TokenManager) *AuthzModule {
	authorizer := usecase.NewAuthorizer(
		gormrepo.NewUserRepository(db),
		gormrepo.NewRoleRepository(db),
		tokenManager,
	)
	return &AuthzModule{Middleware: middleware.NewAuthMiddleware(authorizer)}
}
