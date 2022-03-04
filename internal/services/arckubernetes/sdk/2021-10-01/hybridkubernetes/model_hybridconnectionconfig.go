package hybridkubernetes

type HybridConnectionConfig struct {
	ExpirationTime       *int64  `json:"expirationTime,omitempty"`
	HybridConnectionName *string `json:"hybridConnectionName,omitempty"`
	Relay                *string `json:"relay,omitempty"`
	Token                *string `json:"token,omitempty"`
}
