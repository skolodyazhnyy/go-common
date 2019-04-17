package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
)

// Discard logged records
var Discard = NewWithBackend(&nopBackend{})

// New creates logger with a given backend
func New(fmt string, w io.Writer, f func(R) string) *Logger {
	switch fmt {
	case "json":
		return NewJSON(w)
	default:
		return NewText(w, f)
	}
}

// Logger provides interface for logging
type Logger struct {
	back backend
	rec  R
}

// NewWithBackend creates logger with given backend
func NewWithBackend(b backend) *Logger {
	return &Logger{back: b}
}

// With creates new logger from existing one with predefined parameters
func (l *Logger) With(rec map[string]interface{}) *Logger {
	if l.rec == nil {
		l.rec = R{}
	}

	return &Logger{back: l.back, rec: l.rec.With(rec)}
}

// Channel creates new logger from existing one with `channel` parameter preset
func (l *Logger) Channel(name string) *Logger {
	return l.With(R{"channel": name})
}

func (l *Logger) debug(msg string, rec R) {
	if rec == nil {
		rec = make(R)
	}

	rec["message"] = msg
	rec["level"] = "DEBUG"
	rec["location"] = location(4)

	l.record(rec.With(l.rec))
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
	rec["level"] = "INFO"
	rec["location"] = location(4)

	l.record(rec.With(l.rec))
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
	rec["level"] = "ERROR"
	rec["location"] = location(4)

	l.record(rec.With(l.rec))
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
	rec["level"] = "WARNING"
	rec["location"] = location(4)

	l.record(rec.With(l.rec))
}

func (l *Logger) Warning(msg string, rec map[string]interface{}) {
	l.warning(msg, rec)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.warning(fmt.Sprintf(format, args...), nil)
}

// Fatal signals a fatal error in application (typically during bootstrapping) and stop application execution
// Never use Fatal in places other than application bootstrapping.
func (l *Logger) Fatal(v ...interface{}) {
	l.record(l.rec.With(R{
		"message":  fmt.Sprint(v...),
		"level": "ALERT",
		"location": location(3),
	}))

	os.Exit(-1)
}

func (l *Logger) record(r R) {
	if ctx, ok := r["context"].(context.Context); ok {
		delete(r, "context")
		r = r.With(RecordFromContext(ctx))
	}

	l.back.Record(r)
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
