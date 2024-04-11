package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubRouteTableProperties struct {
	AssociatedConnections  *[]string          `json:"associatedConnections,omitempty"`
	Labels                 *[]string          `json:"labels,omitempty"`
	PropagatingConnections *[]string          `json:"propagatingConnections,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
	Routes                 *[]HubRoute        `json:"routes,omitempty"`
}
