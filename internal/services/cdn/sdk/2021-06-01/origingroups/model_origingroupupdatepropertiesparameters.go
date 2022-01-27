package origingroups

type OriginGroupUpdatePropertiesParameters struct {
	HealthProbeSettings                                   *HealthProbeParameters                       `json:"healthProbeSettings,omitempty"`
	Origins                                               *[]ResourceReference                         `json:"origins,omitempty"`
	ResponseBasedOriginErrorDetectionSettings             *ResponseBasedOriginErrorDetectionParameters `json:"responseBasedOriginErrorDetectionSettings,omitempty"`
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int64                                       `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
}
