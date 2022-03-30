package signalr

type PrivateEndpointConnectionProperties struct {
	GroupIds                          *[]string                          `json:"groupIds,omitempty"`
	PrivateEndpoint                   *PrivateEndpoint                   `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *ProvisioningState                 `json:"provisioningState,omitempty"`
}
