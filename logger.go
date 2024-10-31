package zlogs

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type (
	Logger struct {
		*zerolog.Logger
		Masking ConfigMasking
	}
	Config struct {
		AppName      string
		Level        string        `mapstructure:"level"`
		Masking      ConfigMasking `mapstructure:"masking"`
		CallerEnable bool          `mapstructure:"callerEnable"`
	}

	ConfigMasking struct {
		Enabled         bool     `mapstructure:"enabled"`
		SensitiveFields []string `default:""`
	}
	Event struct {
		*zerolog.Event
	}
)

var (
	sensitiveFields = map[string]struct{}{
		"name": {}, "firstname": {}, "lastname": {}, "cardno": {}, "passport": {},
		"passportid": {}, "passportno": {}, "nationalid": {}, "cid": {},
		"citizen_id": {}, "cvv": {}, "password": {}, "x-api-key": {},
		"authorization": {}, "x-authorization": {},
	}
	appNameKey = "appName"

	loggerInstance *Logger
	once           sync.Once
)

func newStandardLogger() *Logger {
	configDefault := &Config{
		Level: "debug",
		Masking: ConfigMasking{
			Enabled: false,
		},
	}
	return newLogger(configDefault)
}

func newLogger(config *Config) *Logger {
	once.Do(func() {
		logger := zerolog.New(os.Stdout).Hook(&InitHook{
			appName:       config.AppName,
			disableCaller: config.CallerEnable,
		}).With().Timestamp().Logger()

		zerolog.TimeFieldFormat = "2006-01-02T15:04:05-0700"
		zerolog.TimestampFieldName = "timestamp"
		zerolog.LevelFieldName = "severity"
		zerolog.MessageFieldName = "message"

		if lvl, err := zerolog.ParseLevel(config.Level); err != nil {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(lvl)
		}

		for _, field := range config.Masking.SensitiveFields {
			sensitiveFields[strings.ToLower(field)] = struct{}{}
		}

		zlog.Logger = logger
		loggerInstance = &Logger{
			Logger:  &logger,
			Masking: config.Masking,
		}
	})

	return loggerInstance
}

func (l *Logger) maskFields(value map[string]interface{}) map[string]interface{} {
	newData := make(map[string]interface{}, len(value))
	for key, fieldValue := range value {
		newData = l.valueMasking(newData, key, fieldValue)
	}
	return newData
}

func (l *Logger) valueMasking(newData map[string]interface{}, key string, fieldValue interface{}) map[string]interface{} {
	// If the field is sensitive
	if _, exists := sensitiveFields[strings.ToLower(key)]; exists {
		newData[key] = "***"
		return newData
	}
	switch subFieldValue := fieldValue.(type) {
	case map[string]interface{}:
		newData[key] = l.maskFields(subFieldValue)
	case []interface{}:
		newData[key] = l.maskArrayFields(subFieldValue)
	default:
		newData[key] = subFieldValue
	}
	return newData
}

func (l *Logger) maskArrayFields(arrayFieldValue []interface{}) []interface{} {
	for index, value := range arrayFieldValue {
		if valueMap, ok := value.(map[string]interface{}); ok {
			arrayFieldValue[index] = l.maskFields(valueMap)
		}
	}
	return arrayFieldValue
}

func (e *Event) WithField(key string, value interface{}) *Event {
	fields := std.maskFields(map[string]interface{}{
		key: convertStructToFields(value),
	})
	return e.WithFields(fields)
}

func (e *Event) WithFields(fields map[string]interface{}) *Event {
	fields = std.maskFields(fields)
	for key, fieldValue := range fields {
		e.Event = e.Event.Interface(key, fieldValue)
	}
	return e
}

func (e *Event) WithError(err error) *Event {
	e.Event = e.Err(err)
	return e
}

func convertStructToFields(v any) map[string]interface{} {
	marshal, _ := json.Marshal(v)
	fields := make(map[string]interface{})
	_ = json.Unmarshal(marshal, &fields)
	return fields
}
