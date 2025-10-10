package attachednetworkconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttachedNetworkConnectionProperties struct {
	DomainJoinType            *DomainJoinType    `json:"domainJoinType,omitempty"`
	HealthCheckStatus         *HealthCheckStatus `json:"healthCheckStatus,omitempty"`
	NetworkConnectionId       string             `json:"networkConnectionId"`
	NetworkConnectionLocation *string            `json:"networkConnectionLocation,omitempty"`
	ProvisioningState         *ProvisioningState `json:"provisioningState,omitempty"`
}
