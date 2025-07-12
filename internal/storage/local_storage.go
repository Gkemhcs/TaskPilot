package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func NewLocalStorageClient(tempDir, processDir string) *LocalStorageClient {
	return &LocalStorageClient{
		TempDir:    tempDir,
		ProcessDir: processDir,
	}
}

type LocalStorageClient struct {
    TempDir    string // for Upload
    ProcessDir string // for Download
}

func (c *LocalStorageClient) Upload(file multipart.File, filename string)  error {
    fullPath := filepath.Join(c.TempDir, filename)

    out, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    
    defer out.Close()

    _, err = io.Copy(out, file)
    return err
}

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

func (c *LocalStorageClient) Delete(filename string) error {
    _ = os.Remove(filepath.Join(c.TempDir, filename))
    return os.Remove(filepath.Join(c.ProcessDir, filename))
}
