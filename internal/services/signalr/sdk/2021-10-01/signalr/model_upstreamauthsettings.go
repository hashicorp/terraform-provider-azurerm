package signalr

type UpstreamAuthSettings struct {
	ManagedIdentity *ManagedIdentitySettings `json:"managedIdentity,omitempty"`
	Type            *UpstreamAuthType        `json:"type,omitempty"`
}
