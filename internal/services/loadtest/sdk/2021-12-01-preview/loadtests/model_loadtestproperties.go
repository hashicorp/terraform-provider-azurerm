package loadtests

type LoadTestProperties struct {
	DataPlaneURI      *string        `json:"dataPlaneURI,omitempty"`
	Description       *string        `json:"description,omitempty"`
	ProvisioningState *ResourceState `json:"provisioningState,omitempty"`
}
