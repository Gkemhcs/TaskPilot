package project

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/types"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewProjectHandler(logger *logrus.Logger, projectService *ProjectService, taskService types.TaskQueryService) *ProjectHandler {
	return &ProjectHandler{
		logger:           logger,
		projectService:   projectService,
		taskQueryService: taskService,
	}
}

type ProjectHandler struct {
	logger           *logrus.Logger
	projectService   *ProjectService
	taskQueryService types.TaskQueryService
}

func RegisterProjectRoutes(r *gin.RouterGroup, handler *ProjectHandler, jwtManager *auth.JWTManager) {

	projectGroup := r.Group("/projects", middleware.JWTAuthMiddleware(handler.logger, jwtManager))
	{
		projectGroup.POST("/", handler.CreateProject)
		projectGroup.GET("/:id", handler.GetProjectById)
		projectGroup.GET("/", handler.GetProjectsByUserId)
		projectGroup.PUT("/:id", handler.UpdateProject)
		projectGroup.DELETE("/:id", handler.DeleteProject)
		projectGroup.GET("/names/", handler.GetProjectByName)
		projectGroup.GET("/:id/tasks", handler.GetTasksByProjectID)
	}
}

// CreateProject handles the creation of a new project for the authenticated user.
// @Summary      Create a new project
// @Description  Creates a new project for the authenticated user
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project  body      Project  true  "Project creation input"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /api/v1/projects/ [post]
// @Security BearerAuth

func (p *ProjectHandler) CreateProject(c *gin.Context) {

	var project Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		p.logger.Errorf("%v", err)

		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	val, exists := c.Get("userID")
	if !exists {
		p.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusInternalServerError, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		p.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}
	project.User = userID
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	proj, err := p.projectService.CreateProject(ctx, project)
	if errors.Is(err, customErrors.ErrProjectAlreadyExists) {
		p.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		p.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadGateway, err.Error())
		return
	}
	p.logger.Infof("Project named %s created successfully for %d", project.Name, project.User)
	utils.Success(c, http.StatusCreated, map[string]interface{}{
		"data":    proj,
		"message": "successfully created",
		"code":    http.StatusCreated,
	})

}

// GetProjectById retrieves a project by its ID.
// @Summary      Get project by ID
// @Description  Retrieves a project by its unique ID
// @Tags         projects
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/v1/projects/{id} [get]
// @Security BearerAuth
func (p *ProjectHandler) GetProjectById(c *gin.Context) {
	idParam := c.Param("id")
	projectId, err := strconv.Atoi(idParam)
	if err != nil {
		p.logger.Errorf("%v", customErrors.ErrInvalidProjectId)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrInvalidProjectId.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	
	project, err := p.projectService.GetProjectById(ctx, projectId)
	if err != nil {
		p.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p.logger.Infof("Request Succeeded for %d", projectId)
	utils.Success(c, http.StatusOK, map[string]interface{}{
		"data":    project,
		"message": "request succeeded",
		"code":    http.StatusOK,
	})
}

// GetProjectByName retrieves a project by its name for the authenticated user.
// @Summary      Get project by name
// @Description  Retrieves a project by its name for the authenticated user
// @Tags         projects
// @Produce      json
// @Param        name  query     string  true  "Project Name"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Router       /api/v1/projects/names/ [get]
// @Security BearerAuth
func (p *ProjectHandler) GetProjectByName(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		p.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusBadGateway, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		p.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusInternalServerError, "invalid user ID type")
		return
	}
	projectName := c.Query("name")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	project, err := p.projectService.GetProjectByName(ctx, projectName, userID)
	if err != nil {
		p.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, map[string]any{
		"data":    project,
		"message": "request succeeded",
	})

}

// GetProjectsByUserId retrieves all projects for the authenticated user.
// @Summary      Get all projects for user
// @Description  Retrieves all projects for the authenticated user
// @Tags         projects
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/v1/projects/ [get]
// @Security BearerAuth
func (p *ProjectHandler) GetProjectsByUserId(c *gin.Context) {

	val, exists := c.Get("userID")
	if !exists {
		p.logger.Errorf("%v", customErrors.ErrUserIDNotFoundInContext)
		utils.Error(c, http.StatusBadRequest, "unauthenticated: user ID not found")
		return
	}

	userID, ok := val.(int)
	if !ok {
		p.logger.Errorf("%v", customErrors.ErrInvalidUserId)
		utils.Error(c, http.StatusBadRequest, "invalid user ID type")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	projects, err := p.projectService.GetProjectsByUserId(ctx, userID)
	if err != nil {
		p.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p.logger.Infof("Request for Projects succeeded for %d", userID)
	utils.Success(c, http.StatusOK, map[string]interface{}{
		"data":    projects,
		"message": "request suceeeded",
		"code":    http.StatusOK,
	})

}

// UpdateProject updates an existing project by its ID.
// @Summary      Update project
// @Description  Updates an existing project by its ID
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id      path      int     true  "Project ID"
// @Param        project body      Project true  "Project update input"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Router       /api/v1/projects/{id} [put]
// @Security BearerAuth
func (p *ProjectHandler) UpdateProject(c *gin.Context) {

}

// DeleteProject deletes a project by its ID.
// @Summary      Delete project
// @Description  Deletes a project by its unique ID
// @Tags         projects
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/v1/projects/{id} [delete]
// @Security BearerAuth
func (p *ProjectHandler) DeleteProject(c *gin.Context) {
	idParam := c.Param("id")
	projectId, err := strconv.Atoi(idParam)
	if err != nil {
		p.logger.Errorf("%v", customErrors.ErrInvalidProjectId)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrInvalidProjectId.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = p.projectService.DeleteProject(ctx, projectId)
	if err != nil {
		p.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p.logger.Infof("Project %d deleted successfully", projectId)
	utils.Success(c, http.StatusOK, map[string]interface{}{
		"message": "project delete successfully",
		"code":    http.StatusOK,
	})

}

// GetTasksByProjectID retrieves all tasks for a given project ID.
// @Summary      Get tasks by project ID
// @Description  Retrieves all tasks for a given project ID
// @Tags         projects
// @Produce      json
// @Param        id   path      int  true  "Project ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /api/v1/projects/{id}/tasks [get]
// @Security BearerAuth
func (p *ProjectHandler) GetTasksByProjectID(c *gin.Context) {
	id := c.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		p.logger.Errorf("%v", customErrors.ErrInvalidProjectId)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrInvalidProjectId.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	tasks, err := p.taskQueryService.GetTasksByProjectID(ctx, projectID)
	if err != nil {
		p.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p.logger.Infof("%v", tasks)
	if len(tasks) == 0 {
		utils.Success(c, http.StatusOK, map[string]any{
			"data":    "no tasks in the specified project",
			"message": "request succeeded successfully",
		})
		return
	}

	if tasks == nil {
		utils.Success(c, http.StatusOK, map[string]any{
			"data":    []int{},
			"message": "projects are empty",
		})
	}
	utils.Success(c, http.StatusOK, map[string]any{
		"data":    tasks,
		"message": "request succeeded successfully",
	})

}
