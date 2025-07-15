package main

// ProjectImportPayload represents the payload for a project import job message.
// Used for passing job details via RabbitMQ.
type ProjectImportPayload struct {
	JobID    string `json:"job_id"`   // Unique job identifier
	FileName string `json:"filename"` // Name of the file to import
	Type     string `json:"type"`     // Type of import (e.g., "excel")
	UserID   int64  `json:"user_id"`  // User ID associated with the job
}

// ExportJobPayload represents the payload for a project export job message.
// Used for passing job details via RabbitMQ.
type ExportJobPayload struct {
	JobID    string `json:"job_id"`   // Unique job identifier
	Filename string `json:"filename"` // Name of the file to export
	Type     string `json:"type"`     // Type of export (e.g., "excel")
	UserID   int64  `json:"user_id"`  // User ID associated with the job
}
