package project

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewProjectHandler(logger *logrus.Logger, projectService *ProjectService) *ProjectHandler {
	return &ProjectHandler{
		logger:         logger,
		projectService: projectService,
	}
}

type ProjectHandler struct {
	logger         *logrus.Logger
	projectService *ProjectService
}

func RegisterProjectRoutes(r *gin.RouterGroup, handler *ProjectHandler,jwtManager *auth.JWTManager) {

	projectGroup := r.Group("/projects",middleware.JWTAuthMiddleware(handler.logger,jwtManager))
	{
		projectGroup.POST("/", handler.CreateProject)
		projectGroup.GET("/:id", handler.GetProjectById)
		projectGroup.GET("/", handler.GetProjectsByUserId)
		projectGroup.PUT("/:id", handler.UpdateProject)
		projectGroup.DELETE("/:id", handler.DeleteProject)
		projectGroup.GET("/names/",handler.GetProjectByName)
	}
}

func (p *ProjectHandler) CreateProject(c *gin.Context) {

	var project Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		p.logger.Errorf("%v",err)

		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	val, exists := c.Get("userID")
	if !exists {
		p.logger.Errorf("%v",customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		p.logger.Errorf("%v",customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}
	project.User = userID
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	proj, err := p.projectService.CreateProject(ctx, project)
	if err != nil {
		p.logger.Errorf("%v",err)
		utils.Error(c, http.StatusBadGateway, err.Error())
		return
	}
	p.logger.Infof("Project named %s created successfully for %d",project.Name,project.User)
	utils.Success(c, http.StatusCreated, map[string]interface{}{
		"data":    proj,
		"message": "successfully created",
		"code":    http.StatusCreated,
	})

}

func (p *ProjectHandler) GetProjectById(c *gin.Context) {
	idParam := c.Param("id")
	projectId, err := strconv.Atoi(idParam)
	if err != nil {
		p.logger.Errorf("%v",customErrors.ErrInvalidProjectId)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrInvalidProjectId.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	project, err := p.projectService.GetProjectById(ctx, projectId)
	if err != nil {
		p.logger.Errorf("%v",err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p.logger.Infof("Request Succeeded for %d",projectId)
	utils.Success(c, http.StatusOK, map[string]interface{}{
		"data":    project,
		"message": "request succeeded",
		"code":    http.StatusOK,
	})
}

func (p *ProjectHandler) GetProjectByName(c *gin.Context){
	val, exists := c.Get("userID")
	if !exists {
		p.logger.Errorf("%v",customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusBadGateway, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		p.logger.Errorf("%v",customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}
	projectName:=c.Query("name")
	ctx,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	project,err:=p.projectService.GetProjectByName(ctx,projectName,userID)
	if err!=nil{
		p.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return
	}
	utils.Success(c,http.StatusOK,map[string]any{
		"data":project ,
		"message":"request succeeded",
	})

}
func (p *ProjectHandler) GetProjectsByUserId(c *gin.Context) {

	val, exists := c.Get("userID")
	if !exists {
		p.logger.Errorf("%v",customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusBadGateway, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		p.logger.Errorf("%v",customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	projects, err := p.projectService.GetProjectsByUserId(ctx, userID)
	if err != nil {
		p.logger.Errorf("%v",err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p.logger.Infof("Request for Projects succeeded for %d",userID)
	utils.Success(c, http.StatusOK, map[string]interface{}{
		"data":    projects,
		"message": "request suceeeded",
		"code":    http.StatusOK,
	})

}

func (p *ProjectHandler) UpdateProject(c *gin.Context) {

}

func (p *ProjectHandler) DeleteProject(c *gin.Context) {
	idParam := c.Param("id")
	projectId, err := strconv.Atoi(idParam)
	if err != nil {
		p.logger.Errorf("%v",customErrors.ErrInvalidProjectId)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrInvalidProjectId.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = p.projectService.DeleteProject(ctx, projectId)
	if err != nil {
		p.logger.Errorf("%v",err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p.logger.Infof("Project %d deleted successfully",projectId)
	utils.Success(c, http.StatusOK, map[string]interface{}{
		"message": "project delete successfully",
		"code":    http.StatusOK,
	})

}
