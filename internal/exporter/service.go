package exporter

import (
	"context"
	"time"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	exporterdb "github.com/Gkemhcs/taskpilot/internal/exporter/gen"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ExportService coordinates export job creation, publishing, and status tracking.
// It interacts with the database and message queue to manage export jobs for projects and tasks.
type ExportService struct {
	repo             exporterdb.Querier // SQLC-generated DB interface for export jobs
	projectPublisher Publisher          // Publishes project export jobs to RabbitMQ
	taskPublisher    Publisher          // Publishes task export jobs to RabbitMQ
	logger           *logrus.Logger     // Logger for error/info reporting
}

// NewExportService constructs an ExportService with DB, publishers, and logger.
func NewExportService(
	repo exporterdb.Querier, projectPublisher, taskPublisher Publisher,
	logger *logrus.Logger) *ExportService {
	return &ExportService{
		repo:             repo,
		projectPublisher: projectPublisher,
		taskPublisher:    taskPublisher,
		logger:           logger,
	}
}

// ExportProjectExcel creates a new export job for a project's data in Excel format.
// It saves the job in the DB and publishes it to the project export queue.
func (s *ExportService) ExportProjectExcel(ctx context.Context, fileName string, userID int) (string, error) {
	exportID := uuid.New()
	params := exporterdb.CreateExportJobParams{
		ID:         exportID,
		ExportType: exporterdb.ExportTypeProjectExcel,
		UserID:     int32(userID),
	}
	_, err := s.repo.CreateExportJob(ctx, params)
	if err != nil {
		s.logger.Error(err)
		return "", customErrors.ErrCreatingExportJob
	}

	// Prepare and publish the export job message to the queue
	msg := ExportJobMessage{
		JobID:    exportID.String(),
		Filename: fileName,
		Type:     "project_excel",
		UserID:   int64(userID),
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.projectPublisher.PublishExportJob(ctx, msg); err != nil {
		s.logger.Error(err)
		return "", customErrors.ErrWhileEnqueuingExportJob
	}
	s.logger.Info("Project Export job enqueued successfully", "jobID", exportID.String(), "fileName", fileName)
	return exportID.String(), nil
}

// ExportTaskExcel creates a new export job for tasks of a project in Excel format.
// It saves the job in the DB and publishes it to the task export queue.
func (s *ExportService) ExportTaskExcel(ctx context.Context, fileName string, userID int, projectid int) (string, error) {
	exportID := uuid.New()
	params := exporterdb.CreateExportJobParams{
		ID:         exportID,
		ExportType: exporterdb.ExportTypeTaskExcel,
		UserID:     int32(userID),
	}
	_, err := s.repo.CreateExportJob(ctx, params)
	if err != nil {
		return "", customErrors.ErrCreatingExportJob
	}

	// Prepare and publish the export job message to the queue
	msg := ExportJobMessage{
		JobID:     exportID.String(),
		Filename:  fileName,
		Type:      "task_excel",
		UserID:    int64(userID),
		ProjectID: int64(projectid),
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.taskPublisher.PublishExportJob(ctx, msg); err != nil {
		return "", customErrors.ErrWhileEnqueuingExportJob
	}
	s.logger.Info("Task Export job enqueued successfully", "jobID", exportID.String(), "fileName", fileName)
	return exportID.String(), nil
}

// GetExportStatus retrieves the status and details of an export job by job ID and user ID.
// Returns the export job record from the database.
func (s *ExportService) GetExportStatus(ctx context.Context, userId int, jobId string) (*exporterdb.ExportJob, error) {
	uuidID, err := uuid.Parse(jobId)
	if err != nil {
		return nil, customErrors.ErrInvalidJobID
	}
	params := exporterdb.GetExportJobStatusParams{
		ID:     uuidID,
		UserID: int32(userId),
	}
	exportJob, err := s.repo.GetExportJobStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return &exportJob, nil
}
