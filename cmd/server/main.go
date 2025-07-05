package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	"github.com/Gkemhcs/taskpilot/internal/config"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/project"
	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	"github.com/Gkemhcs/taskpilot/internal/user"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// NewServer sets up the Gin router, registers routes, and starts the HTTP server.
// It takes configuration, logger, and database connection as input.
func NewServer(config *config.Config, logger *logrus.Logger, dbConn *sql.DB) error {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	// Create a new Gin router with default middleware (logger, recovery)
	router := gin.Default()

	// Create API v1 group with custom logger middleware
	v1 := router.Group("/api/v1", middleware.LoggerMiddleware(logger))
	// Initialize user service with database connection
	userService := user.NewUserService(userdb.New(dbConn))
	// Initialize JWT manager for authentication
	params := auth.CreateJwtManagerParams{
		AccessTokenDuration:  config.AccessTokenDuration,
		RefreshTokenDuration: config.RefreshTokenDuration,
		AccessTokenKey:       config.JWTAccessTokenSecret,
		RefreshTokenKey:      config.JWTRefreshTokenSecret,
	}
	jwtManager := auth.NewJWTManager(params)

	// Create user handler with service, logger, and JWT manager
	userHandler := user.NewUserHandler(userService, logger, jwtManager)
	// Register user-related routes under /api/v1/users
	user.RegisterRoutes(v1, userHandler)

	projectService:=project.NewProjectService(projectdb.New(dbConn))
	projectHandler := project.NewProjectHandler(logger,projectService)
	project.RegisterProjectRoutes(v1, projectHandler,jwtManager)

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	// Start the HTTP server on the configured host and port
	return router.Run(fmt.Sprintf("%s:%s", config.HOST, config.Port))
}
