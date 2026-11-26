package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGatewayMigrationStatus struct {
	ErrorMessage *string                              `json:"errorMessage,omitempty"`
	Phase        *VirtualNetworkGatewayMigrationPhase `json:"phase,omitempty"`
	State        *VirtualNetworkGatewayMigrationState `json:"state,omitempty"`
}
