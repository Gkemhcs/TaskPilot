package storage

import (
	"context"
)

// StorageFactory returns a StorageClient implementation based on the client type.
// Supports "gcp" for Google Cloud Storage and defaults to local storage otherwise.
// Uses StorageConfig for initialization parameters.
func StorageFactory(ctx context.Context, client string, config StorageConfig) (StorageClient, error) {
	switch client {
	case "gcp":
		gcpClient, err := NewGCPStorageClient(ctx, config.BucketName, config.Prefix, config.ProcessDir)
		if err != nil {
			return nil, err
		}
		return gcpClient, nil
	default:
		localClient, err := NewLocalStorageClient(config.TempDir, config.ProcessDir)
		if err != nil {
			return nil, err
		}
		return localClient, nil
	}
}
