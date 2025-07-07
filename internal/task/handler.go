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
	
	"github.com/Gkemhcs/taskpilot/internal/user"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func NewTaskHandler(taskService TaskService,userService user.UserResolver ,logger *logrus.Logger)*TaskHandler{
return &TaskHandler{
	taskService: taskService,
	userService: userService,
	logger:logger,
}

}

type TaskHandler struct {
	logger *logrus.Logger
	taskService  TaskService
	userService user.UserResolver

}

func RegisterTaskRoutes(router *gin.RouterGroup,taskHandler *TaskHandler,jwtManager *auth.JWTManager){
	taskRouter:=router.Group("/tasks",middleware.JWTAuthMiddleware(taskHandler.logger,jwtManager))
	{
		taskRouter.POST("/",taskHandler.CreateTask)
		taskRouter.GET("/:id",taskHandler.GetTaskByID)
		taskRouter.DELETE("/:id",taskHandler.DeleteTask)

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
	ctx1,cancel:=context.WithTimeout(c.Request.Context(),5*time.Second)
	defer cancel()
	user,err:=t.userService.GetUserByEmail(ctx1,createTaskRequest.AssigneeEmail)
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