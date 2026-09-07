package usecase

import (
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/pkg/auth"

	"github.com/google/uuid"
)

// Authorizer is the single use case behind every authorization decision in the
// HTTP layer. It turns a bearer token into a Principal (authentication) and
// answers permission questions against that Principal (authorization). No
// handler or middleware touches repositories directly.
type Authorizer interface {
	ResolvePrincipal(token string) (*entities.Principal, error)
	CheckPermission(principal *entities.Principal, modelType string, modelID uuid.UUID, permissionName string) (bool, error)
}

type authorizer struct {
	userRepo            repositories.UserRepository
	roleRepo            repositories.RoleRepository
	permissionRepo      repositories.PermissionRepository
	modelPermissionRepo repositories.ModelPermissionRepository
	jwtService          *auth.JWTService
}

func NewAuthorizer(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
	modelPermissionRepo repositories.ModelPermissionRepository,
	jwtService *auth.JWTService,
) Authorizer {
	return &authorizer{
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		permissionRepo:      permissionRepo,
		modelPermissionRepo: modelPermissionRepo,
		jwtService:          jwtService,
	}
}

// ResolvePrincipal validates the bearer token and loads the user together with
// their roles and the permissions those roles grant.
func (a *authorizer) ResolvePrincipal(token string) (*entities.Principal, error) {
	claims, err := a.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	roles, err := a.roleRepo.FindRolesByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]uuid.UUID, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	permissions, err := a.roleRepo.FindPermissionsByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	return &entities.Principal{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

// CheckPermission grants access to superusers, then to principals whose roles
// carry the named permission, and finally falls back to object-level
// model_permissions for the given model.
func (a *authorizer) CheckPermission(principal *entities.Principal, modelType string, modelID uuid.UUID, permissionName string) (bool, error) {
	if principal.IsSuperuser() {
		return true, nil
	}

	permission, err := a.permissionRepo.FindByName(permissionName)
	if err != nil {
		return false, nil // permission does not exist → deny
	}

	if principal.HasPermission(permission.Name) {
		return true, nil
	}

	return a.modelPermissionRepo.CheckPermission(modelType, modelID, permission.ID)
}
