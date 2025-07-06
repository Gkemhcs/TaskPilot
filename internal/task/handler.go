package task

import (
	"net/http"

	"github.com/Gkemhcs/taskpilot/internal/auth"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
	"github.com/Gkemhcs/taskpilot/internal/user"
	"github.com/Gkemhcs/taskpilot/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func NewTaskHandler(taskService TaskService,userService user.IUserService ,logger *logrus.Logger)*TaskHandler{
return &TaskHandler{
	taskService: taskService,
	userService: userService,
	logger:logger,
}

}

type TaskHandler struct {
	logger *logrus.Logger
	taskService  TaskService
	userService user.IUserService

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
	if createTaskRequest.DueDate==""{
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







}

func (t *TaskHandler) GetTaskByID(c *gin.Context){

}

func(t *TaskHandler) GetTasksByProjectID(c *gin.Context){

}

func(t *TaskHandler) DeleteTask(c *gin.Context){

}