package task

import (
	"context"
	"database/sql"
	"errors"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	"github.com/Gkemhcs/taskpilot/internal/task/gen"
)


func NewTaskService(taskRepo taskdb.Querier)*TaskService{
	return &TaskService{
		taskRepository: taskRepo,
	}

}


type 	TaskService struct {
	taskRepository  taskdb.Querier
}

func getStatus(status string)(taskdb.TaskStatus){

	switch status {
	case "in_progress":
		return taskdb.TaskStatusINPROGRESS
	case "done":
		return taskdb.TaskStatusDONE
	default: 
		return taskdb.TaskStatusTODO
	}

}

func getPriority(priority string )(taskdb.TaskPriority){
	switch priority {
	case "critical":
		return taskdb.TaskPriorityCRITICAL
	case "high":
		return taskdb.TaskPriorityHIGH
	case "low":
		return taskdb.TaskPriorityLOW
	default:
		return taskdb.TaskPriorityMEDIUM
	}
}


func (t *TaskService) CreateTask(ctx context.Context,taskInput CreateTaskInput)(*taskdb.Task,error){
	params:=taskdb.CreateTaskParams{
		ProjectID: int64(taskInput.ProjectID),
		Title: taskInput.Title,
		DueDate: sql.NullTime{Time:taskInput.DueDate,Valid:true},
		Description: taskInput.Description,
		AssigneeID:sql.NullInt64{Int64:int64(taskInput.AssigneeID),Valid:true } ,
		Status: getStatus(taskInput.Status),
		Priority: getPriority(taskInput.Priority),
	}
	task,err:=t.taskRepository.CreateTask(ctx,params)
	if  err!=nil{
		return nil,err
	}
	return &task,nil 
}

func(t *TaskService) GetTaskByID(ctx context.Context,taskID int)(*taskdb.Task,error){
	task,err:=t.taskRepository.GetTaskById(ctx,int64(taskID))
	if errors.Is(err,sql.ErrNoRows){
		return nil,customErrors.ErrTaskNotFound
	}
	if err!=nil{
		return nil,err
	}
	return &task,nil

}

func(t *TaskService) DeleteTask(ctx context.Context,taskID int)error{
	err:=t.taskRepository.DeleteTask(ctx,int64(taskID))
	if errors.Is(err, sql.ErrNoRows) {
		return customErrors.ErrTaskNotFound

	}
	if err!=nil{
		return err
	}
	return nil 


}


