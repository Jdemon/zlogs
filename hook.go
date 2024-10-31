package zlogs

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type contextKey int

const (
	TraceID contextKey = iota
	CorrelationID
	RequestID
	CallerSkip
)

type InitHook struct {
	appName       string
	disableCaller bool
}

func (h *InitHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	var callerSkip = 6
	if !h.disableCaller {
		if skip, ok := e.GetCtx().Value(CallerSkip).(int); ok {
			callerSkip = skip
		}
		defer caller(e, callerSkip)
	}

	setEntryData(e, appNameKey, h.appName)
	if e.GetCtx() == nil {
		return
	}
	setEntryData(e, "trace_id", e.GetCtx().Value(TraceID))
	setEntryData(e, "request_id", e.GetCtx().Value(RequestID))
	setEntryData(e, "correlation_id", e.GetCtx().Value(CorrelationID))
}

func setEntryData(e *zerolog.Event, key string, value interface{}) {
	if value != nil && value != "" {
		e.Interface(key, value)
	}
}

func caller(event *zerolog.Event, skip int) *zerolog.Event {
	file, fnc := fileInfo(skip)
	event.Str("file", file)
	event.Str("func", fnc)
	return event
}

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
