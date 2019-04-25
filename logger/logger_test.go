package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TimelessTextFormat(record map[string]interface{}) string {
	channel := record["channel"]
	if channel == nil {
		channel = "main"
	}
	delete(record, "channel")

	severity := record["level"]
	if severity == nil {
		severity = "?"
	}
	delete(record, "level")

	message := record["message"]
	delete(record, "message")

	data, err := json.Marshal(record)
	if err != nil {
		data = []byte(fmt.Sprintf("[!marshalling error: %v!]", err))
	}

	return fmt.Sprintf("%v.%v %v %s", channel, severity, message, data)
}

func ExampleLogger_Info() {
	log := NewText(os.Stdout, TimelessTextFormat)

	log.Info("oops.. something went wrong", nil)
	log.Infof("OK status code is %v", http.StatusOK)

	// Output:
	// main.INFO oops.. something went wrong {"location":"logger_test.go:40"}
	// main.INFO OK status code is 200 {"location":"logger_test.go:41"}
}

func ExampleLogger_With() {
	log := NewText(os.Stdout, TimelessTextFormat)

	// Structured data for a single record
	log.Info("oops.. something went wrong", map[string]interface{}{"foo": "bar"})

	// Structured data for all records
	logWithData := log.With(map[string]interface{}{"foo": "bazzz"})
	logWithData.Infof("OK status code is %v", http.StatusOK)
	logWithData.Infof("NotFound status code is %v", http.StatusNotFound)

	// Output:
	// main.INFO oops.. something went wrong {"foo":"bar","location":"logger_test.go:52"}
	// main.INFO OK status code is 200 {"foo":"bazzz","location":"logger_test.go:56"}
	// main.INFO NotFound status code is 404 {"foo":"bazzz","location":"logger_test.go:57"}
}

func ExampleLogger_Channel() {
	log := NewText(os.Stdout, TimelessTextFormat)

	// Create logger with channel
	dblogger := log.Channel("db")

	// Log to main channel
	log.Infof("OK status code is %v", http.StatusOK)

	// Log to db channel
	dblogger.Infof("NotFound status code is %v", http.StatusNotFound)

	// Output:
	// main.INFO OK status code is 200 {"location":"logger_test.go:72"}
	// db.INFO NotFound status code is 404 {"location":"logger_test.go:75"}
}

func ExampleAppendContext() {
	log := NewText(os.Stdout, TimelessTextFormat)

	ctx := AppendContext(context.Background(), map[string]interface{}{"foo": "bar"})

	// Create logger with channel
	log.Info("Additional data in context", map[string]interface{}{
		"context": ctx,
	})

	// Output:
	// main.INFO Additional data in context {"foo":"bar","location":"logger_test.go:88"}
}

func BenchmarkLogger_Text(b *testing.B) {
	log := NewText(ioutil.Discard, PlainTextFormat)

	for i := 0; i < b.N; i++ {
		log.Info("oops.. something went wrong", map[string]interface{}{"foo": "bar"})
	}
}

func BenchmarkLogger_JSON(b *testing.B) {
	log := NewJSON(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		log.Info("oops.. something went wrong", map[string]interface{}{"foo": "bar"})
	}
}

func BenchmarkTextBackend(b *testing.B) {
	backend := textBackend{writer: ioutil.Discard, formatter: PlainTextFormat}

	for i := 0; i < b.N; i++ {
		backend.Record(map[string]interface{}{"foo": "bar", "message": "foo", "level": "ERROR"})
	}
}

func BenchmarkJSONBackend(b *testing.B) {
	backend := jsonBackend{writer: ioutil.Discard}

	for i := 0; i < b.N; i++ {
		backend.Record(map[string]interface{}{"foo": "bar", "message": "foo", "level": "ERROR"})
	}
}
