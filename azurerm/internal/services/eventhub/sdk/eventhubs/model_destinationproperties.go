package eventhubs

type DestinationProperties struct {
	ArchiveNameFormat        *string `json:"archiveNameFormat,omitempty"`
	BlobContainer            *string `json:"blobContainer,omitempty"`
	StorageAccountResourceId *string `json:"storageAccountResourceId,omitempty"`
}
