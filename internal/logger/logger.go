package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type contextKey string

const (
	requestIDKey  contextKey = "requestID"
	requestURIKey contextKey = "requestURI"
	clientIPKey   contextKey = "clientIP"
)

// Setup configures the global logger with appropriate settings
func Setup(level logrus.Level) {
	// Configure logrus
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return f.Function, filepath.Base(f.File) + ":" + fmt.Sprintf("%d", f.Line)
		},
	})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
}

// WithRequestContext adds request-specific fields to the logger
func WithRequestContext(ctx context.Context) *logrus.Entry {
	logger := logrus.WithContext(ctx)

	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		logger = logger.WithField("request_id", reqID)
	}

	if reqURI, ok := ctx.Value(requestURIKey).(string); ok {
		logger = logger.WithField("uri", reqURI)
	}

	if clientIP, ok := ctx.Value(clientIPKey).(string); ok {
		logger = logger.WithField("client_ip", clientIP)
	}

	return logger
}

// GetRequestID gets the request ID from the context or generates a new one
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return uuid.New().String()
}

// NewRequestContext creates a new context with request information
func NewRequestContext(ctx context.Context, requestURI, clientIP string) context.Context {
	reqID := GetRequestID(ctx)
	ctx = context.WithValue(ctx, requestIDKey, reqID)
	ctx = context.WithValue(ctx, requestURIKey, requestURI)
	ctx = context.WithValue(ctx, clientIPKey, clientIP)
	return ctx
}

// WithFields creates a logger entry with additional fields automatically including caller information
func WithFields(fields logrus.Fields) *logrus.Entry {
	// Add caller information
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fields["caller"] = filepath.Base(file) + ":" + string(rune(line))
	}
	return logrus.WithFields(fields)
}
