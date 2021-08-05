package signalr

type SignalRProperties struct {
	Cors                       *SignalRCorsSettings         `json:"cors,omitempty"`
	ExternalIP                 *string                      `json:"externalIP,omitempty"`
	Features                   *[]SignalRFeature            `json:"features,omitempty"`
	HostName                   *string                      `json:"hostName,omitempty"`
	HostNamePrefix             *string                      `json:"hostNamePrefix,omitempty"`
	NetworkACLs                *SignalRNetworkACLs          `json:"networkACLs,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicPort                 *int64                       `json:"publicPort,omitempty"`
	ServerPort                 *int64                       `json:"serverPort,omitempty"`
	Upstream                   *ServerlessUpstreamSettings  `json:"upstream,omitempty"`
	Version                    *string                      `json:"version,omitempty"`
}
