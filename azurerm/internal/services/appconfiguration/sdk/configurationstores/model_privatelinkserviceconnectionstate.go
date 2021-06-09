package configurationstores

type PrivateLinkServiceConnectionState struct {
	ActionsRequired *ActionsRequired  `json:"actionsRequired,omitempty"`
	Description     *string           `json:"description,omitempty"`
	Status          *ConnectionStatus `json:"status,omitempty"`
}
