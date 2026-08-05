package routefilterrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteFilterRulePropertiesFormat struct {
	Access              Access              `json:"access"`
	Communities         []string            `json:"communities"`
	ProvisioningState   *ProvisioningState  `json:"provisioningState,omitempty"`
	RouteFilterRuleType RouteFilterRuleType `json:"routeFilterRuleType"`
}
