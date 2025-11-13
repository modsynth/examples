package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"task-management-app/internal/config"
	"task-management-app/internal/domain"
	"task-management-app/internal/handler"
	"task-management-app/internal/middleware"
	"task-management-app/internal/repository"
	"task-management-app/internal/service"
	"task-management-app/internal/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schema
	if err := migrateDB(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create WebSocket hub and start it
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, time.Duration(cfg.Auth.JWTExpiration)*time.Minute)
	projectService := service.NewProjectService(projectRepo, userRepo)
	boardService := service.NewBoardService(boardRepo, projectRepo, hub)
	taskService := service.NewTaskService(taskRepo, boardRepo, projectRepo, hub)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	boardHandler := handler.NewBoardHandler(boardService)
	taskHandler := handler.NewTaskHandler(taskService)
	wsHandler := websocket.NewWebSocketHandler(hub)

	// Set gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// WebSocket endpoint (requires auth)
		v1.GET("/ws/:projectId", middleware.AuthMiddleware(cfg.Auth.JWTSecret), wsHandler.HandleConnection)

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.Auth.JWTSecret))
		{
			// Profile routes
			protected.GET("/profile", authHandler.GetProfile)

			// Project routes
			projects := protected.Group("/projects")
			{
				projects.GET("", projectHandler.List)
				projects.POST("", projectHandler.Create)
				projects.GET("/:id", projectHandler.GetByID)
				projects.PUT("/:id", projectHandler.Update)
				projects.DELETE("/:id", projectHandler.Delete)
				projects.POST("/:id/archive", projectHandler.Archive)
				projects.POST("/:id/unarchive", projectHandler.Unarchive)

				// Project members
				projects.GET("/:id/members", projectHandler.GetMembers)
				projects.POST("/:id/members", projectHandler.AddMember)
				projects.DELETE("/:id/members/:memberID", projectHandler.RemoveMember)
				projects.PUT("/:id/members/:memberID/role", projectHandler.UpdateMemberRole)

				// Project online users (WebSocket)
				projects.GET("/:projectId/online-users", wsHandler.GetOnlineUsers)
			}

			// Board routes
			boards := protected.Group("")
			{
				boards.POST("/projects/:projectID/boards", boardHandler.Create)
				boards.GET("/projects/:projectID/boards", boardHandler.ListByProject)
				boards.GET("/boards/:id", boardHandler.GetByID)
				boards.PUT("/boards/:id", boardHandler.Update)
				boards.DELETE("/boards/:id", boardHandler.Delete)
			}

			// Task routes
			tasks := protected.Group("")
			{
				tasks.POST("/boards/:boardID/tasks", taskHandler.Create)
				tasks.GET("/boards/:boardID/tasks", taskHandler.ListByBoard)
				tasks.GET("/tasks/:id", taskHandler.GetByID)
				tasks.PUT("/tasks/:id", taskHandler.Update)
				tasks.DELETE("/tasks/:id", taskHandler.Delete)
				tasks.POST("/tasks/:id/move", taskHandler.Move)

				// Task comments
				tasks.POST("/tasks/:id/comments", taskHandler.AddComment)
				tasks.DELETE("/tasks/:id/comments/:commentID", taskHandler.DeleteComment)

				// Task checklist
				tasks.POST("/tasks/:id/checklist", taskHandler.AddChecklistItem)
				tasks.PUT("/tasks/:id/checklist/:itemID", taskHandler.UpdateChecklistItem)
				tasks.DELETE("/tasks/:id/checklist/:itemID", taskHandler.DeleteChecklistItem)

				// Task labels
				tasks.POST("/tasks/:id/labels", taskHandler.AssignLabels)
			}
		}
	}

	// Start server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s (env: %s)", cfg.Server.Port, cfg.Server.Env)
	log.Printf("WebSocket endpoint: ws://localhost:%s/api/v1/ws/:projectId", cfg.Server.Port)
	log.Printf("Health check: http://localhost:%s/health", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Project{},
		&domain.ProjectMember{},
		&domain.Board{},
		&domain.Task{},
		&domain.Label{},
		&domain.Comment{},
		&domain.Attachment{},
		&domain.ChecklistItem{},
	)
}
