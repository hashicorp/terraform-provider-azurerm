package videoanalyzer

type VideoAnalyzerPropertiesUpdate struct {
	Encryption      *AccountEncryption `json:"encryption,omitempty"`
	Endpoints       *[]Endpoint        `json:"endpoints,omitempty"`
	StorageAccounts *[]StorageAccount  `json:"storageAccounts,omitempty"`
}
