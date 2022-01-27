package endpoints

type DeepCreatedOriginGroupProperties struct {
	HealthProbeSettings                                   *HealthProbeParameters                       `json:"healthProbeSettings,omitempty"`
	Origins                                               []ResourceReference                          `json:"origins"`
	ResponseBasedOriginErrorDetectionSettings             *ResponseBasedOriginErrorDetectionParameters `json:"responseBasedOriginErrorDetectionSettings,omitempty"`
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int64                                       `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
}
