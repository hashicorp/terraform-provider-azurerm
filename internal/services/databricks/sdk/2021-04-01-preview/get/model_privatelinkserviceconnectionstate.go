package get

type PrivateLinkServiceConnectionState struct {
	ActionRequired *string                            `json:"actionRequired,omitempty"`
	Description    *string                            `json:"description,omitempty"`
	Status         PrivateLinkServiceConnectionStatus `json:"status"`
}
