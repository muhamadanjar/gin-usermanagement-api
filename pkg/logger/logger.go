package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// global logger instance
	log *zap.Logger
)

// Config holds logger configuration
type Config struct {
	Level string // debug, info, warn, error
	Mode  string // development, production
}

// Initialize initializes the global logger
func Initialize(cfg Config) error {
	var zapConfig zap.Config

	if cfg.Mode == "production" {
		// Production: JSON output
		zapConfig = zap.NewProductionConfig()
	} else {
		// Development: Console output with colors
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level
	switch cfg.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Build logger
	var err error
	log, err = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	return nil
}

// Sync syncs the logger (should be called before application exits)
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		// Fallback to a default logger if not initialized
		log, _ = zap.NewDevelopment()
	}
	return log
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// Sugar returns a sugared logger for easier use
func Sugar() *zap.SugaredLogger {
	return GetLogger().Sugar()
}
