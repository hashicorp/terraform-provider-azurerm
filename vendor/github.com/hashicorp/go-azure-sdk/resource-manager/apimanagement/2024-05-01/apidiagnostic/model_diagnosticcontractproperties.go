package apidiagnostic

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticContractProperties struct {
	AlwaysLog               *AlwaysLog                  `json:"alwaysLog,omitempty"`
	Backend                 *PipelineDiagnosticSettings `json:"backend,omitempty"`
	Frontend                *PipelineDiagnosticSettings `json:"frontend,omitempty"`
	HTTPCorrelationProtocol *HTTPCorrelationProtocol    `json:"httpCorrelationProtocol,omitempty"`
	LogClientIP             *bool                       `json:"logClientIp,omitempty"`
	LoggerId                string                      `json:"loggerId"`
	Metrics                 *bool                       `json:"metrics,omitempty"`
	OperationNameFormat     *OperationNameFormat        `json:"operationNameFormat,omitempty"`
	Sampling                *SamplingSettings           `json:"sampling,omitempty"`
	Verbosity               *Verbosity                  `json:"verbosity,omitempty"`
}
