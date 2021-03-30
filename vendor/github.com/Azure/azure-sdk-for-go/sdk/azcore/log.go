// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
	"os"
	"time"
)

// LogClassification is used to group entries.  Each group can be toggled on or off.
type LogClassification string

const (
	// LogRequest entries contain information about HTTP requests.
	// This includes information like the URL, query parameters, and headers.
	LogRequest LogClassification = "Request"

	// LogResponse entries contain information about HTTP responses.
	// This includes information like the HTTP status code, headers, and request URL.
	LogResponse LogClassification = "Response"

	// LogRetryPolicy entries contain information specific to the retry policy in use.
	LogRetryPolicy LogClassification = "RetryPolicy"

	// LogLongRunningOperation entries contain information specific to long-running operations.
	// This includes information like polling location, operation state and sleep intervals.
	LogLongRunningOperation LogClassification = "LongRunningOperation"
)

// Listener is the function signature invoked when writing log entries.
// A Listener is required to perform its own synchronization if it's
// expected to be called from multiple Go routines.
type Listener func(LogClassification, string)

// Logger controls which classifications to log and writing to the underlying log.
type Logger struct {
	cls []LogClassification
	lst Listener
}

// SetClassifications is used to control which classifications are written to
// the log.  By default all log classifications are written.
func (l *Logger) SetClassifications(cls ...LogClassification) {
	l.cls = cls
}

// SetListener will set the Logger to write to the specified Listener.
func (l *Logger) SetListener(lst Listener) {
	l.lst = lst
}

// Should returns true if the specified log classification should be written to the log.
// By default all log classifications will be logged.  Call SetClassification() to limit
// the log classifications for logging.
// If no listener has been set this will return false.
// Calling this method is useful when the message to log is computationally expensive
// and you want to avoid the overhead if its log classification is not enabled.
func (l *Logger) Should(cls LogClassification) bool {
	if l.lst == nil {
		return false
	}
	if l.cls == nil || len(l.cls) == 0 {
		return true
	}
	for _, c := range l.cls {
		if c == cls {
			return true
		}
	}
	return false
}

// Write invokes the underlying Listener with the specified classification and message.
// If the classification shouldn't be logged or there is no listener then Write does nothing.
func (l *Logger) Write(cls LogClassification, message string) {
	if l.lst == nil || !l.Should(cls) {
		return
	}
	l.lst(cls, message)
}

// Writef invokes the underlying Listener with the specified classification and formatted message.
// If the classification shouldn't be logged or there is no listener then Writef does nothing.
func (l *Logger) Writef(cls LogClassification, format string, a ...interface{}) {
	if l.lst == nil || !l.Should(cls) {
		return
	}
	l.lst(cls, fmt.Sprintf(format, a...))
}

// for testing purposes
func (l *Logger) resetClassifications() {
	l.cls = nil
}

var logger Logger

// Log returns the process-wide logger.
func Log() *Logger {
	return &logger
}

func init() {
	if cls := os.Getenv("AZURE_SDK_GO_LOGGING"); cls == "all" {
		// cls could be enhanced to support a comma-delimited list of log classifications
		logger.lst = func(cls LogClassification, msg string) {
			// simple console logger, it writes to stderr in the following format:
			// [time-stamp] Classification: message
			fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", time.Now().Format(time.StampMicro), cls, msg)
		}
	}
}
