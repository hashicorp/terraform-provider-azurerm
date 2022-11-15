package frontdoors

type HealthProbeSettingsProperties struct {
	EnabledState      *HealthProbeEnabled         `json:"enabledState,omitempty"`
	HealthProbeMethod *FrontDoorHealthProbeMethod `json:"healthProbeMethod,omitempty"`
	IntervalInSeconds *int64                      `json:"intervalInSeconds,omitempty"`
	Path              *string                     `json:"path,omitempty"`
	Protocol          *FrontDoorProtocol          `json:"protocol,omitempty"`
	ResourceState     *FrontDoorResourceState     `json:"resourceState,omitempty"`
}
