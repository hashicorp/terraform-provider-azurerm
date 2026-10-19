package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingIntentProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	RoutingPolicies   *[]RoutingPolicy   `json:"routingPolicies,omitempty"`
}
