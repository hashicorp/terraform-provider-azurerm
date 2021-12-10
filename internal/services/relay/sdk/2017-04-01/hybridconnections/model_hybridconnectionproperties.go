package hybridconnections

type HybridConnectionProperties struct {
	CreatedAt                   *string `json:"createdAt,omitempty"`
	ListenerCount               *int64  `json:"listenerCount,omitempty"`
	RequiresClientAuthorization *bool   `json:"requiresClientAuthorization,omitempty"`
	UpdatedAt                   *string `json:"updatedAt,omitempty"`
	UserMetadata                *string `json:"userMetadata,omitempty"`
}
