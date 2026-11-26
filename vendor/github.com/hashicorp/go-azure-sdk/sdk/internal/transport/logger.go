// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transport

import (
	"fmt"
	"log"
	"regexp"
)

var sensitiveQueryParamRegexp = regexp.MustCompile(`(?i)([?&](?:sig|code|api-key|token|client_secret)=)([^&\s"']+)`)

type redactingLogger struct{}

func (redactingLogger) Printf(format string, v ...interface{}) {
	log.Print(sensitiveQueryParamRegexp.ReplaceAllString(fmt.Sprintf(format, v...), "${1}REDACTED"))
}

// RedactingLogger returns a retryablehttp.Logger that delegates to the standard
// log package (respecting log.SetOutput and log.Default()'s internal mutex) while
// redacting sensitive query-string parameters (such as SAS sig tokens and secrets).
func RedactingLogger() interface{ Printf(string, ...interface{}) } {
	return redactingLogger{}
}
