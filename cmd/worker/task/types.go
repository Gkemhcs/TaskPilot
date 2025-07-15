package main

// TaskImportPayload represents the payload for a task import job message.
// Used for passing job details via RabbitMQ.
type TaskImportPayload struct {
	JobID    string `json:"job_id"`   // Unique job identifier
	FileName string `json:"filename"` // Name of the file to import
	Type     string `json:"type"`     // Type of import (e.g., "excel")
	UserID   int64  `json:"user_id"`  // User ID associated with the job
}

// ExportJobPayload represents the payload for a task export job message.
// Used for passing job details via RabbitMQ.
type ExportJobPayload struct {
	JobID     string `json:"job_id"`     // Unique job identifier
	Filename  string `json:"filename"`   // Name of the file to export
	Type      string `json:"type"`       // Type of export (e.g., "excel")
	UserID    int64  `json:"user_id"`    // User ID associated with the job
	ProjectID int64  `json:"project_id"` // Project ID for which tasks are exported
}
