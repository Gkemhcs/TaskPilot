package storage

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"cloud.google.com/go/storage"
)

func NewGCPStorageClient(ctx context.Context,bucketName, prefix, processDir string) (*GCPStorageClient, error) {

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCPStorageClient{
		bucketName: bucketName,
		Prefix: prefix,
		LocalDir: processDir,
		client:    client,
	}, nil

}

type GCPStorageClient struct {
    bucketName string
    Prefix     string // e.g., "imports"
    LocalDir   string // where to store downloaded file locally
	client     *storage.Client
}

func (g *GCPStorageClient) Upload(file multipart.File, filename string) error {
    ctx := context.Background()
    objectName := filepath.Join(g.Prefix, filename)

    wc := g.client.Bucket(g.bucketName).Object(objectName).NewWriter(ctx)
    _, err := io.Copy(wc, file)
    if err != nil {
        return  err
    }

    err = wc.Close()
    if err != nil {
        return err 
    }

    return nil
}

func (g *GCPStorageClient) Download(filename string) (string, error) {
    ctx := context.Background()
    objectName := filepath.Join(g.Prefix, filename)
    localPath := filepath.Join(g.LocalDir, filename)

    rc, err := g.client.Bucket(g.bucketName).Object(objectName).NewReader(ctx)
    if err != nil {
        return "", err
    }
    defer rc.Close()

    out, err := os.Create(localPath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    _, err = io.Copy(out, rc)
    return localPath, err
}

func (g *GCPStorageClient) Delete(filename string) error {
    ctx := context.Background()
    objectName := filepath.Join(g.Prefix, filename)

    // Delete blob from GCS
    err := g.client.Bucket(g.bucketName).Object(objectName).Delete(ctx)
    if err != nil {
        return err
    }

    // Delete from local
    return os.Remove(filepath.Join(g.LocalDir, filename))
}
