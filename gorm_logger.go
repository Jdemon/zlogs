package zlogs

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type GORMLogger struct {
	*zerolog.Logger
}

func NewGORMLogger(config *Config) *GORMLogger {
	loggerGORM := zerolog.New(os.Stdout).Hook(&InitHook{
		appName:       config.AppName,
		disableCaller: true,
	}).With().Timestamp().Logger()

	if lvl, err := zerolog.ParseLevel(config.Level); err != nil {
		loggerGORM.Level(zerolog.DebugLevel)
	} else {
		loggerGORM.Level(lvl)
	}
	return &GORMLogger{
		&loggerGORM,
	}
}

func (l *GORMLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *GORMLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	l.Logger.Error().Ctx(ctx).Msgf(msg, opts...)
}

func (l *GORMLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	l.Logger.Warn().Ctx(ctx).Msgf(msg, opts...)
}

func (l *GORMLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	l.Logger.Info().Ctx(ctx).Msgf(msg, opts...)
}

func (l *GORMLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	var event *zerolog.Event
	if err != nil {
		event = l.Logger.Debug().Ctx(ctx)
	} else {
		event = l.Logger.Trace().Ctx(ctx)
	}

	event.Dur(l.getDurationFieldKey(), time.Since(begin))
	sql, rows := f()
	if sql != "" {
		event.Str("sql", sql)
	}
	if rows > -1 {
		event.Int64("rows", rows)
	}
	event.Send()
}

func (l *GORMLogger) getDurationFieldKey() string {
	switch zerolog.DurationFieldUnit {
	case time.Nanosecond:
		return "elapsed_ns"
	case time.Microsecond:
		return "elapsed_us"
	case time.Millisecond:
		return "elapsed_ms"
	case time.Second:
		return "elapsed_s"
	case time.Minute:
		return "elapsed_min"
	case time.Hour:
		return "elapsed_hr"
	default:
		l.Logger.Warn().Msg("Unexpected DurationFieldUnit")
		return "elapsed"
	}
}
