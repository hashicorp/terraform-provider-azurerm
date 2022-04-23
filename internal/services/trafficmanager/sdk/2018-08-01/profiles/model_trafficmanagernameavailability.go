package profiles

type TrafficManagerNameAvailability struct {
	Message       *string `json:"message,omitempty"`
	Name          *string `json:"name,omitempty"`
	NameAvailable *bool   `json:"nameAvailable,omitempty"`
	Reason        *string `json:"reason,omitempty"`
	Type          *string `json:"type,omitempty"`
}
