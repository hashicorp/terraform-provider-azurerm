package accounts

type AddDataLakeStoreWithAccountParameters struct {
	Name       string                      `json:"name"`
	Properties *AddDataLakeStoreProperties `json:"properties,omitempty"`
}
