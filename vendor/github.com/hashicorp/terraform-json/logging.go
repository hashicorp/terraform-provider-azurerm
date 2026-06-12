// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tfjson

import (
	"bytes"
	"encoding/json"
	"time"
)

// LogMessageLevel represents log level
// See https://github.com/hashicorp/go-hclog/blob/v1.6.3/logger.go#L126-L145
type LogMessageLevel string

const (
	// Trace is the most verbose level. Intended to be used for the tracing
	// of actions in code, such as function enters/exits, etc.
	Trace LogMessageLevel = "trace"

	// Debug information for programmer low-level analysis.
	Debug LogMessageLevel = "debug"

	// Info information about steady state operations.
	Info LogMessageLevel = "info"

	// Warn information about rare but handled events.
	Warn LogMessageLevel = "warn"

	// Error information about unrecoverable events.
	Error LogMessageLevel = "error"
)

// LogMessage represents a log message emitted from commands
// which support structured log output.
//
// This is implemented via hashicorp/go-hclog which
// defines the format.
type LogMsg interface {
	Level() LogMessageLevel
	Message() string
	Timestamp() time.Time
}

type baseLogMessage struct {
	Lvl  LogMessageLevel `json:"@level"`
	Msg  string          `json:"@message"`
	Time time.Time       `json:"@timestamp"`
}

type msgType struct {
	// Type represents a message type
	// which is documented at https://developer.hashicorp.com/terraform/internals/machine-readable-ui#message-types
	Type LogMessageType `json:"type"`
}

func (m baseLogMessage) Level() LogMessageLevel {
	return m.Lvl
}

func (m baseLogMessage) Message() string {
	return m.Msg
}

func (m baseLogMessage) Timestamp() time.Time {
	return m.Time
}

// UnknownLogMessage represents a message of unknown type
type UnknownLogMessage struct {
	baseLogMessage
}

func UnmarshalLogMessage(b []byte) (LogMsg, error) {
	d := json.NewDecoder(bytes.NewReader(b))

	mt := msgType{}
	err := d.Decode(&mt)
	if err != nil {
		return nil, err
	}

	v, err := unmarshalByType(mt.Type, b)
	return v, err
}
