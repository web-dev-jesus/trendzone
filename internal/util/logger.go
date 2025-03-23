package util

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger is a simple logging utility
type Logger struct {
	level LogLevel
}

// NewLogger creates a new logger with the specified level
func NewLogger(levelStr string) *Logger {
	var level LogLevel
	switch strings.ToLower(levelStr) {
	case "debug":
		level = DEBUG
	case "info":
		level = INFO
	case "warn":
		level = WARN
	case "error":
		level = ERROR
	default:
		level = INFO
	}

	return &Logger{level: level}
}

// Debug logs debug messages
func (l *Logger) Debug(message string) {
	if l.level <= DEBUG {
		l.log("DEBUG", message)
	}
}

// Info logs informational messages
func (l *Logger) Info(message string) {
	if l.level <= INFO {
		l.log("INFO", message)
	}
}

// Warn logs warning messages
func (l *Logger) Warn(message string) {
	if l.level <= WARN {
		l.log("WARN", message)
	}
}

// Error logs error messages
func (l *Logger) Error(message string) {
	if l.level <= ERROR {
		l.log("ERROR", message)
	}
}

// log formats and writes a log message
func (l *Logger) log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
	fmt.Fprint(os.Stdout, logMessage)
}
