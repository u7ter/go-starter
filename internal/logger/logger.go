package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const requestIDKey contextKey = "request_id"

var log *zap.Logger

// Init initializes the global logger
func Init(level string, isProduction bool) error {
	var config zap.Config

	if isProduction {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Parse log level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Output to stdout for container compatibility
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	log = logger
	return nil
}

// Get returns the global logger instance
func Get() *zap.Logger {
	if log == nil {
		// Fallback logger if Init wasn't called
		log, _ = zap.NewProduction()
	}
	return log
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// FromContext returns a logger with request ID from context if available
func FromContext(ctx context.Context) *zap.Logger {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return Get().With(zap.String("request_id", requestID))
	}
	return Get()
}

// Sync flushes any buffered log entries
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// SetOutput sets the output for testing purposes
func SetOutput(ws zapcore.WriteSyncer) {
	if log != nil {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			ws,
			log.Level(),
		)
		log = zap.New(core)
	}
}

// NewTestLogger creates a logger for testing
func NewTestLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	logger, _ := config.Build()
	return logger
}

// Close closes the logger gracefully
func Close() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
