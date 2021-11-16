package accounts

type UpdateDataLakeStoreAccountParameters struct {
	Properties *UpdateDataLakeStoreAccountProperties `json:"properties,omitempty"`
	Tags       *map[string]string                    `json:"tags,omitempty"`
}
