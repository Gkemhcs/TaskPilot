package storage

import (
	"mime/multipart"
	"time"
)

// StorageClient defines the interface for storage operations used in the application.
// Implementations include GCP and local storage clients, allowing interchangeable usage.
type StorageClient interface {
	// Upload stores a file with the given filename.
	Upload(file multipart.File, filename string) error
	// Download retrieves a file by filename and returns the local path.
	Download(filename string) (string, error)
	// Delete removes the file from storage and local disk.
	Delete(filename string) error
	// GenerateSignedURL returns a signed URL for accessing the file, valid for the given duration.
	GenerateSignedURL(filename string, inExpires time.Duration) (string, error)
}

// StorageConfig holds configuration for initializing storage clients.
// Fields are used for both GCP and local storage implementations.
type StorageConfig struct {
	BucketName string // GCP bucket name
	Prefix     string // Prefix for object names (e.g., "imports")
	ProcessDir string // Directory for processing/downloaded files
	TempDir    string // Directory for temporary/uploaded files
}
