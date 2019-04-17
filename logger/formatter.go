package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

// PlainTextFormat formats structured data into human friendly string without colors
func PlainTextFormat(rec R) string {
	channel := rec["channel"]
	delete(rec, "channel")
	if channel == nil {
		channel = "main"
	}

	severity := rec["level"]
	delete(rec, "level")
	if severity == nil {
		severity = "?"
	}

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

	return fmt.Sprintf(
		"[%v] %v.%v %v %s",
		time.Now().UTC().Format("2006-01-02T15:04:05"),
		channel,
		severity,
		message,
		data,
	)
}

// DefaultTextFormat formats structured data into human friendly string
func DefaultTextFormat(rec R) string {
	channel := rec["channel"]
	delete(rec, "channel")
	if channel == nil {
		channel = "main"
	}

	severity := rec["level"]
	delete(rec, "level")
	if severity == nil {
		severity = "?"
	}

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
