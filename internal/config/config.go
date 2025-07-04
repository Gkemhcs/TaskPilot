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
	viper.SetDefault("ACCESS_TOKEN_DURATION", "15m")
	viper.SetDefault("CONTEXT_TIMEOUT", "10s")

	// Parse token duration from config
	duration, err := time.ParseDuration(viper.GetString("ACCESS_TOKEN_DURATION"))
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
		JWTSecret:           viper.GetString("JWT_SECRET"),
		HOST:                viper.GetString("HOST"),
		AccessTokenDuration: duration,
	}
}
