package zlogs

import (
	"bytes"
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func createTestContext() context.Context {
	ctx := context.WithValue(context.Background(), TraceID, "12345")
	ctx = context.WithValue(ctx, RequestID, "req-67890")
	ctx = context.WithValue(ctx, CorrelationID, "corr-abcde")
	return ctx
}

func TestInitHook_Run(t *testing.T) {
	// Sequential test cases
	testCases := []struct {
		name           string
		disableCaller  bool
		expectedCaller bool
	}{
		{
			name:           "should set entry data and caller info correctly",
			disableCaller:  false,
			expectedCaller: true,
		},
		{
			name:           "should skip caller info when disableCaller is true",
			disableCaller:  true,
			expectedCaller: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			zlogger := zerolog.New(&buf).Hook(&InitHook{
				AppName:       "MyApp",
				DisableCaller: tc.disableCaller,
			})
			zlogger.Level(zerolog.DebugLevel)
			ctx := createTestContext()

			// Zlogger event
			e := zlogger.With().Logger()
			event := e.Error().Ctx(ctx)

			// Log with context
			event.Msg("test log")

			// Assert log output
			logOutput := buf.String()

			if tc.expectedCaller {
				assert.Contains(t, logOutput, `"file":"hook_test.go`) // adjust if needed
				assert.Contains(t, logOutput, `"func":"func1"`)       // adjust if needed
			} else {
				assert.NotContains(t, logOutput, `"file":`)
				assert.NotContains(t, logOutput, `"func":`)
			}

			assert.Contains(t, logOutput, `"trace_id":"12345"`)
			assert.Contains(t, logOutput, `"request_id":"req-67890"`)
			assert.Contains(t, logOutput, `"correlation_id":"corr-abcde"`)
			assert.Contains(t, logOutput, `"MyApp"`)
		})
	}
}
