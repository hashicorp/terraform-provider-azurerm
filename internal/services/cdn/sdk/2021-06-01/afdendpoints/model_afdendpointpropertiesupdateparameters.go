package afdendpoints

type AFDEndpointPropertiesUpdateParameters struct {
	EnabledState *EnabledState `json:"enabledState,omitempty"`
	ProfileName  *string       `json:"profileName,omitempty"`
}
