package privateendpointconnections

type PrivateEndpointConnectionProperties struct {
	GroupIds                          *[]string                  `json:"groupIds,omitempty"`
	PrivateEndpoint                   *PrivateEndpoint           `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *ConnectionState           `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *ResourceProvisioningState `json:"provisioningState,omitempty"`
}
