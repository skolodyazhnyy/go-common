package log

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	Backend backend
	Record  R
}

func (l *Logger) With(rec map[string]interface{}) *Logger {
	if l.Record == nil {
		l.Record = R{}
	}

	return &Logger{Backend: l.Backend, Record: l.Record.With(rec)}
}

func (l *Logger) Channel(name string) *Logger {
	return l.With(R{"channel": name})
}

func (l *Logger) debug(msg string, rec R) {
	if rec == nil {
		rec = make(R)
	}

	rec["message"] = msg
	rec["severity"] = "DEBUG"

	l.record(rec.With(l.Record))
}

func (l *Logger) Debug(msg string, rec map[string]interface{}) {
	l.debug(msg, rec)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.debug(fmt.Sprintf(format, args...), nil)
}

func (l *Logger) info(msg string, rec R) {
	if rec == nil {
		rec = make(R)
	}

	rec["message"] = msg
	rec["severity"] = "INFO"

	l.record(rec.With(l.Record))
}

func (l *Logger) Info(msg string, rec map[string]interface{}) {
	l.info(msg, rec)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...), nil)
}

func (l *Logger) error(msg string, rec R) {
	if rec == nil {
		rec = make(R)
	}

	rec["message"] = msg
	rec["severity"] = "ERROR"
	rec["location"] = location(4)

	l.record(rec.With(l.Record))
}

func (l *Logger) Error(msg string, rec map[string]interface{}) {
	l.error(msg, rec)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.error(fmt.Sprintf(format, args...), nil)
}

func (l *Logger) warning(msg string, rec R) {
	if rec == nil {
		rec = make(R)
	}

	rec["message"] = msg
	rec["severity"] = "WARNING"
	rec["location"] = location(4)

	l.record(rec.With(l.Record))
}

func (l *Logger) Warning(msg string, rec map[string]interface{}) {
	l.warning(msg, rec)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.warning(fmt.Sprintf(format, args...), nil)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.record(l.Record.With(R{
		"message":  fmt.Sprint(v...),
		"severity": "ALERT",
		"location": location(3),
	}))

	os.Exit(-1)
}

func (l *Logger) record(r R) {
	if ctx, ok := r["context"].(context.Context); ok {
		delete(r, "context")
		r = r.With(RecordFromContext(ctx))
	}

	l.Backend.Record(r)
}

func location(skip int) R {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(skip, fpcs)
	if n == 0 {
		return R{}
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return R{}
	}

	file, line := fun.FileLine(fpcs[0] - 1)

	return R{
		"function": path.Base(fun.Name()),
		"file":     path.Base(file),
		"line":     line,
	}
}
