package routingrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingRulePropertiesFormat struct {
	Description       *string                     `json:"description,omitempty"`
	Destination       RoutingRuleRouteDestination `json:"destination"`
	NextHop           RoutingRuleNextHop          `json:"nextHop"`
	ProvisioningState *ProvisioningState          `json:"provisioningState,omitempty"`
	ResourceGuid      *string                     `json:"resourceGuid,omitempty"`
}
