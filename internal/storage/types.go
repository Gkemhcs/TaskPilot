package storage

import (
	"mime/multipart"
)

type StorageClient interface {
	Upload (file multipart.File, filename string) error
	Download (filename string) ( string, error)
	Delete (filename string) error
}



type StorageConfig struct {
	BucketName string
	Prefix     string
	ProcessDir string
	TempDir   string
}


