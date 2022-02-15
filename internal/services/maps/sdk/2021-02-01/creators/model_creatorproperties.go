package creators

type CreatorProperties struct {
	ProvisioningState *string `json:"provisioningState,omitempty"`
	StorageUnits      int64   `json:"storageUnits"`
}
