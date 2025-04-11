package app

func (s *Server) setupRoutes() {
	// Public routes
	s.router.POST("/api/auth/login", s.authHandler.Login)
	s.router.POST("/api/auth/register", s.authHandler.Register)

	// Protected routes
	api := s.router.Group("/api")
	api.Use(s.authMiddleware.RequireAuth())

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.GET("/permissions", s.authHandler.GetUserPermissions)
		auth.POST("/model-permissions", s.authHandler.CreateModelPermission)
		auth.GET("/model-permissions", s.authHandler.GetModelPermissions)
	}

	// User routes
	users := api.Group("/users")
	{
		users.GET("", s.userHandler.GetAllUsers)
		users.POST("", s.userHandler.CreateUser)
		users.GET("/:id", s.userHandler.GetUser)
		users.PUT("/:id", s.userHandler.UpdateUser)
		users.DELETE("/:id", s.userHandler.DeleteUser)
		users.POST("/:id/roles", s.userHandler.AssignRoles)
	}

	// Role routes
	roles := api.Group("/roles")
	{
		roles.GET("", s.roleHandler.GetAllRoles)
		roles.POST("", s.roleHandler.CreateRole)
		roles.GET("/:id", s.roleHandler.GetRole)
		roles.PUT("/:id", s.roleHandler.UpdateRole)
		roles.DELETE("/:id", s.roleHandler.DeleteRole)
		roles.POST("/:id/permissions", s.roleHandler.AssignPermissions)
	}

	// Permission routes
	permissions := api.Group("/permissions")
	{
		permissions.GET("", s.permissionHandler.GetAllPermissions)
		permissions.POST("", s.permissionHandler.CreatePermission)
		permissions.GET("/:id", s.permissionHandler.GetPermission)
		permissions.PUT("/:id", s.permissionHandler.UpdatePermission)
		permissions.DELETE("/:id", s.permissionHandler.DeletePermission)
	}

	// Menu routes
	menus := api.Group("/menus")
	{
		menus.GET("", s.menuHandler.GetAllMenus)
		menus.GET("/active", s.menuHandler.GetActiveMenus)
		menus.POST("", s.menuHandler.CreateMenu)
		menus.GET("/:id", s.menuHandler.GetMenu)
		menus.PUT("/:id", s.menuHandler.UpdateMenu)
		menus.DELETE("/:id", s.menuHandler.DeleteMenu)
	}
}
