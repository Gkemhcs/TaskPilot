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


func NewTaskHandler(taskService TaskService,userService user.UserResolver ,
	logger *logrus.Logger,projectService types.ProjectReader)*TaskHandler{
return &TaskHandler{
	taskService: taskService,
	userService: userService,
	logger:logger,
	projectService:projectService,

}

}

type TaskHandler struct {
	logger *logrus.Logger
	taskService  TaskService
	userService user.UserResolver
	projectService types.ProjectReader
	
}

func RegisterTaskRoutes(router *gin.RouterGroup,taskHandler *TaskHandler,jwtManager *auth.JWTManager){
	taskRouter:=router.Group("/tasks",middleware.JWTAuthMiddleware(taskHandler.logger,jwtManager))
		{
			taskRouter.POST("/",taskHandler.CreateTask)
			taskRouter.GET("/:id",taskHandler.GetTaskByID)
			taskRouter.DELETE("/:id",taskHandler.DeleteTask)
			taskRouter.PATCH("/",taskHandler.UpdateTask)	

	}
}


func(t *TaskHandler) CreateTask(c *gin.Context){

	var createTaskRequest CreateTaskRequest 
	err:=c.ShouldBindJSON(&createTaskRequest)
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return 
	}
	if createTaskRequest.AssigneeEmail==""{
		t.logger.Errorf("%v",customErrors.ErrAssigneeMissingFromBody)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrAssigneeMissingFromBody.Error())
		return 
	}
	if createTaskRequest.DueDate.IsZero(){
		t.logger.Errorf("%v",customErrors.ErrMissingDueDate)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrMissingDueDate.Error())
		return 
	}
	if createTaskRequest.Priority==""{
		createTaskRequest.Priority="medium"
	}
	if createTaskRequest.Status==""{
		createTaskRequest.Status="todo"
	}
	if createTaskRequest.Title==""{
		t.logger.Errorf("%v",customErrors.ErrorTaskTitleMissing)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrorTaskTitleMissing.Error())
		return 	
	}
	if createTaskRequest.ProjectID==0{
		t.logger.Errorf("%v",customErrors.ErrMissingProjectID)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrMissingProjectID.Error())
		return 
	}
	ctx,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	_,err=t.projectService.GetProjectById(ctx,createTaskRequest.ProjectID)
	if errors.Is(err,customErrors.ErrProjectIDNotExist){
		t.logger.Errorf("%v",customErrors.ErrParentProjectIDNotFound)
		utils.Error(c,http.StatusOK,customErrors.ErrParentProjectIDNotFound.Error())
		return 
	
	}
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return 
	}

	ctx1,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	user,err:=t.userService.GetUserByEmail(ctx1,createTaskRequest.AssigneeEmail)
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return
	}
	createTaskInput:=CreateTaskInput{
		ProjectID: createTaskRequest.ProjectID,
		AssigneeID: int(user.ID),
		Title:createTaskRequest.Title,
		Status:createTaskRequest.Status,
		Priority: createTaskRequest.Priority,
		DueDate: createTaskRequest.DueDate,
		Description: createTaskRequest.Description,
	}
	ctx2,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	task,err:=t.taskService.CreateTask(ctx2,createTaskInput)
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return 
	}
	utils.Success(c,http.StatusCreated,map[string]any{
		"data":task ,
		"message":"task created successfully",
	})
}

func (t *TaskHandler) GetTaskByID(c *gin.Context){
	taskID:=c.Param("id")
	id,err:=strconv.Atoi(taskID)
	if err!=nil{
		t.logger.Errorf("%v",customErrors.ErrInvalidTaskID)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrInvalidTaskID.Error())
		return 
	}
	ctx,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	task,err:=t.taskService.GetTaskByID(ctx,id)
	if errors.Is(err,customErrors.ErrTaskNotFound){
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrTaskNotFound.Error())
		return 
	}
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return 
	}
	utils.Success(c,http.StatusOK,map[string]any{
		"data":task ,
		"message":"request succeeded",
	})

}



func(t *TaskHandler) DeleteTask(c *gin.Context){
	taskID:=c.Param("id")
	id,err:=strconv.Atoi(taskID)
	if err!=nil{
		t.logger.Errorf("%v",customErrors.ErrInvalidTaskID)
		utils.Error(c,http.StatusBadRequest,customErrors.ErrInvalidTaskID.Error())
		return 
	}
	ctx,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	err=t.taskService.DeleteTask(ctx,id)
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return 
	}
	utils.Success(c,http.StatusOK,map[string]any{
		"message":"task deleted successfully",
	})


}


func(t *TaskHandler) GetAllTasks(c *gin.Context){
	ctx,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	tasks,err:=t.taskService.GetAllTasks(ctx)
	if err!=nil{
		t.logger.Errorf("%v",err)
		utils.Error(c,http.StatusBadRequest,err.Error())
		return 
	}
	if len(tasks)==0{
		t.logger.Info("no tasks at current time")
		utils.Success(c,http.StatusOK,map[string]any{
			"data":[]any{},
			"message":"request succeeded",
		})
		return 
	}
	t.logger.Info("request succeeded")
	utils.Success(c,http.StatusOK,map[string]any{
		"data":tasks,
		"message":"request succeeded successfully",

	})

}


func(t *TaskHandler) UpdateTask(c *gin.Context){
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}	
	err := t.taskService.UpdateTask(c.Request.Context(), req)
	if err != nil {
		t.logger.Errorf("%v",err)
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
t.logger.Infof("task  %d updated successfully",req.ID)
utils.Success(c,http.StatusOK,map[string]any{
	"message":"task updated successfully",
})
}