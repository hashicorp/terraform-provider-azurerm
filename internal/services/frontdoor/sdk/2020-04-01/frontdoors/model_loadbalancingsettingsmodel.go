package frontdoors

type LoadBalancingSettingsModel struct {
	Id         *string                          `json:"id,omitempty"`
	Name       *string                          `json:"name,omitempty"`
	Properties *LoadBalancingSettingsProperties `json:"properties,omitempty"`
	Type       *string                          `json:"type,omitempty"`
}
