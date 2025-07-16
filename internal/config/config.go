package config

import (
	"errors"
	"log"

	"time"

	"github.com/Gkemhcs/taskpilot/internal/storage"
	"github.com/spf13/viper"
)

// LoadConfig loads application configuration from .env file and environment variables.
// It sets default values and parses durations as needed.
func LoadConfig() (*Config,error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found (ok in prod): %v", err)
	}

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("ACCESS_TOKEN_DURATION", "12h")
	viper.SetDefault("REFRESH_TOKEN_DURATION", "24h")
	viper.SetDefault("CONTEXT_TIMEOUT", "10s")

	// Storage defaults
	viper.SetDefault("STORAGE_TYPE", "local")
	viper.SetDefault("TEMP_DIR", "/tmp")
	viper.SetDefault("PROCESS_DIR", "/tmp/processed")
	viper.SetDefault("GCP_BUCKET", "")
	viper.SetDefault("GCP_PREFIX", "")

	// RabbitMQ
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("PROJECT_IMPORT_QUEUE", "project_import_queue")
	viper.SetDefault("PROJECT_IMPORT_ROUTING_KEY", "project_import_queue")
	viper.SetDefault("PROJECT_IMPORT_EXCHANGE", "")
	viper.SetDefault("TASK_IMPORT_QUEUE", "task_import_queue")
	viper.SetDefault("TASK_IMPORT_ROUTING_KEY", "task_import_queue")
	viper.SetDefault("TASK_IMPORT_EXCHANGE", "")
	viper.SetDefault("PROJECT_EXPORT_QUEUE", "project_export_queue")
	viper.SetDefault("PROJECT_EXPORT_ROUTING_KEY", "project_export_queue")
	viper.SetDefault("PROJECT_EXPORT_EXCHANGE", "")
	viper.SetDefault("TASK_EXPORT_QUEUE", "task_export_queue")
	viper.SetDefault("TASK_EXPORT_ROUTING_KEY", "task_export_queue")
	viper.SetDefault("TASK_EXPORT_EXCHANGE", "")

	accessTokenDuration, err := time.ParseDuration(viper.GetString("ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Fatalf("invalid ACCESS_TOKEN_DURATION: %v", err)
	}
	refreshTokenDuration, err := time.ParseDuration(viper.GetString("REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Fatalf("invalid REFRESH_TOKEN_DURATION: %v", err)
	}
	if viper.GetString("STORAGE_TYPE")=="gcp"{
		if viper.GetString("GOOGLE_APPLICATION_CREDENTIALS")==""{
			return nil,errors.New("PLEASE SET GOOGLE_APPLICATION_CREDENTIALS env pointing to gcp iam service account file")
		}
	}
	return &Config{
		Port:                 viper.GetString("PORT"),
		DBHost:               viper.GetString("DB_HOST"),
		DBPort:               viper.GetString("DB_PORT"),
		DBUser:               viper.GetString("DB_USER"),
		DBPassword:           viper.GetString("DB_PASSWORD"),
		DBName:               viper.GetString("DB_NAME"),
		JWTAccessTokenSecret: viper.GetString("JWT_ACCESS_TOKEN_SECRET"),
		JWTRefreshTokenSecret: viper.GetString("JWT_REFRESH_TOKEN_SECRET"),
		HOST:                 viper.GetString("HOST"),
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
		RedisHost:            viper.GetString("REDIS_HOST"),
		RedisPort:            viper.GetString("REDIS_PORT"),
		StorageType:          viper.GetString("STORAGE_TYPE"),
		StorageConfig: storage.StorageConfig{
			BucketName: viper.GetString("GCP_BUCKET"),
			Prefix:     viper.GetString("GCP_PREFIX"),
			ProcessDir: viper.GetString("PROCESS_DIR"),
			TempDir:    viper.GetString("TEMP_DIR"),
		},

		RabbitMQURL: viper.GetString("RABBITMQ_URL"),

		ProjectPublisher: RabbitMQPublisherConfig{
			QueueName:  viper.GetString("PROJECT_IMPORT_QUEUE"),
			Exchange:   viper.GetString("PROJECT_IMPORT_EXCHANGE"),
			RoutingKey: viper.GetString("PROJECT_IMPORT_ROUTING_KEY"),
		},
		TaskPublisher: RabbitMQPublisherConfig{
			QueueName:  viper.GetString("TASK_IMPORT_QUEUE"),
			Exchange:   viper.GetString("TASK_IMPORT_EXCHANGE"),
			RoutingKey: viper.GetString("TASK_IMPORT_ROUTING_KEY"),
		},
		ProjectExportPublisher: RabbitMQPublisherConfig{
			QueueName:  viper.GetString("PROJECT_EXPORT_QUEUE"),
			Exchange:   viper.GetString("PROJECT_EXPORT_EXCHANGE"),
			RoutingKey: viper.GetString("PROJECT_EXPORT_ROUTING_KEY"),
		},
		TaskExportPublisher: RabbitMQPublisherConfig{
			QueueName:  viper.GetString("TASK_EXPORT_QUEUE"),
			Exchange:   viper.GetString("TASK_EXPORT_EXCHANGE"),
			RoutingKey: viper.GetString("TASK_EXPORT_ROUTING_KEY"),
		},

	},nil 
}