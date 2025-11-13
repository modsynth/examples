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
	"github.com/modsynth/task-management-app/internal/api/handlers"
	"github.com/modsynth/task-management-app/internal/config"
	"github.com/modsynth/task-management-app/internal/domain"
	"github.com/modsynth/task-management-app/internal/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Set gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// Initialize handlers
	wsHandler := handlers.NewWebSocketHandler(hub)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// WebSocket endpoint
		v1.GET("/ws/:projectId", wsHandler.HandleConnection)
		v1.GET("/projects/:projectId/online-users", wsHandler.GetOnlineUsers)

		// Auth routes (placeholder)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
			auth.POST("/login", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
		}

		// Projects routes (placeholder)
		projects := v1.Group("/projects")
		{
			projects.GET("", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
			projects.POST("", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
			projects.GET("/:id", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
		}

		// Tasks routes (placeholder)
		v1.GET("/projects/:projectId/tasks", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
		v1.POST("/projects/:projectId/tasks", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
		v1.GET("/tasks/:id", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
		v1.PUT("/tasks/:id", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
		v1.DELETE("/tasks/:id", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
		v1.PUT("/tasks/:id/move", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"}) })
	}

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s (env: %s)", cfg.Server.Port, cfg.Server.Env)
	log.Printf("WebSocket endpoint: ws://localhost:%s/api/v1/ws/:projectId", cfg.Server.Port)

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
