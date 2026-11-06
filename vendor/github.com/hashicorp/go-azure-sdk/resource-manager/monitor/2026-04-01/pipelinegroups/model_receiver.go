package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Receiver struct {
	Name             string          `json:"name"`
	Otlp             *OtlpReceiver   `json:"otlp,omitempty"`
	Syslog           *SyslogReceiver `json:"syslog,omitempty"`
	TlsConfiguration *string         `json:"tlsConfiguration,omitempty"`
	Type             ReceiverType    `json:"type"`
}
