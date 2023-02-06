package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FrontendEndpointProperties struct {
	CustomHTTPSConfiguration         *CustomHTTPSConfiguration                                         `json:"customHttpsConfiguration,omitempty"`
	CustomHTTPSProvisioningState     *CustomHTTPSProvisioningState                                     `json:"customHttpsProvisioningState,omitempty"`
	CustomHTTPSProvisioningSubstate  *CustomHTTPSProvisioningSubstate                                  `json:"customHttpsProvisioningSubstate,omitempty"`
	HostName                         *string                                                           `json:"hostName,omitempty"`
	ResourceState                    *FrontDoorResourceState                                           `json:"resourceState,omitempty"`
	SessionAffinityEnabledState      *SessionAffinityEnabledState                                      `json:"sessionAffinityEnabledState,omitempty"`
	SessionAffinityTtlSeconds        *int64                                                            `json:"sessionAffinityTtlSeconds,omitempty"`
	WebApplicationFirewallPolicyLink *FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}
