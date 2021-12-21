package accounts

type UpdateDataLakeStoreWithAccountParameters struct {
	Name       string                         `json:"name"`
	Properties *UpdateDataLakeStoreProperties `json:"properties,omitempty"`
}
