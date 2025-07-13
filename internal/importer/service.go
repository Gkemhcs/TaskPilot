// Package importer provides an interface for importing project data from Excel files.
package importer

import (
	"context"
	"mime/multipart"
	"time"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	importerdb "github.com/Gkemhcs/taskpilot/internal/importer/gen"
	"github.com/Gkemhcs/taskpilot/internal/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// internal/import/service.go

type ImportService struct {
    storage      storage.StorageClient      // Local or GCS
    repo         importerdb.Querier          // From sqlc
    projectPublisher    Publisher   // RabbitMQ publisher
    taskPublisher       Publisher   // RabbitMQ publisher
    logger *logrus.Logger
}

func NewImportService(storage storage.StorageClient, 
    repo importerdb.Querier, projectPublisher ,taskPublisher Publisher,
    logger *logrus.Logger) *ImportService {
    return &ImportService{storage: storage,
          repo: repo,
          projectPublisher: projectPublisher, 
          taskPublisher: taskPublisher, 
          logger: logger}
}

func (s *ImportService) ImportProjectExcel(ctx context.Context,file multipart.File, fileName string,userID int) (string, error) {
    err := s.storage.Upload(file, fileName)
    if err != nil {
        s.logger.Errorf("Error uploading file: %v", err)
        return "", customErrors.ErrUploadingFile
    }
    s.logger.Info("File uploaded successfully", "fileName", fileName)

    importID := uuid.New()
    
    params:=importerdb.CreateImportJobParams{
        ID:           importID,
        FilePath:     fileName,
        ImporterType: importerdb.ImportJobTypeProjectExcel,
        Status:       importerdb.ImportJobStatusPending,
        UserID:      int32(userID),
    }
    _,err=s.repo.CreateImportJob(ctx, params)
    if err!=nil{
        return "",customErrors.ErrCreatingImportJob
    }


    // Publish to queue
    msg := ImportJobMessage{
        JobID:  importID.String(),
        Filename:  fileName,
        Type:      "project_excel",
        UserID: int64(userID),
    }
    ctx,cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    if err := s.projectPublisher.PublishImportJob(ctx,msg); err != nil {
        return "", customErrors.ErrWhileEnqueuingImportJob
    }
    s.logger.Info("Project Import job enqueued successfully", "jobID", importID.String(), "fileName", fileName)

    return importID.String(), nil
}


func (s *ImportService) ImportTaskExcel(ctx context.Context,file multipart.File, fileName string,userID int) (string, error) {

    err := s.storage.Upload(file, fileName)
    if err != nil {
        s.logger.Errorf("Error uploading file: %v", err)
        return "", customErrors.ErrUploadingFile
    }
    s.logger.Info("File uploaded successfully", "fileName", fileName)

    importID := uuid.New()
    
    params:=importerdb.CreateImportJobParams{
        ID:           importID,
        FilePath:     fileName,
        ImporterType: importerdb.ImportJobTypeTaskExcel,
        Status:       importerdb.ImportJobStatusPending,
        UserID:      int32(userID),
    }
    _,err=s.repo.CreateImportJob(ctx, params)
    if err!=nil{
        return "",customErrors.ErrCreatingImportJob
    }


    // Publish to queue
    msg := ImportJobMessage{
        JobID:  importID.String(),
        Filename:  fileName,
        Type:      "task_excel",
        UserID: int64(userID),
    }
    ctx,cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    if err := s.taskPublisher.PublishImportJob(ctx,msg); err != nil {
        return "", customErrors.ErrWhileEnqueuingImportJob
    }
    s.logger.Info("Task Import job enqueued successfully", "jobID", importID.String(), "fileName", fileName)

    return importID.String(), nil


}


func(s *ImportService) Getstatus(ctx context.Context,jobId string,userID int)(*importerdb.ImportJob,error){

    uuidID, err := uuid.Parse(jobId)
    if err != nil {
        return nil, customErrors.ErrInvalidJobID
    }

    params := importerdb.GetImportJobParams{
        ID: uuidID,
        UserID: int32(userID),
    }

    job, err := s.repo.GetImportJob(ctx, params)
    if err != nil {
        return nil, err
    }
    return &job, nil
}
