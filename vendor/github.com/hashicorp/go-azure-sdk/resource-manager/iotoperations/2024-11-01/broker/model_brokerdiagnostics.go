package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerDiagnostics struct {
	Logs      *DiagnosticsLogs `json:"logs,omitempty"`
	Metrics   *Metrics         `json:"metrics,omitempty"`
	SelfCheck *SelfCheck       `json:"selfCheck,omitempty"`
	Traces    *Traces          `json:"traces,omitempty"`
}
