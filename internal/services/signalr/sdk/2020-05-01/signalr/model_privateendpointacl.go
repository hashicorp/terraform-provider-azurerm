package signalr

type PrivateEndpointACL struct {
	Allow *[]SignalRRequestType `json:"allow,omitempty"`
	Deny  *[]SignalRRequestType `json:"deny,omitempty"`
	Name  string                `json:"name"`
}
