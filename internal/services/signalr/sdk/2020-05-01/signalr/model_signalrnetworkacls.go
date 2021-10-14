package signalr

type SignalRNetworkACLs struct {
	DefaultAction    *ACLAction            `json:"defaultAction,omitempty"`
	PrivateEndpoints *[]PrivateEndpointACL `json:"privateEndpoints,omitempty"`
	PublicNetwork    *NetworkACL           `json:"publicNetwork,omitempty"`
}
