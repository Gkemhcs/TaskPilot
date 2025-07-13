package utils

import "github.com/sirupsen/logrus"

// NewLogger creates and configures a new logrus.Logger instance.
// The logger outputs logs in JSON format with ISO 8601 timestamps.
func NewLogger() *logrus.Logger {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00", // ISO 8601 format
		PrettyPrint:    true, // Enables pretty printing of logs
	})
	return logger

}
