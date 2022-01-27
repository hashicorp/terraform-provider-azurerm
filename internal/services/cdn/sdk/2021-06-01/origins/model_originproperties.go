package origins

type OriginProperties struct {
	Enabled                    *bool                  `json:"enabled,omitempty"`
	HostName                   string                 `json:"hostName"`
	HttpPort                   *int64                 `json:"httpPort,omitempty"`
	HttpsPort                  *int64                 `json:"httpsPort,omitempty"`
	OriginHostHeader           *string                `json:"originHostHeader,omitempty"`
	Priority                   *int64                 `json:"priority,omitempty"`
	PrivateEndpointStatus      *PrivateEndpointStatus `json:"privateEndpointStatus,omitempty"`
	PrivateLinkAlias           *string                `json:"privateLinkAlias,omitempty"`
	PrivateLinkApprovalMessage *string                `json:"privateLinkApprovalMessage,omitempty"`
	PrivateLinkLocation        *string                `json:"privateLinkLocation,omitempty"`
	PrivateLinkResourceId      *string                `json:"privateLinkResourceId,omitempty"`
	ProvisioningState          *string                `json:"provisioningState,omitempty"`
	ResourceState              *OriginResourceState   `json:"resourceState,omitempty"`
	Weight                     *int64                 `json:"weight,omitempty"`
}
