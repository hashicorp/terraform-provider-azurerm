package privateendpointconnections

type PrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *PrivateEndpoint                            `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState PrivateLinkServiceConnectionState           `json:"privateLinkServiceConnectionState"`
	ProvisioningState                 *PrivateEndpointConnectionProvisioningState `json:"provisioningState,omitempty"`
}
