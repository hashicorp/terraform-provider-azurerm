package accounts

type DataLakeStoreAccountPropertiesBasic struct {
	AccountId         *string                     `json:"accountId,omitempty"`
	CreationTime      *string                     `json:"creationTime,omitempty"`
	Endpoint          *string                     `json:"endpoint,omitempty"`
	LastModifiedTime  *string                     `json:"lastModifiedTime,omitempty"`
	ProvisioningState *DataLakeStoreAccountStatus `json:"provisioningState,omitempty"`
	State             *DataLakeStoreAccountState  `json:"state,omitempty"`
}
