package task

import (
	"context"
	"time"

	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
)

type CreateTaskRequest struct {
	ProjectID     int       `json:"project_id"`
	Title         string    `json:"title"`
	AssigneeEmail string    `json:"assignee_email"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	Priority      string    `json:"priority"`
	DueDate       time.Time `json:"due_date"`
}

type CreateTaskInput struct {
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	AssigneeID  int       `json:"assignee_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	ID          int64      `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
}
type TaskFilterRequest struct {
	ProjectID   *int64     `form:"project_id"`
	AssigneeID  *int64     `form:"assignee_id"`
	Statuses    []string   `form:"statuses"` // comma-separated in URL
	Priority    *string    `form:"priority"`
	DueDateFrom *time.Time `form:"due_date_from" time_format:"2006-01-02T15:04:05Z07:00"`
	DueDateTo   *time.Time `form:"due_date_to"   time_format:"2006-01-02T15:04:05Z07:00"`
	Limit       *int32     `form:"limit"`
	Offset      *int32     `form:"offset"`
}



type BulkTaskService interface{
	CreateTask(ctx context.Context, taskInput CreateTaskInput) (*taskdb.Task, error)
	GetTasksByProjectID(ctx context.Context, projectID int) ([]taskdb.Task, error)
}