// Package logger represents a generic logging interface
package logger

import (
	"context"
	"io"
	"os"
)

// Log is a package level variable, every program should access logging function through "Log"
var Log Logger

// Logger represent common interface for logging function
type Logger interface {
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
}

// SetLogger is the setter for log variable, it should be the only way to assign value to log
func SetLogger(newLogger Logger) {
	Log = newLogger
}

func NewWriter(outputPath string) (io.Writer, error) {
	f, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	return f, err
}

// NewAccessLogWriter returns an implementation of io.Writer interface and can be used
// to write access logs.
func NewAccessLogWriter(accessLogPath string) (io.Writer, error) {
	f, err := os.OpenFile(accessLogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	return f, err
}

// NewErrorLogWriter returns an implementation of io.Writer interface and can be used
// to write error logs.
func NewErrorLogWriter() io.Writer {
	return errorLogWriter{}
}

type errorLogWriter struct{}

func (writer errorLogWriter) Write(p []byte) (n int, err error) {
	if Log == nil { // our logger has not been initialised yet
		return 0, errLoggerNotInitialised
	}

	msgToLog := string(p)
	Log.Error(context.Background(), msgToLog)

	return len(p), nil
}

// NewDebugLogWriter returns an implementation of io.Writer interface and can be used
// to write debug logs.
func NewDebugLogWriter() io.Writer {
	return debugLogWriter{}
}

type debugLogWriter struct{}

func (writer debugLogWriter) Write(p []byte) (n int, err error) {
	if Log == nil { // our logger has not been initialised yet
		return 0, errLoggerNotInitialised
	}

	msgToLog := string(p)
	Log.Debug(context.Background(), msgToLog)

	return len(p), nil
}
