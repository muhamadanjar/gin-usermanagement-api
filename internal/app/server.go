package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/delivery/http/handlers"
	"usermanagement-api/internal/delivery/http/middleware"
	"usermanagement-api/internal/usecase"
	"usermanagement-api/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Server struct {
	router            *gin.Engine
	db                *gorm.DB
	httpServer        *http.Server
	userHandler       *handlers.UserHandler
	roleHandler       *handlers.RoleHandler
	permissionHandler *handlers.PermissionHandler
	menuHandler       *handlers.MenuHandler
	authHandler       *handlers.AuthHandler
	authMiddleware    middleware.AuthMiddleware
	corsMiddleware    middleware.CORSMiddleware
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Initialize() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Connect to database
	s.db = database.ConnectDB()

	// Migrate database
	database.MigrateDB(s.db)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(s.db)
	roleRepo := repositories.NewRoleRepository(s.db)
	permissionRepo := repositories.NewPermissionRepository(s.db)
	menuRepo := repositories.NewMenuRepository(s.db)
	modelPermissionRepo := repositories.NewModelPermissionRepository(s.db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)
	menuUseCase := usecase.NewMenuUseCase(menuRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, roleRepo, modelPermissionRepo)

	// Initialize middleware
	s.authMiddleware = middleware.NewAuthMiddleware(userRepo, roleRepo, permissionRepo, modelPermissionRepo)
	s.corsMiddleware = middleware.NewCORSMiddleware()

	// Initialize handlers
	s.userHandler = handlers.NewUserHandler(userUseCase)
	s.roleHandler = handlers.NewRoleHandler(roleUseCase)
	s.permissionHandler = handlers.NewPermissionHandler(permissionUseCase)
	s.menuHandler = handlers.NewMenuHandler(menuUseCase)
	s.authHandler = handlers.NewAuthHandler(authUseCase)

	// Initialize router
	s.router = gin.Default()

	// Setup routes
	s.setupRoutes()

	// Setup HTTP server
	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	return nil
}

func (s *Server) Run() error {
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Println("Server is running on :8080")
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	return nil
}
