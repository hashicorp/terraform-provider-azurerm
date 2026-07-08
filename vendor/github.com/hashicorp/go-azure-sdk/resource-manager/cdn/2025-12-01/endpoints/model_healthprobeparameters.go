package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthProbeParameters struct {
	ProbeIntervalInSeconds *int64                  `json:"probeIntervalInSeconds,omitempty"`
	ProbePath              *string                 `json:"probePath,omitempty"`
	ProbeProtocol          *ProbeProtocol          `json:"probeProtocol,omitempty"`
	ProbeRequestType       *HealthProbeRequestType `json:"probeRequestType,omitempty"`
}
