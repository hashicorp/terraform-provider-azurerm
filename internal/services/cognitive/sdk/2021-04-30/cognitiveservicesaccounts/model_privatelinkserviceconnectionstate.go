package cognitiveservicesaccounts

type PrivateLinkServiceConnectionState struct {
	ActionsRequired *string                                 `json:"actionsRequired,omitempty"`
	Description     *string                                 `json:"description,omitempty"`
	Status          *PrivateEndpointServiceConnectionStatus `json:"status,omitempty"`
}
