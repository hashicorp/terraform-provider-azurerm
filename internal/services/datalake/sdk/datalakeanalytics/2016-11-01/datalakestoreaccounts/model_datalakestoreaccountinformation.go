package datalakestoreaccounts

type DataLakeStoreAccountInformation struct {
	Id         *string                                    `json:"id,omitempty"`
	Name       *string                                    `json:"name,omitempty"`
	Properties *DataLakeStoreAccountInformationProperties `json:"properties,omitempty"`
	Type       *string                                    `json:"type,omitempty"`
}
