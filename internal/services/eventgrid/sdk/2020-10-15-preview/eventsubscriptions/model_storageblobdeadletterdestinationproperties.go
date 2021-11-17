package eventsubscriptions

type StorageBlobDeadLetterDestinationProperties struct {
	BlobContainerName *string `json:"blobContainerName,omitempty"`
	ResourceId        *string `json:"resourceId,omitempty"`
}
