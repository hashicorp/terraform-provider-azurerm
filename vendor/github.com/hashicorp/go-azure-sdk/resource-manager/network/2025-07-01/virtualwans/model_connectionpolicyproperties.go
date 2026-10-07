package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionPolicyProperties struct {
	AssociatedConnections  *[]string             `json:"associatedConnections,omitempty"`
	EnableInternetSecurity *bool                 `json:"enableInternetSecurity,omitempty"`
	ProvisioningState      *ProvisioningState    `json:"provisioningState,omitempty"`
	RoutingConfiguration   *RoutingConfiguration `json:"routingConfiguration,omitempty"`
}
