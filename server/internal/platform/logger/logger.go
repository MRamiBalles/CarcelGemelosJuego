// Package logger provides structured logging for the game server.
// All actions by "The Twins" (server) should be traceable through this.
package logger

import (
	"log"
	"os"
)

// Logger provides structured logging with context.
type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

// NewLogger creates a new logger instance.
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "[TWINS-INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "[TWINS-WARN] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[TWINS-ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs informational messages.
func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

// Warn logs warning messages.
func (l *Logger) Warn(msg string) {
	l.warnLogger.Println(msg)
}

// Error logs error messages.
func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}

// Event logs a specific game event for "The Twins" oversight.
func (l *Logger) Event(eventType string, actorID string, details string) {
	l.infoLogger.Printf("[EVENT:%s] Actor:%s | %s", eventType, actorID, details)
}
