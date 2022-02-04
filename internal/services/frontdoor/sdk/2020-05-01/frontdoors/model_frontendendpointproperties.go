package frontdoors

type FrontendEndpointProperties struct {
	CustomHttpsConfiguration         *CustomHttpsConfiguration                                         `json:"customHttpsConfiguration,omitempty"`
	CustomHttpsProvisioningState     *CustomHttpsProvisioningState                                     `json:"customHttpsProvisioningState,omitempty"`
	CustomHttpsProvisioningSubstate  *CustomHttpsProvisioningSubstate                                  `json:"customHttpsProvisioningSubstate,omitempty"`
	HostName                         *string                                                           `json:"hostName,omitempty"`
	ResourceState                    *FrontDoorResourceState                                           `json:"resourceState,omitempty"`
	SessionAffinityEnabledState      *SessionAffinityEnabledState                                      `json:"sessionAffinityEnabledState,omitempty"`
	SessionAffinityTtlSeconds        *int64                                                            `json:"sessionAffinityTtlSeconds,omitempty"`
	WebApplicationFirewallPolicyLink *FrontendEndpointUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}
