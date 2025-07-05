package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims defines the custom JWT claims for authenticated users.
type UserClaims struct {
	UserID               int    `json:"user_id"` // Unique user ID
	Username             string `json:"name"`    // Username of the user
	Email                string `json:"email"`
	jwt.RegisteredClaims        // Standard JWT claims (exp, iat, etc.)
}


type CreateJwtManagerParams struct{
	AccessTokenDuration time.Duration
	RefreshTokenDuration time.Duration
	AccessTokenKey  string 
	RefreshTokenKey string 
}
