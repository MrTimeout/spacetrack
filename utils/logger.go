package utils

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

// Configure configures the zap logger from a Config structure
func Configure(cfg Logger) {
	logger = cfg.Tee()
}

// Debug will log a zap.Logger debug message
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info will log a zap.Logger info message
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn will log a zap.Logger warn message
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error will log a zap.Logger error message
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// DPanic will log a zap.Logger dpanic message
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

// Panic will log a zap.Logger panic message
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal will log a zap.Logger fatal message
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
