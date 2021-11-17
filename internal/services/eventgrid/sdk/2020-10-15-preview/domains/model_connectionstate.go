package domains

type ConnectionState struct {
	ActionsRequired *string                    `json:"actionsRequired,omitempty"`
	Description     *string                    `json:"description,omitempty"`
	Status          *PersistedConnectionStatus `json:"status,omitempty"`
}
