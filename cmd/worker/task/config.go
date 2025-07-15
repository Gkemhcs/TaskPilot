package main

import (
	"log"

	"github.com/Gkemhcs/taskpilot/internal/storage"
	"github.com/spf13/viper"
)

// WorkerConfig holds configuration values for the task worker process.
// Includes database, RabbitMQ, storage, and queue settings.
type WorkerConfig struct {
	DBHost          string                // Database host
	DBPort          string                // Database port
	DBUser          string                // Database user
	DBPassword      string                // Database password
	DBName          string                // Database name
	RabbitMQURL     string                // RabbitMQ connection URL
	TaskQueue       string                // Import queue name
	TaskExportQueue string                // Export queue name
	StorageType     string                // Storage type (local/gcp)
	StorageConfig   storage.StorageConfig // Storage configuration
}

// LoadWorkerConfig loads configuration for the task worker from environment variables and .env file.
// Sets sensible defaults and returns a populated WorkerConfig struct.
func LoadWorkerConfig() *WorkerConfig {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Attempt to read from .env file, fallback to environment variables
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️ No .env file found, using env vars: %v", err)
	}

	// Set default values for config keys
	viper.SetDefault("TASK_IMPORT_QUEUE", "task_import_queue")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("STORAGE_TYPE", "local")
	viper.SetDefault("TEMP_DIR", "/tmp")
	viper.SetDefault("PROCESS_DIR", "/tmp/processed")
	viper.SetDefault("GCP_BUCKET", "")
	viper.SetDefault("GCP_PREFIX", "")
	viper.SetDefault("TASK_EXPORT_QUEUE", "task_export_queue")

	// Build storage config
	storageCfg := storage.StorageConfig{
		BucketName: viper.GetString("GCP_BUCKET"),
		Prefix:     viper.GetString("GCP_PREFIX"),
		TempDir:    viper.GetString("TEMP_DIR"),
		ProcessDir: viper.GetString("PROCESS_DIR"),
	}

	// Build and return WorkerConfig
	return &WorkerConfig{
		DBHost:          viper.GetString("DB_HOST"),
		DBPort:          viper.GetString("DB_PORT"),
		DBUser:          viper.GetString("DB_USER"),
		DBPassword:      viper.GetString("DB_PASSWORD"),
		DBName:          viper.GetString("DB_NAME"),
		RabbitMQURL:     viper.GetString("RABBITMQ_URL"),
		TaskQueue:       viper.GetString("TASK_IMPORT_QUEUE"),
		StorageType:     viper.GetString("STORAGE_TYPE"),
		StorageConfig:   storageCfg,
		TaskExportQueue: viper.GetString("TASK_EXPORT_QUEUE"),
	}
}
