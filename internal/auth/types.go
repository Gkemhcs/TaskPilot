package auth

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserID int  `json:"user_id"`
	Username   string `json:"role"`
	jwt.RegisteredClaims
}