package accounts

type UpdateStorageAccountWithAccountParameters struct {
	Name       string                          `json:"name"`
	Properties *UpdateStorageAccountProperties `json:"properties,omitempty"`
}
