package user

import (
	"bytes"
	"encoding/json"
	"errors"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func SetupNewUserHandler()(*UserHandler,*MockUserRepo){

	
	mockUserRepo:=new(MockUserRepo)
	userService:=NewUserService(mockUserRepo)

	logger := logrus.New()
	logger.SetOutput(io.Discard) // Avoid
	
	params:=auth.CreateJwtManagerParams{
		AccessTokenDuration: 10*time.Minute,
		RefreshTokenDuration: 10*time.Hour,
		AccessTokenKey: "qkaniqifiqfi",
		RefreshTokenKey: "fewnfewfnifnif",

	}
	jwtManager:=auth.NewJWTManager(params)
	userHandler:=NewUserHandler(userService,logger,jwtManager)
	return userHandler,mockUserRepo

	
}

func TestCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler,mockRepo:=SetupNewUserHandler()

	testCases:=[]struct{
		testName string
		requestBody map[string]string 
		mockSetup func(string)
		expectedServiceCall bool 
		expectedError error 
	}{
		{
			testName:"Valid User",
			requestBody :map[string]string{
				"name":"gkemhcs",
				"password":"dwfnie2e",
				"email":"gudi@gmail",
			},
			mockSetup: func(password string ){
				hashPassword:=getHashedPassword(password)
				
				mockRepo.On("CreateUser",mock.Anything,mock.Anything).Return(
					userdb.User{
						Name:"gkemhcs",
						HashedPassword:hashPassword ,
						Email: "gudi@gmail",

					},nil,
				)
			},
			expectedServiceCall: true ,
			expectedError: nil,
		},
		{
			testName: "Missing email",
			requestBody: map[string]string{
				"name":"gkemhcs",
				"password":"welcome1234",
			},
			mockSetup: func(password string){

			},
			expectedServiceCall: false,
			expectedError: customErrors.ErrMissingEmail,
		},{
			testName: "missing username",
			requestBody: map[string]string{
				"email":"gudik@gmail",
				"password":"e2ueh2i",
			},
			mockSetup: func(password string ){},
			expectedServiceCall: false,
			expectedError: customErrors.MISSING_USER_NAME,
		},{
			testName: "missing password",
			requestBody: map[string]string{
				"email":"gud@gmail",
				"name": "koti",
			},
			mockSetup: func(password string){},
			expectedServiceCall: false,
			expectedError: customErrors.MISSING_PASSWORD,
		},

	}
		for _ ,tc := range  testCases{


			t.Run(tc.testName,func(t *testing.T){
				mockRepo.Calls=nil
				mockRepo.ExpectedCalls=nil
				password:="efienfinfi"
				pass,ok:=tc.requestBody["password"]
				if ok{
					password=pass
				}
				tc.mockSetup(password)
				body, _ := json.Marshal(tc.requestBody)

				req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				// Set up full group routing like production
				r := gin.Default()
				apiGroup := r.Group("/api/v1")
				userGroup := apiGroup.Group("/users")
				userGroup.POST("/", handler.CreateUser)

				r.ServeHTTP(w, req)
				if tc.expectedError!=nil{
					assert.Error(t,tc.expectedError,w.Code)
				}else{
					assert.Equal(t,http.StatusCreated,w.Code)
				}


				if tc.expectedServiceCall{
					mockRepo.AssertCalled(t,"CreateUser",mock.Anything,mock.Anything)

				}else{
					mockRepo.AssertNotCalled(t,"CreateUser",mock.Anything,mock.Anything)

				}


			})
		}
	




}





func TestLoginHandler(t *testing.T){

	handler,mockRepo:=SetupNewUserHandler()


	testCases:=[]struct{
		testName string
		requestBody map[string]string 
		mockSetup func(string)
		expectedServiceCall bool 
		expectedError error 
	}{
		{
			testName: "Valid User",
			requestBody: map[string]string{
				"name": "koti",
				"email" : "gudi@gmail",
				"password": "gkem23",
			},
			mockSetup: func(password string) {
				mockRepo.On("GetUserByEmail",mock.Anything,"gudi@gmail").Return(
					userdb.User{
						Name: "koti",
						Email: "gudi@gmail",
						HashedPassword: getHashedPassword(password),
					},nil)
			},
			expectedError: nil ,
			expectedServiceCall: true,			
		},
		{
			testName: "Missing User Name",
			requestBody: map[string]string{
				"email":"gudi@gmail",
				"password":"gkem1",
			},
			mockSetup: func(password string){

			},
			expectedServiceCall: false,
			expectedError: customErrors.MISSING_USER_NAME,
		},{
			testName:"Missing Email",
			requestBody: map[string]string{
				"name":"gkemhcs",
				"password":"gkem1234",
			},
			mockSetup: func(password string){},
			expectedServiceCall: false,
			expectedError: customErrors.ErrMissingEmail,
		},
		{
			testName:"Missing Password",
			requestBody: map[string]string{
				"name":"gkemhcs",
				"email":"gudi@gmail",
			},
			mockSetup: func(password string){},
			expectedServiceCall: false,
			expectedError: customErrors.MISSING_PASSWORD,
		},
		{
			testName:"Non Existent User",
			requestBody: map[string]string{
				"name":"gkemhcs",
				"email":"gudi@gmail",
				"password":"gkem3u939",
			},
			mockSetup: func(password string){
				mockRepo.On("GetUserByEmail",mock.Anything,"gudi@gmail").Return(
					userdb.User{},customErrors.USER_NOT_FOUND,
				)
			},
			expectedServiceCall: true,
			expectedError: customErrors.USER_NOT_FOUND,
		},
		{
			testName:"Incorrect password",
			requestBody: map[string]string{
				"name":"gkemhcs",
				"email":"gudi@gmail",
				"password":"gkem3u939",
			},
			mockSetup: func(password string){
				mockRepo.On("GetUserByEmail",mock.Anything,"gudi@gmail").Return(
					userdb.User{
						Name:"gkemhcs",
						Email:"gudi@gmail",
						HashedPassword: getHashedPassword("wrongpassword"),
					},nil,
				)
			},
			expectedServiceCall: true,
			expectedError:customErrors.ErrMismatchedPassword,
			
		},
		{
			testName:"db connection failed",
			requestBody: map[string]string{
				"name":"gkemhcs",
				"email":"gudi@gmail",
				"password":"gkem3u939",
			},
			mockSetup: func(password string){
				mockRepo.On("GetUserByEmail",mock.Anything,"gudi@gmail").Return(
					userdb.User{},errors.New("db error"),
				)
			},
			expectedServiceCall: true,
			expectedError:errors.New("db error"),
			
		},
	}

	for _,tc := range testCases{
			t.Run(tc.testName,func(t *testing.T){

				mockRepo.ExpectedCalls = nil
				mockRepo.Calls = nil
				password:="jooojo"
				pass,ok:=tc.requestBody["password"]
				if ok {
					password=pass
				}


				tc.mockSetup(password)
				body,_:=json.Marshal(tc.requestBody)
			
				req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				// Set up full group routing like production
				r := gin.Default()
				apiGroup := r.Group("/api/v1")
				userGroup := apiGroup.Group("/users")
				userGroup.POST("/login", handler.LoginUser)

				r.ServeHTTP(w, req)

			 
				if tc.expectedError!=nil{
					assert.Error(t,tc.expectedError,w.Code)

				}else{
						assert.Equal(t, http.StatusOK, w.Code)
				}
				
				
			



				if tc.expectedServiceCall{
					mockRepo.AssertCalled(t,"GetUserByEmail",mock.Anything,tc.requestBody["email"])
				}else{
					mockRepo.AssertNotCalled(t,"GetUserByEmail",mock.Anything,tc.requestBody["email"])
				}
				

			})
	}



}


func TestRefreshHandler(t *testing.T){
	handler,_:=SetupNewUserHandler()
	tokenResponse,_:=handler.jwtManager.Generate(123,"koti","eswar@gmail") // to generate access token and refresh token first

	testCases:=[]struct{
		testName string 
		requestBody map[string]string 		
		
		expectedStatus int

	}{
		{
			testName: "Valid Refresh Token",
			requestBody: map[string]string{
				"token":tokenResponse.RefreshToken,
			},

			expectedStatus: http.StatusOK,
		},
		{
			testName: "Invalid Refresh token",
			requestBody:map[string]string{
				"token":"kfwkefnififi",
			},
		
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _,tc := range testCases{
		t.Run(tc.testName,func(t *testing.T){

			body,_:=json.Marshal(tc.requestBody)


			req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Set up full group routing like production
			r := gin.Default()
			apiGroup := r.Group("/api/v1")
			userGroup := apiGroup.Group("/users")
			userGroup.POST("/refresh", handler.GenerateAccessTokenFromRefreshToken)

			r.ServeHTTP(w, req)
		
			assert.Equal(t,tc.expectedStatus,w.Code)
			

		})
	}


}