package account

type AccountPropertiesManagedResources struct {
	EventHubNamespace *string `json:"eventHubNamespace,omitempty"`
	ResourceGroup     *string `json:"resourceGroup,omitempty"`
	StorageAccount    *string `json:"storageAccount,omitempty"`
}
