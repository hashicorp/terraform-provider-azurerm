package wcfrelays

type WcfRelayProperties struct {
	CreatedAt                   *string    `json:"createdAt,omitempty"`
	IsDynamic                   *bool      `json:"isDynamic,omitempty"`
	ListenerCount               *int64     `json:"listenerCount,omitempty"`
	RelayType                   *Relaytype `json:"relayType,omitempty"`
	RequiresClientAuthorization *bool      `json:"requiresClientAuthorization,omitempty"`
	RequiresTransportSecurity   *bool      `json:"requiresTransportSecurity,omitempty"`
	UpdatedAt                   *string    `json:"updatedAt,omitempty"`
	UserMetadata                *string    `json:"userMetadata,omitempty"`
}
