package http

import (
	"usermanagement-api/internal/container"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires the entire HTTP surface onto the router. It is the only
// place that knows how handlers and middleware are combined per resource.
func RegisterRoutes(router *gin.Engine, bc *container.BusinessContainer) {
	// Public routes
	public := router.Group("/")

	public.POST("/auth/login", bc.AuthHandler.Login)
	public.POST("/auth/register", bc.AuthHandler.Register)
	public.POST("/auth/refresh", bc.AuthHandler.Refresh)

	// Protected routes
	api := router.Group("/")
	api.Use(bc.AuthMiddleware.RequireAuth())

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.GET("/permissions", bc.AuthHandler.GetUserPermissions)
		auth.GET("/info", bc.AuthHandler.GetUser)
		auth.POST("/metas", bc.AuthHandler.CreateMeta)
		auth.GET("/metas", bc.AuthHandler.GetUserMeta)
		auth.POST("/change-password", bc.AuthHandler.ChangePassword)
		auth.PUT("/profile", bc.AuthHandler.UpdateProfile)
		auth.GET("/token-history", bc.AuthHandler.GetMyTokenHistory)
	}
	api.GET("/logout", bc.AuthHandler.Logout)

	// User routes
	users := api.Group("/users").Use(bc.AuthMiddleware.RequireRole("admin"))
	{
		users.GET("", bc.UserHandler.GetAllUsers)
		users.POST("", bc.UserHandler.CreateUser)
		users.GET("/:id", bc.UserHandler.GetUser)
		users.PATCH("/:id", bc.UserHandler.UpdateUser)
		users.DELETE("/:id", bc.UserHandler.DeleteUser)
		users.POST("/assign-role/:id", bc.UserHandler.AssignRole)
		users.PUT("/:id/roles", bc.UserHandler.SyncRoles)
		users.POST("/:id/roles", bc.UserHandler.AssignRoles)
		users.GET("/:id/metas", bc.UserHandler.GetUserMeta)
		users.POST("/:id/update-avatar", bc.UserHandler.UpdateUserAvatar)
		users.GET("/:id/token-history", bc.UserHandler.GetUserTokenHistory)
	}

	// Role routes
	roles := api.Group("/roles")
	{
		roles.GET("", bc.RoleHandler.GetAllRoles)
		roles.POST("", bc.RoleHandler.CreateRole)
		roles.GET("/:id", bc.RoleHandler.GetRole)
		roles.PATCH("/:id", bc.RoleHandler.UpdateRole)
		roles.DELETE("/:id", bc.RoleHandler.DeleteRole)
		roles.GET("/:id/members", bc.RoleHandler.GetRoleMembers)
		roles.POST("/:id/permissions", bc.RoleHandler.AssignPermissions)
	}

	// Permission routes
	permissions := api.Group("/permissions")
	{
		permissions.GET("", bc.PermissionHandler.GetAllPermissions)
		permissions.POST("", bc.PermissionHandler.CreatePermission)
		permissions.GET("/:id", bc.PermissionHandler.GetPermission)
		permissions.PATCH("/:id", bc.PermissionHandler.UpdatePermission)
		permissions.DELETE("/:id", bc.PermissionHandler.DeletePermission)
	}

	// Menu routes
	menus := api.Group("/menus")
	{
		menus.GET("", bc.MenuHandler.GetAllMenus)
		menus.GET("/active", bc.MenuHandler.GetActiveMenus)
		menus.POST("", bc.MenuHandler.CreateMenu)
		menus.GET("/:id", bc.MenuHandler.GetMenu)
		menus.PATCH("/:id", bc.MenuHandler.UpdateMenu)
		menus.DELETE("/:id", bc.MenuHandler.DeleteMenu)
		menus.POST("/:id/assign-permissions", bc.MenuHandler.AssignMenuPermissions)
		menus.GET("/permissions", bc.MenuHandler.GetMenuPermissions)
	}

	userMeta := api.Group("/user-meta")
	{
		userMeta.POST("", bc.UserMetaHandler.CreateOrUpdate)
		userMeta.GET("/:user_id", bc.UserMetaHandler.GetAllByUserID)
		userMeta.GET("/:user_id/:key", bc.UserMetaHandler.GetByKey)
		userMeta.DELETE("/:user_id/:key", bc.UserMetaHandler.Delete)
	}

	// Setting routes
	settings := api.Group("/settings").Use(bc.AuthMiddleware.RequireRole("admin"))
	{
		settings.POST("", bc.SettingHandler.CreateOrUpdate)
		settings.GET("", bc.SettingHandler.GetAll)
		settings.GET("/:key", bc.SettingHandler.GetByKey)
		settings.DELETE("/:key", bc.SettingHandler.Delete)
	}

	notifications := api.Group("/notifications")
	notifications.Use(bc.AuthMiddleware.RequireAuth())
	{
		notifications.POST("/send-to-me", bc.NotificationHandler.SendToMe)

		// Admin/Superuser only routes
		adminNotif := notifications.Group("")
		adminNotif.Use(bc.AuthMiddleware.RequireRole("admin", "superuser"))
		{
			adminNotif.POST("/send", bc.NotificationHandler.SendNotification)
		}
	}
}
