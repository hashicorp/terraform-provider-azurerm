package configurationstores

type ConfigurationStoreProperties struct {
	CreationDate               *string                               `json:"creationDate,omitempty"`
	Encryption                 *EncryptionProperties                 `json:"encryption,omitempty"`
	Endpoint                   *string                               `json:"endpoint,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnectionReference `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState                    `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess                  `json:"publicNetworkAccess,omitempty"`
}
