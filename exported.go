package zlogs

import (
	"context"

	"github.com/rs/zerolog"
)

var std = newStandardLogger()

func NewLogger(config *Config) {
	std = newLogger(config)
}

func GetLogger() *Logger {
	return std
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook zerolog.Hook) {
	std.Hook(hook)
}

func AddCallerSkip(ctx context.Context, skip int) context.Context {
	return context.WithValue(ctx, CallerSkip, skip)
}

func Debug() *Event {
	return &Event{std.Logger.Debug()}
}

func Info() *Event {
	return &Event{std.Logger.Info()}
}

func Warn() *Event {
	return &Event{std.Logger.Warn()}
}

func Error() *Event {
	return &Event{std.Logger.Error()}
}

func Fatal() *Event {
	return &Event{std.Logger.Fatal()}
}

func Panic() *Event {
	return &Event{std.Logger.Panic()}
}
