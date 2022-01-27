package origingroups

type OriginGroupProperties struct {
	HealthProbeSettings                                   *HealthProbeParameters                       `json:"healthProbeSettings,omitempty"`
	Origins                                               []ResourceReference                          `json:"origins"`
	ProvisioningState                                     *string                                      `json:"provisioningState,omitempty"`
	ResourceState                                         *OriginGroupResourceState                    `json:"resourceState,omitempty"`
	ResponseBasedOriginErrorDetectionSettings             *ResponseBasedOriginErrorDetectionParameters `json:"responseBasedOriginErrorDetectionSettings,omitempty"`
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int64                                       `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
}
