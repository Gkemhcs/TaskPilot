package task

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

func(t *TaskService) GetTasksByProjectID(ctx context.Context,projectID int)([]taskdb.Task,error){
	tasks,err:=t.taskRepository.GetTasksByProjectId(ctx,int64(projectID))
	if err!=nil{
		return nil,err
	}
	return tasks,nil 
}
func(t *TaskService) DeleteTask(ctx context.Context,taskID int)error{
	rows,err:=t.taskRepository.DeleteTask(ctx,int64(taskID))
	if rows==0 {
		return customErrors.ErrTaskNotFound

	}
	if err!=nil{
		return err
	}
	return nil 


}


func(t *TaskService) GetAllTasks(ctx context.Context)([]taskdb.Task,error){

	tasks,err:=t.taskRepository.GetAllTasks(ctx)
	if err!=nil{
		return nil,err
	}
	return tasks,nil
}


func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}


func(t *TaskService)UpdateTask(ctx context.Context,req UpdateTaskRequest)(error){
	updateParams:=taskdb.UpdateTaskParams{
		ID: req.ID,
		Title: sql.NullString{
			String: deref(req.Title),
			Valid:  req.Title != nil,
		},
		Description: sql.NullString{
			String: deref(req.Description),
			Valid:  req.Description != nil,
		},
		DueDate: sql.NullTime{
			Time:  derefTime(req.DueDate),
			Valid: req.DueDate != nil,
		},
		Status: taskdb.NullTaskStatus{
			TaskStatus: taskdb.TaskStatus(deref(req.Status)),
			Valid:      req.Status != nil,
		},
		Priority: taskdb.NullTaskPriority{
			TaskPriority: taskdb.TaskPriority(deref(req.Priority)),
			Valid:        req.Priority != nil,
		},
	}
	rows,err:=t.taskRepository.UpdateTask(ctx,updateParams)
	if rows==0{
		return customErrors.ErrTaskNotFound
	}
	if err!=nil{
		return err
	}
	return nil 
}


