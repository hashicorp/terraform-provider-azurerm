package managedapis

type ConnectionParameter struct {
	OAuthSettings *ApiOAuthSettings        `json:"oAuthSettings,omitempty"`
	Type          *ConnectionParameterType `json:"type,omitempty"`
}
