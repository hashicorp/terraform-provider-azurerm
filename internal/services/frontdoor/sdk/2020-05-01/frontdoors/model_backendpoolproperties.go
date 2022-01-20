package frontdoors

type BackendPoolProperties struct {
	Backends              *[]Backend              `json:"backends,omitempty"`
	HealthProbeSettings   *SubResource            `json:"healthProbeSettings,omitempty"`
	LoadBalancingSettings *SubResource            `json:"loadBalancingSettings,omitempty"`
	ResourceState         *FrontDoorResourceState `json:"resourceState,omitempty"`
}
