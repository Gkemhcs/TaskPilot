package workers


package config

import (
    "os"
)

type Config struct {
    DB struct {
        DSN string
    }
    RabbitMQ struct {
        URL string
    }
    Storage struct {
        Type       string // "gcp" or "local"
        BucketName string
        Prefix     string
        TempDir    string
        ProcessDir string
    }
}






package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig loads application configuration from .env file and environment variables.
// It sets default values and parses durations as needed.
func LoadConfig() *Config {
	viper.SetConfigFile(".env") // support for .env file
	viper.AutomaticEnv()        // also read from real env vars

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found (ok in prod): %v", err)
	}

	// Set default values for configuration
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("ACCESS_TOKEN_DURATION", "12h")
	viper.SetDefault("REFRESH_TOKEN_DURATION", "24h")
	viper.SetDefault("CONTEXT_TIMEOUT", "10s")

	// Parse token duration from config
	accessTokenDuration, err := time.ParseDuration(viper.GetString("ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Fatalf("invalid duration: %v", err)
	}
	refreshTokenDuration, err := time.ParseDuration(viper.GetString("REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Fatalf("invalid duration: %v", err)
	}

	// Return a Config struct populated with values from config/env
	return &Config{
		Port:                viper.GetString("PORT"),
		DBHost:              viper.GetString("DB_HOST"),
		DBPort:              viper.GetString("DB_PORT"),
		DBUser:              viper.GetString("DB_USER"),
		DBPassword:          viper.GetString("DB_PASSWORD"),
		DBName:              viper.GetString("DB_NAME"),
		JWTAccessTokenSecret: viper.GetString("JWT_ACCESS_TOKEN_SECRET"),
		JWTRefreshTokenSecret:viper.GetString("JWT_REFRESH_TOKEN_SECRET"),
		HOST:                viper.GetString("HOST"),
		AccessTokenDuration: accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
		RedisHost: 		 viper.GetString("REDIS_HOST"),
		RedisPort: 		 viper.GetString("REDIS_PORT"),
	}}


func Load() *Config {
    c := &Config{}

    c.DB.DSN = os.Getenv("DB_DSN")
    c.RabbitMQ.URL = os.Getenv("RABBITMQ_URL")
    c.Storage.Type = os.Getenv("STORAGE_TYPE")
    c.Storage.BucketName = os.Getenv("GCP_BUCKET_NAME")
    c.Storage.Prefix = os.Getenv("GCP_PREFIX")
    c.Storage.TempDir = os.Getenv("LOCAL_TEMP_DIR")
    c.Storage.ProcessDir = os.Getenv("LOCAL_PROCESS_DIR")

    return c
}
