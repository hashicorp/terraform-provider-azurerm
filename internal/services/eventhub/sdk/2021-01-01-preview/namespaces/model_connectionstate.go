package namespaces

type ConnectionState struct {
	Description *string                      `json:"description,omitempty"`
	Status      *PrivateLinkConnectionStatus `json:"status,omitempty"`
}
