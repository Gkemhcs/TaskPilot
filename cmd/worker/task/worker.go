package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Gkemhcs/taskpilot/internal/exporter"
	exporterdb "github.com/Gkemhcs/taskpilot/internal/exporter/gen"
	"github.com/Gkemhcs/taskpilot/internal/importer"
	importerdb "github.com/Gkemhcs/taskpilot/internal/importer/gen"

	"github.com/Gkemhcs/taskpilot/internal/storage"
	"github.com/Gkemhcs/taskpilot/internal/task"
)

type TaskWorker struct {
	Importer   importer.Importer
	Exporter   exporter.Exporter
	Storage    storage.StorageClient
	TaskSvc    task.BulkTaskService
	ImportRepo importerdb.Querier
	ExportRepo exporterdb.Querier
	Logger     *logrus.Logger
	Headers    []string // expected headers for validation
	SheetName  string
	LocalDir   string
}

func NewTaskWorker(
	importer importer.Importer,
	exporter exporter.Exporter,
	storage storage.StorageClient,
	taskSvc task.BulkTaskService,
	importRepo importerdb.Querier,
	exportRepo exporterdb.Querier,
	logger *logrus.Logger,
	headers []string,
	sheetName string,
	localDir string,
) *TaskWorker {
	return &TaskWorker{
		Importer:   importer,
		Exporter:   exporter,
		Storage:    storage,
		TaskSvc:    taskSvc,
		ImportRepo: importRepo,
		ExportRepo: exportRepo,
		Logger:     logger,
		Headers:    headers,
		SheetName:  sheetName,
		LocalDir:   localDir,
	}
}
func (w *TaskWorker) StartConsuming(ch *amqp.Channel, queueName string) error {
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	w.Logger.Infof("📥 Listening on queue: %s", queueName)

	for msg := range msgs {
		func() {
			var payload TaskImportPayload
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				w.Logger.Errorf("❌ Invalid message format: %v", err)
				msg.Nack(false, false)
				return
			}

			ctx := context.Background()
			w.Logger.Infof("📦 Job Received: %s", payload.JobID)

			localPath, err := w.Storage.Download(payload.FileName)

			// ✅ Always clean up the file (even on failure)
			defer func() {
				if err := w.Storage.Delete(payload.FileName); err != nil {
					w.Logger.Warnf("⚠️ Failed to delete file: %s", payload.FileName)
				}
			}()

			if err != nil {
				errMsg := fmt.Sprintf("Download failed: %v", err)
				w.Logger.Error(errMsg)
				w.ImportRepo.UpdateImportJobStatus(ctx, importerdb.UpdateImportJobStatusParams{
					ID:           uuid.MustParse(payload.JobID),
					UserID:       int32(payload.UserID),
					Status:       importerdb.ImportJobStatusFailed,
					ErrorMessage: sql.NullString{String: errMsg, Valid: true},
				})
				msg.Nack(false, false)
				return
			}

			err = w.Importer.Import(localPath, w.Headers, int(payload.UserID))
			if err != nil {
				errMsg := fmt.Sprintf("Import failed: %v", err)
				w.Logger.Error(errMsg)
				w.ImportRepo.UpdateImportJobStatus(ctx, importerdb.UpdateImportJobStatusParams{
					ID:           uuid.MustParse(payload.JobID),
					UserID:       int32(payload.UserID),
					Status:       importerdb.ImportJobStatusFailed,
					ErrorMessage: sql.NullString{String: errMsg, Valid: true},
				})
				msg.Nack(false, false)
				return
			}

			// ✅ Success path
			w.ImportRepo.UpdateImportJobStatus(ctx, importerdb.UpdateImportJobStatusParams{
				ID:           uuid.MustParse(payload.JobID),
				UserID:       int32(payload.UserID),
				Status:       importerdb.ImportJobStatusCompleted,
				ErrorMessage: sql.NullString{Valid: false},
			})
			msg.Ack(false)
			w.Logger.Infof("✅ Job Completed: %s", payload.JobID)
		}()
	}

	return nil
}

func (w *TaskWorker) StartConsumingExport(ch *amqp.Channel, queueName string) error {
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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			projects, err := w.TaskSvc.GetTasksByProjectID(ctx, int(payload.ProjectID))
			if err != nil {
				w.failExport(ctx, payload, err)
				msg.Nack(false, false)
				return
			}

			for _, t := range projects {
				row := []any{t.ID,t.ProjectID, t.Title,t.AssigneeID.Int64, t.Description, t.Status,t.Priority,t.DueDate.Time,t.CreatedAt,t.UpdatedAt}
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
			url, err := w.Storage.GenerateSignedURL(payload.Filename, 10*time.Minute)
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


func (w *TaskWorker) failExport(ctx context.Context, payload ExportJobPayload, err error) {
	w.ExportRepo.UpdateExportJobStatus(ctx, exporterdb.UpdateExportJobStatusParams{
		ID:           uuid.MustParse(payload.JobID),
		UserID:       int32(payload.UserID),
		Status:       exporterdb.ExportJobStatusFailed,
		ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
	})
	w.Logger.Errorf("❌ Export job failed: %s | %v", payload.JobID, err)
}
