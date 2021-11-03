package accounts

type StorageAccountInformation struct {
	Id         *string                              `json:"id,omitempty"`
	Name       *string                              `json:"name,omitempty"`
	Properties *StorageAccountInformationProperties `json:"properties,omitempty"`
	Type       *string                              `json:"type,omitempty"`
}
