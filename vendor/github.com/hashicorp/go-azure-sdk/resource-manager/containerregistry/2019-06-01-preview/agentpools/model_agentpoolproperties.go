package agentpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolProperties struct {
	Count                          *int64             `json:"count,omitempty"`
	Os                             *OS                `json:"os,omitempty"`
	ProvisioningState              *ProvisioningState `json:"provisioningState,omitempty"`
	Tier                           *string            `json:"tier,omitempty"`
	VirtualNetworkSubnetResourceId *string            `json:"virtualNetworkSubnetResourceId,omitempty"`
}
