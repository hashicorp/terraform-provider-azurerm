package storageaccounts

type StorageContainer struct {
	Id         *string                     `json:"id,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	Properties *StorageContainerProperties `json:"properties,omitempty"`
	Type       *string                     `json:"type,omitempty"`
}
