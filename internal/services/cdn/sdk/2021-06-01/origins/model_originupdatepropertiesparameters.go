package origins

type OriginUpdatePropertiesParameters struct {
	Enabled                    *bool   `json:"enabled,omitempty"`
	HostName                   *string `json:"hostName,omitempty"`
	HttpPort                   *int64  `json:"httpPort,omitempty"`
	HttpsPort                  *int64  `json:"httpsPort,omitempty"`
	OriginHostHeader           *string `json:"originHostHeader,omitempty"`
	Priority                   *int64  `json:"priority,omitempty"`
	PrivateLinkAlias           *string `json:"privateLinkAlias,omitempty"`
	PrivateLinkApprovalMessage *string `json:"privateLinkApprovalMessage,omitempty"`
	PrivateLinkLocation        *string `json:"privateLinkLocation,omitempty"`
	PrivateLinkResourceId      *string `json:"privateLinkResourceId,omitempty"`
	Weight                     *int64  `json:"weight,omitempty"`
}
