package signalr

type NetworkACL struct {
	Allow *[]SignalRRequestType `json:"allow,omitempty"`
	Deny  *[]SignalRRequestType `json:"deny,omitempty"`
}
