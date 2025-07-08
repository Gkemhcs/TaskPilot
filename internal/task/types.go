package task

import  "time"
	



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


type UpdateTaskRequest struct {
	ID          int64      `json:"id" binding:"required"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
}


