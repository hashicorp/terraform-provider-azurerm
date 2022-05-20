package privateendpointconnection

type PrivateLinkServiceConnectionState struct {
	ActionsRequired *string `json:"actionsRequired,omitempty"`
	Description     *string `json:"description,omitempty"`
	Status          *Status `json:"status,omitempty"`
}
