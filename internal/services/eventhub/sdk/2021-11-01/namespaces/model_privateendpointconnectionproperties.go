package namespaces

type PrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *PrivateEndpoint           `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *ConnectionState           `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *EndPointProvisioningState `json:"provisioningState,omitempty"`
}
