package accounts

type NameAvailabilityInformation struct {
	Message       *string `json:"message,omitempty"`
	NameAvailable *bool   `json:"nameAvailable,omitempty"`
	Reason        *string `json:"reason,omitempty"`
}
