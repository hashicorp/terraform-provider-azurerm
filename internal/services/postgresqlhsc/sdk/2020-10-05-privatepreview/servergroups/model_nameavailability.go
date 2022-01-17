package servergroups

type NameAvailability struct {
	Message       *string `json:"message,omitempty"`
	Name          *string `json:"name,omitempty"`
	NameAvailable *bool   `json:"nameAvailable,omitempty"`
	Type          *string `json:"type,omitempty"`
}
