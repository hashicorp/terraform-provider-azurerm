package labplan

type ConnectionProfile struct {
	ClientRdpAccess *ConnectionType `json:"clientRdpAccess,omitempty"`
	ClientSshAccess *ConnectionType `json:"clientSshAccess,omitempty"`
	WebRdpAccess    *ConnectionType `json:"webRdpAccess,omitempty"`
	WebSshAccess    *ConnectionType `json:"webSshAccess,omitempty"`
}
