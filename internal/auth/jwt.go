package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager handles creation and verification of JWT tokens.
type JWTManager struct {
	
	accessTokenSecretKey     string    
	refreshTokenSecretKey string 
	refreshTokenDuration time.Duration    // Secret key used to sign tokens
	accessTokenDuration time.Duration // Duration for which the token is valid
}

// NewJWTManager creates a new JWTManager with the given secret key and token duration.
func NewJWTManager(params CreateJwtManagerParams) *JWTManager {
	return &JWTManager{
		accessTokenSecretKey: params.AccessTokenKey,
		refreshTokenDuration: params.RefreshTokenDuration,
		accessTokenDuration: params.AccessTokenDuration,
		refreshTokenSecretKey: params.RefreshTokenKey,
	}
}

// Generate creates a new JWT token for a user with the given userID and username.

func (j *JWTManager) Generate(userID int, username string,email string ) (string, error) {
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		Email : email ,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenDuration)), // Set token expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),                      // Set token issue time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessTokenSecretKey)) // Sign and return the token
}

// Verify parses and validates a JWT token string and returns the claims if valid.
func (j *JWTManager) Verify(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(j.accessTokenSecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
