package zlogs

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/net/context"
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

// testLogger is a struct used for logging test information with a Password field that may contain sensitive data.
type testLogger struct {
	Password string `json:"password"`
}

// BenchmarkLoggerMasking benchmarks the logging process with context and field masking. It measures the performance of logging with sensitive data masking.
func BenchmarkLoggerMasking(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.WithValue(context.Background(), TraceID, "trace-id-value")
		Debug().WithField("data", data).Ctx(ctx).Msg("benchmark test")
	}
}

// BenchmarkOriginalZeroLog benchmarks the original zerolog setup by logging a debug-level message with context.
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
			Password: "password",
		}).Msg("benchmark test")
	}
}

// TestMaskingFunctions tests the masking functionality of the logger.
func TestMaskingFunctions(t *testing.T) {
	logger := newLogger(&Config{
		Level: "debug",
		Masking: MaskingConfig{
			Enabled:         false,
			SensitiveFields: []string{"cc_number"},
		},
	})

	t.Run("testMaskFields", func(t *testing.T) {
		cases := []struct {
			name  string
			input map[string]interface{}
			want  map[string]interface{}
		}{
			{"AllSensitive", map[string]interface{}{"password": "secret", "cc_number": "1234-5678-9876-5432"}, map[string]interface{}{"password": redactedValue, "cc_number": redactedValue}},
			{"MixedFields", map[string]interface{}{"password": "secret", "username": "john_doe"}, map[string]interface{}{"password": redactedValue, "username": "john_doe"}},
			{"NoSensitive", map[string]interface{}{"username": "john_doe"}, map[string]interface{}{"username": "john_doe"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := logger.maskFields(tc.input)
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})

	t.Run("testValueMasking", func(t *testing.T) {
		cases := []struct {
			name   string
			result map[string]interface{}
			key    string
			value  interface{}
			want   map[string]interface{}
		}{
			{"SensitiveField", map[string]interface{}{}, "password", "secret", map[string]interface{}{"password": redactedValue}},
			{"NestedSensitiveField", map[string]interface{}{}, "details", map[string]interface{}{"password": "secret"}, map[string]interface{}{"details": map[string]interface{}{"password": redactedValue}}},
			{"NonSensitiveField", map[string]interface{}{}, "username", "john_doe", map[string]interface{}{"username": "john_doe"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := logger.valueMasking(tc.result, tc.key, tc.value)
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})

	t.Run("testMaskArrayFields", func(t *testing.T) {
		cases := []struct {
			name  string
			input []interface{}
			want  []interface{}
		}{
			{"ArrayWithSensitive", []interface{}{map[string]interface{}{"password": "secret"}}, []interface{}{map[string]interface{}{"password": redactedValue}}},
			{"ArrayWithMixed", []interface{}{map[string]interface{}{"password": "secret"}, map[string]interface{}{"username": "john_doe"}}, []interface{}{map[string]interface{}{"password": redactedValue}, map[string]interface{}{"username": "john_doe"}}},
			{"ArrayWithNonSensitive", []interface{}{map[string]interface{}{"username": "john_doe"}}, []interface{}{map[string]interface{}{"username": "john_doe"}}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := logger.maskArrayFields(tc.input)
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})

	t.Run("testConvertStructToFields", func(t *testing.T) {
		type SampleStruct struct {
			Name  string
			Email string
		}
		cases := []struct {
			name  string
			input SampleStruct
			want  map[string]interface{}
		}{
			{"BasicStruct", SampleStruct{Name: "John Doe", Email: "john@doe.com"}, map[string]interface{}{"Name": "John Doe", "Email": "john@doe.com"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got := logger.ConvertStructToFields(tc.input)
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			})
		}
	})
}

// TestEventFunctions tests functionalities of Event methods such as WithField, WithFields, and WithError.
func TestEventFunctions(t *testing.T) {
	logger = newStandardLogger()
	event := &Event{
		Event: zlog.Info(),
	}

	t.Run("testWithField", func(t *testing.T) {
		cases := []struct {
			name  string
			key   string
			value interface{}
		}{
			{"SensitiveField", "password", "secret"},
			{"NonSensitiveField", "username", "john_doe"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				modifiedEvent := event.WithField(tc.key, tc.value)
				if modifiedEvent.Event == nil {
					t.Fatal("Expected modified event to be non-nil")
				}
				// Additional assertions to verify the field
			})
		}
	})

	t.Run("testWithFields", func(t *testing.T) {
		cases := []struct {
			name   string
			fields map[string]interface{}
		}{
			{"SingleSensitiveField", map[string]interface{}{"password": "secret"}},
			{"SingleNonSensitiveField", map[string]interface{}{"username": "john_doe"}},
			{"MixedFields", map[string]interface{}{"password": "secret", "username": "john_doe"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				modifiedEvent := event.WithFields(tc.fields)
				if modifiedEvent.Event == nil {
					t.Fatal("Expected modified event to be non-nil")
				}
				// Additional assertions to verify the fields
			})
		}
	})

	t.Run("testWithError", func(t *testing.T) {
		cases := []struct {
			name string
			err  error
		}{
			{"BasicError", errors.New("sample error")},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				modifiedEvent := event.WithError(tc.err)
				if modifiedEvent.Event == nil {
					t.Fatal("Expected modified event to be non-nil")
				}
				// Additional assertions to verify the error
			})
		}
	})
}

// TestLoggerFunctions runs a series of subtests to verify the functionality of different logger initialization and utility functions.
func TestLoggerFunctions(t *testing.T) {
	t.Run("testNewStandardLogger", func(t *testing.T) {
		logger := newStandardLogger()
		if logger == nil {
			t.Error("expected Logger instance, got nil")
		}
	})

	t.Run("testNewLogger", func(t *testing.T) {
		cases := []struct {
			name   string
			config *Config
		}{
			{"DebugLevel", &Config{Level: "debug"}},
			{"InfoLevel", &Config{Level: "info"}},
			{"WithMasking", &Config{Level: "debug", Masking: MaskingConfig{Enabled: true, SensitiveFields: []string{"password"}}}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				logger := newLogger(tc.config)
				if logger == nil {
					t.Error("expected Logger instance, got nil")
				}
			})
		}
	})

	t.Run("testInitZerologLogger", func(t *testing.T) {
		cases := []struct {
			name   string
			config *Config
			level  zerolog.Level
		}{
			{"DefaultDebugLevel", &Config{Level: "debug", Masking: MaskingConfig{Enabled: true, SensitiveFields: []string{"password"}}}, zerolog.DebugLevel},
			{"InfoLevel", &Config{Level: "info"}, zerolog.InfoLevel},
			{"InvalidLevelFallbackToDebug", &Config{Level: "invalid"}, zerolog.DebugLevel},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				zerologLogger := initZerologLogger(tc.config)
				if zerologLogger.GetLevel() != tc.level {
					t.Errorf("expected %v level, got %v", tc.level, zerologLogger.GetLevel())
				}
			})
		}
	})

	t.Run("testIsSensitiveField", func(t *testing.T) {
		sensitiveFields["password"] = struct{}{}
		cases := []struct {
			name     string
			field    string
			expected bool
		}{
			{"SensitiveField", "password", true},
			{"NonSensitiveField", "username", false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := isSensitiveField(tc.field); got != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
			})
		}
	})

	t.Run("testIsPrimitiveType", func(t *testing.T) {
		cases := []struct {
			name     string
			value    interface{}
			expected bool
		}{
			{"BoolType", true, true},
			{"IntType", 42, true},
			{"StringType", "test", true},
			{"StructType", struct{}{}, false},
			{"MapType", map[string]interface{}{}, false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got := isPrimitiveType(tc.value); got != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
			})
		}
	})
}
