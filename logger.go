package zlogs

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// redactedValue is the constant used to replace sensitive data fields with a masked value to prevent information leakage.
const redactedValue = "***"

// Logger wraps zerolog.Logger and includes ConfigMasking for masking purposes.
type (
	Logger struct {
		*zerolog.Logger
		Masking MaskingConfig
	}
	Config struct {
		AppName      string
		Level        string
		Masking      MaskingConfig
		CallerEnable bool
	}
	MaskingConfig struct {
		Enabled         bool
		SensitiveFields []string
	}
	Event struct {
		*zerolog.Event
	}
)

// sensitiveFields contains a set of keys that are considered sensitive and need to be masked in logs or data processing.
var (
	sensitiveFields = map[string]struct{}{
		"name": {}, "firstname": {}, "lastname": {}, "cardno": {}, "passport": {},
		"passportid": {}, "passportno": {}, "nationalid": {}, "cid": {},
		"citizen_id": {}, "cvc": {}, "password": {}, "x-api-key": {},
		"authorization": {}, "x-authorization": {},
	}
	appNameKey = "appName"
	logger     *Logger
)

// newStandardLogger initializes a standard Logger instance with default configuration for level "debug" and masking disabled.
func newStandardLogger() *Logger {
	defaultConfig := &Config{
		Level: "debug",
		Masking: MaskingConfig{
			Enabled: false,
		},
	}
	return newLogger(defaultConfig)
}

// newLogger initializes the global logger instance with the provided configuration.
func newLogger(config *Config) *Logger {
	zlog.Logger = initZerologLogger(config)
	logger = &Logger{
		Logger:  &zlog.Logger,
		Masking: config.Masking,
	}

	return logger
}

// initZerologLogger initializes and configures a zerolog.Logger instance based on the provided configuration.
func initZerologLogger(config *Config) zerolog.Logger {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "severity"
	zerolog.MessageFieldName = "message"
	if level, err := zerolog.ParseLevel(config.Level); err != nil {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
	for _, field := range config.Masking.SensitiveFields {
		sensitiveFields[strings.ToLower(field)] = struct{}{}
	}
	return zerolog.New(os.Stdout).Hook(&InitHook{
		AppName:       config.AppName,
		DisableCaller: config.CallerEnable,
	}).With().Timestamp().Logger().Level(zerolog.GlobalLevel())
}

// maskFields processes a map to mask sensitive fields based on the Logger configuration.
func (l *Logger) maskFields(value map[string]interface{}) map[string]interface{} {
	newData := make(map[string]interface{}, len(value))
	for key, fieldValue := range value {
		newData = l.valueMasking(newData, key, fieldValue)
	}
	return newData
}

// valueMasking masks sensitive fields in a nested map or array structure, otherwise it retains the original field value.
func (l *Logger) valueMasking(result map[string]interface{}, key string, value interface{}) map[string]interface{} {
	if isSensitiveField(key) {
		result[key] = redactedValue
		return result
	}
	switch v := value.(type) {
	case map[string]interface{}:
		result[key] = l.maskFields(v)
	case []interface{}:
		result[key] = l.maskArrayFields(v)
	default:
		result[key] = v
	}
	return result
}

// isSensitiveField checks if a given field name is considered sensitive based on a predefined list of sensitive fields.
func isSensitiveField(field string) bool {
	_, exists := sensitiveFields[strings.ToLower(field)]
	return exists
}

// maskArrayFields iterates over an array of interface values and applies field masking to any map elements.
func (l *Logger) maskArrayFields(array []interface{}) []interface{} {
	for i, value := range array {
		if valueMap, ok := value.(map[string]interface{}); ok {
			array[i] = l.maskFields(valueMap)
		}
	}
	return array
}

// WithField adds a key-value pair to the event, converting and masking the value if necessary, and returns the updated event.
func (e *Event) WithField(key string, value interface{}) *Event {
	if !isPrimitiveType(value) {
		value = logger.ConvertStructToFields(value)
	}
	return e.WithFields(logger.maskFields(map[string]interface{}{key: value}))
}

// isPrimitiveType determines if the given value is of a primitive Go type that does not require conversion.
func isPrimitiveType(value interface{}) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return true
	default:
		return false
	}
}

// WithFields adds multiple fields to the event, masking sensitive data, and returns the updated event.
func (e *Event) WithFields(fields map[string]interface{}) *Event {
	fields = logger.maskFields(fields)
	for key, fieldValue := range fields {
		e.Event = e.Event.Interface(key, fieldValue)
	}
	return e
}

// WithError attaches an error to the Event and returns the modified Event.
func (e *Event) WithError(err error) *Event {
	e.Event = e.Err(err)
	return e
}

// ConvertStructToFields converts a struct to a map of string keys and interface{} values using JSON marshaling.
func (l *Logger) ConvertStructToFields(v any) map[string]interface{} {
	data, _ := json.Marshal(v)
	fields := make(map[string]interface{})
	_ = json.Unmarshal(data, &fields)
	return fields
}
