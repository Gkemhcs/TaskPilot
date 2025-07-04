package config

import "time"
type Config struct {
	
	Port                string 
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	HOST 				string 
	JWTSecret           string
	AccessTokenDuration time.Duration
}