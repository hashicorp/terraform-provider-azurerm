package frontdoors

type HealthProbeSettingsModel struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *HealthProbeSettingsProperties `json:"properties,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
