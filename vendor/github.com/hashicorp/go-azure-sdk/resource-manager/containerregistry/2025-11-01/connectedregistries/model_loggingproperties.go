package connectedregistries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoggingProperties struct {
	AuditLogStatus *AuditLogStatus `json:"auditLogStatus,omitempty"`
	LogLevel       *LogLevel       `json:"logLevel,omitempty"`
}
