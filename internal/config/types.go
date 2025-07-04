package config

import "time"

// Config holds all configuration values for the application.
type Config struct {
	Port                string        // Port for the HTTP server
	DBHost              string        // Database host
	DBPort              string        // Database port
	DBUser              string        // Database user
	DBPassword          string        // Database password
	DBName              string        // Database name
	HOST                string        // Host for the HTTP server
	JWTSecret           string        // Secret key for JWT signing
	AccessTokenDuration time.Duration // Duration for JWT access tokens
}
