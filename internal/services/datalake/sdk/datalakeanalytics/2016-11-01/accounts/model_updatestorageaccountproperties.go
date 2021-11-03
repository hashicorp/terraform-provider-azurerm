package accounts

type UpdateStorageAccountProperties struct {
	AccessKey *string `json:"accessKey,omitempty"`
	Suffix    *string `json:"suffix,omitempty"`
}
