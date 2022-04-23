package connections

type ConnectionStatusDefinition struct {
	Error  *ConnectionError `json:"error,omitempty"`
	Status *string          `json:"status,omitempty"`
	Target *string          `json:"target,omitempty"`
}
