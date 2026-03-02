// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tfjson

import (
	"bytes"
	"encoding/json"
)

type LogMessageType string

const (
	MessageTypeVersion    LogMessageType = "version"
	MessageTypeLog        LogMessageType = "log"
	MessageTypeDiagnostic LogMessageType = "diagnostic"
)

// allLogMessageTypes is a slice containing all recognised message types
// to be passed into cmp.AllowUnexported
var allLogMessageTypes = []any{
	VersionLogMessage{},
	LogMessage{},
	DiagnosticLogMessage{},
	UnknownLogMessage{},

	// query
	ListStartMessage{},
	ListResourceFoundMessage{},
	ListCompleteMessage{},
}

func unmarshalByType(t LogMessageType, b []byte) (LogMsg, error) {
	d := json.NewDecoder(bytes.NewReader(b))

	// decode numbers as json.Number to avoid losing precision
	d.UseNumber()

	switch t {

	// generic
	case MessageTypeVersion:
		v := VersionLogMessage{}
		return v, d.Decode(&v)
	case MessageTypeLog:
		v := LogMessage{}
		return v, d.Decode(&v)
	case MessageTypeDiagnostic:
		v := DiagnosticLogMessage{}
		return v, d.Decode(&v)

	// query
	case MessageListStart:
		v := ListStartMessage{}
		return v, d.Decode(&v)
	case MessageListResourceFound:
		v := ListResourceFoundMessage{}
		return v, d.Decode(&v)
	case MessageListComplete:
		v := ListCompleteMessage{}
		return v, d.Decode(&v)
	}

	v := UnknownLogMessage{}
	return v, d.Decode(&v)
}
