package log

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type backend interface {
	Record(R)
}

// New creates logger with a given backend
func New(fmt string, w io.Writer, f func(R) string) *Logger {
	switch fmt {
	case "json":
		return NewJSON(w)
	default:
		return NewText(w, f)
	}
}

// NewText creates logger with text backend
func NewText(w io.Writer, f func(R) string) *Logger {
	if f == nil {
		f = DefaultTextFormat
	}

	return NewLogger(&TextBackend{Writer: w, Formatter: f})
}

// NewJSON creates logger with JSON backend
func NewJSON(w io.Writer) *Logger {
	return NewLogger(&JSONBackend{Writer: w})
}

// JSONBackend writes logs as JSON objects
type JSONBackend struct {
	Writer io.Writer
}

// Record an entry
func (b *JSONBackend) Record(rec R) {
	if rec == nil {
		rec = R{}
	}

	rec["timestamp"] = time.Now().UTC()

	json.NewEncoder(b.Writer).Encode(rec)
}

// TextBackend writes logs as formatted text
type TextBackend struct {
	Writer    io.Writer
	Formatter func(R) string
}

// DefaultTextFormat formats structured data into human friendly string
func DefaultTextFormat(rec R) string {
	channel := rec["channel"]
	if channel == nil {
		channel = "main"
	}
	delete(rec, "channel")

	severity := rec["severity"]
	if severity == nil {
		severity = "?"
	}
	delete(rec, "severity")

	message := rec["message"]
	delete(rec, "message")

	data, err := json.Marshal(rec)
	if err != nil {
		data = []byte(fmt.Sprintf("[!marshalling error: %v!]", err))
	}

	return fmt.Sprintf("%v %v.%v %v %s", time.Now().Format(time.RFC3339Nano), channel, severity, message, data)
}

// Record an entry
func (b *TextBackend) Record(rec R) {
	if rec == nil {
		rec = R{}
	}

	fmt.Fprintln(b.Writer, b.Formatter(rec))
}
