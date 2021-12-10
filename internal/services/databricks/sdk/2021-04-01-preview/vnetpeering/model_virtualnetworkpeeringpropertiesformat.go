package vnetpeering

type VirtualNetworkPeeringPropertiesFormat struct {
	AllowForwardedTraffic     *bool                                                          `json:"allowForwardedTraffic,omitempty"`
	AllowGatewayTransit       *bool                                                          `json:"allowGatewayTransit,omitempty"`
	AllowVirtualNetworkAccess *bool                                                          `json:"allowVirtualNetworkAccess,omitempty"`
	DatabricksAddressSpace    *AddressSpace                                                  `json:"databricksAddressSpace,omitempty"`
	DatabricksVirtualNetwork  *VirtualNetworkPeeringPropertiesFormatDatabricksVirtualNetwork `json:"databricksVirtualNetwork,omitempty"`
	PeeringState              *PeeringState                                                  `json:"peeringState,omitempty"`
	ProvisioningState         *PeeringProvisioningState                                      `json:"provisioningState,omitempty"`
	RemoteAddressSpace        *AddressSpace                                                  `json:"remoteAddressSpace,omitempty"`
	RemoteVirtualNetwork      VirtualNetworkPeeringPropertiesFormatRemoteVirtualNetwork      `json:"remoteVirtualNetwork"`
	UseRemoteGateways         *bool                                                          `json:"useRemoteGateways,omitempty"`
}
