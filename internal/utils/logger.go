package utils

import "github.com/sirupsen/logrus"



func NewLogger() *logrus.Logger{

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00", // ISO 8601 format
	})
	return logger

}
