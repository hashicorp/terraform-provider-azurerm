package accounts

type DataLakeAnalyticsAccountPropertiesBasic struct {
	AccountId         *string                         `json:"accountId,omitempty"`
	CreationTime      *string                         `json:"creationTime,omitempty"`
	Endpoint          *string                         `json:"endpoint,omitempty"`
	LastModifiedTime  *string                         `json:"lastModifiedTime,omitempty"`
	ProvisioningState *DataLakeAnalyticsAccountStatus `json:"provisioningState,omitempty"`
	State             *DataLakeAnalyticsAccountState  `json:"state,omitempty"`
}
