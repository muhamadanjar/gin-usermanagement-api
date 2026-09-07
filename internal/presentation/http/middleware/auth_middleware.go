package middleware

import (
	"net/http"
	"strings"
	"usermanagement-api/domain/entities"
	"usermanagement-api/internal/application/usecase"
	"usermanagement-api/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware is a thin HTTP adapter. All authentication/authorization
// logic lives in the Authorizer use case; this package only translates HTTP
// concerns (headers, status codes, context keys).
type AuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
	RequirePermission(modelType string, modelID uuid.UUID, permissionName string) gin.HandlerFunc
	RequireRole(roles ...string) gin.HandlerFunc
	RequireSuperuser() gin.HandlerFunc
}

type authMiddleware struct {
	authorizer usecase.Authorizer
}

func NewAuthMiddleware(authorizer usecase.Authorizer) AuthMiddleware {
	return &authMiddleware{authorizer: authorizer}
}

func (m *authMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrTokenMissing})
			return
		}

		principal, err := m.authorizer.ResolvePrincipal(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}

		if !principal.User.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			return
		}

		c.Set(constants.AccessToken, token)
		c.Set(constants.PrincipalKey, principal)
		c.Set(constants.UserIDKey, principal.User.ID)
		c.Next()
	}
}

func (m *authMiddleware) RequirePermission(modelType string, modelID uuid.UUID, permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal, ok := Principal(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}

		// The route parameter overrides the static model ID for object-level checks.
		if idParam := c.Param("id"); idParam != "" {
			if parsed, err := uuid.Parse(idParam); err == nil {
				modelID = parsed
			}
		}

		allowed, err := m.authorizer.CheckPermission(principal, modelType, modelID, permissionName)
		if err != nil || !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			return
		}
		c.Next()
	}
}

func (m *authMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal, ok := Principal(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}
		if !principal.IsSuperuser() && !principal.HasRole(roles...) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrForbidden})
			return
		}
		c.Next()
	}
}

func (m *authMiddleware) RequireSuperuser() gin.HandlerFunc {
	return func(c *gin.Context) {
		principal, ok := Principal(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}
		if !principal.IsSuperuser() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "superuser access required"})
			return
		}
		c.Next()
	}
}

func Principal(c *gin.Context) (*entities.Principal, bool) {
	v, ok := c.Get(constants.PrincipalKey)
	if !ok {
		return nil, false
	}
	p, ok := v.(*entities.Principal)
	return p, ok
}

func bearerToken(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
