package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancingRule struct {
	BackendPort      int64         `json:"backendPort"`
	FrontendPort     int64         `json:"frontendPort"`
	LoadDistribution *string       `json:"loadDistribution,omitempty"`
	ProbePort        *int64        `json:"probePort,omitempty"`
	ProbeProtocol    ProbeProtocol `json:"probeProtocol"`
	ProbeRequestPath *string       `json:"probeRequestPath,omitempty"`
	Protocol         Protocol      `json:"protocol"`
}
