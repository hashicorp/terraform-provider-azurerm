package eventhubs

type DestinationProperties struct {
	ArchiveNameFormat        *string `json:"archiveNameFormat,omitempty"`
	BlobContainer            *string `json:"blobContainer,omitempty"`
	DataLakeAccountName      *string `json:"dataLakeAccountName,omitempty"`
	DataLakeFolderPath       *string `json:"dataLakeFolderPath,omitempty"`
	DataLakeSubscriptionId   *string `json:"dataLakeSubscriptionId,omitempty"`
	StorageAccountResourceId *string `json:"storageAccountResourceId,omitempty"`
}
