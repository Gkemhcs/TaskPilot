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
	TaskQueue     string
	TaskExportQueue string 
	StorageType   string
	StorageConfig storage.StorageConfig
}

func LoadWorkerConfig() *WorkerConfig {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️ No .env file found, using env vars: %v", err)
	}

	// ---------- Set Defaults ----------
	viper.SetDefault("TASK_IMPORT_QUEUE", "task_import_queue")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("STORAGE_TYPE", "local")
	viper.SetDefault("TEMP_DIR", "/tmp")
	viper.SetDefault("PROCESS_DIR", "/tmp/processed")
	viper.SetDefault("GCP_BUCKET", "")
	viper.SetDefault("GCP_PREFIX", "")
	viper.SetDefault("TASK_EXPORT_QUEUE","task_export_queue")

	storageCfg := storage.StorageConfig{
		BucketName: viper.GetString("GCP_BUCKET"),
		Prefix:     viper.GetString("GCP_PREFIX"),
		TempDir:    viper.GetString("TEMP_DIR"),
		ProcessDir: viper.GetString("PROCESS_DIR"),
	}

	return &WorkerConfig{
		DBHost:        viper.GetString("DB_HOST"),
		DBPort:        viper.GetString("DB_PORT"),
		DBUser:        viper.GetString("DB_USER"),
		DBPassword:    viper.GetString("DB_PASSWORD"),
		DBName:        viper.GetString("DB_NAME"),
		RabbitMQURL:   viper.GetString("RABBITMQ_URL"),
		TaskQueue:     viper.GetString("TASK_IMPORT_QUEUE"),
		StorageType:   viper.GetString("STORAGE_TYPE"),
		StorageConfig: storageCfg,
		TaskExportQueue: viper.GetString("TASK_EXPORT_QUEUE"),
	}
}
