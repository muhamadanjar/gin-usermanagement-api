package app

func (s *Server) setupRoutes() {
	bc := s.businessContainer

	// Public routes
	public := s.router.Group("/")

	public.POST("/auth/login", bc.AuthHandler.Login)
	public.POST("/auth/register", bc.AuthHandler.Register)

	// Protected routes
	api := s.router.Group("/")
	api.Use(bc.AuthMiddleware.RequireAuth())

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.GET("/permissions", bc.AuthHandler.GetUserPermissions)
		auth.POST("/model-permissions", bc.AuthHandler.CreateModelPermission)
		auth.GET("/model-permissions", bc.AuthHandler.GetModelPermissions)
		auth.GET("/info", bc.AuthHandler.GetUser)
		auth.POST("/metas", bc.AuthHandler.CreateMeta)
		auth.GET("/metas", bc.AuthHandler.GetUserMeta)
	}

	// User routes
	users := api.Group("/users").Use(bc.AuthMiddleware.RequireRole("admin"))
	{
		users.GET("", bc.UserHandler.GetAllUsers)
		users.POST("", bc.UserHandler.CreateUser)
		users.GET("/:id", bc.UserHandler.GetUser)
		users.PUT("/:id", bc.UserHandler.UpdateUser)
		users.DELETE("/:id", bc.UserHandler.DeleteUser)
		users.POST("/:id/roles", bc.UserHandler.AssignRoles)
	}

	// Role routes
	roles := api.Group("/roles")
	{
		roles.GET("", bc.RoleHandler.GetAllRoles)
		roles.POST("", bc.RoleHandler.CreateRole)
		roles.GET("/:id", bc.RoleHandler.GetRole)
		roles.PUT("/:id", bc.RoleHandler.UpdateRole)
		roles.DELETE("/:id", bc.RoleHandler.DeleteRole)
		roles.POST("/:id/permissions", bc.RoleHandler.AssignPermissions)
	}

	// Permission routes
	permissions := api.Group("/permissions")
	{
		permissions.GET("", bc.PermissionHandler.GetAllPermissions)
		permissions.POST("", bc.PermissionHandler.CreatePermission)
		permissions.GET("/:id", bc.PermissionHandler.GetPermission)
		permissions.PUT("/:id", bc.PermissionHandler.UpdatePermission)
		permissions.DELETE("/:id", bc.PermissionHandler.DeletePermission)
	}

	// Menu routes
	menus := api.Group("/menus")
	{
		menus.GET("", bc.MenuHandler.GetAllMenus)
		menus.GET("/active", bc.MenuHandler.GetActiveMenus)
		menus.POST("", bc.MenuHandler.CreateMenu)
		menus.GET("/:id", bc.MenuHandler.GetMenu)
		menus.PUT("/:id", bc.MenuHandler.UpdateMenu)
		menus.DELETE("/:id", bc.MenuHandler.DeleteMenu)
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
