package exporter

import "context"

// Exporter defines the interface for exporting data to a file (e.g., Excel, CSV).
// Implementations should handle opening a file, adding rows, and saving the file to disk.
type Exporter interface {
	// Open initializes the export file with the given filename.
	Open(filename string) error
	// AddRow appends a row of data to the export file.
	AddRow(row []any) error
	// Save writes the export file to the specified local directory and returns the file path.
	Save(localDir string) (string, error)
}

// Publisher abstracts the logic for publishing export jobs to a message queue (e.g., RabbitMQ).
// Used to trigger asynchronous export processing.
type Publisher interface {
	// PublishExportJob sends an export job message to the queue for background processing.
	PublishExportJob(ctx context.Context, job ExportJobMessage) error
}

// ExportTaskRequest represents the request payload for exporting tasks of a specific project.
type ExportTaskRequest struct {
	ProjectID int `json:"project_id"`
}

// ExportJobMessage is the message structure sent to the queue for an export job.
// Contains metadata needed for processing and tracking the export.
type ExportJobMessage struct {
	JobID     string `json:"job_id"`     // Unique identifier for the export job
	Filename  string `json:"filename"`   // Name of the file to be generated
	Type      string `json:"type"`       // Type of export (e.g., "project_excel", "task_excel")
	UserID    int64  `json:"user_id"`    // ID of the user requesting the export
	ProjectID int64  `json:"project_id"` // ID of the project (if applicable)
}
