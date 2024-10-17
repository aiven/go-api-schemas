// Package util is the package that contains the utility functions.
package util

import (
	"log"
	"os"
)

const (
	// logSeparator is the separator for the log.
	logSeparator = " "

	// logFlag is the flag for the log.
	logFlag = log.LstdFlags | log.LUTC | log.Lmsgprefix
)

const (
	// ServiceTypesFilename is the name of the service types file.
	ServiceTypesFilename = "service_types.yml"

	// IntegrationTypesFilename is the name of the integration types file.
	IntegrationTypesFilename = "integration_types.yml"

	// IntegrationEndpointTypesFilename is the name of the integration endpoint types file.
	IntegrationEndpointTypesFilename = "integration_endpoint_types.yml"
)

// Logger is a struct that holds the loggers for the application.
type Logger struct {
	// Info is the logger for info messages.
	Info *log.Logger

	// Error is the logger for error messages.
	Error *log.Logger
}

// EnvMap is a type for a map of environment variables.
type EnvMap map[string]string

// SetupLogger sets up the logger.
func SetupLogger(logger *Logger) {
	logger.Info = log.New(os.Stdout, "[INFO]"+logSeparator, logFlag)
	logger.Error = log.New(os.Stderr, "[ERROR]"+logSeparator, logFlag)
}

// Ref returns the reference (pointer) of the provided value.
func Ref[T any](v T) *T {
	return &v
}
