package configurationstores

type PrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *PrivateEndpoint                  `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState"`
	ProvisioningState                 *ProvisioningState                `json:"provisioningState,omitempty"`
}
