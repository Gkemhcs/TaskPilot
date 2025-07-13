package importer

import (
	"context"
	
)

type Importer interface {
	Import(path string , headers []string,userID int) error
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
