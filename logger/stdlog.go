package logger

import (
	"fmt"
	"io"
)

// NewStdLogger creates new instance of standard logger
func NewStdLogger(log *Logger) *StdLogger {
	return &StdLogger{
		log: log,
	}
}

// StdLogger implements interface of go standard library logger
type StdLogger struct {
	log *Logger
}

func (l *StdLogger) SetOutput(w io.Writer) {
}

func (l *StdLogger) Output(calldepth int, s string) error {
	return nil
}

func (l *StdLogger) Printf(format string, v ...interface{}) {
	l.log.Debugf(format, v...)
}

func (l *StdLogger) Print(v ...interface{}) {
	l.log.Debug(fmt.Sprint(v...), nil)
}

func (l *StdLogger) Println(v ...interface{}) {
	l.log.Debug(fmt.Sprint(v...), nil)
}

func (l *StdLogger) Fatal(v ...interface{}) {
	l.log.Fatal(v...)
}

func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	l.log.Fatal(fmt.Sprintf(format, v...))
}

func (l *StdLogger) Fatalln(v ...interface{}) {
	l.log.Fatal(v...)
}

func (l *StdLogger) Panic(v ...interface{}) {
	l.log.Fatal(v...)
}

func (l *StdLogger) Panicf(format string, v ...interface{}) {
	l.log.Fatal(fmt.Sprintf(format, v...))
}

func (l *StdLogger) Panicln(v ...interface{}) {
	l.log.Fatal(v...)
}

func (l *StdLogger) Flags() int {
	return 0
}

func (l *StdLogger) SetFlags(flag int) {
}

func (l *StdLogger) Prefix() string {
	return ""
}

func (l *StdLogger) SetPrefix(prefix string) {
}
