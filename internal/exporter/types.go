package exporter

import "context"


type Exporter interface {
	
	Open(filename string)error 
	AddRow(row []any)error 
	Save(localDir string) (string, error) 

}


type Publisher interface {	
	PublishExportJob(ctx context.Context, job ExportJobMessage) error
}



type ExportTaskRequest struct {
	ProjectID int `json:"project_id"`
}

type ExportJobMessage struct {
	JobID    string `json:"job_id"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	UserID   int64    `json:"user_id"`
	ProjectID int64  `json:"project_id"`
}


