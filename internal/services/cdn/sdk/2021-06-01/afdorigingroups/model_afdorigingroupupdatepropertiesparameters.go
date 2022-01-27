package afdorigingroups

type AFDOriginGroupUpdatePropertiesParameters struct {
	HealthProbeSettings                                   *HealthProbeParameters                       `json:"healthProbeSettings,omitempty"`
	LoadBalancingSettings                                 *LoadBalancingSettingsParameters             `json:"loadBalancingSettings,omitempty"`
	ProfileName                                           *string                                      `json:"profileName,omitempty"`
	ResponseBasedAfdOriginErrorDetectionSettings          *ResponseBasedOriginErrorDetectionParameters `json:"responseBasedAfdOriginErrorDetectionSettings,omitempty"`
	SessionAffinityState                                  *EnabledState                                `json:"sessionAffinityState,omitempty"`
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int64                                       `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
}
