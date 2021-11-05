package privateendpointconnections

type PrivateEndpointConnectionProperties struct {
	GroupIds                          *[]string                                   `json:"groupIds,omitempty"`
	PrivateEndpoint                   *PrivateEndpoint                            `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState PrivateLinkServiceConnectionState           `json:"privateLinkServiceConnectionState"`
	ProvisioningState                 *PrivateEndpointConnectionProvisioningState `json:"provisioningState,omitempty"`
}
