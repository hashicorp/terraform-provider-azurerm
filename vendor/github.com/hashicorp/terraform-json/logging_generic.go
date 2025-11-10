// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tfjson

import "github.com/hashicorp/go-version"

// VersionLogMessage represents information about the Terraform version
// and the version of the schema used for the following messages.
// This is a message of type "version".
type VersionLogMessage struct {
	baseLogMessage
	Terraform *version.Version `json:"terraform"`
	UI        *version.Version `json:"ui"`
}

// LogMessage represents a generic human-readable log line
// This is a message of type "log"
type LogMessage struct {
	baseLogMessage
}

// DiagnosticLogMessage represents diagnostic warning or error message.
// This is a message of type "diagnostic"
type DiagnosticLogMessage struct {
	baseLogMessage
	Diagnostic `json:"diagnostic"`
}
