package user

import (
	"context"
	"net/http"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	"github.com/Gkemhcs/taskpilot/internal/errors"

	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewUserHandler(userService IUserService, logger *logrus.Logger, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
		jwtManager:  jwtManager,
	}

}

func RegisterRoutes(rg *gin.RouterGroup, handler *UserHandler) {
	userGroup := rg.Group("/users")
	{
		userGroup.POST("/", handler.CreateUser)
		userGroup.POST("/login", handler.LoginUser)
	}
}

type UserHandler struct {
	logger      *logrus.Logger
	userService IUserService
	jwtManager  *auth.JWTManager
}

func (u *UserHandler) CreateUser(c *gin.Context) {

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Errorf("incorrect request body %s , %v", user.Name, err)
		utils.Error(c, http.StatusBadRequest, errors.INCORRECT_REQUEST_BODY.Error())

	}
	if user.Name == "" {
		u.logger.Errorf("missing name in  request body %s , %v", user.Name, errors.MISSING_USER_NAME)
		utils.Error(c, http.StatusBadRequest, errors.MISSING_USER_NAME.Error())

	}
	if user.Password == "" {
		u.logger.Errorf("missing password in request body %s , %v", user.Name, errors.MISSING_PASSWORD)
		utils.Error(c, http.StatusBadRequest, errors.MISSING_PASSWORD.Error())

	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := u.userService.CreateUser(ctx, user.Name, user.Password)
	if err != nil {
		u.logger.Errorf("error while creating the user %v",err)
		utils.Error(c, http.StatusBadRequest, err.Error())

	}
	u.logger.Infof("user creation successful")
	utils.Success(c, http.StatusCreated, map[string]interface{}{
		"message":     "user crated",
		"status_code": http.StatusCreated,
	})
}

func (u *UserHandler) LoginUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Errorf("incorrect request body %s , %v", user.Name, err)
		utils.Error(c, http.StatusBadRequest, errors.INCORRECT_REQUEST_BODY.Error())
		return 
	}
	if user.Name == "" {
		u.logger.Errorf("missing name in  request body %s , %v", user.Name, errors.MISSING_USER_NAME)
		utils.Error(c, http.StatusBadRequest, errors.MISSING_USER_NAME.Error())
		return
	}
	if user.Password == "" {
		u.logger.Errorf("missing password in  request body %s , %v", user.Name, errors.MISSING_PASSWORD)
		utils.Error(c, http.StatusBadRequest, errors.MISSING_PASSWORD.Error())
		return 
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userInfo, err := u.userService.LoginUser(ctx, user.Name, user.Password)
	if err != nil {
		u.logger.Errorf("unable to login  %v",err )
		utils.Error(c, http.StatusBadRequest, err.Error())
		return 

	}
	tokenString, err := u.jwtManager.Generate(int(userInfo.ID), userInfo.Name)
	if err != nil {
		u.logger.Errorf("error while generating the token %v",err)
		
		utils.Error(c, http.StatusBadRequest, err.Error())
		return 
	}
	u.logger.Infof("%s logged in and token generated successfully", user.Name)
	utils.Success(c, http.StatusOK, map[string]interface{}{

		"token": tokenString,

		"message": "login successful",
	})

}
