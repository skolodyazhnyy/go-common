package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func TimelessTextFormat(rec R) string {
	channel := rec["channel"]
	if channel == nil {
		channel = "main"
	}
	delete(rec, "channel")

	severity := rec["level"]
	if severity == nil {
		severity = "?"
	}
	delete(rec, "level")

	message := rec["message"]
	delete(rec, "message")

	// remove location so it does not screw up tests
	delete(rec, "location")

	data, err := json.Marshal(rec)
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
	// main.INFO oops.. something went wrong {}
	// main.INFO OK status code is 200 {}
}

func ExampleLogger_With() {
	log := NewText(os.Stdout, TimelessTextFormat)

	// Structured data for a single record
	log.Info("oops.. something went wrong", R{"foo": "bar"})

	// Structured data for all records
	logWithData := log.With(R{"foo": "bazzz"})
	logWithData.Infof("OK status code is %v", http.StatusOK)
	logWithData.Infof("NotFound status code is %v", http.StatusNotFound)

	// Output:
	// main.INFO oops.. something went wrong {"foo":"bar"}
	// main.INFO OK status code is 200 {"foo":"bazzz"}
	// main.INFO NotFound status code is 404 {"foo":"bazzz"}
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
	// main.INFO OK status code is 200 {}
	// db.INFO NotFound status code is 404 {}
}

func ExampleLogger_Context() {
	log := NewText(os.Stdout, TimelessTextFormat)

	ctx := AppendContext(context.Background(), map[string]interface{}{"foo": "bar"})

	// Create logger with channel
	log.Info("Additional data in context", map[string]interface{}{
		"context": ctx,
	})

	// Output:
	// main.INFO Additional data in context {"foo":"bar"}
}
