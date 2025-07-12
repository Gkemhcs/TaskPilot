package storage

import (
	"context"
)

func StorageFactory(ctx context.Context,client string, config StorageConfig) (StorageClient, error) {

	switch client {
	case "gcp":
		gcpClient, err := NewGCPStorageClient(ctx, config.BucketName, config.Prefix, config.ProcessDir)
		if err != nil {
			return nil, err
		}
		return gcpClient, nil

	default:
		localClient := NewLocalStorageClient(config.TempDir, config.ProcessDir)
		return localClient, nil
	}
}
