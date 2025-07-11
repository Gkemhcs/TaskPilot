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
		userGroup.POST("/refresh", handler.GenerateAccessTokenFromRefreshToken)
	}
}

type UserHandler struct {
	logger      *logrus.Logger
	userService IUserService
	jwtManager  *auth.JWTManager
}

// CreateUser handles user registration
// @Summary      Register a new user
// @Description  Creates a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "User registration input"
// @Success      201   {object}  utils.SuccessResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      500   {object}  utils.ErrorResponse
// @Router       /api/v1/users/ [post]
func (u *UserHandler) CreateUser(c *gin.Context) {

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Errorf("incorrect request body %s , %v", user.Name, err)
		utils.Error(c, http.StatusBadRequest, errors.INCORRECT_REQUEST_BODY.Error())
		return
	}
	if user.Email == "" {
		u.logger.Errorf("missing email in  request body %s , %v", user.Name, errors.ErrMissingEmail)
		utils.Error(c, http.StatusBadRequest, errors.ErrMissingEmail.Error())
		return
	}
	if user.Name == "" {
		u.logger.Errorf("missing name in  request body %s , %v", user.Name, errors.MISSING_USER_NAME)
		utils.Error(c, http.StatusBadRequest, errors.MISSING_USER_NAME.Error())
		return
	}
	if user.Password == "" {
		u.logger.Errorf("missing password in request body %s , %v", user.Name, errors.MISSING_PASSWORD)
		utils.Error(c, http.StatusBadRequest, errors.MISSING_PASSWORD.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)

	defer cancel()

	err := u.userService.CreateUser(ctx, user.Name, user.Email, user.Password)
	if err != nil {
		u.logger.Errorf("error while creating the user %v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	u.logger.Infof("user creation successful")
	utils.Success(c, http.StatusCreated, map[string]interface{}{
		"message":     "user created",
		"status_code": http.StatusCreated,
	})
}

func (u *UserHandler) LoginUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Errorf("incorrect request body %s , %v", user.Email, err)
		utils.Error(c, http.StatusBadRequest, errors.INCORRECT_REQUEST_BODY.Error())
		return
	}
	if user.Email == "" {
		u.logger.Errorf("missing email in  request body %s , %v", user.Name, errors.ErrMissingEmail)
		utils.Error(c, http.StatusBadRequest, errors.ErrMissingEmail.Error())
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	userInfo, err := u.userService.LoginUser(ctx, user.Email, user.Password)
	if err != nil {
		u.logger.Errorf("unable to login  %v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return

	}
	jwtTokenResponse, err := u.jwtManager.Generate(int(userInfo.ID), userInfo.Name, userInfo.Email)
	if err != nil {
		u.logger.Errorf("error while generating the token %v", err)

		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	u.logger.Infof("%s logged in and token generated successfully", user.Email)
	utils.Success(c, http.StatusOK, map[string]interface{}{

		"tokens": jwtTokenResponse,

		"message": "login successful",
	})

}

func (u *UserHandler) GenerateAccessTokenFromRefreshToken(c *gin.Context) {
	var request RefreshRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errors.INCORRECT_REQUEST_BODY.Error())
		return
	}
	accessToken, err := u.jwtManager.GenerateAccessTokenFromRefresh(request.RefreshToken)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, map[string]any{

		"token": accessToken,
	})

}
