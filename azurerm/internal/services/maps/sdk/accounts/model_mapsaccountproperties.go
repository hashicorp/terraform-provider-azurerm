package accounts

type MapsAccountProperties struct {
	DisableLocalAuth  *bool   `json:"disableLocalAuth,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
	UniqueId          *string `json:"uniqueId,omitempty"`
}
