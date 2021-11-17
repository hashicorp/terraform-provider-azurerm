package exports

type ExportDeliveryDestination struct {
	Container      string  `json:"container"`
	ResourceId     *string `json:"resourceId,omitempty"`
	RootFolderPath *string `json:"rootFolderPath,omitempty"`
	SasToken       *string `json:"sasToken,omitempty"`
	StorageAccount *string `json:"storageAccount,omitempty"`
}
