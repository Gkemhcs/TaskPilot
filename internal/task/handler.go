
package task

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

	"github.com/Gkemhcs/taskpilot/internal/user"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewTaskHandler(taskService TaskService, userService user.UserResolver,
	logger *logrus.Logger, projectService types.ProjectReader) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		userService:    userService,
		logger:         logger,
		projectService: projectService,
	}

}

type TaskHandler struct {
	logger         *logrus.Logger
	taskService    TaskService
	userService    user.UserResolver
	projectService types.ProjectReader
}

func RegisterTaskRoutes(router *gin.RouterGroup, taskHandler *TaskHandler, jwtManager *auth.JWTManager) {
	taskRouter := router.Group("/tasks", middleware.JWTAuthMiddleware(taskHandler.logger, jwtManager))
	{
		taskRouter.POST("/", taskHandler.CreateTask)
		taskRouter.GET("/:id", taskHandler.GetTaskByID)
		taskRouter.DELETE("/:id", taskHandler.DeleteTask)
		taskRouter.PATCH("/:id", taskHandler.UpdateTask)
		taskRouter.GET("/filter", taskHandler.FilterTasks)

	}
}

// @Summary      Create a new task
// @Description  Creates a new task for the authenticated user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      CreateTaskRequest true  "Task creation input"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/tasks/ [post]
// @Security BearerAuth
func (t *TaskHandler) CreateTask(c *gin.Context) {

	var createTaskRequest CreateTaskRequest
	err := c.ShouldBindJSON(&createTaskRequest)
	if err != nil {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if createTaskRequest.AssigneeEmail == "" {
		t.logger.Errorf("%v", customErrors.ErrAssigneeMissingFromBody)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrAssigneeMissingFromBody.Error())
		return
	}
	if createTaskRequest.DueDate.IsZero() {
		t.logger.Errorf("%v", customErrors.ErrMissingDueDate)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrMissingDueDate.Error())
		return
	}
	if createTaskRequest.Priority == "" {
		createTaskRequest.Priority = "medium"
	}
	if createTaskRequest.Status == "" {
		createTaskRequest.Status = "todo"
	}
	if createTaskRequest.Title == "" {
		t.logger.Errorf("%v", customErrors.ErrorTaskTitleMissing)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrorTaskTitleMissing.Error())
		return
	}
	if createTaskRequest.ProjectID == 0 {
		t.logger.Errorf("%v", customErrors.ErrMissingProjectID)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrMissingProjectID.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	_, err = t.projectService.GetProjectById(ctx, createTaskRequest.ProjectID)
	if errors.Is(err, customErrors.ErrProjectIDNotExist) {
		t.logger.Errorf("%v", customErrors.ErrParentProjectIDNotFound)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrParentProjectIDNotFound.Error())
		return

	}
	if err != nil {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx1, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	user, err := t.userService.GetUserByEmail(ctx1, createTaskRequest.AssigneeEmail)
	if err != nil {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	createTaskInput := CreateTaskInput{
		ProjectID:   createTaskRequest.ProjectID,
		AssigneeID:  int(user.ID),
		Title:       createTaskRequest.Title,
		Status:      createTaskRequest.Status,
		Priority:    createTaskRequest.Priority,
		DueDate:     createTaskRequest.DueDate,
		Description: createTaskRequest.Description,
	}
	ctx2, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	task, err := t.taskService.CreateTask(ctx2, createTaskInput)
	if err != nil {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, map[string]any{
		"data":    task,
		"message": "task created successfully",
	})
}

// @Summary      Get task by ID
// @Description  Retrieves a task by its unique ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/tasks/{id} [get]
// @Security BearerAuth
func (t *TaskHandler) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	id, err := strconv.Atoi(taskID)
	if err != nil {
		t.logger.Errorf("%v", customErrors.ErrInvalidTaskID)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrInvalidTaskID.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	task, err := t.taskService.GetTaskByID(ctx, id)
	if errors.Is(err, customErrors.ErrTaskNotFound) {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrTaskNotFound.Error())
		return
	}
	if err != nil {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, map[string]any{
		"data":    task,
		"message": "request succeeded",
	})

}

// @Summary      Delete task
// @Description  Deletes a task by its unique ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/tasks/{id} [delete]
// @Security BearerAuth
func (t *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	id, err := strconv.Atoi(taskID)
	if err != nil {
		t.logger.Errorf("%v", customErrors.ErrInvalidTaskID)
		utils.Error(c, http.StatusBadRequest, customErrors.ErrInvalidTaskID.Error())
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = t.taskService.DeleteTask(ctx, id)
	if err != nil {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, map[string]any{
		"message": "task deleted successfully",
	})

}

// @Summary      Get all tasks
// @Description  Retrieves all tasks for the authenticated user
// @Tags         tasks
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/tasks/ [get]
// @Security BearerAuth
func (t *TaskHandler) GetAllTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	tasks, err := t.taskService.GetAllTasks(ctx)
	if err != nil {
		t.logger.Errorf("%v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(tasks) == 0 {
		t.logger.Info("no tasks at current time")
		utils.Success(c, http.StatusOK, map[string]any{
			"data":    []any{},
			"message": "request succeeded",
		})
		return
	}
	t.logger.Info("request succeeded")
	utils.Success(c, http.StatusOK, map[string]any{
		"data":    tasks,
		"message": "request succeeded successfully",
	})

}

// @Summary      Update task
// @Description  Updates an existing task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      UpdateTaskRequest true  "Task update input"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/tasks/ [patch]
// @Security BearerAuth
func (t *TaskHandler) UpdateTask(c *gin.Context) {
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	id:=c.Param("id")
	taskID,err:=strconv.Atoi(id)
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrInvalidTaskID.Error())
		return 

	}
	req.ID=int64(taskID)

	err = t.taskService.UpdateTask(c.Request.Context(), req)
	if err != nil {
		t.logger.Errorf(" error is %v", err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	t.logger.Infof("task  %d updated successfully", req.ID)
	utils.Success(c, http.StatusOK, map[string]any{
		"message": "task updated successfully",
	})
}

// @Summary      Filter tasks
// @Description  Filters tasks based on query parameters
// @Tags         tasks
// @Produce      json
// @Param        filter query TaskFilterRequest false "Task filter parameters"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/v1/tasks/filter [get]
// @Security BearerAuth
func (h *TaskHandler) FilterTasks(c *gin.Context) {
	var req TaskFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Errorf("Invalid filter query params: %v", err)
		utils.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	tasks, err := h.taskService.FilterTasks(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorf("Failed to fetch filtered tasks: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Could not filter tasks")
		return
	}

	utils.Success(c, http.StatusOK, tasks)
}
