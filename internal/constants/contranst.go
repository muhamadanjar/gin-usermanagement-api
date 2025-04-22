package constants

// Context keys
const (
	UserIDKey      = "userID"
	UserRolesKey   = "userRoles"
	PermissionsKey = "permissions"
	AccessToken    = "access_token"
	RefreshToken   = "refresh_token"
)

// Authentication errors
const (
	ErrInvalidCredentials = "invalid credentials"
	ErrTokenExpired       = "token expired"
	ErrTokenInvalid       = "invalid token"
	ErrTokenMissing       = "token missing"
	ErrUnauthorized       = "unauthorized"
	ErrForbidden          = "forbidden"
)

const SettingsCacheKey = "global_settings"

// ModelTypes for ModelPermission
const (
	ModelTypeRole = "role"
	ModelTypeMenu = "menu"
)
