package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FrontDoorProperties struct {
	BackendPools          *[]BackendPool                `json:"backendPools,omitempty"`
	BackendPoolsSettings  *BackendPoolsSettings         `json:"backendPoolsSettings,omitempty"`
	Cname                 *string                       `json:"cname,omitempty"`
	EnabledState          *FrontDoorEnabledState        `json:"enabledState,omitempty"`
	FriendlyName          *string                       `json:"friendlyName,omitempty"`
	FrontdoorId           *string                       `json:"frontdoorId,omitempty"`
	FrontendEndpoints     *[]FrontendEndpoint           `json:"frontendEndpoints,omitempty"`
	HealthProbeSettings   *[]HealthProbeSettingsModel   `json:"healthProbeSettings,omitempty"`
	LoadBalancingSettings *[]LoadBalancingSettingsModel `json:"loadBalancingSettings,omitempty"`
	ProvisioningState     *string                       `json:"provisioningState,omitempty"`
	ResourceState         *FrontDoorResourceState       `json:"resourceState,omitempty"`
	RoutingRules          *[]RoutingRule                `json:"routingRules,omitempty"`
	RulesEngines          *[]RulesEngine                `json:"rulesEngines,omitempty"`
}
