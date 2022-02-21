package proxy

type ServicesNameAvailabilityInfo struct {
	Message       *string                          `json:"message,omitempty"`
	NameAvailable *bool                            `json:"nameAvailable,omitempty"`
	Reason        *ServiceNameUnavailabilityReason `json:"reason,omitempty"`
}
