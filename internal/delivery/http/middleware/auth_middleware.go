package middleware

import (
	"net/http"
	"strings"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/constants"
	"usermanagement-api/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
	RequirePermission(modelType string, modelID uuid.UUID, permissionName string) gin.HandlerFunc
	RequireRole(roles ...string) gin.HandlerFunc
	RequireSuperuser() gin.HandlerFunc
}

type authMiddleware struct {
	userRepo            repositories.UserRepository
	roleRepo            repositories.RoleRepository
	permissionRepo      repositories.PermissionRepository
	modelPermissionRepo repositories.ModelPermissionRepository
}

func NewAuthMiddleware(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
	modelPermissionRepo repositories.ModelPermissionRepository,
) AuthMiddleware {
	return &authMiddleware{
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		permissionRepo:      permissionRepo,
		modelPermissionRepo: modelPermissionRepo,
	}
}

func (m *authMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrTokenMissing})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		// Format: "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrTokenInvalid})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate the token
		claims, err := auth.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrTokenInvalid})
			c.Abort()
			return
		}

		// Get user from database
		user, err := m.userRepo.FindByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			c.Abort()
			return
		}

		// Get user roles
		roles, err := m.roleRepo.FindRolesByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user roles"})
			c.Abort()
			return
		}

		// Get role IDs
		var roleIDs []uuid.UUID
		for _, role := range roles {
			roleIDs = append(roleIDs, role.ID)
		}

		// Get permissions
		permissions, err := m.roleRepo.FindPermissionsByRoleIDs(roleIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user permissions"})
			c.Abort()
			return
		}

		c.Set(constants.AccessToken, tokenString)

		// Store user ID in the context
		c.Set(constants.UserIDKey, user.ID)

		// Store user roles in the context
		c.Set(constants.UserRolesKey, roles)

		// Store user permissions in the context
		c.Set(constants.PermissionsKey, permissions)

		c.Next()
	}
}

func (m *authMiddleware) RequirePermission(modelType string, modelID uuid.UUID, permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, exists := c.Get(constants.UserIDKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}
		user, err := m.userRepo.FindByID(userUUID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		// If user is superuser, allow access immediately
		if user.IsSuperuser {
			c.Next()
			return
		}

		// Get roles from context
		roles, exists := c.Get(constants.UserRolesKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		// Get permission by name
		permission, err := m.permissionRepo.FindByName(permissionName)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			c.Abort()
			return
		}

		// Check if user has permission
		hasPermission := false

		// Option 1: Check through model permissions
		for range roles.([]*entities.Role) {
			hasAccess, err := m.modelPermissionRepo.CheckPermission(modelType, modelID, permission.ID)
			if err == nil && hasAccess {
				hasPermission = true
				break
			}
		}

		// Option 2: Check through role permissions
		if !hasPermission {
			userRoles := roles.([]*entities.Role)
			for _, role := range userRoles {
				// Check if the role has the required permission
				for _, rolePermission := range role.Permissions {
					if rolePermission.ID == permission.ID {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *authMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(constants.UserIDKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		userUUID, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		user, err := m.userRepo.FindByID(userUUID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		// If user is superuser, allow access immediately
		if user.IsSuperuser {
			c.Next()
			return
		}

		// Get user roles from context
		userRoles, exists := c.Get(constants.UserRolesKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		userRolesList := userRoles.([]*entities.Role)
		hasRequiredRole := false

		for _, userRole := range userRolesList {
			for _, requiredRole := range roles {
				if strings.EqualFold(userRole.Name, requiredRole) {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *authMiddleware) RequireSuperuser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userInterface, exists := c.Get(constants.UserIDKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			c.Abort()
			return
		}

		user := userInterface.(*entities.User)

		// Check if user is superuser
		if !user.IsSuperuser {
			c.JSON(http.StatusForbidden, gin.H{"error": "superuser access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
