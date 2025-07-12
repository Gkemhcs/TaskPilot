package importer

import (
	"context"
	"mime/multipart"
)

type Importer interface {
	Import(file multipart.File) error
}

type ImportJobMessage struct {
	JobID    string `json:"job_id"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	UserID   int64    `json:"user_id"`
}

type Publisher interface {
	PublishImportJob(ctx context.Context, job ImportJobMessage) error
}
