package vnetpeering

type VirtualNetworkPeering struct {
	Id         *string                               `json:"id,omitempty"`
	Name       *string                               `json:"name,omitempty"`
	Properties VirtualNetworkPeeringPropertiesFormat `json:"properties"`
	Type       *string                               `json:"type,omitempty"`
}
