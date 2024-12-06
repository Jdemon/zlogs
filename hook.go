package zlogs

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

// contextKey is a custom type used for defining unique keys for context values in Go applications.
type contextKey int

// TraceID is used as a context key for tracing the execution path of a request.
// CorrelationID is used as a context key for grouping and correlating related log entries.
// RequestID is used as a context key for uniquely identifying an individual request.
// CallerSkip is used as a context key for controlling the skip level in caller information retrieval.
const (
	TraceID contextKey = iota
	CorrelationID
	RequestID
	CallerSkip
)

// InitHook is a type used to initialize log entries with application-specific data and optionally include caller information.
type InitHook struct {
	AppName       string
	DisableCaller bool
}

// Run sets entry data and caller information into the zerolog event. It retrieves context values and updates the event accordingly.
func (h *InitHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	var callerSkip = 5
	if !h.DisableCaller {
		if skip, ok := e.GetCtx().Value(CallerSkip).(int); ok {
			callerSkip = skip
		}
		defer caller(e, callerSkip)
	}

	setEntryData(e, appNameKey, h.AppName)
	if e.GetCtx() == nil {
		return
	}
	setEntryData(e, "trace_id", e.GetCtx().Value(TraceID))
	setEntryData(e, "request_id", e.GetCtx().Value(RequestID))
	setEntryData(e, "correlation_id", e.GetCtx().Value(CorrelationID))
}

// setEntryData adds a key-value pair to the given zerolog event if the value is neither nil nor an empty string.
func setEntryData(e *zerolog.Event, key string, value interface{}) {
	if value != nil && value != "" {
		e.Interface(key, value)
	}
}

// caller enriches the provided zerolog.Event by adding file and function name information based on the stack skip level.
func caller(event *zerolog.Event, skip int) *zerolog.Event {
	file, fnc := fileInfo(skip)
	event.Str("file", file)
	event.Str("func", fnc)
	return event
}

// fileInfo returns the file name and line number of the caller located `skip` frames up the call stack,
// and the name of the function where the call occurred.
func fileInfo(skip int) (string, string) {
	var funcName string
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}

		funcName = runtime.FuncForPC(pc).Name()
		funcName = funcName[strings.LastIndex(funcName, ".")+1:]
	}

	return fmt.Sprintf("%s:%d", file, line), funcName
}
