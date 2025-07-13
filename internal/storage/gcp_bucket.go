package storage

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
)

func NewGCPStorageClient(ctx context.Context, bucketName, prefix, processDir string) (*GCPStorageClient, error) {
	// Ensure ProcessDir exists

	if err := os.MkdirAll(processDir, os.ModePerm); err != nil {
		return nil, err
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	// Load credentials from GOOGLE_APPLICATION_CREDENTIALS env var
	credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsPath == "" {
		return nil, customErrors.ErrGoogleApplicationCredentialsNotSet
	}

	saEmail, saKey, err := loadServiceAccount(credsPath)
	if err != nil {
		return nil, customErrors.ErrGoogleApplicationCredentialsNotSet
	}

	return &GCPStorageClient{
		bucketName:     bucketName,
		Prefix:         prefix,
		LocalDir:       processDir,
		client:         client,
		PrivateKey:     saKey,
		GoogleAccessID: saEmail,
	}, nil

}

func loadServiceAccount(path string) (string, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, customErrors.ErrLoadingServiceAccountFile
	}
	var sa serviceAccount
	if err := json.Unmarshal(data, &sa); err != nil {
		return "", nil, customErrors.ErrInvalidServiceAccountFile
	}
	return sa.ClientEmail, []byte(sa.PrivateKey), nil
}

type serviceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

type GCPStorageClient struct {
	bucketName     string
	Prefix         string // e.g., "imports"
	LocalDir       string // where to store downloaded file locally
	client         *storage.Client
	GoogleAccessID string
	PrivateKey     []byte
}

func (g *GCPStorageClient) Upload(file multipart.File, filename string) error {
	ctx := context.Background()
	objectName := filepath.Join(g.Prefix, filename)

	wc := g.client.Bucket(g.bucketName).Object(objectName).NewWriter(ctx)
	wc.ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	_, err := io.Copy(wc, file)
	if err != nil {
		return err
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

func (g *GCPStorageClient) GenerateSignedURL(filename string, expiresIn time.Duration) (string, error) {
	objectName := filepath.Join(g.Prefix, filename)

	url, err := storage.SignedURL(g.bucketName, objectName, &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(expiresIn),
		GoogleAccessID: g.GoogleAccessID,
		PrivateKey:     g.PrivateKey,
	})
	if err != nil {
		return "", customErrors.ErrGeneratingSignedURL
	}
	return url, nil
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
