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

// Discard logged records
var Discard = NewLogger(&nopBackend{})

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

	if e, ok := rec["error"]; ok {
		message = fmt.Sprint(message, ": ", e)
		delete(rec, "error")
	}

	data, err := json.Marshal(rec)
	if err != nil {
		data = []byte(fmt.Sprintf("[!marshalling error: %v!]", err))
	}

	color := "1;37"
	switch severity {
	case "ERROR":
		color = "1;31"
	case "INFO":
		color = "1;34"
	}

	return fmt.Sprintf(
		"\033[2;37m[%v]\033[0m \033[%vm%v.%v\033[0m \033[37m%v\033[0m \033[2;37m%s\033[0m",
		time.Now().UTC().Format("2006-01-02T15:04:05"),
		color,
		channel,
		severity,
		message,
		data,
	)
}

// Record an entry
func (b *TextBackend) Record(rec R) {
	if rec == nil {
		rec = R{}
	}

	fmt.Fprintln(b.Writer, b.Formatter(rec))
}

type nopBackend struct {
}

// Record an entry
func (b *nopBackend) Record(rec R) {
}
