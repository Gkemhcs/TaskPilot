package config

import (
	"time"

	"github.com/Gkemhcs/taskpilot/internal/storage"
)

type RabbitMQPublisherConfig struct {
	QueueName  string
	Exchange   string
	RoutingKey string
}

type Config struct {
	Port                 string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	HOST                 string
	JWTAccessTokenSecret string
	JWTRefreshTokenSecret string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	RedisHost            string
	RedisPort            string
	StorageType         string // e.g., "local", "gcp"
	// Importer-related configs
	StorageConfig     storage.StorageConfig
	RabbitMQURL       string
	ProjectPublisher  RabbitMQPublisherConfig
	TaskPublisher     RabbitMQPublisherConfig
	ProjectExportPublisher RabbitMQPublisherConfig
	TaskExportPublisher RabbitMQPublisherConfig

}