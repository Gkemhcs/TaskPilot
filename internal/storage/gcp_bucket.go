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

// NewGCPStorageClient initializes a GCPStorageClient for interacting with Google Cloud Storage.
// Ensures the local processing directory exists, loads credentials, and sets up the client.
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

// loadServiceAccount loads service account credentials from a JSON file.
// Returns the client email and private key for signing URLs.
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

// serviceAccount represents the relevant fields from a GCP service account JSON file.
type serviceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

// GCPStorageClient implements StorageClient for Google Cloud Storage.
// Handles file upload, download, deletion, and signed URL generation.
type GCPStorageClient struct {
	bucketName     string          // GCP bucket name
	Prefix         string          // Prefix for object names (e.g., "imports")
	LocalDir       string          // Local directory for downloaded files
	client         *storage.Client // GCP storage client
	GoogleAccessID string          // Service account email for signed URLs
	PrivateKey     []byte          // Private key for signed URLs
}

// Upload uploads a file to the GCP bucket under the specified prefix and filename.
// Sets the content type for Excel files.
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

// Download retrieves a file from the GCP bucket and saves it locally.
// Returns the local file path.
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

// GenerateSignedURL creates a signed URL for accessing the file in GCP bucket.
// The URL is valid for the specified duration.
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

// Delete removes the file from both GCP bucket and local disk.
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
