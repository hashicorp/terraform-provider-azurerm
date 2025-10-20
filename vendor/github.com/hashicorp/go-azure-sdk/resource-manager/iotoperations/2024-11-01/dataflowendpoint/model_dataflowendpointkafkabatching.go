package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointKafkaBatching struct {
	LatencyMs   *int64           `json:"latencyMs,omitempty"`
	MaxBytes    *int64           `json:"maxBytes,omitempty"`
	MaxMessages *int64           `json:"maxMessages,omitempty"`
	Mode        *OperationalMode `json:"mode,omitempty"`
}
