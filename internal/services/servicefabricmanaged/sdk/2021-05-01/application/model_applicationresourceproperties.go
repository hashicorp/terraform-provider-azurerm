package application

type ApplicationResourceProperties struct {
	ManagedIdentities *[]ApplicationUserAssignedIdentity `json:"managedIdentities,omitempty"`
	Parameters        *map[string]string                 `json:"parameters,omitempty"`
	ProvisioningState *string                            `json:"provisioningState,omitempty"`
	UpgradePolicy     *ApplicationUpgradePolicy          `json:"upgradePolicy,omitempty"`
	Version           *string                            `json:"version,omitempty"`
}
