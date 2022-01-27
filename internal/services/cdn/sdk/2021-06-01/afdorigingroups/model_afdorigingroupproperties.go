package afdorigingroups

type AFDOriginGroupProperties struct {
	DeploymentStatus                                      *DeploymentStatus                            `json:"deploymentStatus,omitempty"`
	HealthProbeSettings                                   *HealthProbeParameters                       `json:"healthProbeSettings,omitempty"`
	LoadBalancingSettings                                 *LoadBalancingSettingsParameters             `json:"loadBalancingSettings,omitempty"`
	ProfileName                                           *string                                      `json:"profileName,omitempty"`
	ProvisioningState                                     *AfdProvisioningState                        `json:"provisioningState,omitempty"`
	ResponseBasedAfdOriginErrorDetectionSettings          *ResponseBasedOriginErrorDetectionParameters `json:"responseBasedAfdOriginErrorDetectionSettings,omitempty"`
	SessionAffinityState                                  *EnabledState                                `json:"sessionAffinityState,omitempty"`
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int64                                       `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
}
