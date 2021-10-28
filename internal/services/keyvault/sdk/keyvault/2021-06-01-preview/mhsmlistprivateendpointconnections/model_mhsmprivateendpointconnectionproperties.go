package mhsmlistprivateendpointconnections

type MHSMPrivateEndpointConnectionProperties struct {
	PrivateEndpoint                   *MHSMPrivateEndpoint                        `json:"privateEndpoint,omitempty"`
	PrivateLinkServiceConnectionState *MHSMPrivateLinkServiceConnectionState      `json:"privateLinkServiceConnectionState,omitempty"`
	ProvisioningState                 *PrivateEndpointConnectionProvisioningState `json:"provisioningState,omitempty"`
}
