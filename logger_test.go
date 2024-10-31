package zlogs

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type testLogger struct {
	Password string `json:"password"`
}

func BenchmarkLoggerMasking(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.WithValue(context.Background(), TraceID, "trace-id-value")
		Debug().WithField("data", testLogger{
			Password: "P@ssw0rd",
		}).Ctx(ctx).Msg("benchmark test")
	}
}

func BenchmarkOriginalZeroLog(b *testing.B) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05-0700"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "severity"
	zerolog.MessageFieldName = "message"
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	zlog.Logger = logger
	for i := 0; i < b.N; i++ {
		ctx := context.WithValue(context.Background(), TraceID, "trace-id-value")
		logger.Debug().Ctx(ctx).Interface("data", testLogger{
			Password: "P@ssw0rd",
		}).Msg("benchmark test")
	}
}
