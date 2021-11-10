package storageaccounts

type AddStorageAccountProperties struct {
	AccessKey string  `json:"accessKey"`
	Suffix    *string `json:"suffix,omitempty"`
}
