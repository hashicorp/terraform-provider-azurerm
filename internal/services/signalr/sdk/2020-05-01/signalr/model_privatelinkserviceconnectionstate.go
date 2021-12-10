package signalr

type PrivateLinkServiceConnectionState struct {
	ActionsRequired *string                             `json:"actionsRequired,omitempty"`
	Description     *string                             `json:"description,omitempty"`
	Status          *PrivateLinkServiceConnectionStatus `json:"status,omitempty"`
}
