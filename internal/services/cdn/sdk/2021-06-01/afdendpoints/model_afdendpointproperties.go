package afdendpoints

type AFDEndpointProperties struct {
	DeploymentStatus  *DeploymentStatus     `json:"deploymentStatus,omitempty"`
	EnabledState      *EnabledState         `json:"enabledState,omitempty"`
	HostName          *string               `json:"hostName,omitempty"`
	ProfileName       *string               `json:"profileName,omitempty"`
	ProvisioningState *AfdProvisioningState `json:"provisioningState,omitempty"`
}
