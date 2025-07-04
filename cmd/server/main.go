package server

import (
	"fmt"
	"net/http"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	"github.com/Gkemhcs/taskpilot/internal/config"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/user"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NewServer sets up the Gin router, registers routes, and starts the HTTP server.
// It takes configuration, logger, and database connection as input.
func NewServer(config *config.Config, logger *logrus.Logger, dbConn *userdb.Queries) error {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	// Create a new Gin router with default middleware (logger, recovery)
	router := gin.Default()

	// Create API v1 group with custom logger middleware
	v1 := router.Group("/api/v1", middleware.LoggerMiddleware(logger))
	// Initialize user service with database connection
	userService := user.NewUserService(dbConn)
	// Initialize JWT manager for authentication
	jwtManager := auth.NewJWTManager(config.JWTSecret, config.AccessTokenDuration)

	// Create user handler with service, logger, and JWT manager
	userHandler := user.NewUserHandler(userService, logger, jwtManager)
	// Register user-related routes under /api/v1/users
	user.RegisterRoutes(v1, userHandler)

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	// Start the HTTP server on the configured host and port
	return router.Run(fmt.Sprintf("%s:%s", config.HOST, config.Port))
}
