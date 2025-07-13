package main

import (
	"log"

	"github.com/Gkemhcs/taskpilot/internal/storage"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RabbitMQURL   string
	ProjectQueue  string
	ProjectExportQueue string 
	StorageType   string
	StorageConfig storage.StorageConfig
}

func LoadWorkerConfig() *WorkerConfig {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// ---------- Set Defaults ----------
	viper.SetDefault("TEMP_DIR", "/tmp")
	viper.SetDefault("PROCESS_DIR", "/tmp/processed")

	viper.SetDefault("STORAGE_TYPE", "local")
	viper.SetDefault("GCP_BUCKET", "")
	viper.SetDefault("GCP_PREFIX", "")

	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("PROJECT_IMPORT_QUEUE", "project_import_queue")
	viper.SetDefault("PROJECT_EXPORT_QUEUE", "project_export_queue")

	// ---------- Load from .env (optional) ----------
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, relying on environment variables: %v", err)
	}

	// ---------- Return WorkerConfig ----------
	return &WorkerConfig{
		DBHost:       viper.GetString("DB_HOST"),
		DBPort:       viper.GetString("DB_PORT"),
		DBUser:       viper.GetString("DB_USER"),
		DBPassword:   viper.GetString("DB_PASSWORD"),
		DBName:       viper.GetString("DB_NAME"),
		RabbitMQURL:  viper.GetString("RABBITMQ_URL"),
		ProjectQueue: viper.GetString("PROJECT_IMPORT_QUEUE"),
		ProjectExportQueue: viper.GetString("PROJECT_EXPORT_QUEUE"),
		StorageType:  viper.GetString("STORAGE_TYPE"),
		StorageConfig: storage.StorageConfig{
			BucketName: viper.GetString("GCP_BUCKET"),
			Prefix:     viper.GetString("GCP_PREFIX"),
			TempDir:    viper.GetString("TEMP_DIR"),
			ProcessDir: viper.GetString("PROCESS_DIR"),
		},
	}
}
