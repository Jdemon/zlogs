package zlogs

import (
	"context"

	"github.com/rs/zerolog"
)

// std is the default instance of Logger configured with default settings.
var std = newStandardLogger()

// NewLogger initializes the global logger instance with the provided configuration.
func NewLogger(config *Config) {
	std = newLogger(config)
}

// GetLogger returns the standard logger instance used for logging in the application.
func GetLogger() *Logger {
	return std
}

// AddHook adds a hook to the default logger instance.
func AddHook(hook zerolog.Hook) {
	std.Hook(hook)
}

// AddCallerSkip returns a new context.Context that carries the specified caller skip value.
func AddCallerSkip(ctx context.Context, skip int) context.Context {
	return context.WithValue(ctx, CallerSkip, skip)
}

// Debug logs a message at the debug level and returns an Event object to further customize the log entry.
func Debug() *Event {
	return &Event{std.Logger.Debug()}
}

// Info creates a new event at the info log level.
func Info() *Event {
	return &Event{std.Logger.Info()}
}

// Warn creates a log event at the warn level using the standard logger.
func Warn() *Event {
	return &Event{std.Logger.Warn()}
}

// Error returns a new Event with the logging level set to 'error'.
func Error() *Event {
	return &Event{std.Logger.Error()}
}

// Fatal creates an Event with fatal log level
func Fatal() *Event {
	return &Event{std.Logger.Fatal()}
}

// Panic creates a new logging event at the panic level using the standard logger and returns an Event pointer.
func Panic() *Event {
	return &Event{std.Logger.Panic()}
}
