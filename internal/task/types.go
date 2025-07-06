package task

import (
	"time"

	"golang.org/x/text/date"
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