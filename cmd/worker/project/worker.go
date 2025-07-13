package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"os"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Gkemhcs/taskpilot/internal/exporter"
	exporterdb "github.com/Gkemhcs/taskpilot/internal/exporter/gen"
	"github.com/Gkemhcs/taskpilot/internal/importer"
	importerdb "github.com/Gkemhcs/taskpilot/internal/importer/gen"
	"github.com/Gkemhcs/taskpilot/internal/project"
	"github.com/Gkemhcs/taskpilot/internal/storage"
)

type ProjectWorker struct {
	Importer      importer.Importer
	Exporter      exporter.Exporter
	Storage       storage.StorageClient
	ProjectSvc    project.ProjectService // unified service for import & export
	ImportRepo    importerdb.Querier
	ExportRepo    exporterdb.Querier
	Logger        *logrus.Logger
	Headers       []string
	SheetName     string
	LocalDir      string
}

func NewProjectWorker(
	importer importer.Importer,
	exporter *exporter.ExcelExporter,
	storage storage.StorageClient,
	projectSvc project.ProjectService,
	importRepo importerdb.Querier,
	exportRepo exporterdb.Querier,
	logger *logrus.Logger,
	headers []string,
	sheetName string,
	localDir string,
) *ProjectWorker {
	return &ProjectWorker{
		Importer:   importer,
		Exporter:   exporter,
		Storage:    storage,
		ProjectSvc: projectSvc,
		ImportRepo: importRepo,
		ExportRepo: exportRepo,
		Logger:     logger,
		Headers:    headers,
		SheetName:  sheetName,
		LocalDir:   localDir,
	}
}





func (w *ProjectWorker) StartConsuming(ch *amqp.Channel, queueName string) error {
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	w.Logger.Infof("📥 Listening for import jobs on queue: %s", queueName)
	for msg := range msgs {
		go func(msg amqp.Delivery) {
			var payload ProjectImportPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				w.Logger.Errorf("❌ Invalid import message format: %v", err)
				msg.Nack(false, false)
				return
			}
			ctx := context.Background()
			w.Logger.Infof("📦 Import Job Received: %s", payload.JobID)
			localPath, err := w.Storage.Download(payload.FileName)
			defer func() {
				if err := w.Storage.Delete(payload.FileName); err != nil {
					w.Logger.Warnf("⚠️ Failed to delete: %s", payload.FileName)
				}
			}()
			if err != nil {
				w.failImport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			err = w.Importer.Import(localPath, w.Headers, int(payload.UserID))
			if err != nil {
				w.failImport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			w.ImportRepo.UpdateImportJobStatus(ctx, importerdb.UpdateImportJobStatusParams{
				ID:           uuid.MustParse(payload.JobID),
				UserID:       int32(payload.UserID),
				Status:       importerdb.ImportJobStatusCompleted,
				ErrorMessage: sql.NullString{Valid: false},
			})
			msg.Ack(false)
			w.Logger.Infof("✅ Import Job Completed: %s", payload.JobID)
		}(msg)
	}
	return nil
}

func (w *ProjectWorker) StartConsumingExport(ch *amqp.Channel, queueName string) error {
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	w.Logger.Infof("📤 Listening for export jobs on queue: %s", queueName)
	for msg := range msgs {
		go func(msg amqp.Delivery) {
			var payload ExportJobPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				w.failExport(context.Background(), payload, err)
				msg.Nack(false, false)
				return
			}
			ctx := context.Background()
			w.Logger.Infof("📦 Export Job Received: %s", payload.JobID)
			if err := w.Exporter.Open(payload.Filename); err != nil {
				w.failExport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			
			ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
			defer cancel()
			projects, err := w.ProjectSvc.GetProjectsByUserId(ctx, int(payload.UserID))
			if err != nil {
				w.failExport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			for _, p := range projects {
				row := []any{p.ID, p.Name, p.Description.String, p.Color.ProjectColor, p.CreatedAt, p.UpdatedAt}
				if err := w.Exporter.AddRow(row); err != nil {
					w.failExport(ctx, payload, err)
					msg.Nack(false, false)
					return
				}
			}
			localPath, err := w.Exporter.Save(w.LocalDir)
			if err != nil {
				w.failExport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			f, err := os.Open(localPath)
			if err != nil {
				w.failExport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			defer f.Close()
			if err := w.Storage.Upload(f, payload.Filename); err != nil {
				w.failExport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			_ = os.Remove(localPath)
			url, err := w.Storage.GenerateSignedURL(payload.Filename,10*time.Minute)
			if err != nil {
				w.failExport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}
			w.ExportRepo.UpdateExportJobURL(ctx, exporterdb.UpdateExportJobURLParams{
				ID:     uuid.MustParse(payload.JobID),
				UserID: int32(payload.UserID),
				Url:    sql.NullString{String: url, Valid: true},
			})
			msg.Ack(false)
			w.Logger.Infof("✅ Export Job Completed: %s", payload.JobID)
		}(msg)
	}
	return nil
}

func (w *ProjectWorker) failImport(ctx context.Context, payload ProjectImportPayload, err error) {
	w.ImportRepo.UpdateImportJobStatus(ctx, importerdb.UpdateImportJobStatusParams{
		ID:           uuid.MustParse(payload.JobID),
		UserID:       int32(payload.UserID),
		Status:       importerdb.ImportJobStatusFailed,
		ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
	})
	w.Logger.Errorf("❌ Import job failed: %s | %v", payload.JobID, err)
}

func (w *ProjectWorker) failExport(ctx context.Context, payload ExportJobPayload, err error) {
	w.ExportRepo.UpdateExportJobStatus(ctx, exporterdb.UpdateExportJobStatusParams{
		ID:           uuid.MustParse(payload.JobID),
		UserID:       int32(payload.UserID),
		Status:       exporterdb.ExportJobStatusFailed,
		ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
	})
	w.Logger.Errorf("❌ Export job failed: %s | %v", payload.JobID, err)
}
