package http

import (
	"log"
	"zeneye-gateway/internal/adapter/http/handlers"
	"zeneye-gateway/internal/adapter/http/middlewares"
	"zeneye-gateway/internal/adapter/repository/postgres"
	"zeneye-gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter sets up the Gin router with routes and middleware.
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Middleware setup for logging and rate limiting
	router.Use(middlewares.LoggingMiddleware())
	router.Use(middlewares.RateLimitingMiddleware())

	logger.LogInfo("SetupRouter", "Initializing routes", "Setting up health check route", "")
	router.GET("/health", handlers.HealthCheck)

	logger.LogInfo("SetupRouter", "Initializing routes", "Setting up public routes", "")
	// Public routes
	router.POST("/login", handlers.Login(db))
	router.POST("/refresh-token", handlers.RefreshToken(db))

	logger.LogInfo("SetupRouter", "Initializing routes", "Setting up superadmin routes", "")
	// Superadmin Routes
	superadminGroup := router.Group("/superadmin")
	{
		superadminGroup.GET("/check", handlers.CheckSuperadmin(db))
		superadminGroup.POST("/create", handlers.CreateSuperadmin(db))
	}

	logger.LogInfo("SetupRouter", "Initializing routes", "Setting up protected routes", "")
	// Protected routes for gateway-handled APIs
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middlewares.AuthMiddleware())
	{
		// User Routes
		userGroup := protectedRoutes.Group("/users")
		{
			userGroup.POST("/", handlers.CreateUser(db))
			userGroup.PATCH("/:id", handlers.EditUser(db))
			userGroup.DELETE("/:id", handlers.DeleteUser(db))
			userGroup.GET("/:id", handlers.GetUser(db))
			userGroup.GET("/", handlers.ListUsers(db))
		}
	}

	// Initialize user repository for middleware
	userRepo := postgres.NewUserRepository(db)

	// Casting to concrete type *postgres.UserRepository
	postgresUserRepo, ok := userRepo.(*postgres.UserRepository)
	if !ok {
		log.Fatalf("Failed to assert userRepo to *postgres.UserRepository")
	}

	// Protected routes for microservices (no need for handlers, will be handled by middleware)
	microserviceRoutes := router.Group("/")
	microserviceRoutes.Use(middlewares.AuthMiddleware())
	microserviceRoutes.Use(middlewares.MicroserviceRoutingMiddleware(postgresUserRepo))
	{
		// Admin Management Routes
		adminGroup := microserviceRoutes.Group("/admin-management")
		{
			adminGroup.GET("/get-all")
		}
	}

	return router
}
