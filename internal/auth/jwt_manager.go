// GenerateAccessTokenFromRefresh generates a new access token using a valid refresh token.
// It verifies the provided refresh token, extracts the user claims, and issues a new access token
// for the user. Returns the new access token string or an error if the refresh token is invalid.
package auth

import (
	"errors"
	"time"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
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

// Generate creates a new JWT Access token for a user with the given userID and username.
func (j *JWTManager)GenerateAccessToken(userID int,username string ,email string )(string,error){
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


// Generate creates a new JWT  Refreshtoken for a user with the given userID and username.
func (j *JWTManager)GenerateRefreshToken(userID int,username string ,email string )(string,error){
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		Email : email ,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenDuration)), // Set token expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),                      // Set token issue time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshTokenSecretKey)) // Sign and return the token
}

// Generate creates a new JWT Access token and Refresh token for a user with the given userID and username.
func (j *JWTManager) Generate(userID int, username string,email string ) (*GenerateJwtResponse, error) {
	refreshToken,err:=j.GenerateRefreshToken(userID,username,email)
	if err!=nil{
		return nil,err 
	}
	accessToken,err:=j.GenerateAccessToken(userID,username,email)
	if err!=nil{
		return nil,err
	}
	return &GenerateJwtResponse{
		RefreshToken: refreshToken,
		AccessToken: accessToken,
	},err


}



// GenerateAccessTokenFromRefresh generates a new access token from a valid refresh token.
// It verifies the refresh token, extracts user claims, and issues a new access token.
// Returns the new access token string or an error if the refresh token is invalid.
func (j *JWTManager) GenerateAccessTokenFromRefresh(refreshToken string)(string,error){
	userClaims,err:=j.VerifyRefreshToken(refreshToken)
	if err!=nil{
		return "",err
	}
	accessToken,err:=j.GenerateAccessToken(userClaims.UserID,userClaims.Username,userClaims.Email)
	if err!=nil{
		return "",err
	}
	return accessToken,nil 
}


// VerifyRefreshToken  verifies the jwt refresh token is valid  and returns claims if it is valid
// It checks the signing method and expiration, returning an error if the token is invalid.
// It returns the user claims if the token is valid.
func(j *JWTManager) VerifyRefreshToken(refreshToken string)(*UserClaims,error){

	token, err := jwt.ParseWithClaims(
		refreshToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(j.refreshTokenSecretKey), nil
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
// Verify parses and validates a JWT token string and returns the claims if valid.
// It checks the signing method and expiration, returning an error if the token is invalid.
// If the token is valid, it returns the user claims contained in the token.
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
	if errors.Is(err,jwt.ErrTokenExpired){
		return nil,customErrors.ErrTokenExpired
	}

	if err != nil {
		return nil, err
	}
	
	claims, ok := token.Claims.(*UserClaims)



	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	

	return claims, nil
}
