package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// NewLocalStorageClient initializes a LocalStorageClient for local file operations.
// Ensures both temp and process directories exist.
func NewLocalStorageClient(tempDir, processDir string) (*LocalStorageClient, error) {
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Ensure ProcessDir exists
	if err := os.MkdirAll(processDir, os.ModePerm); err != nil {
		return nil, err
	}

	return &LocalStorageClient{
		TempDir:    tempDir,
		ProcessDir: processDir,
	}, nil
}

// LocalStorageClient implements StorageClient for local disk operations.
// Handles file upload, download, deletion, and signed URL generation locally.
type LocalStorageClient struct {
	TempDir    string // Directory for uploaded files
	ProcessDir string // Directory for processed/downloaded files
}

// Upload saves the uploaded file to the TempDir with the given filename.
func (c *LocalStorageClient) Upload(file multipart.File, filename string) error {
	fullPath := filepath.Join(c.TempDir, filename)

	out, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

// Download copies a file from TempDir to ProcessDir and returns the destination path.
func (c *LocalStorageClient) Download(filename string) (string, error) {
	sourcePath := filepath.Join(c.TempDir, filename)
	destPath := filepath.Join(c.ProcessDir, filename)

	in, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return destPath, err
}

// GenerateSignedURL returns the local file path as a "signed URL" for local storage.
// No actual signing is performed; this is for interface compatibility.
func (c *LocalStorageClient) GenerateSignedURL(filename string, inExpires time.Duration) (string, error) {
	destPath := filepath.Join(c.ProcessDir, filename)
	return destPath, nil
}

// Delete removes the file from both TempDir and ProcessDir.
func (c *LocalStorageClient) Delete(filename string) error {
	_ = os.Remove(filepath.Join(c.TempDir, filename))
	return os.Remove(filepath.Join(c.ProcessDir, filename))
}
