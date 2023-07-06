// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"fmt"
	"log"
)

var _ Logger = ConsoleLogger{}

// ConsoleLogger provides a Logger implementation which writes the log messages
// to StdOut - in Terraform's perspective that's proxied via the Plugin SDK
type ConsoleLogger struct{}

// Info prints out a message prefixed with `[INFO]` verbatim
func (ConsoleLogger) Info(message string) {
	log.Printf("[INFO] %s", message)
}

// Infof prints out a message prefixed with `[INFO]` formatted
// with the specified arguments
func (l ConsoleLogger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn prints out a message prefixed with `[WARN]` formatted verbatim
func (l ConsoleLogger) Warn(message string) {
	log.Printf("[WARN] %s", message)
}

// Warnf prints out a message prefixed with `[WARN]` formatted
// with the specified arguments
func (l ConsoleLogger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}
