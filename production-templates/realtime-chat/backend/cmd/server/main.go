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

	"realtime-chat/internal/config"
	"realtime-chat/internal/domain"
	"realtime-chat/internal/handler"
	"realtime-chat/internal/middleware"
	"realtime-chat/internal/repository"
	"realtime-chat/internal/service"
	"realtime-chat/internal/websocket"
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
	roomRepo := repository.NewRoomRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, time.Duration(cfg.Auth.JWTExpiration)*time.Minute)
	roomService := service.NewRoomService(roomRepo, userRepo, messageRepo, hub)
	messageService := service.NewMessageService(messageRepo, roomRepo, userRepo, hub)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	roomHandler := handler.NewRoomHandler(roomService)
	messageHandler := handler.NewMessageHandler(messageService)
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
			"status":         "ok",
			"time":           time.Now().Unix(),
			"active_rooms":   hub.GetRoomCount(),
			"active_clients": hub.GetClientCount(),
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

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.Auth.JWTSecret))
		{
			// Profile routes
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)

			// User search
			protected.GET("/users/search", authHandler.SearchUsers)

			// Room routes
			rooms := protected.Group("/rooms")
			{
				rooms.GET("", roomHandler.GetUserRooms)
				rooms.POST("", roomHandler.Create)
				rooms.GET("/:id", roomHandler.GetByID)
				rooms.PUT("/:id", roomHandler.Update)
				rooms.DELETE("/:id", roomHandler.Delete)
				rooms.POST("/:id/archive", roomHandler.Archive)
				rooms.POST("/:id/leave", roomHandler.LeaveRoom)

				// Room participants
				rooms.GET("/:id/participants", roomHandler.GetParticipants)
				rooms.POST("/:id/participants", roomHandler.AddParticipant)
				rooms.DELETE("/:id/participants/:userId", roomHandler.RemoveParticipant)

				// Unread count and mark as read
				rooms.GET("/:id/unread", roomHandler.GetUnreadCount)
				rooms.POST("/:id/read", roomHandler.MarkAsRead)
			}

			// Direct message
			protected.POST("/direct", roomHandler.GetOrCreateDirectRoom)

			// Message routes
			messages := protected.Group("/rooms/:roomId/messages")
			{
				messages.GET("", messageHandler.GetRoomMessages)
				messages.POST("", messageHandler.Send)
			}

			protected.GET("/messages/:id", messageHandler.GetByID)
			protected.PUT("/messages/:id", messageHandler.Update)
			protected.DELETE("/messages/:id", messageHandler.Delete)

			// Message reactions
			protected.POST("/messages/:id/reactions", messageHandler.AddReaction)
			protected.DELETE("/messages/:id/reactions", messageHandler.RemoveReaction)

			// Read receipts
			protected.POST("/messages/:id/read", messageHandler.MarkAsRead)

			// Typing indicator
			protected.POST("/rooms/:roomId/typing", messageHandler.SendTypingIndicator)

			// WebSocket endpoint (requires auth)
			protected.GET("/ws/:roomId", wsHandler.HandleConnection)

			// WebSocket stats
			protected.GET("/ws/stats", wsHandler.GetStats)
			protected.GET("/rooms/:roomId/online", wsHandler.GetOnlineUsers)
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
	log.Printf("WebSocket endpoint: ws://localhost:%s/api/v1/ws/:roomId", cfg.Server.Port)
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
		&domain.Room{},
		&domain.Participant{},
		&domain.Message{},
		&domain.MessageReaction{},
		&domain.ReadReceipt{},
	)
}
