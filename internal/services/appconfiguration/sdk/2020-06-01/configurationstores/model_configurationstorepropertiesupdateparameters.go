package configurationstores

type ConfigurationStorePropertiesUpdateParameters struct {
	Encryption          *EncryptionProperties `json:"encryption,omitempty"`
	PublicNetworkAccess *PublicNetworkAccess  `json:"publicNetworkAccess,omitempty"`
}
