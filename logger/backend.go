package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type backend interface {
	Record(map[string]interface{})
}

// NewText creates logger with text backend
func NewText(w io.Writer, f func(map[string]interface{}) string) *Logger {
	if f == nil {
		f = DefaultTextFormat
	}

	return NewWithBackend(&textBackend{writer: w, formatter: f})
}

// NewJSON creates logger with JSON backend
func NewJSON(w io.Writer) *Logger {
	return NewWithBackend(&jsonBackend{writer: w})
}

// jsonBackend writes logs as JSON objects
type jsonBackend struct {
	writer io.Writer
}

// Record an entry
func (b *jsonBackend) Record(rec map[string]interface{}) {
	if rec == nil {
		rec = make(map[string]interface{})
	}

	rec["time"] = time.Now().UTC()

	//nolint:errcheck
	json.NewEncoder(b.writer).Encode(rec)
}

// textBackend writes logs as formatted text
type textBackend struct {
	writer    io.Writer
	formatter func(map[string]interface{}) string
}

// Record an entry
func (b *textBackend) Record(rec map[string]interface{}) {
	if rec == nil {
		rec = make(map[string]interface{})
	}

	//nolint:errcheck
	fmt.Fprintln(b.writer, b.formatter(rec))
}

// nopBackend discards all messages
type nopBackend struct {
}

// Record an entry
func (b *nopBackend) Record(rec map[string]interface{}) {
}
