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
func New(fmt string, w io.Writer) *Logger {
	switch fmt {
	case "plain":
		return NewText(w, PlainTextFormat)
	case "text":
		return NewText(w, DefaultTextFormat)
	default:
		return NewJSON(w)
	}
}

// Logger provides interface for logging
type Logger struct {
	back backend
	rec  map[string]interface{}
}

// NewWithBackend creates logger with given backend
func NewWithBackend(b backend) *Logger {
	return &Logger{back: b}
}

// With creates new logger from existing one with predefined parameters
func (l *Logger) With(rec map[string]interface{}) *Logger {
	return &Logger{back: l.back, rec: merge(l.rec, rec)}
}

// Channel creates new logger from existing one with `channel` parameter preset
func (l *Logger) Channel(name string) *Logger {
	return l.With(map[string]interface{}{"channel": name})
}

// Debug message
func (l *Logger) Debug(msg string, record map[string]interface{}) {
	l.record("DEBUG", msg, record)
}

// Debugf message in prtinf format
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.record("DEBUG", fmt.Sprintf(format, args...), nil)
}

// Info message
func (l *Logger) Info(msg string, record map[string]interface{}) {
	l.record("INFO", msg, record)
}

// Infof message in prtinf format
func (l *Logger) Infof(format string, args ...interface{}) {
	l.record("INFO", fmt.Sprintf(format, args...), nil)
}

// Error message
func (l *Logger) Error(msg string, record map[string]interface{}) {
	l.record("ERROR", msg, record)
}

// Errorf message in prtinf format
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.record("ERROR", fmt.Sprintf(format, args...), nil)
}

// Warning message
func (l *Logger) Warning(msg string, record map[string]interface{}) {
	l.record("WARNING", msg, record)
}

// Warningf message in prtinf format
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.record("WARNING", fmt.Sprintf(format, args...), nil)
}

// Fatal signals a fatal error in application (typically during bootstrapping) and stop application execution
// Never use Fatal in places other than application bootstrapping.
func (l *Logger) Fatal(v ...interface{}) {
	l.record("ALERT", fmt.Sprint(v...), nil)
	os.Exit(-1)
}

func (l *Logger) record(level, message string, record map[string]interface{}) {
	var fromContext map[string]interface{}

	if ctx, ok := record["context"].(context.Context); ok {
		delete(record, "context")
		fromContext = RecordFromContext(ctx)
	}

	l.back.Record(merge(l.rec, fromContext, record, map[string]interface{}{
		"message":  message,
		"level":    level,
		"location": location(4),
	}))
}

func location(skip int) string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(skip, fpcs)
	if n == 0 {
		return ""
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return ""
	}

	file, line := fun.FileLine(fpcs[0] - 1)

	return fmt.Sprintf("%v:%v", path.Base(file), line)
}
