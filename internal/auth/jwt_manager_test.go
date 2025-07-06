package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	"github.com/stretchr/testify/assert"
	
)

var jwtManager *auth.JWTManager



func getAccessTokenString(userID int,email string,userName string )(string){
	tokenString,err:=jwtManager.GenerateAccessToken(userID,userName,email)
	if err!=nil{
		return  ""
	}
	return tokenString 
}

func getRefreshTokenString(userID int,email string,userName string )(string){
	tokenString,err:=jwtManager.GenerateRefreshToken(userID,userName,email)
	if err!=nil{
		return  ""
	}
	return tokenString 
}



func TestMain(m *testing.M) {
	params := auth.CreateJwtManagerParams{
		AccessTokenDuration:  2 * time.Minute,
		RefreshTokenDuration: 24 * time.Hour,
		AccessTokenKey:       "ffj3if3ifjeifnefn",
		RefreshTokenKey:      "2ijij2ifi2fj32ii2ji",
	}

	jwtManager = auth.NewJWTManager(params)
	os.Exit(m.Run())
}


func TestGenerate(t *testing.T) {
	

	userID := 101
	username := "gkemhcs"
	email := "mani@gkemhcs.com"

	tokens, err := jwtManager.Generate(userID, username, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)

	// Now decode and verify
	claims, err := jwtManager.Verify(tokens.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)

	// And check refresh token too
	refreshClaims, err := jwtManager.VerifyRefreshToken(tokens.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, refreshClaims.UserID)
}

func TestGenerateAccessToken(t *testing.T){

	userID := 101
	username := "gkemhcs"
	email := "mani@gkemhcs.com"
	_,err:=jwtManager.GenerateAccessToken(userID,username,email)
	assert.NoError(t, err)
}




func TestGenerateAccessTokenFromRefresh(t *testing.T){
	testCases:=[]struct{
		name string 
		refreshToken string 
		expectedError bool 
	}{{
		"Valid Refresh Token",
		getRefreshTokenString(123,"gudi@gmail","eswar"),
		false ,
	},
	{
		"Invalid Refresh Token",
		"dknfinfknk",
		true,
	}}
	for _,test := range testCases {
		
		_,err:=jwtManager.GenerateAccessTokenFromRefresh(test.refreshToken)
		if test.expectedError {
			assert.Error(t, err, test.name)
		} else {
			assert.NoError(t, err, test.name)
		}
	} 
}
func TestVerify(t *testing.T){

	testCases := []struct{
		name string 
		tokenString string 
		expectedError bool
	}{
		{
			"Valid Access Token",
			getAccessTokenString(1234,"gudi@gmail","koti"),
			false,
		},
		{
			"Invalid Access Token",
			"nfnfi3f3mf",
			true,
		},
	}
	for _, test := range testCases {
		_, err := jwtManager.Verify(test.tokenString)
		if test.expectedError {
			assert.Error(t, err, test.name)
		} else {
			assert.NoError(t, err, test.name)
		}
	}


}

func TestVerifyRefreshToken(t *testing.T){

	testCases := []struct{
		name string 
		tokenString string 
		expectedError bool
	}{
		{
			"Valid Access Token",
			getRefreshTokenString(1234,"gudi@gmail","koti"),
			false,
		},
		{
			"Invalid Access Token",
			"nfnfi3f3mf",
			true,
		},
	}
	for _, test := range testCases {
		_, err := jwtManager.VerifyRefreshToken(test.tokenString)
		if test.expectedError {
			assert.Error(t, err, test.name)
		} else {
			assert.NoError(t, err, test.name)
		}
	}

	
}