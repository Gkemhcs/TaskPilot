package config

import "time"

type StorageConfig struct {
	BucketName string
	Prefix     string
	ProcessDir string
	TempDir   string

}

type RabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Queue    string
}

type Config struct {
	Port                string        // Port for the HTTP server
	DBHost              string        // Database host
	DBPort              string        // Database port
	DBUser              string        // Database user
	DBPassword          string        // Database password
	DBName              string        // Database name
	HOST                string        // Host for the HTTP server
	JWTAccessTokenSecret           string        // Secret key for JWT Access Key signing
	JWTRefreshTokenSecret 	string  // Secret key for JWT Refresh Key signing
	AccessTokenDuration time.Duration // Duration for JWT access tokens
	RefreshTokenDuration time.Duration // Duration for JWT refresh tokens
	RedisHost           string        // Redis host for caching
	RedisPort           string        // Redis port for caching
	StorageConfig       StorageConfig // Configuration for storage client
}

