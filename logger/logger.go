package logger

import (
	"context"
	"fmt"
	"log"
)

// Logger defines the interface for structured logging
type Logger interface {
	Info(msg string)
	Error(msg string)
	V(level int) Logger
	Extra(key, value string) Logger
}

// basicLogger provides a simple logger implementation
type basicLogger struct {
	extras map[string]string
}

// NewOCMLogger creates a new logger from context
func NewOCMLogger(ctx context.Context) Logger {
	return &basicLogger{
		extras: make(map[string]string),
	}
}

// Info logs an info message
func (l *basicLogger) Info(msg string) {
	log.Printf("[INFO] %s %v", msg, l.extras)
}

// Error logs an error message
func (l *basicLogger) Error(msg string) {
	log.Printf("[ERROR] %s %v", msg, l.extras)
}

// V returns a logger with verbosity level
func (l *basicLogger) V(level int) Logger {
	return l // For simplicity, ignore verbosity in this minimal implementation
}

// Extra adds structured data to the logger
func (l *basicLogger) Extra(key, value string) Logger {
	newLogger := &basicLogger{
		extras: make(map[string]string),
	}
	for k, v := range l.extras {
		newLogger.extras[k] = v
	}
	newLogger.extras[key] = value
	return newLogger
}

// UpdateAdvisoryLockCountMetric is a placeholder for metrics
func UpdateAdvisoryLockCountMetric(lockType interface{}, status string) {
	// Placeholder - in a real implementation, this would update metrics
	fmt.Printf("Advisory lock metric: %v %s\n", lockType, status)
}

// UpdateAdvisoryLockDurationMetric is a placeholder for metrics
func UpdateAdvisoryLockDurationMetric(lockType interface{}, status string, startTime interface{}) {
	// Placeholder - in a real implementation, this would update metrics
	fmt.Printf("Advisory lock duration metric: %v %s\n", lockType, status)
}