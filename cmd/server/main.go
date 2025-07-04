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


func NewServer(config *config.Config, logger *logrus.Logger,dbConn *userdb.Queries) error {
	router := gin.Default() // Includes Logger and Recovery middleware
	gin.SetMode(gin.ReleaseMode)
   
    v1 := router.Group("/api/v1",middleware.LoggerMiddleware(logger))
    userService:=user.NewUserService(dbConn)
    jwtManager:=auth.NewJWTManager(config.JWTSecret,config.AccessTokenDuration)

    userHandler:=user.NewUserHandler(userService,logger,jwtManager)
    user.RegisterRoutes(v1,userHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return router.Run(fmt.Sprintf("%s:%s",config.HOST,config.Port))
	// Start server on port 8080
}
