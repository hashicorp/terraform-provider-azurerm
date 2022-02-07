package storageaccounts

type UpdateStorageAccountProperties struct {
	AccessKey *string `json:"accessKey,omitempty"`
	Suffix    *string `json:"suffix,omitempty"`
}
