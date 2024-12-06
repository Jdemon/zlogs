package zlogs_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"gitlab-dev.tripperpix.com/vb-mobile-backend/channel-api/batch/batch-card-apply-task/pkg/zlogs"
)

var data = map[string]any{
	"data": map[string]any{
		"password":      "P@ssw0rd",
		"mobile_number": "0909263742",
		"id":            "112132321312",
		"firstname":     "John",
		"lastName":      "Doe",
	},
	"credit_card": "4231234512341234",
}

// TestNewLogger verifies that a new logger instance is correctly created and retrievable. The function configures the logger
// with a debug level and masking enabled, ensuring the logger instance is not nil.
func TestNewLogger(t *testing.T) {
	config := &zlogs.Config{
		Level:   "debug",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	assert.NotNil(t, zlogs.GetLogger())
}

// TestGetLogger tests the GetLogger function to ensure a non-nil logger instance is returned after initialization with configuration.
func TestGetLogger(t *testing.T) {
	config := &zlogs.Config{
		Level:   "debug",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	logger := zlogs.GetLogger()
	assert.NotNil(t, logger)
}

// TestAddHook ensures that a hook can be added to the zlogs logger without causing a panic.
func TestAddHook(t *testing.T) {
	config := &zlogs.Config{
		Level:   "debug",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	hook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {})

	// Note: zlogs.Logger does not expose hooks directly; this test verifies no panic occurs
	assert.NotPanics(t, func() {
		zlogs.AddHook(hook)
	})
}

// TestAddCallerSkip checks if the AddCallerSkip function correctly sets the caller skip value in the context.
func TestAddCallerSkip(t *testing.T) {
	testCases := []struct {
		skip         int
		expectedSkip int
	}{
		{0, 0},
		{1, 1},
		{5, 5},
	}

	for _, tc := range testCases {
		ctx := context.Background()
		newCtx := zlogs.AddCallerSkip(ctx, tc.skip)
		val := newCtx.Value(zlogs.CallerSkip).(int)
		assert.Equal(t, tc.expectedSkip, val)
	}
}

// TestDebug verifies that the debug logging level is initialized correctly and an event is created.
func TestDebug(t *testing.T) {
	config := &zlogs.Config{
		Level:   "debug",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	event := zlogs.Debug()
	assert.NotNil(t, event)
	event.WithField("data", data).Msg("test message")
}

// TestInfo verifies that the Info log level creates a non-nil log event and initializes the logger with proper configuration.
func TestInfo(t *testing.T) {
	config := &zlogs.Config{
		Level:   "info",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	event := zlogs.Info()
	assert.NotNil(t, event)
	event.WithField("data", data).Msg("test message")
}

// TestWarn verifies the functionality of logging an event at the warn level with masking enabled.
func TestWarn(t *testing.T) {
	config := &zlogs.Config{
		Level:   "warn",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	event := zlogs.Warn()
	assert.NotNil(t, event)
	event.WithField("data", data).Msg("test message")
}

// TestError tests the Error function of the zlogs package to ensure it returns a non-nil event.
func TestError(t *testing.T) {
	config := &zlogs.Config{
		Level:   "error",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	event := zlogs.Error()
	assert.NotNil(t, event)
	event.WithField("data", data).Msg("test message")
}

// TestFatal verifies the creation of a fatal log event using the zlogs package.
func TestFatal(t *testing.T) {
	config := &zlogs.Config{
		Level:   "fatal",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	event := zlogs.Fatal()
	assert.NotNil(t, event)
}

// TestPanic configures a logger with panic level and verifies that a new panic event is created successfully.
func TestPanic(t *testing.T) {
	config := &zlogs.Config{
		Level:   "panic",
		Masking: zlogs.MaskingConfig{Enabled: true},
	}
	zlogs.NewLogger(config)
	event := zlogs.Panic()
	assert.NotNil(t, event)
}
