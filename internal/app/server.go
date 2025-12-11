package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"usermanagement-api/internal/container"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router            *gin.Engine
	httpServer        *http.Server
	appContainer      *AppContainer
	businessContainer *container.BusinessContainer
}

// NewServer creates a new server instance with containers
func NewServer(appContainer *AppContainer, businessContainer *container.BusinessContainer) *Server {
	return &Server{
		appContainer:      appContainer,
		businessContainer: businessContainer,
	}
}

// Initialize sets up the router and routes
func (s *Server) Initialize() error {
	// Initialize router
	s.router = gin.Default()
	s.router.Use(s.businessContainer.CORSMiddleware.SetupCORS())

	// Setup routes
	s.setupRoutes()

	// Setup HTTP server using config from AppContainer
	addr := fmt.Sprintf("%s:%d", s.appContainer.Config.Server.Host, s.appContainer.Config.Server.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
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
	addr := fmt.Sprintf("%s:%d", s.appContainer.Config.Server.Host, s.appContainer.Config.Server.Port)
	log.Printf("Server is running on %s", addr)
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	return nil
}
