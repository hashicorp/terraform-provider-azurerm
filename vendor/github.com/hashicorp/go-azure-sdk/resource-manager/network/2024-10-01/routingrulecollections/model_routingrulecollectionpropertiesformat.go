package routingrulecollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingRuleCollectionPropertiesFormat struct {
	AppliesTo                  []NetworkManagerRoutingGroupItem `json:"appliesTo"`
	Description                *string                          `json:"description,omitempty"`
	DisableBgpRoutePropagation *DisableBgpRoutePropagation      `json:"disableBgpRoutePropagation,omitempty"`
	ProvisioningState          *ProvisioningState               `json:"provisioningState,omitempty"`
	ResourceGuid               *string                          `json:"resourceGuid,omitempty"`
}
