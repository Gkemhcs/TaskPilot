package middleware

import (
	"net/http"
	"strings"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	"github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// JWTAuthMiddleware checks for valid JWT token and sets user ID in context
func JWTAuthMiddleware(logger *logrus.Logger, jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.WithFields(logrus.Fields{
				"path":   c.FullPath(),
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			}).Warn("Missing Authorization header")

			utils.Error(c,http.StatusUnauthorized,errors.MISSING_AUTHORIZATION_HEADER.Error())
			return 
		
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.WithFields(logrus.Fields{
				"path":   c.FullPath(),
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			}).Warn("Invalid Authorization header format")
			utils.Error(c,http.StatusUnauthorized,errors.INVALID_TOKEN_FORMAT.Error())
			return 	
			
			
		}

		tokenString := parts[1]
		userClaims, err := jwtManager.Verify(tokenString)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"path":   c.FullPath(),
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
				"error":  err.Error(),
			}).Warn("Token verification failed")
			utils.Error(c,http.StatusUnauthorized,errors.INVALID_TOKEN.Error())
			return 
		}
		

		// Attach user ID to request context
		c.Set("userID", int(userClaims.UserID))

		c.Next()

	}
}
