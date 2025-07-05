package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// JWTAuthMiddleware checks for a valid JWT token in the Authorization header.
// If valid, it sets the user ID in the request context for downstream handlers.
func JWTAuthMiddleware(logger *logrus.Logger, jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Log and return error if Authorization header is missing
			logger.WithFields(logrus.Fields{
				"path":   c.FullPath(),
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			}).Warn("Missing Authorization header")
			utils.Error(c, http.StatusUnauthorized, customErrors.MISSING_AUTHORIZATION_HEADER.Error())
			return
		}

		// Split the header to extract the token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// Log and return error if header format is invalid
			logger.WithFields(logrus.Fields{
				"path":   c.FullPath(),
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
			}).Warn("Invalid Authorization header format")
			utils.Error(c, http.StatusUnauthorized, customErrors.INVALID_TOKEN_FORMAT.Error())
			return
		}

		tokenString := parts[1]
		// Verify the JWT token
		userClaims, err := jwtManager.Verify(tokenString)
		if errors.Is(err,customErrors.ErrTokenExpired){
			logger.Errorf("%v",customErrors.ErrTokenExpired)
			utils.Error(c,http.StatusBadRequest,customErrors.ErrTokenExpired.Error())
			return 
		}
		if err != nil {
			// Log and return error if token verification fails
			logger.WithFields(logrus.Fields{
				"path":   c.FullPath(),
				"method": c.Request.Method,
				"ip":     c.ClientIP(),
				"error":  err.Error(),
			}).Warn("Token verification failed")
			utils.Error(c, http.StatusUnauthorized, customErrors.INVALID_TOKEN.Error())
			return
		}

		// Attach user ID to request context for downstream handlers
		c.Set("userID", int(userClaims.UserID))
		c.Set("userName",userClaims.Username)
		c.Set("email",userClaims.Email)
		c.Next()
	}
}
