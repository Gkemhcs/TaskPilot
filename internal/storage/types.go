package storage

import (
	"mime/multipart"
	"time"
)


// for imports
type StorageClient interface {
	Upload (file multipart.File, filename string) error
	Download (filename string) ( string, error)
	Delete (filename string) error
	GenerateSignedURL(filename string, inExpires time.Duration) (string, error)
	
}






type StorageConfig struct {
	BucketName string
	Prefix     string
	ProcessDir string
	TempDir   string
}


