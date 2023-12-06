package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentInstance struct {
	DiscoveryStatus *string `json:"discoveryStatus,omitempty"`
	Name            *string `json:"name,omitempty"`
	Reason          *string `json:"reason,omitempty"`
	StartTime       *string `json:"startTime,omitempty"`
	Status          *string `json:"status,omitempty"`
	Zone            *string `json:"zone,omitempty"`
}
