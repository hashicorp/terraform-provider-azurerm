package rulesets

type RuleSetProperties struct {
	DeploymentStatus  *DeploymentStatus     `json:"deploymentStatus,omitempty"`
	ProfileName       *string               `json:"profileName,omitempty"`
	ProvisioningState *AfdProvisioningState `json:"provisioningState,omitempty"`
}
