package log

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

	return fmt.Sprintf("%v.%v %v %s", channel, severity, message, data)
}

func ExampleLogger_Info() {
	logger := NewText(os.Stdout, TimelessTextFormat)

	logger.Info("oops.. something went wrong", nil)
	logger.Infof("OK status code is %v", http.StatusOK)

	// Output:
	// main.INFO oops.. something went wrong {}
	// main.INFO OK status code is 200 {}
}

func ExampleLogger_With() {
	logger := NewText(os.Stdout, TimelessTextFormat)

	// Structured data for a single record
	logger.Info("oops.. something went wrong", R{"foo": "bar"})

	// Structured data for all records
	loggerWithData := logger.With(R{"foo": "bazzz"})
	loggerWithData.Infof("OK status code is %v", http.StatusOK)
	loggerWithData.Infof("NotFound status code is %v", http.StatusNotFound)

	// Output:
	// main.INFO oops.. something went wrong {"foo":"bar"}
	// main.INFO OK status code is 200 {"foo":"bazzz"}
	// main.INFO NotFound status code is 404 {"foo":"bazzz"}
}

func ExampleLogger_Channel() {
	logger := NewText(os.Stdout, TimelessTextFormat)

	// Create logger with channel
	dblogger := logger.Channel("db")

	// Log to main channel
	logger.Infof("OK status code is %v", http.StatusOK)

	// Log to db channel
	dblogger.Infof("NotFound status code is %v", http.StatusNotFound)

	// Output:
	// main.INFO OK status code is 200 {}
	// db.INFO NotFound status code is 404 {}
}

func ExampleLogger_Context() {
	logger := NewText(os.Stdout, TimelessTextFormat)

	ctx := AppendContext(context.Background(), map[string]interface{}{"foo": "bar"})

	// Create logger with channel
	logger.Info("Additional data in context", map[string]interface{}{
		"context": ctx,
	})

	// Output:
	// main.INFO Additional data in context {"foo":"bar"}
}
