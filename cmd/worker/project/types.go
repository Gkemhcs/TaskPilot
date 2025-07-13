package main


type ProjectImportPayload struct {
	JobID    string `json:"job_id"`
	FileName string `json:"filename"`
	Type     string `json:"type"`
	UserID   int64    `json:"user_id"`
}

type ExportJobPayload struct {
	JobID    string `json:"job_id"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	UserID   int64  `json:"user_id"`
}