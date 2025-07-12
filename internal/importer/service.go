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
)

// internal/import/service.go

type ImportService struct {
    storage      storage.StorageClient      // Local or GCS
    repo         importerdb.Querier          // From sqlc
    publisher    Publisher   // RabbitMQ publisher
}

func NewImportService(storage storage.StorageClient, repo importerdb.Querier, publisher Publisher) *ImportService {
    return &ImportService{storage: storage, repo: repo, publisher: publisher}
}

func (s *ImportService) ImportProjectExcel(ctx context.Context,file multipart.File, fileName string,userID int) (string, error) {
     err := s.storage.Upload(file, fileName)
    if err != nil {
        return "", customErrors.ErrUploadingFile
    }

    importID := uuid.New()
    
    params:=importerdb.CreateImportJobParams{
        ID:           importID,
        FilePath:     fileName,
        ImporterType: importerdb.ImportJobTypeProjectExcel,
        Status:       importerdb.ImportJobStatusPending,
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
    if err := s.publisher.PublishImportJob(ctx,msg); err != nil {
        return "", customErrors.ErrWhileEnqueuingImportJob
    }

    return importID.String(), nil
}

