package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type backend interface {
	Record(R)
}

// NewText creates logger with text backend
func NewText(w io.Writer, f func(R) string) *Logger {
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
func (b *jsonBackend) Record(rec R) {
	if rec == nil {
		rec = R{}
	}

	rec["time"] = time.Now().UTC()

	json.NewEncoder(b.writer).Encode(rec)
}

// textBackend writes logs as formatted text
type textBackend struct {
	writer    io.Writer
	formatter func(R) string
}

// Record an entry
func (b *textBackend) Record(rec R) {
	if rec == nil {
		rec = R{}
	}

	fmt.Fprintln(b.writer, b.formatter(rec))
}

// nopBackend discards all messages
type nopBackend struct {
}

// Record an entry
func (b *nopBackend) Record(rec R) {
}
