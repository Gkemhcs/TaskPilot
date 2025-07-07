package task

import (
	"context"
	"time"

	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
)



type CreateTaskRequest struct {
   ProjectID int `json:"project_id"`
   Title string `json:"title"`
   AssigneeEmail string `json:"assignee_email"`
   Description string `json:"description"`
   Status string `json:"status"`
   Priority  string `json:"priority"`
   DueDate  time.Time `json:"due_date"`
}

type CreateTaskInput struct {
	ProjectID int `json:"project_id"`
   Title string `json:"title"`
   AssigneeID int `json:"assignee_id"`
   Description string `json:"description"`
   Status string `json:"status"`
   Priority  string `json:"priority"`
   DueDate  time.Time `json:"due_date"`
}


type TaskQueryService interface {
   GetTasksByProjectID(ctx context.Context,projectID int)([]taskdb.Task,error)
}