package sdk

import (
	"fmt"
	"log"
)

// ConsoleLogger provides a Logger implementation which writes the log messages
// to StdOut - in Terraform's perspective that's proxied via the Plugin SDK
type ConsoleLogger struct {
}

// Info prints out a message prefixed with `[INFO]` verbatim
func (ConsoleLogger) Info(message string) {
	log.Print(fmt.Sprintf("[INFO] %s", message))
}

// Infof prints out a message prefixed with `[INFO]` formatted
// with the specified arguments
func (l ConsoleLogger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn prints out a message prefixed with `[WARN]` formatted verbatim
func (l ConsoleLogger) Warn(message string) {
	log.Print(fmt.Sprintf("[WARN] %s", message))
}

// Warnf prints out a message prefixed with `[WARN]` formatted
// with the specified arguments
func (l ConsoleLogger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}
